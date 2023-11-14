package client

import (
	"context"
	"fmt"
	"time"

	"github.com/shurcooL/githubv4"

	"github.com/skarlso/caretaker/pkg/logger"
)

// Repository https://docs.github.com/en/graphql/reference/objects#repository
type Repository struct {
	PullRequests struct {
		PageInfo struct {
			EndCursor   githubv4.String
			HasNextPage bool
		} // should not be needed as we don't have more than 100 OPEN PRs.
		Nodes []PullRequest
	} `graphql:"pullRequests(first: 100, states: OPEN)"`
}

// PullRequest https://docs.github.com/en/graphql/reference/objects#pullrequest
type PullRequest struct {
	ID        githubv4.String
	Number    githubv4.Int
	UpdatedAt githubv4.Date
	Labels    struct {
		Nodes []struct {
			Name githubv4.String
		}
	} `graphql:"labels(first: 50)"` // We can't use Label with name because that fails if the label is not there
	ClosingIssuesReferences struct {
		Nodes    []Issue
		PageInfo struct {
			EndCursor   githubv4.String
			HasNextPage bool
		} // should not be needed because we don't reference hundreds of issues per pull request.
	} `graphql:"closingIssuesReferences(first: 10)"`
}

// ProjectV2Item https://docs.github.com/en/graphql/reference/objects#projectv2item
type ProjectV2Item struct {
	ID githubv4.String
}

// Issue https://docs.github.com/en/graphql/reference/objects#issue
type Issue struct {
	ID         githubv4.ID
	Closed     githubv4.Boolean
	Title      githubv4.String
	Number     githubv4.Int
	ProjectsV2 struct {
		Nodes []ProjectV2
	} `graphql:"projectsV2(first: 1)"` // we assume an issue is only part of a single project
	ProjectItems struct {
		TotalCount githubv4.Int
		Nodes      []ProjectV2Item
	} `graphql:"projectItems(first: 1)"` // there should be only one card associated with this issue
}

// Label https://docs.github.com/en/graphql/reference/objects#label
type Label struct{}

type Field struct{}

// ProjectV2 https://docs.github.com/en/graphql/reference/objects#projectv2
type ProjectV2 struct {
	Title  githubv4.String
	ID     githubv4.String
	Number githubv4.Int
	Field  struct { // TODO: Make this optional. Makes the query a bit not nice.
		ProjectV2SingleSelectField struct {
			ID      githubv4.String
			Options []struct {
				ID githubv4.String
			} `graphql:"options(names: [$statusName])"`
		} `graphql:"... on ProjectV2SingleSelectField"`
	} `graphql:"field(name: \"Status\")"` // gather the selection options for the Status field
}

// GraphQLClient hides the GitHub GraphQL library.
type GraphQLClient interface {
	Query(ctx context.Context, q any, variables map[string]any) error
	Mutate(ctx context.Context, m any, input githubv4.Input, variables map[string]any) error
}

// Client defines the capabilities of Caretaker.
//
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -o fakes/client.go . Client
type Client interface {
	AddLabel(ctx context.Context, label string, id githubv4.String) error
	AssignIssueToProject(ctx context.Context, issueNumber, projectNumber int) error
	LeaveComment(ctx context.Context, prID githubv4.String, comment string) error
	RemoveLabel(ctx context.Context, label string, id githubv4.String) error
	PullRequests(ctx context.Context) ([]PullRequest, error)
	PullRequest(ctx context.Context, prNumber int) (PullRequest, error)
	UpdateIssueStatus(ctx context.Context, issue Issue) error
	Issue(ctx context.Context, issueNumber int) (Issue, error)
}

// Options are for Caretaker's functionality.
type Options struct {
	Repo             string
	Owner            string
	TargetStatusName string
	Interval         time.Duration
	ScanLabel        string
	IsOrganization   bool
}

// Caretaker defines the main Caretaker capabilities.
type Caretaker struct {
	Options

	gclient GraphQLClient
	log     logger.Logger
}

// NewCaretaker creates a new Caretaker with an available GitHub GraphQL client.
func NewCaretaker(log logger.Logger, gc GraphQLClient, opts Options) *Caretaker {
	return &Caretaker{
		Options: opts,

		log:     log,
		gclient: gc,
	}
}

// Make sure Caretaker implements Client.
var _ Client = &Caretaker{}

func (c *Caretaker) AddLabel(ctx context.Context, label string, id githubv4.String) error {
	labelID, err := c.queryLabelID(ctx, label)
	if err != nil {
		return err
	}

	var addLabel struct {
		AddLabel struct {
			Labelable struct {
				Labels struct {
					TotalCount githubv4.Int
				}
			}
		} `graphql:"addLabelsToLabelable(input: $input)"`
	}

	input := githubv4.AddLabelsToLabelableInput{
		LabelableID: id,
		LabelIDs:    []githubv4.ID{labelID},
	}
	if err := c.gclient.Mutate(ctx, &addLabel, input, nil); err != nil {
		return fmt.Errorf("failed to add label on object: %w", err)
	}

	c.log.Debug("added label to pull request")

	return nil
}

func (c *Caretaker) AssignIssueToProject(ctx context.Context, issueNumber, projectNumber int) error {
	var getIssueQuery struct {
		Repository struct {
			Issue Issue `graphql:"issue(number: $issue)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	issueValues := map[string]any{
		"owner":      githubv4.String(c.Owner),
		"name":       githubv4.String(c.Repo),
		"issue":      githubv4.Int(issueNumber),
		"statusName": githubv4.String(""),
	}

	if err := c.gclient.Query(ctx, &getIssueQuery, issueValues); err != nil {
		return fmt.Errorf("failed to find issue with number %d: %w", issueNumber, err)
	}

	if c.IsOrganization {
		return c.assignToOrganization(ctx, &getIssueQuery.Repository.Issue, projectNumber)
	}

	return c.assignToUser(ctx, &getIssueQuery.Repository.Issue, projectNumber)
}

func (c *Caretaker) RemoveLabel(ctx context.Context, label string, id githubv4.String) error {
	labelID, err := c.queryLabelID(ctx, label)
	if err != nil {
		return err
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

	input := githubv4.RemoveLabelsFromLabelableInput{
		LabelableID: id,
		LabelIDs:    []githubv4.ID{labelID},
	}
	if err := c.gclient.Mutate(ctx, &removeLabel, input, nil); err != nil {
		return fmt.Errorf("failed to remove label from object: %w", err)
	}

	c.log.Debug("removed label from pull request")

	return nil
}

func (c *Caretaker) PullRequests(ctx context.Context) ([]PullRequest, error) {
	var queryPullRequests struct {
		Repository Repository `graphql:"repository(owner: $owner, name: $name)"`
	}

	variables := map[string]any{
		"owner":      githubv4.String(c.Owner),
		"name":       githubv4.String(c.Repo),
		"statusName": githubv4.String(c.TargetStatusName),
	}

	if err := c.gclient.Query(ctx, &queryPullRequests, variables); err != nil {
		return nil, fmt.Errorf("failed to list all pull requests: %w", err)
	}

	return queryPullRequests.Repository.PullRequests.Nodes, nil
}

func (c *Caretaker) PullRequest(ctx context.Context, prNumber int) (PullRequest, error) {
	var queryPullRequests struct {
		Repository struct {
			PullRequest PullRequest `graphql:"pullRequest(number: $pullNumber)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	variables := map[string]any{
		"owner":      githubv4.String(c.Owner),
		"name":       githubv4.String(c.Repo),
		"statusName": githubv4.String(c.TargetStatusName),
		"pullNumber": githubv4.Int(prNumber),
	}

	if err := c.gclient.Query(ctx, &queryPullRequests, variables); err != nil {
		return PullRequest{}, fmt.Errorf("failed to get pull requests: %w", err)
	}

	return queryPullRequests.Repository.PullRequest, nil
}

func (c *Caretaker) Issue(ctx context.Context, issueNumber int) (Issue, error) {
	var queryPullRequests struct {
		Repository struct {
			Issue Issue `graphql:"issue(number: $issueNumber)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	variables := map[string]any{
		"owner":      githubv4.String(c.Owner),
		"name":       githubv4.String(c.Repo),
		"statusName": githubv4.String(c.TargetStatusName),
		"pullNumber": githubv4.Int(issueNumber),
	}

	if err := c.gclient.Query(ctx, &queryPullRequests, variables); err != nil {
		return Issue{}, fmt.Errorf("failed to get issue: %w", err)
	}

	return queryPullRequests.Repository.Issue, nil
}

func (c *Caretaker) UpdateIssueStatus(ctx context.Context, issue Issue) error {
	if issue.Closed {
		c.log.Log("issue already closed, skip")

		return nil
	}

	if len(issue.ProjectsV2.Nodes) != 1 {
		c.log.Log("issues that are attached to more than one project are not supported ATM")

		return nil
	}

	// mutateIssueStatus sets the Status of an Issue to the desired option.
	var mutateIssueStatus struct {
		UpdateProjectV2ItemFieldValue struct {
			ProjectV2Item struct {
				ID githubv4.String
			} `graphql:"projectV2Item"` // value is case-sensitive and the default is projectV2item which is wrong.
		} `graphql:"updateProjectV2ItemFieldValue(input: $input)"`
	}

	project := issue.ProjectsV2.Nodes[0]

	if l := len(project.Field.ProjectV2SingleSelectField.Options); l != 1 {
		return fmt.Errorf("incorrect number of options found for name %s; want 1; got: %d", c.TargetStatusName, l)
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

	if err := c.gclient.Mutate(ctx, &mutateIssueStatus, input, nil); err != nil {
		return fmt.Errorf("failed to mutate issue: %w", err)
	}

	return nil
}

func (c *Caretaker) LeaveComment(ctx context.Context, prID githubv4.String, comment string) error {
	var leaveComment struct {
		AddComment struct {
			Subject struct {
				ID githubv4.String
			}
		} `graphql:"addComment(input: $input)"`
	}

	input := githubv4.AddCommentInput{
		SubjectID: prID,
		Body:      githubv4.String(comment),
	}

	if err := c.gclient.Mutate(ctx, &leaveComment, input, nil); err != nil {
		return fmt.Errorf("failed to leave comment on object: %w", err)
	}

	c.log.Debug("added comment with ID %s", leaveComment.AddComment.Subject.ID)

	return nil
}

func (c *Caretaker) queryLabelID(ctx context.Context, label string) (githubv4.String, error) {
	variables := map[string]any{
		"owner": githubv4.String(c.Owner),
		"name":  githubv4.String(c.Repo),
		"query": githubv4.String(label),
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

	if err := c.gclient.Query(ctx, &queryLabelID, variables); err != nil {
		return "", fmt.Errorf("failed to query for label id: %w", err)
	}

	if len(queryLabelID.Repository.Labels.Nodes) != 1 {
		return "", fmt.Errorf("expected a single label to be returned, got: %d", len(queryLabelID.Repository.Labels.Nodes))
	}

	return queryLabelID.Repository.Labels.Nodes[0].ID, nil
}

func (c *Caretaker) assignToOrganization(ctx context.Context, issue *Issue, projectNumber int) error {
	projectValues := map[string]any{
		"login":      githubv4.String(c.Owner),
		"number":     githubv4.Int(projectNumber),
		"statusName": githubv4.String(""),
	}

	var projectQuery struct {
		Organization struct {
			ProjectV2 ProjectV2 `graphql:"projectV2(number: $number)"`
		} `graphql:"organization(login: $login)"`
	}

	if err := c.gclient.Query(ctx, &projectQuery, projectValues); err != nil {
		return fmt.Errorf("failed to find project with number %d for owner %s: %w", projectNumber, c.Owner, err)
	}

	return c.assignIssueToProject(ctx, projectQuery.Organization.ProjectV2, issue)
}

func (c *Caretaker) assignToUser(ctx context.Context, issue *Issue, projectNumber int) error {
	projectValues := map[string]any{
		"login":      githubv4.String(c.Owner),
		"number":     githubv4.Int(projectNumber),
		"statusName": githubv4.String(""),
	}

	var projectQuery struct {
		User struct {
			ProjectV2 ProjectV2 `graphql:"projectV2(number: $number)"`
		} `graphql:"user(login: $login)"`
	}

	if err := c.gclient.Query(ctx, &projectQuery, projectValues); err != nil {
		return fmt.Errorf("failed to find project with number %d for owner %s: %w", projectNumber, c.Owner, err)
	}

	return c.assignIssueToProject(ctx, projectQuery.User.ProjectV2, issue)
}

func (c *Caretaker) assignIssueToProject(ctx context.Context, project ProjectV2, issue *Issue) error {
	c.log.Log(
		"assigning issue number %d with title %s to project number %d and title %s",
		issue.Number,
		issue.Title,
		project.Number,
		project.Title,
	)

	// mutateIssue sets the Status of an Issue to the desired option.
	var addProjectV2ItemByID struct {
		AddProjectV2ItemById struct { //nolint:stylecheck,revive // this needs to be Id.
			Item ProjectV2Item
		} `graphql:"addProjectV2ItemById(input: $input)"`
	}

	input := githubv4.AddProjectV2ItemByIdInput{
		ProjectID: project.ID,
		ContentID: issue.ID,
	}

	if err := c.gclient.Mutate(ctx, &addProjectV2ItemByID, input, nil); err != nil {
		return fmt.Errorf("failed to assign issue to project: %w", err)
	}

	return nil
}
