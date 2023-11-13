package stale

import (
	"context"
	"fmt"
	"testing"

	"github.com/shurcooL/githubv4"
	"github.com/stretchr/testify/assert"

	"github.com/skarlso/caretaker/pkg/logger"
)

func TestChecker_Check(t *testing.T) {
	tests := []struct {
		name      string
		setupMock func() Client
		repo      string
		owner     string
		wantErr   assert.ErrorAssertionFunc
	}{
		{
			name: "",
			setupMock: func() Client {
				return &mockClient{}
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				if err != nil {
					return false
				}

				return true
			},
		},
	}
	for _, tt := range tests {
		log := logger.VerboseLogger{}
		t.Run(tt.name, func(t *testing.T) {
			m := tt.setupMock()

			c := &Checker{
				Options: Options{},
				client:  m,
				log:     &log,
			}
			err := c.Check(context.Background())

			tt.wantErr(t, err, fmt.Sprintf("check stale pull requests"))
		})
	}
}

type mockClient struct{}

var _ Client = &mockClient{}

func (m *mockClient) Query(ctx context.Context, q any, variables map[string]any) error {
	return nil
}

func (m *mockClient) Mutate(ctx context.Context, mut any, input githubv4.Input, variables map[string]any) error {
	return nil
}
