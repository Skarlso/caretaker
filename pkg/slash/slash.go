package slash

import (
	"context"
	"fmt"
	"strings"

	"github.com/shurcooL/githubv4"

	"github.com/skarlso/caretaker/pkg/client"
)

type Command interface {
	// Execute runs the respective comment. The ID can be obtained through GitHub action's context: github.
	// cmd can be used for further parsing arguments to the command.
	// GraphQL object https://docs.github.com/en/graphql/reference/objects#issuecomment
	Execute(ctx context.Context, pullNumber int, actor, cmd string) error
}

type Slash struct {
	supportedCommands map[string]Command
	client            client.Client
}

func NewSlashHandler(client client.Client) *Slash {
	return &Slash{
		supportedCommands: make(map[string]Command),
		client:            client,
	}
}

func (s *Slash) RegisterHandler(key string, cmd Command) {
	s.supportedCommands[key] = cmd
}

// Run runs a command parsed from a comment body.
// Every line is examined to be a possible command.
func (s *Slash) Run(ctx context.Context, pullNumber int, actor, commentID, commentBody string) error {
	split := strings.Split(commentBody, "\n")

	seen := false

	for _, cmd := range split {
		if cmd == "" || !strings.HasPrefix(cmd, "/") {
			continue
		}

		// add an eye if we found at least ONE command that can be executed.
		if !seen {
			if err := s.client.AddReaction(ctx, commentID, githubv4.ReactionContentEyes); err != nil {
				return fmt.Errorf("failed to add reaction to comment: %w", err)
			}

			seen = true
		}

		handler, ok := s.supportedCommands[cmd]
		if !ok {
			return fmt.Errorf("command handler not registered for command %s", cmd)
		}

		if err := handler.Execute(ctx, pullNumber, actor, cmd); err != nil {
			return fmt.Errorf("failed to run command: %w", err)
		}
	}

	// add a thumbs up if all commands ran successfully
	if err := s.client.AddReaction(ctx, commentID, githubv4.ReactionContentThumbsUp); err != nil {
		return fmt.Errorf("failed to add reaction to comment: %w", err)
	}

	// Once we are done, we approve.
	return nil
}
