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

package users

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	apiutil "github.com/superseriousbusiness/gotosocial/internal/api/util"
	"github.com/superseriousbusiness/gotosocial/internal/gtserror"
	"github.com/superseriousbusiness/gotosocial/internal/paging"
)

// OutboxGETHandler swagger:operation GET /users/{username}/outbox s2sOutboxGet
//
// Get the public outbox collection for an actor.
//
// Note that the response will be a Collection with a page as `first`, as shown below, if `page` is `false`.
//
// If `page` is `true`, then the response will be a single `CollectionPage` without the wrapping `Collection`.
//
// HTTP signature is required on the request.
//
//	---
//	tags:
//	- s2s/federation
//
//	produces:
//	- application/activity+json
//
//	parameters:
//	-
//		name: username
//		type: string
//		description: Username of the account.
//		in: path
//		required: true
//	-
//		name: page
//		type: boolean
//		description: Return response as a CollectionPage.
//		in: query
//		default: false
//	-
//		name: min_id
//		type: string
//		description: Minimum ID of the next status, used for paging.
//		in: query
//	-
//		name: max_id
//		type: string
//		description: Maximum ID of the next status, used for paging.
//		in: query
//
//	responses:
//		'200':
//			in: body
//			schema:
//				"$ref": "#/definitions/swaggerCollection"
//		'400':
//			description: bad request
//		'401':
//			description: unauthorized
//		'403':
//			description: forbidden
//		'404':
//			description: not found
func (m *Module) OutboxGETHandler(c *gin.Context) {
	// usernames on our instance are always lowercase
	requestedUsername := strings.ToLower(c.Param(UsernameKey))
	if requestedUsername == "" {
		err := errors.New("no username specified in request")
		apiutil.ErrorHandler(c, gtserror.NewErrorBadRequest(err, err.Error()), m.processor.InstanceGetV1)
		return
	}

	contentType, err := apiutil.NegotiateAccept(c, apiutil.ActivityPubOrHTMLHeaders...)
	if err != nil {
		apiutil.ErrorHandler(c, gtserror.NewErrorNotAcceptable(err, err.Error()), m.processor.InstanceGetV1)
		return
	}

	if contentType == string(apiutil.TextHTML) {
		// This isn't an ActivityPub request;
		// redirect to the user's profile.
		c.Redirect(http.StatusSeeOther, "/@"+requestedUsername)
		return
	}

	page, errWithCode := paging.ParseIDPage(c,
		1,  // min limit
		80, // max limit
		0,  // default = disabled
	)
	if errWithCode != nil {
		apiutil.ErrorHandler(c, errWithCode, m.processor.InstanceGetV1)
		return
	}

	resp, errWithCode := m.processor.Fedi().OutboxGet(c.Request.Context(), requestedUsername, page)
	if errWithCode != nil {
		apiutil.ErrorHandler(c, errWithCode, m.processor.InstanceGetV1)
		return
	}

	apiutil.JSONType(c, http.StatusOK, contentType, resp)
}
