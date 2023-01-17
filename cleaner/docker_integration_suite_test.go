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
)

// --------------------------- TEST SUITE -----------------
type DockerTestSuite struct {
	suite.Suite
	LocalRegistry DockerLocalRegistry
	RegistryID    string
	Docker        Docker
}

func (suite *DockerTestSuite) SetupSuite() {
	dockerRegistryContainer, registryID, docker := SetupDockerSocket()
	if len(registryID) > 0 {
		suite.LocalRegistry = dockerRegistryContainer
		suite.RegistryID = registryID
		suite.Docker = docker
	} else {
		assert.FailNow(suite.T(), "Initialization failed %s", registryID)
	}
}

func (suite *DockerTestSuite) TearDownSuite() {
	registryID := suite.LocalRegistry.GetRegistryRunningID()
	if len(registryID) > 0 {
		DockerTearDown(suite.LocalRegistry)
	} else {
		suite.LocalRegistry.StopRegistry()
	}
	purged, err := suite.Docker.PurgeContainer("", REGISTRY)
	if err != nil {
		logrus.Errorf("Error during purged container in TearDown SUite %t", err)
	}
	logrus.Infof("Purged containers %t", purged)
}
