package scan

import (
	"context"
	"fmt"
	"time"

	"github.com/skarlso/caretaker/pkg/client"
	"github.com/skarlso/caretaker/pkg/logger"
)

type Scanner struct {
	interval        time.Duration
	scanLabel       string
	client          client.Client
	log             logger.Logger
	disableComments bool
}

func NewScanner(log logger.Logger, client client.Client, interval time.Duration, scanLabel string, disableComments bool) *Scanner {
	return &Scanner{
		log:             log,
		client:          client,
		interval:        interval,
		scanLabel:       scanLabel,
		disableComments: disableComments,
	}
}

// Scan checks if any associated issues should be moved into a different column based on this PR.
func (c *Scanner) Scan(ctx context.Context) error {
	pullRequests, err := c.client.PullRequests(ctx)
	if err != nil {
		return fmt.Errorf("failed to list pull requests: %w", err)
	}

	now := time.Now()
loop:
	for _, pr := range pullRequests {
		pr := pr

		for _, label := range pr.Labels.Nodes {
			if string(label.Name) == c.scanLabel {
				c.log.Log("pull request with number %d already processed", pr.Number)

				continue loop
			}
		}

		// If the last action ( any action ) on the Pull Request is after now, skip it.
		if pr.UpdatedAt.Add(c.interval).After(now) {
			continue
		}

		if len(pr.ClosingIssuesReferences.Nodes) == 0 {
			c.log.Log("pull request with number %d doesn't have any issues associated with it", pr.Number)

			continue
		}

		for _, issue := range pr.ClosingIssuesReferences.Nodes {
			issue := issue
			if err := c.client.UpdateIssueStatus(ctx, issue); err != nil {
				return fmt.Errorf("failed to mutate issue: %w", err)
			}

			c.log.Debug("issue number %d successfully mutated", issue.Number)
		}

		if err := c.client.AddLabel(ctx, c.scanLabel, pr.ID); err != nil {
			return fmt.Errorf("failed to add label to processed entity: %w", err)
		}

		if !c.disableComments {
			if err := c.client.LeaveComment(ctx, pr.ID, "Pull request successfully processed by Caretaker."); err != nil {
				c.log.Log("failed to leave comment on pull request %d with error: %s", pr.Number, err)
				// we continue as everything else seemed to have worked and a comment shouldn't stop the flow
			}
		}
	}

	return nil
}
