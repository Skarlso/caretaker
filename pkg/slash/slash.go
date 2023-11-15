package slash

import (
	"context"
	"fmt"
	"strings"
)

type Command interface {
	// Execute runs the respective comment. The ID can be obtained through GitHub action's context: github.
	// GraphQL object https://docs.github.com/en/graphql/reference/objects#issuecomment
	Execute(ctx context.Context, pullNumber int, actor, commentBody string) error
}

type Slash struct {
	supportedCommands map[string]Command
}

func NewSlashHandler() *Slash {
	return &Slash{
		supportedCommands: make(map[string]Command),
	}
}

func (s *Slash) RegisterHandler(key string, cmd Command) {
	s.supportedCommands[key] = cmd
}

// Run runs a command parsed from a comment body. The comment MUST NOT contain anything else but the command.
func (s *Slash) Run(ctx context.Context, pullNumber int, actor, commentBody string) error {
	split := strings.Split(commentBody, "\n")
	for _, cmd := range split {
		if cmd == "" || !strings.HasPrefix(cmd, "/") {
			return fmt.Errorf("invalid comment format, expected comment to start with / but got: %s", cmd)
		}

		handler, ok := s.supportedCommands[cmd]
		if !ok {
			return fmt.Errorf("command handler not registered for command %s", cmd)
		}

		if err := handler.Execute(ctx, pullNumber, actor, cmd); err != nil {
			return fmt.Errorf("failed to run command: %w", err)
		}
	}

	return nil
}
