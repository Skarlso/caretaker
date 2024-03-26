package scanproject

import (
	"context"
	"fmt"
	"time"

	"github.com/shurcooL/githubv4"

	"github.com/skarlso/caretaker/pkg/client"
	"github.com/skarlso/caretaker/pkg/logger"
)

type Options struct {
	ProjectNumber   int
	Interval        time.Duration
	DisableComments bool
	FromStatus      string
	ToStatus        string
}

type Scanner struct {
	Options

	client client.Client
	log    logger.Logger
}

func NewScanner(log logger.Logger, client client.Client, opts Options) *Scanner {
	return &Scanner{
		log:     log,
		client:  client,
		Options: opts,
	}
}

// ScanIssues checks if any issues of a project should be moved into a different column based on ToStatus.
func (c *Scanner) ScanIssues(ctx context.Context) error {
	now := time.Now()

	items, err := c.client.ProjectItems(ctx, c.ProjectNumber)
	if err != nil {
		return err
	}

	for _, i := range items {
		item := i
		if v := item.FieldValueByName.ProjectV2SingleSelectField.Name; v != githubv4.String(c.FromStatus) {
			c.log.Log("skipping issue %s; status %s doesn't match with %s", item.Content.Issue.Title, v, c.FromStatus)

			continue
		}

		// If the last action ( any action ) on the issue is after now, skip it.
		if item.UpdatedAt.Add(c.Interval).After(now) {
			c.log.Log(
				"issue %s has been last updated at %s which doesn't exceed the interval %s",
				item.Content.Issue.Title,
				item.UpdatedAt.Format(time.RFC3339),
				c.Interval,
			)

			continue
		}

		if _, err := c.client.UpdateIssueStatus(
			ctx,
			item.Content.Issue,
			githubv4.String(c.ToStatus),
			c.ProjectNumber,
		); err != nil {
			return fmt.Errorf("failed to update issue: %w", err)
		}
	}

	return nil
}
