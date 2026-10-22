package mermaid_test

import (
	"testing"

	"github.com/capcom6/mergram-tg-bot/internal/bot/mermaid"
)

func TestImageFileName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		data     []byte
		expected string
	}{
		{
			name:     "jpeg",
			data:     []byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 0x4A, 0x46},
			expected: "diagram.jpg",
		},
		{
			name:     "png",
			data:     []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A},
			expected: "diagram.png",
		},
		{
			name:     "unknown content",
			data:     []byte("not an image"),
			expected: "diagram.jpg",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := mermaid.ImageFileNameForTest(tt.data); got != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}
