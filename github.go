package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/pkg/errors"
)

type github struct {
	Token string
}

func (gh *github) auth(req *http.Request) {
	if gh.Token != "" {
		req.Header.Set("Authorization", "token "+gh.Token)
	}
}

func (gh *github) DispatchRepositoryEvent(ctx context.Context, org, repo string, eventType string) error {
	endpoint := fmt.Sprintf("https://api.github.com/repos/%s/%s/dispatches",
		org,
		repo,
	)

	payload, err := json.Marshal(struct {
		EventType string `json:"event_type"`
	}{
		EventType: eventType,
	})

	if err != nil {
		return errors.Wrap(err, "marshal payload")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(payload))
	if err != nil {
		return errors.Wrap(err, "build request")
	}

	gh.auth(req)
	req.Header.Set("Accept", "application/vnd.github.everest-preview+json")
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "do request")
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return errors.Wrap(err, "read body")
	}

	log.Printf("%s", string(body))

	if res.StatusCode != http.StatusNoContent {
		var apiErr struct {
			Message string `json:"message"`
		}

		if err := json.Unmarshal(body, &apiErr); err != nil {
			return errors.Wrap(err, "unmarshal error")
		}

		return errors.New(apiErr.Message)
	}

	return nil
}

func (gh *github) GetLastCommitSHA(ctx context.Context, org, repo, branch string) (string, error) {
	endpoint := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits/%s",
		org,
		repo,
		branch,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return "", errors.Wrap(err, "build request")
	}

	gh.auth(req)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "do request")
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", errors.Wrap(err, "read body")
	}

	if res.StatusCode != http.StatusOK {
		var apiErr struct {
			Message string `json:"message"`
		}

		if err := json.Unmarshal(body, &apiErr); err != nil {
			return "", errors.Wrap(err, "unmarshal error")
		}

		return "", errors.New(apiErr.Message)
	}

	var response struct {
		SHA string `json:"sha"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return "", errors.Wrap(err, "umarshal body")
	}

	return response.SHA, nil
}
