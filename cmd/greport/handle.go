package main

import (
	"fmt"

	gitreport "github.com/vanhtuan0409/git-report"
	cli "gopkg.in/urfave/cli.v1"
)

func generateReport(c *cli.Context) error {
	configPath := gitreport.GetDefaultConfigPath()
	config, err := gitreport.ReadConfigFromFile(configPath)
	if err != nil {
		panic(err)
	}

	resultChan := make(chan string)
	errChan := make(chan error)
	for _, repoPath := range config.Repos {
		go func(path string) {
			gitClient := gitreport.NewGitClient(path)
			result, err := gitClient.Log(&gitreport.LogOption{
				Authors:           config.FilterEmail,
				Limit:             5,
				FetchAllBranch:    true,
				FilterMergeCommit: true,
			})
			if err != nil {
				errChan <- fmt.Errorf("Cannot fetch git commits from url: %s. Original Error:\n%s", path, err.Error())
				return
			}

			generator := gitreport.NewReportGenerator()
			report := generator.GenerateFromCommits(result)
			resultChan <- report
		}(repoPath)
	}

	for i := 0; i < len(config.Repos); i++ {
		select {
		case result := <-resultChan:
			fmt.Println(result)
		case err := <-errChan:
			fmt.Println(err.Error())
		}
	}

	return nil
}
