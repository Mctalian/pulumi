// Copyright 2016-2022, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build !xplatform_acceptance

package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//nolint:paralleltest // changes directory for process
func TestCreatingPolicyPackWithPromptedName(t *testing.T) {
	skipIfShortOrNoPulumiAccessToken(t)

	tempdir := tempProjectDir(t)
	chdir(t, tempdir)

	args := newPolicyArgs{
		templateNameOrURL: "aws-javascript",
	}

	err := runNewPolicyPack(context.Background(), args)
	assert.NoError(t, err)

	assert.FileExists(t, filepath.Join(tempdir, "PulumiPolicy.yaml"))
	assert.FileExists(t, filepath.Join(tempdir, "index.js"))
}

//nolint:paralleltest // changes directory for process
func TestInvalidPolicyPackTemplateName(t *testing.T) {
	skipIfShortOrNoPulumiAccessToken(t)

	// A template that will never exist.
	const nonExistantTemplate = "this-is-not-the-template-youre-looking-for"

	t.Run("RemoteTemplateNotFound", func(t *testing.T) {
		tempdir := tempProjectDir(t)
		chdir(t, tempdir)

		args := newPolicyArgs{
			templateNameOrURL: nonExistantTemplate,
		}

		err := runNewPolicyPack(context.Background(), args)
		assert.Error(t, err)
		assertNotFoundError(t, err)
	})

	t.Run("LocalTemplateNotFound", func(t *testing.T) {
		tempdir := tempProjectDir(t)
		chdir(t, tempdir)

		args := newPolicyArgs{
			generateOnly:      true,
			offline:           true,
			templateNameOrURL: nonExistantTemplate,
		}

		err := runNewPolicyPack(context.Background(), args)
		assert.Error(t, err)
		assertNotFoundError(t, err)
	})
}

func skipIfShortOrNoPulumiAccessToken(t *testing.T) {
	_, ok := os.LookupEnv("PULUMI_ACCESS_TOKEN")
	if !ok {
		t.Skipf("Skipping: PULUMI_ACCESS_TOKEN is not set")
	}
	if testing.Short() {
		t.Skip("Skipped in short test run")
	}
}

func chdir(t *testing.T, dir string) {
	cwd, err := os.Getwd()
	assert.NoError(t, err)
	assert.NoError(t, os.Chdir(dir)) // Set directory
	t.Cleanup(func() {
		assert.NoError(t, os.Chdir(cwd)) // Restore directory
		restoredDir, err := os.Getwd()
		assert.NoError(t, err)
		assert.Equal(t, cwd, restoredDir)
	})
}

func tempProjectDir(t *testing.T) string {
	t.Helper()

	dir := filepath.Join(t.TempDir(), genUniqueName(t))
	require.NoError(t, os.MkdirAll(dir, 0o700))
	return dir
}

func genUniqueName(t *testing.T) string {
	t.Helper()

	var bs [8]byte
	_, err := rand.Read(bs[:])
	require.NoError(t, err)

	return "test-" + hex.EncodeToString(bs[:])
}
