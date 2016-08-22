package main

import (
	"encoding/json"
	"io"

	potentials "github.com/artisresistance/githubpotentials"
)

type config struct {
	Token           string
	OutputPath      string
	OutCount        int
	FetchPagesCount int
}

func loadConfig(r io.ReadCloser) (config, error) {
	defer r.Close()
	conf := config{}

	err := json.NewDecoder(r).Decode(&conf)

	return conf, err
}

type result struct {
	Metadata       meta
	ByCommits      []potentials.Repository
	ByStars        []potentials.Repository
	ByContributors []potentials.Repository
}

type meta struct {
	UpdatedUnix int64
	APICalls    int
	Errors      int
	DurationSec int
	ResetUnix   int64
}

func (r result) Write(wc io.WriteCloser) error {
	err := json.NewEncoder(wc).Encode(r)
	return err
}