package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/shurcooL/githubv4"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"

	"github.com/skarlso/caretaker/pkg/client"
	"github.com/skarlso/caretaker/pkg/logger"
	"github.com/skarlso/caretaker/pkg/scan"
)

func CreateScanCommand(rootArgs *rootArgsStruct) *cobra.Command {
	scanCmd := &cobra.Command{
		Use:   "scan",
		Short: "Marks issues belonging to PRs that are over 24 hours old as in review",
	}

	scanCmd.RunE = scanRunE(rootArgs)

	return scanCmd
}

func scanRunE(rootArgs *rootArgsStruct) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: rootArgs.token},
		)
		tc := oauth2.NewClient(ctx, ts)
		gclient := githubv4.NewClient(tc)

		// setup logger
		var log logger.Logger = &logger.QuiteLogger{}
		if rootArgs.verbose {
			log = &logger.VerboseLogger{}
		}

		log.Log("running scan command")

		interval, err := time.ParseDuration(rootArgs.scanInterval)
		if err != nil {
			return fmt.Errorf("failed to parse interval: %w", err)
		}

		client := client.NewCaretaker(log, gclient, client.Options{
			Repo:  rootArgs.repo,
			Owner: rootArgs.owner,
		})
		scanner := scan.NewScanner(log, client, scan.Options{
			Interval:        interval,
			ScanLabel:       rootArgs.pullRequestProcessedLabel,
			DisableComments: rootArgs.disableComments != "",
			StatusName:      rootArgs.statusOption,
		})

		return scanner.Scan(ctx)
	}
}
