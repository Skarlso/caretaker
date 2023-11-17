package cmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/shurcooL/githubv4"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"

	"github.com/skarlso/caretaker/pkg/client"
	"github.com/skarlso/caretaker/pkg/logger"
	"github.com/skarlso/caretaker/pkg/pullrequestupdated"
)

// CreatePullRequestUpdatedCommand gets the issue and updates its status to the desired Status.
func CreatePullRequestUpdatedCommand(rootArgs *rootArgsStruct) *cobra.Command {
	pullRequestUpdatedCmd := &cobra.Command{
		Use:   "pull-request-updated",
		Short: "Execute this command if there was a pull request update event to update any connecting issues.",
	}

	pullRequestUpdatedCmd.RunE = pullRequestUpdatedRunE(rootArgs)

	return pullRequestUpdatedCmd
}

func pullRequestUpdatedRunE(rootArgs *rootArgsStruct) func(cmd *cobra.Command, args []string) error {
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

		log.Log("running pull request updated command")

		prNumber, err := strconv.Atoi(rootArgs.pullRequestNumber)
		if err != nil {
			return fmt.Errorf("failed to convert pull request number: %w", err)
		}

		client := client.NewCaretaker(log, gclient, client.Options{
			Repo:  rootArgs.repo,
			Owner: rootArgs.owner,
		})
		updater := pullrequestupdated.NewUpdater(log, client, pullrequestupdated.Options{
			PullRequestNumber: prNumber,
			StatusName:        rootArgs.statusOption,
			ScanLabel:         rootArgs.pullRequestProcessedLabel,
		})

		return updater.PullRequestUpdated(ctx)
	}
}
