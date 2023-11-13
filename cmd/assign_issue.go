package cmd

import "github.com/spf13/cobra"

func CreateAssignIssueCommand(rootArgs *rootArgsStruct) *cobra.Command {
	createIssueCmd := &cobra.Command{
		Use:   "assign-issue",
		Short: "Assigns an issue created in this repository to a specific project.",
	}

	return createIssueCmd
}
