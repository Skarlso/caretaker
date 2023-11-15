package assign

import (
	"context"
	"fmt"

	"github.com/skarlso/caretaker/pkg/client"
	"github.com/skarlso/caretaker/pkg/slash"
)

// Command defines the command this handler understands.
const Command = "/assign"

type Handler struct {
	client client.Client
}

func NewHandler(client client.Client) *Handler {
	return &Handler{
		client: client,
	}
}

var _ slash.Command = &Handler{}

// Execute assigns the pull request and all related issues to the user.
func (h *Handler) Execute(ctx context.Context, pullNumber int, actor string, _ ...string) error {
	pr, err := h.client.PullRequest(ctx, pullNumber)
	if err != nil {
		return fmt.Errorf("failed to get related pull request: %w", err)
	}

	user, err := h.client.User(ctx, actor)
	if err != nil {
		return fmt.Errorf("failed to fetch user: %w", err)
	}

	// assign user to the PR
	if err := h.client.AssignUserToAssignable(ctx, user.ID, pr.ID); err != nil {
		return fmt.Errorf("failed to assign user to pull request: %w", err)
	}

	// assign user to all related issues
	for _, issue := range pr.ClosingIssuesReferences.Nodes {
		if err := h.client.AssignUserToAssignable(ctx, user.ID, issue.ID); err != nil {
			return fmt.Errorf("failed to assign user to issue: %w", err)
		}
	}

	return nil
}
