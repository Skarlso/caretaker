package moveissue

import (
	"context"
	"fmt"

	"github.com/skarlso/caretaker/pkg/client"
	"github.com/skarlso/caretaker/pkg/logger"
)

type Options struct {
	PullRequestNumber int
	StatusName        string
	StaleLabel        string
}

type Mover struct {
	Options

	client client.Client
	log    logger.Logger
}

func NewMoveIssueAction(log logger.Logger, client client.Client, opts Options) *Mover {
	return &Mover{
		log:     log,
		client:  client,
		Options: opts,
	}
}

// Move moves issues into a specific status on a given Pull Request.
func (c *Mover) Move(ctx context.Context) error {
	pr, err := c.client.PullRequest(ctx, c.PullRequestNumber)
	if err != nil {
		return fmt.Errorf("failed to get pull request: %w", err)
	}

	if len(pr.ClosingIssuesReferences.Nodes) == 0 {
		c.log.Log("pull request with number %d doesn't have any issues associated with it", pr.Number)

		return nil
	}

	for _, issue := range pr.ClosingIssuesReferences.Nodes {
		issue := issue
		if err := c.client.MutateIssue(ctx, issue); err != nil {
			return fmt.Errorf("failed to mutate issue: %w", err)
		}

		c.log.Debug("issue number %d successfully mutated", issue.Number)
	}

	if err := c.client.RemoveLabel(ctx, c.StaleLabel, pr.ID); err != nil {
		return fmt.Errorf("failed to remove label from entity: %w", err)
	}

	if err := c.client.LeaveComment(
		ctx,
		pr.ID,
		fmt.Sprintf("Update detected, any open associated issue has been transfer to %s.", c.StatusName),
	); err != nil {
		// we continue as everything else seemed to have worked and a comment shouldn't stop the flow
		c.log.Log("failed to leave comment on pull request %d with error: %s", pr.Number, err)
	}

	return nil
}
