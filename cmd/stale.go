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
	"github.com/skarlso/caretaker/pkg/stale"
)

func CreateStaleCommand(rootArgs *rootArgsStruct) *cobra.Command {
	staleCmd := &cobra.Command{
		Use:   "stale",
		Short: "Marks issues belonging to PRs that are over 24 hours old as in review",
	}

	staleCmd.RunE = staleRunE(rootArgs)

	return staleCmd
}

func staleRunE(rootArgs *rootArgsStruct) func(cmd *cobra.Command, args []string) error {
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

		log.Log("running stale command")

		interval, err := time.ParseDuration(rootArgs.staleInterval)
		if err != nil {
			return fmt.Errorf("failed to parse interval: %w", err)
		}

		client := client.NewCaretaker(log, gclient, client.Options{
			Repo:       rootArgs.repo,
			Owner:      rootArgs.owner,
			StatusName: rootArgs.statusOption,
			Interval:   interval,
			StaleLabel: rootArgs.pullRequestProcessedLabel,
		})
		checker := stale.NewStaleChecker(log, client, interval, rootArgs.pullRequestProcessedLabel)

		return checker.Check(ctx)
	}
}
