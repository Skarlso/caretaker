package cmd

import (
	"context"

	"github.com/shurcooL/githubv4"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"

	"github.com/skarlso/caretaker/pkg/assignissue"
	"github.com/skarlso/caretaker/pkg/client"
	"github.com/skarlso/caretaker/pkg/logger"
)

func CreateAssignIssueCommand(rootArgs *rootArgsStruct) *cobra.Command {
	createIssueCmd := &cobra.Command{
		Use:   "assign-issue",
		Short: "Assigns an issue created in this repository to a specific project.",
	}

	createIssueCmd.RunE = assignIssueRunE(rootArgs)

	return createIssueCmd
}

func assignIssueRunE(rootArgs *rootArgsStruct) func(cmd *cobra.Command, args []string) error {
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

		log.Log("running assign command")

		client := client.NewCaretaker(log, gclient, client.Options{
			Repo:           rootArgs.repo,
			Owner:          rootArgs.owner,
			StatusName:     rootArgs.statusOption,
			Interval:       rootArgs.staleInterval,
			StaleLabel:     rootArgs.pullRequestProcessedLabel,
			IsOrganization: rootArgs.isOrganization != "",
		})
		assigner := assignissue.NewAssignIssueAction(log, client, assignissue.Options{
			ProjectNumber: rootArgs.projectNumber,
			IssueNumber:   rootArgs.issueNumber,
		})

		return assigner.Assign(ctx)
	}
}
