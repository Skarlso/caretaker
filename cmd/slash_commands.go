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
	"github.com/skarlso/caretaker/pkg/slash"
	"github.com/skarlso/caretaker/pkg/slash/assign"
	"github.com/skarlso/caretaker/pkg/slash/status"
)

// CreateSlashCommand defines a command that handles comments made by users on objects that Caretaker tracks.
// These commands run on pull requests and their respective issues.
// This will create a reaction of a +1 on the comment once it's done.
func CreateSlashCommand(rootArgs *rootArgsStruct) *cobra.Command {
	scanCmd := &cobra.Command{
		Use:   "slash",
		Short: "Marks issues belonging to PRs that are over 24 hours old as in review",
	}

	scanCmd.RunE = slashRunE(rootArgs)

	return scanCmd
}

func slashRunE(rootArgs *rootArgsStruct) func(cmd *cobra.Command, args []string) error {
	return func(_ *cobra.Command, _ []string) error {
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

		client := client.NewCaretaker(log, gclient, client.Options{
			Repo:           rootArgs.repo,
			Owner:          rootArgs.owner,
			IsOrganization: rootArgs.isOrganization != "",
		})

		assignHandler := assign.NewHandler(client)
		statusHandler := status.NewHandler(client)
		s := slash.NewSlashHandler(client)
		s.RegisterHandler(assign.Command, assignHandler)
		s.RegisterHandler(status.Command, statusHandler)
		s.RegisterHandler(slash.Help, s)

		prNumber, err := strconv.Atoi(rootArgs.pullRequestNumber)
		if err != nil {
			return fmt.Errorf("failed to convert pull number: %w", err)
		}

		return s.Run(ctx, prNumber, rootArgs.actor, rootArgs.commentID, rootArgs.commentBody)
	}
}
