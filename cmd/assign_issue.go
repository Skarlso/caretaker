package cmd

import (
	"context"
	"fmt"
	"strconv"
	"time"

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

		interval, err := time.ParseDuration(rootArgs.staleInterval)
		if err != nil {
			return fmt.Errorf("failed to parse interval: %w", err)
		}

		projectNumber, err := strconv.Atoi(rootArgs.projectNumber)
		if err != nil {
			return fmt.Errorf("failed to convert pull request number: %w", err)
		}

		issueNumber, err := strconv.Atoi(rootArgs.issueNumber)
		if err != nil {
			return fmt.Errorf("failed to convert issue number: %w", err)
		}

		client := client.NewCaretaker(log, gclient, client.Options{
			Repo:           rootArgs.repo,
			Owner:          rootArgs.owner,
			StatusName:     rootArgs.statusOption,
			Interval:       interval,
			StaleLabel:     rootArgs.pullRequestProcessedLabel,
			IsOrganization: rootArgs.isOrganization != "",
		})
		assigner := assignissue.NewAssignIssueAction(log, client, assignissue.Options{
			ProjectNumber: projectNumber,
			IssueNumber:   issueNumber,
		})

		return assigner.Assign(ctx)
	}
}
