//go:build integration_podman

/*
 * Copyright 2022 Red Hat, Inc. and/or its affiliates.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package cleaner

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestPodmanIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(PodmanTestSuite))
}

// --------------------------- TESTS -----------------

func (suite *PodmanTestSuite) TestImagesOperationsOnPodmanRegistryForTest() {
	registryContainer, err := GetRegistryContainer()
	assert.NotNil(suite.T(), registryContainer)
	assert.Nil(suite.T(), err)
	repos, err := registryContainer.GetRepositories()
	assert.Nil(suite.T(), err)
	assert.True(suite.T(), len(repos) == 0)
	_, err = suite.Podman.PullImage(TEST_IMAGE)
	assert.Nil(suite.T(), err, "Pull image failed")
	assert.Nil(suite.T(), suite.Podman.TagImage(TEST_IMAGE_TAG, TEST_IMAGE_LOCAL_TAG, TEST_REGISTRY_REPO), "Tag image failed")
	assert.Nil(suite.T(), suite.Podman.PushImage(TEST_IMAGE_LOCAL_TAG, REGISTRY_CONTAINER_URL_FROM_DOCKER_SOCKET, "", ""), "Push image in the Podman container failed")
	//give the time to update the registry status
	time.Sleep(2 * time.Second)
	repos, err = registryContainer.GetRepositories()
	assert.Nil(suite.T(), err)
	assert.True(suite.T(), len(repos) == 1)

	digest, erroDIgest := registryContainer.Connection.ManifestDigest(TEST_IMAGE, LATEST_TAG)
	assert.Nil(suite.T(), erroDIgest)
	assert.NotNil(suite.T(), digest)
	assert.NotNil(suite.T(), registryContainer.DeleteImage(TEST_IMAGE, LATEST_TAG), "Delete Image not allowed")
}
