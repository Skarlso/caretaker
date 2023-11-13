package cmd

import (
	"context"

	"github.com/shurcooL/githubv4"
	"github.com/skarlso/caretaker/pkg/logger"
	"github.com/skarlso/caretaker/pkg/moveissue"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

// CreateMoveIssueCommand gets the issue and updates its status to the desired Status.
func CreateMoveIssueCommand(rootArgs *rootArgsStruct) *cobra.Command {
	moveIssueCmd := &cobra.Command{
		Use:   "move-issue",
		Short: "Moves an issue into a specific column location on a Project Board.",
	}

	moveIssueCmd.RunE = moveIssueRunE(rootArgs)

	return moveIssueCmd
}

func moveIssueRunE(rootArgs *rootArgsStruct) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: rootArgs.token},
		)
		tc := oauth2.NewClient(ctx, ts)
		client := githubv4.NewClient(tc)

		// setup logger
		var log logger.Logger = &logger.QuiteLogger{}
		if rootArgs.verbose {
			log = &logger.VerboseLogger{}
		}

		log.Log("running move issue command")

		mover := moveissue.NewMoveIssueAction(log, client, moveissue.Options{
			Repo:              rootArgs.repo,
			Owner:             rootArgs.owner,
			StatusName:        rootArgs.statusOption,
			StaleLabel:        rootArgs.pullRequestProcessedLabel,
			PullRequestNumber: rootArgs.pullRequestNumber,
		})

		return mover.Move(ctx)
	}
}
