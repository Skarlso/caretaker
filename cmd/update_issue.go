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
	"github.com/skarlso/caretaker/pkg/updateissue"
)

func CreateUpdateIssueCommand(rootArgs *rootArgsStruct) *cobra.Command {
	updateIssueCmd := &cobra.Command{
		Use:   "update-issue",
		Short: "Sets an issue into the desired state. Can be used regardless of any actions.",
	}

	updateIssueCmd.RunE = updateIssueRunE(rootArgs)

	return updateIssueCmd
}

func updateIssueRunE(rootArgs *rootArgsStruct) func(cmd *cobra.Command, args []string) error {
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

		log.Log("running update command")

		issueNumber, err := strconv.Atoi(rootArgs.issueNumber)
		if err != nil {
			return fmt.Errorf("failed to convert issue number: %w", err)
		}

		client := client.NewCaretaker(log, gclient, client.Options{
			Repo:             rootArgs.repo,
			Owner:            rootArgs.owner,
			TargetStatusName: rootArgs.statusOption,
			IsOrganization:   rootArgs.isOrganization != "",
		})
		updater := updateissue.NewUpdateIssueAction(log, client, updateissue.Options{
			IssueNumber: issueNumber,
		})

		return updater.Update(ctx)
	}
}
