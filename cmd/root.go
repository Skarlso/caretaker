package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

type rootArgsStruct struct {
	token                     string
	owner                     string
	repo                      string
	authorName                string
	authorEmail               string
	verbose                   bool
	pullRequestNumber         int
	issueNumber               int
	projectNumber             string
	statusOption              string
	staleInterval             time.Duration
	pullRequestProcessedLabel string
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
	flag.DurationVar(&rootArgs.staleInterval, "stale-interval", 24*time.Hour, "--stale-interval defines after how long duration a pull request should be considered stale")
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
	flag.BoolVarP(&rootArgs.verbose, "verbose", "v", false, "--verbose|-v if enabled, will output extra debug information")
	flag.IntVar(
		&rootArgs.pullRequestNumber,
		"pull-request-number",
		0,
		"--pull-request-number is the number of the pull request currently inspected")
	flag.IntVar(
		&rootArgs.issueNumber,
		"issue-number",
		0,
		"--issue-number the number of the issue currently inspected")
	flag.StringVar(
		&rootArgs.projectNumber,
		"project-number",
		"",
		"--project-number the number of the project to add a created issue to")
	flag.StringVar(
		&rootArgs.statusOption,
		"status-option",
		"",
		"--status-option is the status to set an issue to")
	flag.StringVar(
		&rootArgs.pullRequestProcessedLabel,
		"pull-request-processed-label",
		"caretaker-processed",
		"--pull-request-processed-label is the label used to signify that Caretaker already processed this pull request. This label is removed on update.",
	)

	markFlagAsRequired(rootCmd, "token")
	markFlagAsRequired(rootCmd, "owner")
	markFlagAsRequired(rootCmd, "repo")

	staleCmd := CreateStaleCommand(rootArgs)
	createIssueCmd := CreateCreateIssueCommand(rootArgs)
	moveIssueCmd := CreateMoveIssueCommand(rootArgs)
	assignIssueCmd := CreateAssignIssueCommand(rootArgs)
	rootCmd.AddCommand(staleCmd, createIssueCmd, moveIssueCmd, assignIssueCmd)

	return rootCmd
}

func markFlagAsRequired(cmd *cobra.Command, flag string) {
	if err := cmd.MarkPersistentFlagRequired(flag); err != nil {
		fmt.Printf("failed to mark %s flag as required", flag)
		os.Exit(1)
	}
}
