package scanproject

import (
	"context"
	"time"

	"github.com/skarlso/caretaker/pkg/client"
	"github.com/skarlso/caretaker/pkg/logger"
)

type Options struct {
	ProjectNumber   int
	Interval        time.Duration
	DisableComments bool
	FromStatus      string
	ToStatus        string
}

type Scanner struct {
	Options

	client client.Client
	log    logger.Logger
}

func NewScanner(log logger.Logger, client client.Client, opts Options) *Scanner {
	return &Scanner{
		log:     log,
		client:  client,
		Options: opts,
	}
}

// ScanIssues checks if any issues of a project should be moved into a different column based on ToStatus.
func (c *Scanner) ScanIssues(ctx context.Context) error {

	return nil
}
