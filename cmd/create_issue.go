package cmd

import "github.com/spf13/cobra"

func CreateCreateIssueCommand(rootArgs *rootArgsStruct) *cobra.Command {
	createIssueCmd := &cobra.Command{
		Use:   "create-issue",
		Short: "Creates an issue and links it to a project board.",
	}

	return createIssueCmd
}
