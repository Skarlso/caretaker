package stale

import (
	"context"
	"fmt"
	"testing"

	"github.com/skarlso/caretaker/pkg/client"
	"github.com/skarlso/caretaker/pkg/client/fakes"
	"github.com/stretchr/testify/assert"

	"github.com/skarlso/caretaker/pkg/logger"
)

func TestChecker_Check(t *testing.T) {
	tests := []struct {
		name      string
		setupMock func() client.Client
		repo      string
		owner     string
		wantErr   assert.ErrorAssertionFunc
	}{
		{
			name: "",
			setupMock: func() client.Client {
				return &fakes.FakeClient{}
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
				client: m,
				log:    &log,
			}
			err := c.Check(context.Background())

			tt.wantErr(t, err, fmt.Sprintf("check stale pull requests"))
		})
	}
}
