package cli

import (
	"strings"
	"testing"
)

func TestValidateCreateTaskInput(t *testing.T) {
	tests := []struct {
		name    string
		title   string
		content string
		desc    string
		wantErr bool
	}{
		{
			name:    "all fields below limit",
			title:   strings.Repeat("a", 499),
			content: strings.Repeat("b", 499),
			desc:    strings.Repeat("c", 499),
			wantErr: false,
		},
		{
			name:    "title over limit",
			title:   strings.Repeat("a", 500),
			content: "ok",
			desc:    "ok",
			wantErr: true,
		},
		{
			name:    "content over limit with multibyte chars",
			title:   "ok",
			content: strings.Repeat("你", 500),
			desc:    "ok",
			wantErr: true,
		},
		{
			name:    "desc over limit",
			title:   "ok",
			content: "ok",
			desc:    strings.Repeat("d", 500),
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validateCreateTaskInput(tc.title, tc.content, tc.desc)
			if tc.wantErr && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("expected nil error, got %v", err)
			}
		})
	}
}
