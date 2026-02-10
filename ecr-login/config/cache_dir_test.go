// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//	http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package config

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestGetCacheDir(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping XDG tests on Windows")
	}

	// Save original env vars
	oldAwsCacheDir := os.Getenv("AWS_ECR_CACHE_DIR")
	oldXdgCacheHome := os.Getenv("XDG_CACHE_HOME")
	oldHome := os.Getenv("HOME")

	// Cleanup after tests
	defer func() {
		if oldAwsCacheDir != "" {
			os.Setenv("AWS_ECR_CACHE_DIR", oldAwsCacheDir)
		} else {
			os.Unsetenv("AWS_ECR_CACHE_DIR")
		}
		if oldXdgCacheHome != "" {
			os.Setenv("XDG_CACHE_HOME", oldXdgCacheHome)
		} else {
			os.Unsetenv("XDG_CACHE_HOME")
		}
		if oldHome != "" {
			os.Setenv("HOME", oldHome)
		} else {
			os.Unsetenv("HOME")
		}
	}()

	tests := []struct {
		name           string
		awsCacheDir    string
		xdgCacheHome   string
		home           string
		expectedSuffix string
		description    string
	}{
		{
			name:           "AWS_ECR_CACHE_DIR takes highest priority",
			awsCacheDir:    "/custom/cache/dir",
			xdgCacheHome:   "/xdg/cache",
			home:           "/home/user",
			expectedSuffix: "/custom/cache/dir",
			description:    "Environment variable override",
		},
		{
			name:           "XDG_CACHE_HOME is used when set",
			awsCacheDir:    "",
			xdgCacheHome:   "/xdg/cache",
			home:           "/home/user",
			expectedSuffix: "/xdg/cache/docker-credential-ecr-login",
			description:    "XDG spec compliance",
		},
		{
			name:           "Fallback to ~/.cache when XDG not set",
			awsCacheDir:    "",
			xdgCacheHome:   "",
			home:           "/home/user",
			expectedSuffix: "/home/user/.cache/docker-credential-ecr-login",
			description:    "Default XDG path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.awsCacheDir != "" {
				os.Setenv("AWS_ECR_CACHE_DIR", tt.awsCacheDir)
			} else {
				os.Unsetenv("AWS_ECR_CACHE_DIR")
			}
			if tt.xdgCacheHome != "" {
				os.Setenv("XDG_CACHE_HOME", tt.xdgCacheHome)
			} else {
				os.Unsetenv("XDG_CACHE_HOME")
			}
			if tt.home != "" {
				os.Setenv("HOME", tt.home)
			} else {
				os.Unsetenv("HOME")
			}

			result := GetCacheDir()
			expected := tt.expectedSuffix

			if result != expected {
				t.Errorf("%s: GetCacheDir() = %v, want %v", tt.description, result, expected)
			}
		})
	}
}

func TestGetCacheDirWindows(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Skipping Windows-specific tests")
	}

	oldAwsCacheDir := os.Getenv("AWS_ECR_CACHE_DIR")
	defer func() {
		if oldAwsCacheDir != "" {
			os.Setenv("AWS_ECR_CACHE_DIR", oldAwsCacheDir)
		} else {
			os.Unsetenv("AWS_ECR_CACHE_DIR")
		}
	}()

	// Test Windows behavior with env override
	os.Setenv("AWS_ECR_CACHE_DIR", "C:\\custom\\cache")
	result := GetCacheDir()
	if result != "C:\\custom\\cache" {
		t.Errorf("Windows with AWS_ECR_CACHE_DIR: got %v, want C:\\custom\\cache", result)
	}

	os.Unsetenv("AWS_ECR_CACHE_DIR")
	result = GetCacheDir()
	if result != "~/.ecr" {
		t.Errorf("Windows without AWS_ECR_CACHE_DIR: got %v, want ~/.ecr", result)
	}
}

func TestGetLogDir(t *testing.T) {
	cacheDir := "/test/cache"
	oldCacheDir := os.Getenv("AWS_ECR_CACHE_DIR")

	os.Setenv("AWS_ECR_CACHE_DIR", cacheDir)
	defer func() {
		if oldCacheDir != "" {
			os.Setenv("AWS_ECR_CACHE_DIR", oldCacheDir)
		} else {
			os.Unsetenv("AWS_ECR_CACHE_DIR")
		}
	}()

	result := GetLogDir()
	expected := filepath.Join(cacheDir, "log")

	if result != expected {
		t.Errorf("GetLogDir() = %v, want %v", result, expected)
	}
}
