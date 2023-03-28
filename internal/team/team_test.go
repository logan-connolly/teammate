package team

import (
	"testing"
)

func TestMemoryApplication(t *testing.T) {
	t.Run("Intialize team app", func(t *testing.T) {
		_, err := NewMemoryApplication()
		if err != nil {
			t.Errorf("Did not expect to get error: %v", err)
		}
	})
}
