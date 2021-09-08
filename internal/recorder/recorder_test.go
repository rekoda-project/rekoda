package recorder

import (
	"testing"

	"github.com/stretchr/testify/assert"
	//	"github.com/wmw9/rekoda/internal/recorder"
)

func TestAddOnline(t *testing.T) {
	r := New()
	r.AddOnline("test")
	assert.Equal(t, "test", r.Online[0])

}

func TestIsOnline(t *testing.T) {
	r := New()
	r.AddOnline("test")
	assert.Equal(t, true, r.IsOnline("test"))
}

func TestRemoveOnline(t *testing.T) {
	r := New()
	r.AddOnline("test")
	r.RemoveOnline("test")
	assert.Equal(t, false, r.IsOnline("test"))
}
