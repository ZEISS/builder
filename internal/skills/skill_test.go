package skills_test

import (
	"strings"
	"testing"

	"github.com/zeiss/builder/internal/skills"

	"github.com/stretchr/testify/require"
)

func TestParserSkill(t *testing.T) {
	tests := []struct {
		name        string
		description string
		markdown    string
		expected    struct {
			Name          string         `yaml:"name"`
			Description   string         `yaml:"description"`
			License       string         `yaml:"license"`
			Compatibility string         `yaml:"compatibility"`
			Metadata      map[string]any `yaml:"metadata"`
			AllowedTools  string         `yaml:"allowed-tools"`
			Body          string         `yaml:"body"`
		}
		wantErr     bool
		expectedErr error
	}{
		{
			name:        "test skill",
			description: "test skill description",
			markdown: `
---
name: pdf-processing
description: Extract text and tables from PDF files, fill forms, merge documents.
license: Apache-2.0
compatibility: Linux, macOS
metadata:
  author: example-org
  version: "1.0"
allowed-tools: pdftk pdftotext
---
BODY`,
			expected: struct {
				Name          string         `yaml:"name"`
				Description   string         `yaml:"description"`
				License       string         `yaml:"license"`
				Compatibility string         `yaml:"compatibility"`
				Metadata      map[string]any `yaml:"metadata"`
				AllowedTools  string         `yaml:"allowed-tools"`
				Body          string         `yaml:"body"`
			}{
				Name:          "pdf-processing",
				Description:   "Extract text and tables from PDF files, fill forms, merge documents.",
				License:       "Apache-2.0",
				Compatibility: "Linux, macOS",
				Metadata: map[string]any{
					"author":  "example-org",
					"version": "1.0",
				},
				AllowedTools: "pdftk pdftotext",
				Body:         "BODY",
			},
			wantErr:     false,
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := skills.NewParser()
			s, err := p.Parse(strings.NewReader(tt.markdown))

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.expected.Body, s.Body())
			require.Equal(t, tt.expected.Name, s.Name())
			require.Equal(t, tt.expected.Description, s.Description())
			require.Equal(t, tt.expected.License, s.License())
			require.Equal(t, tt.expected.Compatibility, s.Compatibility())
			require.Equal(t, tt.expected.Metadata, s.Metadata())
			require.Equal(t, tt.expected.AllowedTools, s.AllowedTools())
		})
	}
}
