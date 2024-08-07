package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// All of these are string to conform to GitHub's map[string]string actions.yaml.
type rootArgsStruct struct {
	token                     string
	owner                     string
	repo                      string
	authorName                string
	authorEmail               string
	verbose                   bool
	pullRequestNumber         string
	issueNumber               string
	projectNumber             string
	statusOption              string
	scanInterval              string
	pullRequestProcessedLabel string
	isOrganization            string
	disableComments           string
	commentBody               string
	commentID                 string
	actor                     string
	fromStatusOption          string
	moveClosed                string
}

func CreateRootCommand() *cobra.Command {
	rootArgs := &rootArgsStruct{}

	rootCmd := &cobra.Command{
		Use:   "root",
		Short: "Dependabot bundler action",
	}

	flag := rootCmd.PersistentFlags()

	// Server Configs
	flag.StringVar(&rootArgs.token, "token", "", "--token github token")
	flag.StringVar(&rootArgs.owner, "owner", "", "--owner github organization / owner")
	flag.StringVar(&rootArgs.repo, "repo", "", "--repo github repository")
	flag.StringVar(
		&rootArgs.scanInterval,
		"scan-interval",
		"24h",
		"--scan-interval defines after how long duration a pull request should be considered scan",
	)
	flag.StringVar(
		&rootArgs.authorName,
		"author-name",
		"Github Action",
		"--author-name name of the committer, default to Github Action")
	flag.StringVar(
		&rootArgs.authorEmail,
		"author-email",
		"41898282+github-actions[bot]@users.noreply.github.com",
		"--author-email email address of the committer, defaults to github action's email address")
	flag.BoolVarP(
		&rootArgs.verbose,
		"verbose",
		"v",
		false,
		"--verbose|-v if enabled, will output extra debug information",
	)
	flag.StringVar(
		&rootArgs.pullRequestNumber,
		"pull-request-number",
		"0",
		"--pull-request-number is the number of the pull request currently inspected")
	flag.StringVar(
		&rootArgs.issueNumber,
		"issue-number",
		"0",
		"--issue-number the number of the issue currently inspected")
	flag.StringVar(
		&rootArgs.projectNumber,
		"project-number",
		"0",
		"--issue-number the number of the project to add a created issue to")
	flag.StringVar(
		&rootArgs.statusOption,
		"status-option",
		"",
		"--status-option is the status to set an issue to")
	flag.StringVar(
		&rootArgs.fromStatusOption,
		"from-status-option",
		"",
		"--from-status-option the status to transition from")
	flag.StringVar(
		&rootArgs.pullRequestProcessedLabel,
		"pull-request-processed-label",
		"caretaker-processed",
		"--pull-request-processed-label label used to mark pull request as processed. This label is removed on update.",
	)
	flag.StringVar(
		&rootArgs.isOrganization,
		"is-organization",
		"",
		"--is-organization=true is defined if the user is an organization",
	)
	flag.StringVar(
		&rootArgs.disableComments,
		"disable-comments",
		"",
		"--disable-comments=true is used to disable caretaker commenting back on PRs",
	)
	flag.StringVar(
		&rootArgs.commentBody,
		"comment-body",
		"",
		"--comment-body:/test the body of the comment as passed from the github action context",
	)
	flag.StringVar(
		&rootArgs.commentID,
		"comment-id",
		"",
		"--comment-id:IC_ the node_id of the comment that triggered this action",
	)
	flag.StringVar(
		&rootArgs.actor,
		"actor",
		"",
		"--actor is the username of the actor who performed the action",
	)
	flag.StringVar(
		&rootArgs.moveClosed,
		"move-closed",
		"",
		"--move-closed will edit closed issues, by default closed issues are not moved",
	)

	markFlagAsRequired(rootCmd, "token")
	markFlagAsRequired(rootCmd, "owner")

	scanCmd := CreateScanCommand(rootArgs)
	pullRequestUpdatedCmd := CreatePullRequestUpdatedCommand(rootArgs)
	assignIssueCmd := CreateAssignIssueCommand(rootArgs)
	slashCommandCmd := CreateSlashCommand(rootArgs)
	updateIssueCommandCmd := CreateUpdateIssueCommand(rootArgs)
	scanProjectCmd := CreateScanProjectCommand(rootArgs)
	rootCmd.AddCommand(
		scanCmd,
		pullRequestUpdatedCmd,
		assignIssueCmd,
		slashCommandCmd,
		updateIssueCommandCmd,
		scanProjectCmd,
	)

	return rootCmd
}

func markFlagAsRequired(cmd *cobra.Command, flag string) {
	if err := cmd.MarkPersistentFlagRequired(flag); err != nil {
		fmt.Printf("failed to mark %s flag as required", flag)
		os.Exit(1)
	}
}
