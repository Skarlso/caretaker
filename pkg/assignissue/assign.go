package assignissue

import (
	"context"
	"fmt"

	"github.com/skarlso/caretaker/pkg/client"
	"github.com/skarlso/caretaker/pkg/logger"
)

type Options struct {
	ProjectNumber int
	IssueNumber   int
}

type Assigner struct {
	Options

	client client.Client
	log    logger.Logger
}

func NewAssignIssueAction(log logger.Logger, client client.Client, opts Options) *Assigner {
	return &Assigner{
		log:     log,
		client:  client,
		Options: opts,
	}
}

// Assign assigns an issue to project.
func (c *Assigner) Assign(ctx context.Context) error {
	if err := c.client.AssignIssueToProject(ctx, c.IssueNumber, c.ProjectNumber); err != nil {
		return fmt.Errorf("failed to assign issue to project: %w", err)
	}

	return nil
}
