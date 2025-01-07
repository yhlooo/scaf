package streams

import (
	"testing"
)

func TestNewGRPCStreamClientConnection(t *testing.T) {
	_ = NewGRPCStreamClientConnection("test", nil)
}
