package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/pkg/errors"
)

var revision = "unknown"

func main() {

	cfg := parseConfig()

	cfg.Calculate()

	_, err := flags.ParseArgs(&cfg, os.Args)
	if err != nil {
		log.Fatalf("parse flags: %v", err)
	}

	log.Printf("start github-actions-upstream-watch, version %s, state=%s", revision, cfg.State)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sig := make(chan os.Signal, 1)

		signal.Notify(sig, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		<-sig

		cancel()
	}()

	if err := run(ctx, cfg); err != nil {
		log.Fatalf("runtime error: %v", err)
	}
}

func parseConfig() config {
	cfg := config{}
	parser := flags.NewParser(&cfg, flags.Default)

	if _, err := parser.Parse(); err != nil {
		switch flagsErr := err.(type) {
		case flags.ErrorType:
			if flagsErr == flags.ErrHelp {
				os.Exit(0)
			}
			os.Exit(1)
		default:
			os.Exit(1)
		}
	}

	return cfg
}

func run(ctx context.Context, cfg config) error {
	ticker := time.NewTicker(cfg.Period)
	defer ticker.Stop()

	gh := &github{Token: cfg.Token}
	state := &state{Path: cfg.State}

	do := func() {
		if err := step(ctx, gh, state, cfg); err != nil {
			log.Printf("error in iteration: %v", err)
		}
	}

	do()
	for {
		select {
		case <-ticker.C:
			do()
		case <-ctx.Done():
			log.Println("shutdown")
			return nil
		}
	}
}

func step(ctx context.Context, gh *github, state *state, cfg config) error {
	log.Printf("check for new commits in %s/%s@%s", cfg.SrcOrg, cfg.SrcRepo, cfg.SrcBranch)

	remoteVersion, err := gh.GetLastCommitSHA(ctx, cfg.SrcOrg, cfg.SrcRepo, cfg.SrcBranch)
	if err != nil {
		return errors.Wrap(err, "get last commit sha")
	}

	localVersion, err := state.Get()
	if err != nil {
		return errors.Wrap(err, "get local commit sha")
	}

	// if local version is unknown, store current and go next
	if localVersion == nil {
		log.Printf("local version is undefined, store %s as current version", remoteVersion)
		if err := state.Set(remoteVersion); err != nil {
			return errors.Wrap(err, "set last version")
		}
		return nil
	}

	if remoteVersion == *localVersion {
		return nil
	}

	log.Printf("find new commit, current local: %s, current remote: %s", *localVersion, remoteVersion)

	log.Printf("dispatch event '%s' to %s/%s", cfg.DstEvent, cfg.DstOrg, cfg.DstRepo)
	if err := gh.DispatchRepositoryEvent(ctx, cfg.DstOrg, cfg.DstRepo, cfg.DstEvent); err != nil {
		return errors.Wrap(err, "dispatch event")
	}

	log.Printf("store %s as last commit to %s", remoteVersion, cfg.State)
	if err := state.Set(remoteVersion); err != nil {
		return errors.Wrap(err, "set last version")
	}

	return nil
}
