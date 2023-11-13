package client

import (
	"context"

	"github.com/shurcooL/githubv4"
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

// Issue https://docs.github.com/en/graphql/reference/objects#issue
type Issue struct {
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

// Label https://docs.github.com/en/graphql/reference/objects#label
type Label struct{}

// ProjectV2 https://docs.github.com/en/graphql/reference/objects#projectv2
type ProjectV2 struct{}

type GraphQLClient interface {
	Query(ctx context.Context, q any, variables map[string]any) error
	Mutate(ctx context.Context, m any, input githubv4.Input, variables map[string]any) error
}

type Client interface {
	AddLabel()
	RemoveLabel()
	PullRequests()
	PullRequest()
	MutateIssue()
}
