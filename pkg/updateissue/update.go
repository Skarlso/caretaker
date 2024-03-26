package updateissue

import (
	"context"
	"fmt"

	"github.com/shurcooL/githubv4"

	"github.com/skarlso/caretaker/pkg/client"
	"github.com/skarlso/caretaker/pkg/logger"
)

type Options struct {
	ProjectNumber int
	IssueNumber   int
	FromStatus    string
	ToStatus      string
}

type Assigner struct {
	Options

	client client.Client
	log    logger.Logger
}

func NewUpdateIssueAction(log logger.Logger, client client.Client, opts Options) *Assigner {
	return &Assigner{
		log:     log,
		client:  client,
		Options: opts,
	}
}

// Update an issue status.
func (c *Assigner) Update(ctx context.Context) error {
	issue, err := c.client.Issue(ctx, c.IssueNumber)
	if err != nil {
		return fmt.Errorf("failed to fetch issue: %w", err)
	}

	var projectItem *client.ProjectV2Item

	for _, item := range issue.ProjectItems.Nodes {
		item := item
		if item.Project.Number == githubv4.Int(c.ProjectNumber) {
			projectItem = &item

			break
		}
	}

	if projectItem == nil {
		return fmt.Errorf(
			"project with id %d not found in the list of project associated to issue number %d",
			c.ProjectNumber,
			c.IssueNumber)
	}

	if c.FromStatus != "" &&
		projectItem.FieldValueByName.ProjectV2SingleSelectField.Name != githubv4.String(c.FromStatus) {
		c.log.Log(
			"issue with number %d ignored as the from status %s did not match with set status %s",
			c.IssueNumber,
			projectItem.FieldValueByName.ProjectV2SingleSelectField.Name,
			c.FromStatus,
		)

		return nil
	}

	if _, err := c.client.UpdateIssueStatus(ctx, issue, githubv4.String(c.ToStatus), c.ProjectNumber); err != nil {
		return fmt.Errorf("failed to update issue with number %d to status %s: %w", c.IssueNumber, c.ToStatus, err)
	}

	return nil
}
