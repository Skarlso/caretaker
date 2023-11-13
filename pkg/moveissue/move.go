package moveissue

import (
	"context"
	"fmt"

	"github.com/shurcooL/githubv4"
	"github.com/skarlso/caretaker/pkg/logger"
)

type Client interface {
	Query(ctx context.Context, q any, variables map[string]any) error
	Mutate(ctx context.Context, m any, input githubv4.Input, variables map[string]any) error
}

// queryPullRequestIssues gets all related open issues for this pull request.
var queryPullRequestIssues struct {
	Repository struct {
		PullRequest struct {
			ID        githubv4.String
			Number    githubv4.Int
			UpdatedAt githubv4.Date
			Labels    struct {
				Nodes []struct {
					Name githubv4.String
				}
			} `graphql:"labels(first: 50)"` // We can't use Label with name because that fails if the label is not there
			ClosingIssuesReferences struct {
				Nodes []struct {
					Closed     githubv4.Boolean
					Title      githubv4.String
					Number     githubv4.Int
					ProjectsV2 struct {
						Nodes []struct {
							Title githubv4.String
							ID    githubv4.String
							Field struct {
								ProjectV2SingleSelectField struct {
									ID      githubv4.String
									Options []struct {
										ID githubv4.String
									} `graphql:"options(names: [$statusName])"`
								} `graphql:"... on ProjectV2SingleSelectField"`
							} `graphql:"field(name: \"Status\")"`
						}
					} `graphql:"projectsV2(first: 1)"` // we assume an issue is only part of a single project
					ProjectItems struct {
						TotalCount githubv4.Int
						Nodes      []struct {
							ID githubv4.String
						}
					} `graphql:"projectItems(first: 1)"` // there should be only one card associated with this issue.
				}
				PageInfo struct {
					EndCursor   githubv4.String
					HasNextPage bool
				} // should not be needed as we like... 1 issue that this thing closes.
			} `graphql:"closingIssuesReferences(first: 10)"`
		} `graphql:"pullRequest(number: $pullNumber)"`
	} `graphql:"repository(owner: $owner, name: $name)"`
}

// mutateIssueStatus sets the Status of an Issue to the desired option.
var mutateIssueStatus struct {
	UpdateProjectV2ItemFieldValue struct {
		ProjectV2Item struct {
			ID githubv4.String
		} `graphql:"projectV2Item"` // important because the return value is case-sensitive and the default is projectV2item which is wrong.
	} `graphql:"updateProjectV2ItemFieldValue(input: $input)"`
}

type Options struct {
	Repo              string
	Owner             string
	PullRequestNumber int
	StatusName        string
	StaleLabel        string
}

type Mover struct {
	Options

	client Client
	log    logger.Logger
}

func NewMoveIssueAction(log logger.Logger, client *githubv4.Client, opts Options) *Mover {
	return &Mover{
		log:     log,
		client:  client,
		Options: opts,
	}
}

// Move moves issues into a specific status on a given Pull Request.
func (c *Mover) Move(ctx context.Context) error {
	variables := map[string]any{
		"owner":      githubv4.String(c.Owner),
		"name":       githubv4.String(c.Repo),
		"statusName": githubv4.String(c.StatusName),
		"pullNumber": githubv4.Int(c.PullRequestNumber),
	}
	if err := c.client.Query(ctx, &queryPullRequestIssues, variables); err != nil {
		return fmt.Errorf("failed to query for issues: %w", err)
	}
	pr := queryPullRequestIssues.Repository.PullRequest

	if len(pr.ClosingIssuesReferences.Nodes) == 0 {
		c.log.Log("pull request with number %d doesn't have any issues associated with it", pr.Number)
		return nil
	}

	for _, issue := range pr.ClosingIssuesReferences.Nodes {
		if issue.Closed {
			c.log.Log("issue already closed, skip")
			continue
		}

		if len(issue.ProjectsV2.Nodes) != 1 {
			c.log.Log("issues that are attached to more than one project are not supported ATM")
			continue
		}

		project := issue.ProjectsV2.Nodes[0]

		if l := len(project.Field.ProjectV2SingleSelectField.Options); l != 1 {
			return fmt.Errorf("incorrect number of options found for name %s; want 1; got: %d", c.StatusName, l)
		}

		c.log.Debug("associated issue number %d and title %s on project: %s", issue.Number, issue.Title, project.Title)

		projectItem := issue.ProjectItems.Nodes[0]
		option := project.Field.ProjectV2SingleSelectField.Options[0]

		input := githubv4.UpdateProjectV2ItemFieldValueInput{
			ProjectID: githubv4.NewString(project.ID),
			ItemID:    githubv4.NewString(projectItem.ID),
			FieldID:   githubv4.NewString(project.Field.ProjectV2SingleSelectField.ID),
			Value: githubv4.ProjectV2FieldValue{
				SingleSelectOptionID: githubv4.NewString(option.ID),
			},
		}

		if err := c.client.Mutate(ctx, &mutateIssueStatus, input, nil); err != nil {
			return fmt.Errorf("failed to mutate issue: %w", err)
		}

		c.log.Debug("issue number %d successfully mutated", issue.Number)
	}

	if err := c.RemoveLabel(ctx, c.StaleLabel, pr.ID); err != nil {
		return fmt.Errorf("failed to remove label from entity: %w", err)
	}

	if err := c.LeaveComment(ctx, pr.ID, fmt.Sprintf("Update detected, any open associated issue has been transfer to %s.", c.StatusName)); err != nil {
		c.log.Log("failed to leave comment on pull request %d with error: %s", pr.Number, err)
		// we continue as everything else seemed to have worked and a comment shouldn't stop the flow
	}

	return nil
}

var leaveComment struct {
	AddComment struct {
		Subject struct {
			ID githubv4.String
		}
	} `graphql:"addComment(input: $input)"`
}

func (c *Mover) LeaveComment(ctx context.Context, prID githubv4.String, comment string) error {
	input := githubv4.AddCommentInput{
		SubjectID: prID,
		Body:      githubv4.String(comment),
	}
	if err := c.client.Mutate(ctx, &leaveComment, input, nil); err != nil {
		return fmt.Errorf("failed to leave comment on object: %w", err)
	}

	c.log.Debug("added comment with ID %s", leaveComment.AddComment.Subject.ID)
	return nil
}

var removeLabel struct {
	RemoveLabel struct {
		Labelable struct {
			Labels struct {
				TotalCount githubv4.Int
			}
		}
	} `graphql:"removeLabelsFromLabelable(input: $input)"`
}

func (c *Mover) RemoveLabel(ctx context.Context, label string, id githubv4.String) error {
	labelID, err := c.QueryLabelID(ctx, label)
	if err != nil {
		return err
	}

	input := githubv4.RemoveLabelsFromLabelableInput{
		LabelableID: id,
		LabelIDs:    []githubv4.ID{labelID},
	}
	if err := c.client.Mutate(ctx, &removeLabel, input, nil); err != nil {
		return fmt.Errorf("failed to remove label from object: %w", err)
	}

	c.log.Debug("removed label from pull request")
	return nil
}

var queryLabelID struct {
	Repository struct {
		Labels struct {
			Nodes []struct {
				ID githubv4.String
			}
		} `graphql:"labels(first: 1, query: $query)"` // There Can Be Only One!
	} `graphql:"repository(owner: $owner, name: $name)"`
}

func (c *Mover) QueryLabelID(ctx context.Context, label string) (githubv4.String, error) {
	variables := map[string]any{
		"owner": githubv4.String(c.Owner),
		"name":  githubv4.String(c.Repo),
		"query": githubv4.String(label),
	}

	if err := c.client.Query(ctx, &queryLabelID, variables); err != nil {
		return "", fmt.Errorf("failed to query for label id: %w", err)
	}

	if len(queryLabelID.Repository.Labels.Nodes) != 1 {
		return "", fmt.Errorf("expected a single label to be returned, got: %d", len(queryLabelID.Repository.Labels.Nodes))
	}

	return queryLabelID.Repository.Labels.Nodes[0].ID, nil
}
