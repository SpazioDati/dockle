package cache

import (
	"testing"

	"github.com/SpazioDati/dockle/pkg/log"
	"github.com/SpazioDati/dockle/pkg/types"
	deckodertypes "github.com/goodwithtech/deckoder/types"
)

func init() {
	// Initialize logger for tests
	log.InitLogger(false, false)
}

func TestAssess(t *testing.T) {
	tests := []struct {
		name          string
		fileMap       deckodertypes.FileMap
		expectedCount int
		expectedFiles []string
	}{
		{
			name: "detects suspicious directories anywhere in filesystem",
			fileMap: deckodertypes.FileMap{
				"root/.cache/pip/http-v2/1/2/3/file": {},
				"root/.aws/credentials":              {},
				"/root/.git/config":                  {},
				"home/user/.npm/registry":            {},
				".cache/direct":                      {},
			},
			expectedCount: 5,
			expectedFiles: []string{"root/.cache/pip/http-v2/1/2/3", "root/.aws", "/root/.git", "home/user/.npm", ".cache"},
		},
		{
			name: "detects required files by basename",
			fileMap: deckodertypes.FileMap{
				"root/Dockerfile":        {},
				"app/docker-compose.yml": {},
				"home/.vimrc":            {},
				"project/.DS_Store":      {},
			},
			expectedCount: 4,
			expectedFiles: []string{"root/Dockerfile", "app/docker-compose.yml", "home/.vimrc", "project/.DS_Store"},
		},
		{
			name: "deduplicates directories with multiple files",
			fileMap: deckodertypes.FileMap{
				"root/.cache/file1": {},
				"root/.cache/file2": {},
				"root/.cache/file3": {},
			},
			expectedCount: 1,
			expectedFiles: []string{"root/.cache"},
		},
		{
			name: "ignores uncontrollable directories",
			fileMap: deckodertypes.FileMap{
				"app/node_modules/.cache/file": {},
				"lib/vendor/.git/config":       {},
			},
			expectedCount: 0,
			expectedFiles: []string{},
		},
		{
			name: "handles edge cases",
			fileMap: deckodertypes.FileMap{
				"usr/bin/app":     {},
				"etc/config.conf": {},
			},
			expectedCount: 0,
			expectedFiles: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset detectedDir map before each test
			detectedDir = map[string]struct{}{}

			assessor := CacheAssessor{}
			assessments, err := assessor.Assess(tt.fileMap)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(assessments) != tt.expectedCount {
				t.Errorf("expected %d assessments, got %d", tt.expectedCount, len(assessments))
				for i, a := range assessments {
					t.Logf("  [%d] %s", i, a.Filename)
				}
			}

			// Verify expected files are present
			assessmentMap := make(map[string]*types.Assessment)
			for _, a := range assessments {
				assessmentMap[a.Filename] = a
			}

			for _, expectedFile := range tt.expectedFiles {
				if _, found := assessmentMap[expectedFile]; !found {
					t.Errorf("expected assessment for '%s' not found", expectedFile)
				}
			}

			// Verify all have correct code
			for _, a := range assessments {
				if a.Code != types.InfoDeletableFiles {
					t.Errorf("wrong code for %s: got %v, want %v", a.Filename, a.Code, types.InfoDeletableFiles)
				}
			}
		})
	}
}

func TestInIgnoreDir(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		expected bool
	}{
		{
			name:     "File in node_modules",
			filename: "app/node_modules/.cache/file",
			expected: true,
		},
		{
			name:     "File in vendor",
			filename: "lib/vendor/.git/config",
			expected: true,
		},
		{
			name:     "File not in ignore dir",
			filename: "app/.cache/file",
			expected: false,
		},
		{
			name:     "File with node_modules in name but not as directory",
			filename: "app/my-node_modules-file.txt",
			expected: false, // Only matches "node_modules/" with trailing slash
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := inIgnoreDir(tt.filename)
			if result != tt.expected {
				t.Errorf("inIgnoreDir(%s) = %v, expected %v", tt.filename, result, tt.expected)
			}
		})
	}
}
