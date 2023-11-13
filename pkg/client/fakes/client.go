package fakes

import (
	"context"

	"github.com/shurcooL/githubv4"
	"github.com/skarlso/caretaker/pkg/client"
)

type FakeClient struct{}

func (m *FakeClient) AddLabel(ctx context.Context, label string, id githubv4.String) error {
	//TODO implement me
	panic("implement me")
}

func (m *FakeClient) RemoveLabel(ctx context.Context, label string, id githubv4.String) error {
	//TODO implement me
	panic("implement me")
}

func (m *FakeClient) PullRequests(ctx context.Context) ([]client.PullRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (m *FakeClient) PullRequest(ctx context.Context, prNumber int) (client.PullRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (m *FakeClient) MutateIssue(ctx context.Context, issue client.Issue) error {
	//TODO implement me
	panic("implement me")
}

func (m *FakeClient) LeaveComment(ctx context.Context, prID githubv4.String, comment string) error {
	//TODO implement me
	panic("implement me")
}

var _ client.Client = &FakeClient{}
