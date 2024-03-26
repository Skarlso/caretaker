package pullrequestupdated

import (
	"context"
	"fmt"

	"github.com/shurcooL/githubv4"

	"github.com/skarlso/caretaker/pkg/client"
	"github.com/skarlso/caretaker/pkg/logger"
)

type Options struct {
	PullRequestNumber int
	StatusName        string
	ScanLabel         string
	NoComment         bool
}

type Updater struct {
	Options

	client client.Client
	log    logger.Logger
}

func NewUpdater(log logger.Logger, client client.Client, opts Options) *Updater {
	return &Updater{
		log:     log,
		client:  client,
		Options: opts,
	}
}

// PullRequestUpdated moves issues into a specific status on a given Pull Request.
func (c *Updater) PullRequestUpdated(ctx context.Context) error {
	pr, err := c.client.PullRequest(ctx, c.PullRequestNumber)
	if err != nil {
		return fmt.Errorf("failed to get pull request: %w", err)
	}

	if len(pr.ClosingIssuesReferences.Nodes) == 0 {
		c.log.Log("pull request with number %d doesn't have any issues associated with it", pr.Number)

		return nil
	}

	var updated bool

	for _, issue := range pr.ClosingIssuesReferences.Nodes {
		issue := issue

		// issue is not part of a project, skip it
		if issue.ProjectItems.TotalCount == 0 {
			c.log.Log("issue with number %d is not part of any project, skipping", issue.Number)

			continue
		}

		// if any of its project items is not in the desired state, we'll update it.
		updated, err = c.client.UpdateIssueStatus(ctx, issue, githubv4.String(c.StatusName), -1)
		if err != nil {
			return fmt.Errorf("failed to mutate issue: %w", err)
		}

		c.log.Debug("issue number %d successfully mutated", issue.Number)
	}

	if err := c.client.RemoveLabel(ctx, c.ScanLabel, pr.ID); err != nil {
		return fmt.Errorf("failed to remove label from entity: %w", err)
	}

	// if there was no update performed, don't leave a comment
	if !c.NoComment && updated {
		if err := c.client.LeaveComment(
			ctx,
			pr.ID,
			fmt.Sprintf("Update detected, any open associated issue has been transfer to %s.", c.StatusName),
		); err != nil {
			// we continue as everything else seemed to have worked and a comment shouldn't stop the flow
			c.log.Log("failed to leave comment on pull request %d with error: %s", pr.Number, err)
		}
	}

	return nil
}
