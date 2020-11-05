package main

import (
	"fmt"
	"time"
)

type config struct {
	SrcOrg    string `long:"src-org" description:"organization of source github repo" required:"true" env:"SRC_ORG"`
	SrcRepo   string `long:"src-repo" description:"source github repo name" required:"true" env:"SRC_REPO"`
	SrcBranch string `long:"src-branch" description:"source github repo branch name" default:"master" env:"SRC_BRANCH"`

	DstOrg   string `long:"dst-org" description:"organization of destination repository" required:"true" env:"DST_ORG"`
	DstRepo  string `long:"dst-repo" description:"destionation repository name" required:"true" env:"DST_REPO"`
	DstEvent string `long:"dst-event" description:"event type to dispatch" default:"upstream_commit"`

	Token string `long:"token" description:"github token to trigger build and fetch last commits from src" required:"true" env:"TOKEN"`

	State  string        `long:"state" description:"path to file for store state" env:"STATE"`
	Period time.Duration `long:"period" description:"delay between check of new commits" default:"1s" env:"PERIOD"`
}

func (cfg *config) Calculate() {
	if cfg.State == "" {
		cfg.State = fmt.Sprintf("%s-%s-%s.sha", cfg.SrcOrg, cfg.SrcRepo, cfg.SrcBranch)
	}
}
