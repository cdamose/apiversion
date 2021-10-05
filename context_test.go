package apiversion

import (
	"context"
	"testing"
)

func TestContext(t *testing.T) {
	v1 := new(Version)
	ctx := NewContext(context.Background(), v1)
	v2 := FromContext(ctx)
	if v1 != v2 {
		t.Fatal("failed to get version from context")
	}
}
