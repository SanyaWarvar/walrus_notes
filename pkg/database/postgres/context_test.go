package postgres

import (
	"context"
	"testing"
)

func TestContextManager_InjectAndExtractNilTx(t *testing.T) {
	tm := NewContextManager()

	ctx := tm.InjectTx(context.Background(), nil)
	extracted := tm.ExtractTx(ctx)

	if extracted != nil {
		t.Errorf("ExtractTx() = %v, want nil", extracted)
	}
}

func TestContextManager_ExtractTx_Empty(t *testing.T) {
	tm := NewContextManager()

	if got := tm.ExtractTx(context.Background()); got != nil {
		t.Errorf("ExtractTx() = %v, want nil", got)
	}
}
