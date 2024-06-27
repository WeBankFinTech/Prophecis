package v2

import (
	"context"
	"testing"
)

func Test1(t *testing.T) {
	ctx := context.WithValue(context.Background(), "X-Request-Id", "xxxyyy111")
	log := GetLogger(ctx)
	log.Infof("hello, world")
	log.Debugf("hello, world")
}
