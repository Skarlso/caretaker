package review

import (
	"context"
	"fmt"

	"github.com/shurcooL/githubv4"

	"github.com/skarlso/caretaker/pkg/client"
	"github.com/skarlso/caretaker/pkg/slash"
)

const (
	// Command defines the command this handler understands.
	Command = "/review"
	Status  = "In Review"

	statusKey = "status"
)

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
func (h *Handler) Execute(ctx context.Context, pullNumber int, _ string, args ...string) error {
	pr, err := h.client.PullRequest(ctx, pullNumber)
	if err != nil {
		return fmt.Errorf("failed to get related pull request: %w", err)
	}

	status := Status

	if len(args) > 0 {
		argMap, err := slash.ConvertArgs(args...)
		if err != nil {
			return fmt.Errorf("failed to convert arguments to command: %w", err)
		}

		if v, ok := argMap[statusKey]; ok {
			status = v
		}
	}

	// assign user to all related issues
	for _, issue := range pr.ClosingIssuesReferences.Nodes {
		if err := h.client.UpdateIssueStatus(ctx, issue, githubv4.String(status)); err != nil {
			return fmt.Errorf("failed to update issue into desired state %s: %w", status, err)
		}
	}

	return nil
}

func (h *Handler) Help() string {
	return "- `/review` set all attached issues to \"In Progress\", " +
		"to override the status use: status=\"Custom In Review\" as command argument"
}
