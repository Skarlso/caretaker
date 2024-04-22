package cmd

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/shurcooL/githubv4"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"

	"github.com/skarlso/caretaker/pkg/client"
	"github.com/skarlso/caretaker/pkg/logger"
	"github.com/skarlso/caretaker/pkg/scanproject"
)

func CreateScanProjectCommand(rootArgs *rootArgsStruct) *cobra.Command {
	scanProjectCmd := &cobra.Command{
		Use:   "scan-project",
		Short: "Scans a project for out-dated project items and moves them into the specified column.",
	}

	scanProjectCmd.RunE = scanProjectRunE(rootArgs)

	return scanProjectCmd
}

func scanProjectRunE(rootArgs *rootArgsStruct) func(cmd *cobra.Command, args []string) error {
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

		log.Log("running scan issues command")

		projectNumber, err := strconv.Atoi(rootArgs.projectNumber)
		if err != nil {
			return fmt.Errorf("failed to convert project number: %w", err)
		}

		interval, err := time.ParseDuration(rootArgs.scanInterval)
		if err != nil {
			return fmt.Errorf("failed to parse interval: %w", err)
		}

		caretaker := client.NewCaretaker(log, gclient, client.Options{
			Repo:           rootArgs.repo,
			Owner:          rootArgs.owner,
			IsOrganization: rootArgs.isOrganization != "",
			MoveClosed:     rootArgs.moveClosed != "",
		})
		scanner := scanproject.NewScanner(log, caretaker, scanproject.Options{
			ProjectNumber: projectNumber,
			FromStatus:    rootArgs.fromStatusOption,
			ToStatus:      rootArgs.statusOption,
			Interval:      interval,
		})

		return scanner.ScanIssues(ctx)
	}
}
