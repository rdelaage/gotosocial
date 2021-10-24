/*
   GoToSocial
   Copyright (C) 2021 GoToSocial Authors admin@gotosocial.org

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package typeutils_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/go-fed/activity/streams"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type InternalToASTestSuite struct {
	TypeUtilsTestSuite
}

func (suite *InternalToASTestSuite) TestAccountToAS() {
	testAccount := suite.testAccounts["local_account_1"] // take zork for this test

	asPerson, err := suite.typeconverter.AccountToAS(context.Background(), testAccount)
	suite.NoError(err)

	ser, err := streams.Serialize(asPerson)
	suite.NoError(err)

	bytes, err := json.Marshal(ser)
	suite.NoError(err)

	fmt.Println(string(bytes))
	// TODO: write assertions here, rn we're just eyeballing the output
}

func (suite *InternalToASTestSuite) TestOutboxToASCollection() {
	testAccount := suite.testAccounts["admin_account"]
	ctx := context.Background()

	collection, err := suite.typeconverter.OutboxToASCollection(ctx, testAccount.OutboxURI)
	suite.NoError(err)

	ser, err := streams.Serialize(collection)
	assert.NoError(suite.T(), err)

	bytes, err := json.Marshal(ser)
	suite.NoError(err)

	/*
		we want this:
		{
			"@context": "https://www.w3.org/ns/activitystreams",
			"first": "http://localhost:8080/users/admin/outbox?page=true",
			"id": "http://localhost:8080/users/admin/outbox",
			"type": "OrderedCollection"
		}
	*/

	suite.Equal(`{"@context":"https://www.w3.org/ns/activitystreams","first":"http://localhost:8080/users/admin/outbox?page=true","id":"http://localhost:8080/users/admin/outbox","type":"OrderedCollection"}`, string(bytes))
}

func (suite *InternalToASTestSuite) TestStatusesToASOutboxPage() {
	testAccount := suite.testAccounts["admin_account"]
	ctx := context.Background()

	// get public statuses from testaccount
	statuses, err := suite.db.GetAccountStatuses(ctx, testAccount.ID, 30, true, "", "", false, false, true)
	suite.NoError(err)

	page, err := suite.typeconverter.StatusesToASOutboxPage(ctx, testAccount.OutboxURI, "", "", statuses)
	suite.NoError(err)

	ser, err := streams.Serialize(page)
	assert.NoError(suite.T(), err)

	bytes, err := json.Marshal(ser)
	suite.NoError(err)

	/*

		we want this:

		{
			"@context": "https://www.w3.org/ns/activitystreams",
			"id": "http://localhost:8080/users/admin/outbox?page=true",
			"next": "http://localhost:8080/users/admin/outbox?page=true&max_id=01F8MH75CBF9JFX4ZAD54N0W0R",
			"orderedItems": [
				{
					"actor": "http://localhost:8080/users/admin",
					"cc": "http://localhost:8080/users/admin/followers",
					"id": "http://localhost:8080/users/admin/statuses/01F8MHAAY43M6RJ473VQFCVH37/activity",
					"object": "http://localhost:8080/users/admin/statuses/01F8MHAAY43M6RJ473VQFCVH37",
					"published": "2021-10-20T12:36:45Z",
					"to": "https://www.w3.org/ns/activitystreams#Public",
					"type": "Create"
				},
				{
					"actor": "http://localhost:8080/users/admin",
					"cc": "http://localhost:8080/users/admin/followers",
					"id": "http://localhost:8080/users/admin/statuses/01F8MH75CBF9JFX4ZAD54N0W0R/activity",
					"object": "http://localhost:8080/users/admin/statuses/01F8MH75CBF9JFX4ZAD54N0W0R",
					"published": "2021-10-20T11:36:45Z",
					"to": "https://www.w3.org/ns/activitystreams#Public",
					"type": "Create"
				}
			],
			"partOf": "http://localhost:8080/users/admin/outbox",
			"prev": "http://localhost:8080/users/admin/outbox?page=true&min_id=01F8MHAAY43M6RJ473VQFCVH37",
			"type": "OrderedCollectionPage"
		}
	*/

	suite.Equal(`{"@context":"https://www.w3.org/ns/activitystreams","id":"http://localhost:8080/users/admin/outbox?page=true","next":"http://localhost:8080/users/admin/outbox?page=true\u0026max_id=01F8MH75CBF9JFX4ZAD54N0W0R","orderedItems":[{"actor":"http://localhost:8080/users/admin","cc":"http://localhost:8080/users/admin/followers","id":"http://localhost:8080/users/admin/statuses/01F8MHAAY43M6RJ473VQFCVH37/activity","object":"http://localhost:8080/users/admin/statuses/01F8MHAAY43M6RJ473VQFCVH37","published":"2021-10-20T12:36:45Z","to":"https://www.w3.org/ns/activitystreams#Public","type":"Create"},{"actor":"http://localhost:8080/users/admin","cc":"http://localhost:8080/users/admin/followers","id":"http://localhost:8080/users/admin/statuses/01F8MH75CBF9JFX4ZAD54N0W0R/activity","object":"http://localhost:8080/users/admin/statuses/01F8MH75CBF9JFX4ZAD54N0W0R","published":"2021-10-20T11:36:45Z","to":"https://www.w3.org/ns/activitystreams#Public","type":"Create"}],"partOf":"http://localhost:8080/users/admin/outbox","prev":"http://localhost:8080/users/admin/outbox?page=true\u0026min_id=01F8MHAAY43M6RJ473VQFCVH37","type":"OrderedCollectionPage"}`, string(bytes))
}

func TestInternalToASTestSuite(t *testing.T) {
	suite.Run(t, new(InternalToASTestSuite))
}
