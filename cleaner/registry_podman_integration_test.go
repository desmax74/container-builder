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

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// --------------------------- TEST SUITE -----------------

func TestRegistryPodmanIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(PodmanTestSuite))
}

// -------------------------------------- TESTS -----------------------------

func (suite *PodmanTestSuite) TestRegistry() {
	assert.Truef(suite.T(), suite.RegistryID != "", "Registry not started")
	assert.Truef(suite.T(), suite.LocalRegistry.IsRegistryImagePresent(), "Registry image not present")
	assert.Truef(suite.T(), suite.LocalRegistry.IsRegistryRunning(), "Registry container not running")
}

func (suite *PodmanTestSuite) TestPullTagPush() {

	assert.Truef(suite.T(), suite.RegistryID != "", "Registry not started")
	registryContainer, err := GetRegistryContainer()
	assert.Nil(suite.T(), err)
	repos := CheckRepositoriesSize(suite.T(), 0, registryContainer)
	logrus.Info("Empty Repo Size = ", len(repos))

	result := podmanPullTagPushOnRegistryContainer(suite)
	logrus.Info("Pull Tag and Push Image on Registry Container successful = ", result)

	// Give some time to the registry to refresh status
	time.Sleep(2 * time.Second)
	repos = CheckRepositoriesSize(suite.T(), 1, registryContainer)
	logrus.Info("Repo Size after pull image = ", len(repos))
}

func podmanPullTagPushOnRegistryContainer(suite *PodmanTestSuite) bool {
	podmanSocketConn, errSock := GetPodmanConnection()
	if errSock != nil {
		assert.FailNow(suite.T(), "Cant get podman socket")
	}
	p := Podman{Connection: podmanSocketConn}
	if !suite.LocalRegistry.IsRegistryImagePresent() {
		_, err := p.PullImage(TEST_IMAGE_TAG)
		if err != nil {
			assert.Fail(suite.T(), "Pull Image Failed", err)
			return false
		}
	}
	logrus.Info("Pull image")

	err := p.TagImage(TEST_REPO+TEST_IMAGE_TAG, LATEST_TAG, TEST_REGISTRY_REPO+TEST_IMAGE)
	if err != nil {
		assert.Fail(suite.T(), "Tag Image Failed", err)
		return false
	}
	logrus.Info("Tag image")

	err = p.PushImage(TEST_IMAGE_LOCAL_TAG, TEST_IMAGE_LOCAL_TAG, "", "")
	if err != nil {
		assert.Fail(suite.T(), "Push Image Failed", err)
		return false
	}
	logrus.Info("Push image")
	return true
}
