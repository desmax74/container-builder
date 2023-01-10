//go:build integration_docker

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
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

func TestDockerIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(DockerTestSuite))
}

// --------------------------- TESTS -----------------

func (suite *DockerTestSuite) TestImagesOperationsOnDockerRegistryForTest() {
	registryContainer, err := GetRegistryContainer()
	assert.NotNil(suite.T(), registryContainer)
	assert.Nil(suite.T(), err)
	repos, err := registryContainer.GetRepositories()
	initialSize := len(repos)
	assert.Nil(suite.T(), err)

	pullErr := suite.Docker.PullImage(TEST_IMAGE + ":" + LATEST_TAG)
	if pullErr != nil {
		logrus.Infof("Pull Error:%s", pullErr)
	}
	assert.Nil(suite.T(), pullErr, "Pull image failed")
	time.Sleep(2 * time.Second) // Needed on CI
	assert.True(suite.T(), suite.LocalRegistry.IsImagePresent(TEST_IMAGE), "Test image not found in the registry after the pull")
	tagErr := suite.Docker.TagImage(TEST_IMAGE, TEST_IMAGE_LOCAL_TAG)
	if tagErr != nil {
		logrus.Infof("Tag Error:%s", tagErr)
	}

	assert.Nil(suite.T(), tagErr, "Tag image failed")
	//time.Sleep(2 * time.Second) // Needed on CI
	pushErr := suite.Docker.PushImage(TEST_IMAGE_LOCAL_TAG, REGISTRY_CONTAINER_URL_FROM_DOCKER_SOCKET, "", "")
	if pushErr != nil {
		logrus.Infof("Push Error:%s", pushErr)
	}

	assert.Nil(suite.T(), pushErr, "Push image in the DOcker container failed")
	//give the time to update the registry status
	//time.Sleep(2 * time.Second)
	repos, err = registryContainer.GetRepositories()
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), repos)
	assert.True(suite.T(), len(repos) == initialSize+1)

	digest, erroDIgest := registryContainer.Connection.ManifestDigest(TEST_IMAGE, LATEST_TAG)
	assert.Nil(suite.T(), erroDIgest)
	assert.NotNil(suite.T(), digest)
	assert.NotNil(suite.T(), registryContainer.DeleteImage(TEST_IMAGE, LATEST_TAG), "Delete Image not allowed")
}
