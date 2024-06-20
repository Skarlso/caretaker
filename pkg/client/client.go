package client

import (
	"context"
	"fmt"

	"github.com/shurcooL/githubv4"

	"github.com/skarlso/caretaker/pkg/logger"
)

const itemPerPage = 100

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
	ID        githubv4.ID
	Number    githubv4.Int
	UpdatedAt githubv4.Date
	Closed    githubv4.Boolean
	Title     githubv4.String
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
	ProjectsV2   ProjectsV2   `graphql:"projectsV2(first: 10)"`
	ProjectItems ProjectItems `graphql:"projectItems(first: 20)"`
}

func (p PullRequest) GetTitle() githubv4.String {
	return p.Title
}

func (p PullRequest) GetID() githubv4.ID {
	return p.ID
}

func (p PullRequest) GetNumber() githubv4.Int {
	return p.Number
}

func (p PullRequest) GetProjectsV2() ProjectsV2 {
	return p.ProjectsV2
}

func (p PullRequest) GetProjectItems() ProjectItems {
	return p.ProjectItems
}

func (p PullRequest) IsClosed() bool {
	return bool(p.Closed)
}

// ProjectV2Item https://docs.github.com/en/graphql/reference/objects#projectv2item
type ProjectV2Item struct {
	ID               githubv4.String
	Project          ProjectV2
	FieldValueByName struct {
		ProjectV2SingleSelectField struct {
			Name githubv4.String
		} `graphql:"... on ProjectV2ItemFieldSingleSelectValue"`
	} `graphql:"fieldValueByName(name: \"Status\")"`
}

// The fields need to be none pointer types to unmarshal. Hence, we need to
// check types of the returned items.
const (
	IssueType       = "ISSUE"
	PullRequestType = "PULL_REQUEST"
)

// ProjectV2ItemWithIssueContent https://docs.github.com/en/graphql/reference/objects#projectv2item
type ProjectV2ItemWithIssueContent struct {
	ID        githubv4.String
	Project   ProjectV2
	Type      githubv4.String
	UpdatedAt githubv4.Date
	Content   struct {
		Issue       Issue       `graphql:"... on Issue"`
		PullRequest PullRequest `graphql:"... on PullRequest"`
	}
	FieldValueByName struct {
		ProjectV2SingleSelectField struct {
			Name githubv4.String
		} `graphql:"... on ProjectV2ItemFieldSingleSelectValue"`
	} `graphql:"fieldValueByName(name: \"Status\")"`
}

type ProjectsV2 struct {
	Nodes []ProjectV2
}

type ProjectItems struct {
	TotalCount githubv4.Int
	Nodes      []ProjectV2Item
}

// Issue https://docs.github.com/en/graphql/reference/objects#issue
type Issue struct {
	ID           githubv4.ID
	Closed       githubv4.Boolean
	Title        githubv4.String
	Number       githubv4.Int
	ProjectsV2   ProjectsV2   `graphql:"projectsV2(first: 10)"`
	ProjectItems ProjectItems `graphql:"projectItems(first: 20)"`
}

func (i Issue) GetTitle() githubv4.String {
	return i.Title
}

func (i Issue) GetID() githubv4.ID {
	return i.ID
}

func (i Issue) GetNumber() githubv4.Int {
	return i.Number
}

func (i Issue) GetProjectsV2() ProjectsV2 {
	return i.ProjectsV2
}

func (i Issue) GetProjectItems() ProjectItems {
	return i.ProjectItems
}

func (i Issue) IsClosed() bool {
	return bool(i.Closed)
}

// Comment https://docs.github.com/en/graphql/reference/objects#issuecomment
type Comment struct {
	ID          githubv4.ID
	Body        githubv4.String
	BodyText    githubv4.String
	Issue       Issue
	PullRequest PullRequest
}

// User https://docs.github.com/en/graphql/reference/objects#user
type User struct {
	ID githubv4.ID
}

// Label https://docs.github.com/en/graphql/reference/objects#label
type Label struct{}

// ProjectV2 https://docs.github.com/en/graphql/reference/objects#projectv2
type ProjectV2 struct {
	Title  githubv4.String
	ID     githubv4.String
	Number githubv4.Int
	Field  struct {
		ProjectV2SingleSelectField struct {
			ID      githubv4.String
			Options []struct {
				ID   githubv4.String
				Name githubv4.String
			}
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
	AddLabel(ctx context.Context, label string, id githubv4.ID) error
	RemoveLabel(ctx context.Context, label string, id githubv4.ID) error
	AssignIssueToProject(ctx context.Context, issueNumber, projectNumber int) error // Consider combining these two
	AssignUserToAssignable(ctx context.Context, userID, objectID githubv4.ID) error
	AddReaction(ctx context.Context, objectID githubv4.ID, reaction githubv4.ReactionContent) error
	LeaveComment(ctx context.Context, prID githubv4.ID, comment string) error
	PullRequests(ctx context.Context) ([]PullRequest, error)
	PullRequest(ctx context.Context, prNumber int) (PullRequest, error)
	Issue(ctx context.Context, issueNumber int) (Issue, error)
	ProjectItems(
		ctx context.Context,
		projectNumber int,
	) ([]ProjectV2ItemWithIssueContent, error)
	UpdateIssueStatus(ctx context.Context, issue GenericIssue, statusName githubv4.String, projectNumber int) (bool, error)
	User(ctx context.Context, username string) (User, error)
}

// Options are for Caretaker's functionality.
type Options struct {
	Repo           string
	Owner          string
	IsOrganization bool
	MoveClosed     bool
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

func (c *Caretaker) AddReaction(ctx context.Context, objectID githubv4.ID, reaction githubv4.ReactionContent) error {
	var addReaction struct {
		AddReaction struct {
			Subject struct {
				ID githubv4.ID
			}
		} `graphql:"addReaction(input: $input)"`
	}

	input := githubv4.AddReactionInput{
		SubjectID: objectID,
		Content:   reaction,
	}

	if err := c.gclient.Mutate(ctx, &addReaction, input, nil); err != nil {
		return fmt.Errorf("failed to add reaction on object: %w", err)
	}

	return nil
}

func (c *Caretaker) AddLabel(ctx context.Context, label string, id githubv4.ID) error {
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
		"owner": githubv4.String(c.Owner),
		"name":  githubv4.String(c.Repo),
		"issue": githubv4.Int(issueNumber),
	}

	if err := c.gclient.Query(ctx, &getIssueQuery, issueValues); err != nil {
		return fmt.Errorf("failed to find issue with number %d: %w", issueNumber, err)
	}

	if c.IsOrganization {
		return c.assignToOrganization(ctx, &getIssueQuery.Repository.Issue, projectNumber)
	}

	return c.assignToUser(ctx, &getIssueQuery.Repository.Issue, projectNumber)
}

func (c *Caretaker) AssignUserToAssignable(ctx context.Context, userID, objectID githubv4.ID) error {
	var addAssigneesToAssignable struct {
		AddAssigneesToAssignable struct {
			ClientMutationID githubv4.ID `graphql:"clientMutationId"`
		} `graphql:"addAssigneesToAssignable(input: $input)"`
	}

	input := githubv4.AddAssigneesToAssignableInput{
		AssignableID: objectID,
		AssigneeIDs:  []githubv4.ID{userID},
	}

	if err := c.gclient.Mutate(ctx, &addAssigneesToAssignable, input, nil); err != nil {
		return fmt.Errorf("failed to assign user to object: %w", err)
	}

	return nil
}

func (c *Caretaker) RemoveLabel(ctx context.Context, label string, id githubv4.ID) error {
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
		"owner": githubv4.String(c.Owner),
		"name":  githubv4.String(c.Repo),
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
		"pullNumber": githubv4.Int(prNumber),
	}

	if err := c.gclient.Query(ctx, &queryPullRequests, variables); err != nil {
		return PullRequest{}, fmt.Errorf("failed to get pull requests: %w", err)
	}

	return queryPullRequests.Repository.PullRequest, nil
}

func (c *Caretaker) Issue(ctx context.Context, issueNumber int) (Issue, error) {
	var queryIssue struct {
		Repository struct {
			Issue Issue `graphql:"issue(number: $number)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	variables := map[string]any{
		"owner":  githubv4.String(c.Owner),
		"name":   githubv4.String(c.Repo),
		"number": githubv4.Int(issueNumber),
	}

	if err := c.gclient.Query(ctx, &queryIssue, variables); err != nil {
		return Issue{}, fmt.Errorf("failed to get issue: %w", err)
	}

	return queryIssue.Repository.Issue, nil
}

type PageInfo struct {
	EndCursor   githubv4.String
	HasNextPage githubv4.Boolean
}

type projectQueryForUser struct {
	Entity struct {
		ProjectV2 struct {
			Items struct {
				Nodes    []ProjectV2ItemWithIssueContent
				PageInfo PageInfo
			} `graphql:"items(first: $first, after: $after)"`
		} `graphql:"projectV2(number: $number)"`
	} `graphql:"user(login: $login)"`
}

func (p *projectQueryForUser) Content() []ProjectV2ItemWithIssueContent {
	return p.Entity.ProjectV2.Items.Nodes
}

func (p *projectQueryForUser) PageInfo() PageInfo {
	return p.Entity.ProjectV2.Items.PageInfo
}

type projectQueryForOrganization struct {
	Entity struct {
		ProjectV2 struct {
			Items struct {
				Nodes    []ProjectV2ItemWithIssueContent
				PageInfo PageInfo
			} `graphql:"items(first: $first, after: $after)"`
		} `graphql:"projectV2(number: $number)"`
	} `graphql:"organization(login: $login)"`
}

func (p *projectQueryForOrganization) Content() []ProjectV2ItemWithIssueContent {
	return p.Entity.ProjectV2.Items.Nodes
}

func (p *projectQueryForOrganization) PageInfo() PageInfo {
	return p.Entity.ProjectV2.Items.PageInfo
}

// query unifies the two query types for user and organization.
type query interface {
	Content() []ProjectV2ItemWithIssueContent
	PageInfo() PageInfo
}

func (c *Caretaker) ProjectItems(ctx context.Context, projectNumber int) ([]ProjectV2ItemWithIssueContent, error) {
	var projectQuery query = &projectQueryForUser{}

	if c.IsOrganization {
		projectQuery = &projectQueryForOrganization{}
	}

	projectValues := map[string]any{
		"login":  githubv4.String(c.Owner),
		"number": githubv4.Int(projectNumber),
		"first":  githubv4.Int(itemPerPage),
		"after":  (*githubv4.String)(nil),
	}

	var result []ProjectV2ItemWithIssueContent

	for {
		if err := c.gclient.Query(ctx, projectQuery, projectValues); err != nil {
			return nil, fmt.Errorf("failed to get issue: %w", err)
		}

		result = append(result, projectQuery.Content()...)

		if !projectQuery.PageInfo().HasNextPage {
			break
		}

		projectValues["after"] = githubv4.NewString(projectQuery.PageInfo().EndCursor)
	}

	return result, nil
}

func (c *Caretaker) User(ctx context.Context, name string) (User, error) {
	var user struct {
		User User `graphql:"user(login: $name)"`
	}

	variables := map[string]any{
		"name": githubv4.String(name),
	}

	if err := c.gclient.Query(ctx, &user, variables); err != nil {
		return User{}, fmt.Errorf("failed to get user: %w", err)
	}

	return user.User, nil
}

type GenericIssue interface {
	GetTitle() githubv4.String
	GetID() githubv4.ID
	GetNumber() githubv4.Int
	GetProjectsV2() ProjectsV2
	GetProjectItems() ProjectItems
	IsClosed() bool
}

func (c *Caretaker) UpdateIssueStatus(
	ctx context.Context,
	issue GenericIssue,
	statusName githubv4.String,
	projectNumber int,
) (bool, error) {
	if issue.IsClosed() && !c.MoveClosed {
		c.log.Log("issue %s already closed, skip", issue.GetTitle())

		return false, nil
	}

	// mutateIssueStatus sets the Status of an Issue to the desired option.
	var mutateIssueStatus struct {
		UpdateProjectV2ItemFieldValue struct {
			ProjectV2Item struct {
				ID githubv4.String
			} `graphql:"projectV2Item"` // value is case-sensitive and the default is projectV2item which is wrong.
		} `graphql:"updateProjectV2ItemFieldValue(input: $input)"`
	}

	var updated bool

	for _, project := range issue.GetProjectsV2().Nodes {
		c.log.Debug("issue number %d and title %s on project: %s", issue.GetNumber(), issue.GetTitle(), project.Title)

		if projectNumber > 0 && int(project.Number) != projectNumber {
			c.log.Log("skipping project number %d as it wasn't requested for update", project.Number)

			continue
		}

		var projectItem ProjectV2Item

		// Select the right project item for the project we are checking.
		for _, i := range issue.GetProjectItems().Nodes {
			if i.Project.ID == project.ID {
				projectItem = i

				break
			}
		}

		// If there are no project items for this project that belong to this issue, it means
		// that the issue is not assigned to this project.
		if projectItem.ID == "" {
			c.log.Log("ProjectItem not found for project with number %d; skipping", project.Number)

			continue
		}

		if projectItem.FieldValueByName.ProjectV2SingleSelectField.Name == statusName {
			c.log.Log("ProjectItem already in request status, skipping mutation")

			continue
		}

		var option githubv4.String

		for _, o := range project.Field.ProjectV2SingleSelectField.Options {
			if o.Name == statusName {
				option = o.ID

				break
			}
		}

		// This project might not have the same statuses configured. We skip setting it in that case.
		// Note, we are doing this because an issue can be assigned to multiple projects.
		if option == "" {
			c.log.Log("status with name %s not found for project %d, skipping setting it", statusName, project.Number)

			continue
		}

		input := githubv4.UpdateProjectV2ItemFieldValueInput{
			ProjectID: githubv4.NewString(project.ID),
			ItemID:    githubv4.NewString(projectItem.ID),
			FieldID:   githubv4.NewString(project.Field.ProjectV2SingleSelectField.ID),
			Value: githubv4.ProjectV2FieldValue{
				SingleSelectOptionID: githubv4.NewString(option),
			},
		}

		if err := c.gclient.Mutate(ctx, &mutateIssueStatus, input, nil); err != nil {
			return false, fmt.Errorf("failed to mutate issue: %w", err)
		}

		c.log.Log("updated status on issue %s with number %d", issue.GetTitle(), issue.GetNumber())

		updated = true
	}

	return updated, nil
}

func (c *Caretaker) LeaveComment(ctx context.Context, prID githubv4.ID, comment string) error {
	var leaveComment struct {
		AddComment struct {
			Subject struct {
				ID githubv4.ID
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

func (c *Caretaker) queryLabelID(ctx context.Context, label string) (githubv4.ID, error) {
	variables := map[string]any{
		"owner": githubv4.String(c.Owner),
		"name":  githubv4.String(c.Repo),
		"query": githubv4.String(label),
	}

	var queryLabelID struct {
		Repository struct {
			Labels struct {
				Nodes []struct {
					ID githubv4.ID
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
		"login":  githubv4.String(c.Owner),
		"number": githubv4.Int(projectNumber),
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
