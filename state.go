package main

import (
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
)

type state struct {
	Path string
}

func (s *state) Set(sha string) error {
	if err := ioutil.WriteFile(s.Path, []byte(sha), os.ModePerm); err != nil {
		return errors.Wrap(err, "write file")
	}

	return nil
}

func (s *state) Get() (*string, error) {
	data, err := ioutil.ReadFile(s.Path)
	if os.IsNotExist(err) {
		return nil, nil
	} else if err != nil {
		return nil, errors.Wrap(err, "read file")
	}

	sha := string(data)

	return &sha, nil
}
