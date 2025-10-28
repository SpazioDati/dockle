package cache

import (
	"testing"
)

func TestIsUnderRootReqDir(t *testing.T) {
	tests := []struct {
		name     string
		dirName  string
		expected bool
	}{
		{
			name:     "Directory is root/.aws",
			dirName:  "root/.aws",
			expected: true,
		},
		{
			name:     "Directory is /root/.git",
			dirName:  "/root/.git",
			expected: true,
		},
		{
			name:     "Directory deep under root/.cache",
			dirName:  "root/.cache/pip/http-v2/1/2/3",
			expected: true,
		},
		{
			name:     "Directory not under root",
			dirName:  "/home/user/.aws",
			expected: false,
		},
		{
			name:     "Directory in root but not in reqDir",
			dirName:  "/root/somedir",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isUnderRootReqDir(tt.dirName, reqDirs)
			if result != tt.expected {
				t.Errorf("isUnderRootReqDir(%s) = %v, expected %v", tt.dirName, result, tt.expected)
			}
		})
	}
}

