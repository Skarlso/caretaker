package updateissue

import (
	"context"
	"fmt"

	"github.com/skarlso/caretaker/pkg/client"
	"github.com/skarlso/caretaker/pkg/logger"
)

type Options struct {
	IssueNumber int
}

type Updater struct {
	Options

	client client.Client
	log    logger.Logger
}

func NewUpdateIssueAction(log logger.Logger, client client.Client, opts Options) *Updater {
	return &Updater{
		log:     log,
		client:  client,
		Options: opts,
	}
}

// Update Updates an issue to project.
func (c *Updater) Update(ctx context.Context) error {
	issue, err := c.client.Issue(ctx, c.IssueNumber)
	if err != nil {
		return fmt.Errorf("failed to fetch issue: %w", err)
	}

	if err := c.client.UpdateIssueStatus(ctx, issue); err != nil {
		return fmt.Errorf("failed to update issue: %w", err)
	}

	return nil
}
