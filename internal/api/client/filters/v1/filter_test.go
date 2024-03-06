// GoToSocial
// Copyright (C) GoToSocial Authors admin@gotosocial.org
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package v1_test

import (
	"github.com/stretchr/testify/suite"
	filtersV1 "github.com/superseriousbusiness/gotosocial/internal/api/client/filters/v1"
	"github.com/superseriousbusiness/gotosocial/internal/config"
	"github.com/superseriousbusiness/gotosocial/internal/db"
	"github.com/superseriousbusiness/gotosocial/internal/email"
	"github.com/superseriousbusiness/gotosocial/internal/federation"
	"github.com/superseriousbusiness/gotosocial/internal/filter/visibility"
	"github.com/superseriousbusiness/gotosocial/internal/gtsmodel"
	"github.com/superseriousbusiness/gotosocial/internal/media"
	"github.com/superseriousbusiness/gotosocial/internal/processing"
	"github.com/superseriousbusiness/gotosocial/internal/state"
	"github.com/superseriousbusiness/gotosocial/internal/storage"
	"github.com/superseriousbusiness/gotosocial/internal/typeutils"
	"github.com/superseriousbusiness/gotosocial/testrig"
	"testing"
)

type FiltersTestSuite struct {
	suite.Suite
	db           db.DB
	storage      *storage.Driver
	mediaManager *media.Manager
	federator    *federation.Federator
	processor    *processing.Processor
	emailSender  email.Sender
	sentEmails   map[string]string
	state        state.State

	// standard suite models
	testTokens         map[string]*gtsmodel.Token
	testClients        map[string]*gtsmodel.Client
	testApplications   map[string]*gtsmodel.Application
	testUsers          map[string]*gtsmodel.User
	testAccounts       map[string]*gtsmodel.Account
	testStatuses       map[string]*gtsmodel.Status
	testFilters        map[string]*gtsmodel.Filter
	testFilterKeywords map[string]*gtsmodel.FilterKeyword
	testFilterStatuses map[string]*gtsmodel.FilterStatus

	// module being tested
	filtersModule *filtersV1.Module
}

func (suite *FiltersTestSuite) SetupSuite() {
	suite.testTokens = testrig.NewTestTokens()
	suite.testClients = testrig.NewTestClients()
	suite.testApplications = testrig.NewTestApplications()
	suite.testUsers = testrig.NewTestUsers()
	suite.testAccounts = testrig.NewTestAccounts()
	suite.testStatuses = testrig.NewTestStatuses()
	suite.testFilters = testrig.NewTestFilters()
	suite.testFilterKeywords = testrig.NewTestFilterKeywords()
	suite.testFilterStatuses = testrig.NewTestFilterStatuses()
}

func (suite *FiltersTestSuite) SetupTest() {
	suite.state.Caches.Init()
	testrig.StartNoopWorkers(&suite.state)

	testrig.InitTestConfig()
	config.Config(func(cfg *config.Configuration) {
		cfg.WebAssetBaseDir = "../../../../../web/assets/"
		cfg.WebTemplateBaseDir = "../../../../../web/templates/"
	})
	testrig.InitTestLog()

	suite.db = testrig.NewTestDB(&suite.state)
	suite.state.DB = suite.db
	suite.storage = testrig.NewInMemoryStorage()
	suite.state.Storage = suite.storage

	testrig.StartTimelines(
		&suite.state,
		visibility.NewFilter(&suite.state),
		typeutils.NewConverter(&suite.state),
	)

	suite.mediaManager = testrig.NewTestMediaManager(&suite.state)
	suite.federator = testrig.NewTestFederator(&suite.state, testrig.NewTestTransportController(&suite.state, testrig.NewMockHTTPClient(nil, "../../../../../testrig/media")), suite.mediaManager)
	suite.sentEmails = make(map[string]string)
	suite.emailSender = testrig.NewEmailSender("../../../../../web/template/", suite.sentEmails)
	suite.processor = testrig.NewTestProcessor(&suite.state, suite.federator, suite.emailSender, suite.mediaManager)
	suite.filtersModule = filtersV1.New(suite.processor)

	testrig.StandardDBSetup(suite.db, nil)
	testrig.StandardStorageSetup(suite.storage, "../../../../../testrig/media")
}

func (suite *FiltersTestSuite) TearDownTest() {
	testrig.StandardDBTeardown(suite.db)
	testrig.StandardStorageTeardown(suite.storage)
	testrig.StopWorkers(&suite.state)
}

func TestFiltersTestSuite(t *testing.T) {
	suite.Run(t, new(FiltersTestSuite))
}
