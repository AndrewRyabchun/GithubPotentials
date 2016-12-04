package main

import (
	"fmt"
	potentials "github.com/artisresistance/githubpotentials"
	"log"
	"os"
	"sync"
	"time"
)

const configPath = "config.json"

var conf config
var startTime time.Time

func init() {
	startTime = time.Now()

	file, err := os.Open(configPath)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	conf, err = loadConfig(file)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}

func main() {
	updatedFrom := time.Now().AddDate(0, 0, -2)

	logger := log.New(os.Stdout, "ghp: ", log.Ltime|log.Lshortfile)
	client := potentials.New(conf.Token, updatedFrom, logger)

	errCount := 0
	onError := func(err error) {
		errCount++
		fmt.Fprintln(os.Stderr, err.Error())
	}

	it := client.Search(conf.FetchPagesCount, onError)
	demuxed := client.
		CountStats(it, onError).
		Split(3)

	joiner := new(sync.WaitGroup)
	collected := make([]potentials.RepositoryCollection, 3)
	for i, in := range demuxed {
		joiner.Add(1)
		go func(i int, in potentials.RepositoryChannel) {
			defer joiner.Done()
			criteria := potentials.SortCriteria(i)
			repositories := in.FilterZeroStats(criteria).
				Dump(onError).
				Sort(criteria).
				Trim(conf.OutCount)
			collected[criteria] = repositories
		}(i, in)
	}
	joiner.Wait()

	remained, reset, err := client.APIRates()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	out := result{
		Metadata: meta{
			Updated:          time.Now(),
			Errors:           errCount,
			APICallsRemained: remained,
			Reset:            reset,
			DurationSec:      int(time.Since(startTime).Seconds()),
		},
		ByCommits:      collected[potentials.CommitsCriteria],
		ByStars:        collected[potentials.StarsCriteria],
		ByContributors: collected[potentials.ContributorsCriteria],
	}

	err = writeToFile(out, conf.OutputPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func writeToFile(result result, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	err = result.Write(file)
	if err != nil {
		return err
	}
	err = file.Sync()
	return err
}
