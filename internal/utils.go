package internal

import (
	"time"

	"github.com/cenkalti/backoff"
)

func Retry(op func() error, maxWait time.Duration) error {
	if maxWait == 0 {
		maxWait = time.Minute
	}

	bo := backoff.NewExponentialBackOff()
	bo.MaxInterval = time.Second * 5
	bo.MaxElapsedTime = maxWait
	return backoff.Retry(op, bo)
}
