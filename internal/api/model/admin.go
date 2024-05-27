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

package model

// AdminAccountInfo models the admin view of an account's details.
//
// swagger:model adminAccountInfo
type AdminAccountInfo struct {
	// The ID of the account in the database.
	// example: 01GQ4PHNT622DQ9X95XQX4KKNR
	ID string `json:"id"`
	// The username of the account.
	// example: dril
	Username string `json:"username"`
	// The domain of the account.
	// Null for local accounts.
	// example: example.org
	Domain *string `json:"domain"`
	// When the account was first discovered. (ISO 8601 Datetime)
	// example: 2021-07-30T09:20:25+00:00
	CreatedAt string `json:"created_at"`
	// The email address associated with the account.
	// Empty string for remote accounts or accounts with
	// no known email address.
	// example: someone@somewhere.com
	Email string `json:"email"`
	// The IP address last used to login to this account.
	// Null if not known.
	// example: 192.0.2.1
	IP *string `json:"ip"`
	// All known IP addresses associated with this account.
	// NOT IMPLEMENTED (will always be empty array).
	// example: []
	IPs []interface{} `json:"ips"`
	// The locale of the account. (ISO 639 Part 1 two-letter language code)
	// example: en
	Locale string `json:"locale"`
	// The reason given when signing up.
	// Null if no reason / remote account.
	// example: Pleaaaaaaaaaaaaaaase!!
	InviteRequest *string `json:"invite_request"`
	// The current role of the account.
	Role AccountRole `json:"role"`
	// Whether the account has confirmed their email address.
	Confirmed bool `json:"confirmed"`
	// Whether the account is currently approved.
	Approved bool `json:"approved"`
	// Whether the account is currently disabled.
	Disabled bool `json:"disabled"`
	// Whether the account is currently silenced
	Silenced bool `json:"silenced"`
	// Whether the account is currently suspended.
	Suspended bool `json:"suspended"`
	// User-level information about the account.
	Account *Account `json:"account"`
	// The ID of the application that created this account.
	CreatedByApplicationID string `json:"created_by_application_id,omitempty"`
	// The ID of the account that invited this user
	InvitedByAccountID string `json:"invited_by_account_id,omitempty"`
}

// AdminReport models the admin view of a report.
//
// swagger:model adminReport
type AdminReport struct {
	// ID of the report.
	// example: 01FBVD42CQ3ZEEVMW180SBX03B
	ID string `json:"id"`
	// Whether an action has been taken by an admin in response to this report.
	// example: false
	ActionTaken bool `json:"action_taken"`
	// If an action was taken, at what time was this done? (ISO 8601 Datetime)
	// Will be null if not set / no action yet taken.
	// example: 2021-07-30T09:20:25+00:00
	ActionTakenAt *string `json:"action_taken_at"`
	// Under what category was this report created?
	// example: spam
	Category string `json:"category"`
	// Comment submitted when the report was created.
	// Will be empty if no comment was submitted.
	// example: This person has been harassing me.
	Comment string `json:"comment"`
	// Bool to indicate that report should be federated to remote instance.
	// example: true
	Forwarded bool `json:"forwarded"`
	// The date when this report was created (ISO 8601 Datetime).
	// example: 2021-07-30T09:20:25+00:00
	CreatedAt string `json:"created_at"`
	// Time of last action on this report (ISO 8601 Datetime).
	// example: 2021-07-30T09:20:25+00:00
	UpdatedAt string `json:"updated_at"`
	// The account that created the report.
	Account *AdminAccountInfo `json:"account"`
	// Account that was reported.
	TargetAccount *AdminAccountInfo `json:"target_account"`
	// The account assigned to handle the report.
	// Null if no account assigned.
	AssignedAccount *AdminAccountInfo `json:"assigned_account"`
	// Account that took admin action (if any).
	// Null if no action (yet) taken.
	ActionTakenByAccount *AdminAccountInfo `json:"action_taken_by_account"`
	// Array of  statuses that were submitted along with this report.
	// Will be empty if no status IDs were submitted with the report.
	Statuses []*Status `json:"statuses"`
	// Array of rules that were broken according to this report.
	// Will be empty if no rule IDs were submitted with the report.
	Rules []*InstanceRule `json:"rules"`
	// If an action was taken, what comment was made by the admin on the taken action?
	// Will be null if not set / no action yet taken.
	// example: Account was suspended.
	ActionTakenComment *string `json:"action_taken_comment"`
}

// AdminReportResolveRequest can be submitted along with a POST to /api/v1/admin/reports/{id}/resolve
//
// swagger:ignore
type AdminReportResolveRequest struct {
	// Comment to show to the creator of the report when an admin marks it as resolved.
	ActionTakenComment *string `form:"action_taken_comment" json:"action_taken_comment" xml:"action_taken_comment"`
}

// AdminEmoji models the admin view of a custom emoji.
//
// swagger:model adminEmoji
type AdminEmoji struct {
	Emoji
	// The ID of the emoji.
	// example: 01GEM7SFDZ7GZNRXFVZ3X4E4N1
	ID string `json:"id"`
	// True if this emoji has been disabled by an admin action.
	// example: false
	Disabled bool `json:"disabled"`
	// The domain from which the emoji originated. Only defined for remote domains, otherwise key will not be set.
	//
	// example: example.org
	Domain string `json:"domain,omitempty"`
	// Time when the emoji image was last updated.
	// example: 2022-10-05T09:21:26.419Z
	UpdatedAt string `json:"updated_at"`
	// The total file size taken up by the emoji in bytes, including static and animated versions.
	// example: 69420
	TotalFileSize int `json:"total_file_size"`
	// The MIME content type of the emoji.
	// example: image/png
	ContentType string `json:"content_type"`
	// The ActivityPub URI of the emoji.
	// example: https://example.org/emojis/016T5Q3SQKBT337DAKVSKNXXW1
	URI string `json:"uri"`
}

// AdminActionRequest models a request
// for an admin action to be performed.
//
// swagger:ignore
type AdminActionRequest struct {
	// Category of the target entity.
	Category string `form:"-" json:"-" xml:"-"`
	// Type of admin action to take. One of disable, silence, suspend.
	Type string `form:"type" json:"type" xml:"type"`
	// Text describing why an action was taken.
	Text string `form:"text" json:"text" xml:"text"`
	// ID of the target entity.
	TargetID string `form:"-" json:"-" xml:"-"`
}

// AdminActionResponse models the server
// response to an admin action.
//
// swagger:model adminActionResponse
type AdminActionResponse struct {
	// Internal ID of the action.
	//
	// example: 01H9QG6TZ9W5P0402VFRVM17TH
	ActionID string `json:"action_id"`
}

// MediaCleanupRequest models admin media cleanup parameters
//
// swagger:parameters mediaCleanup
type MediaCleanupRequest struct {
	// Number of days of remote media to keep. Native values will be treated as 0.
	// If value is not specified, the value of media-remote-cache-days in the server config will be used.
	RemoteCacheDays *int `form:"remote_cache_days" json:"remote_cache_days" xml:"remote_cache_days"`
}

// AdminSendTestEmailRequest models a test email send request (woah).
type AdminSendTestEmailRequest struct {
	// Email address to send the test email to.
	Email string `form:"email" json:"email"`
	// Optional message to include in the test email.
	Message string `form:"message" json:"message"`
}

type AdminInstanceRule struct {
	ID        string `json:"id"`         // id of this item in the database
	CreatedAt string `json:"created_at"` // when was item created
	UpdatedAt string `json:"updated_at"` // when was item last updated
	Text      string `json:"text"`       // text content of the rule
}

// DebugAPUrlResponse provides detailed debug
// information for an AP URL dereference request.
//
// swagger:model debugAPUrlResponse
type DebugAPUrlResponse struct {
	// Remote AP URL that was requested.
	RequestURL string `json:"request_url"`
	// HTTP headers used in the outgoing request.
	RequestHeaders map[string][]string `json:"request_headers"`
	// HTTP headers returned from the remote instance.
	ResponseHeaders map[string][]string `json:"response_headers"`
	// HTTP response code returned from the remote instance.
	ResponseCode int `json:"response_code"`
	// Body returned from the remote instance.
	// Will be stringified bytes; may be JSON,
	// may be an error, may be both!
	ResponseBody string `json:"response_body"`
}

// AdminGetAccountsRequest models a request
// to get an admin view of one or more
// accounts using given parameters.
//
// swagger:ignore
type AdminGetAccountsRequest struct {
	// Filter for `local` or `remote` accounts.
	Origin string
	// Filter for `active`, `pending`, `disabled`,
	// `silenced`, or `suspended` accounts.
	Status string
	// Filter for accounts with staff perms
	// (users that can manage reports).
	Permissions string
	// Filter for users with these roles.
	RoleIDs []string
	// Lookup users invited by the account with this ID.
	InvitedBy string
	// Search for the given username.
	Username string
	// Search for the given display name.
	DisplayName string
	// Filter by the given domain.
	ByDomain string
	// Lookup a user with this email.
	Email string
	// Lookup users with this IP address.
	IP string
	// API version to use for this request (1 or 2).
	// Set internally, not by callers.
	APIVersion int
}

// AdminAccountRejectRequest models a
// request to deny a new account sign-up.
//
// swagger:ignore
type AdminAccountRejectRequest struct {
	// Comment to leave on why the account was denied.
	// The comment will be visible to admins only.
	PrivateComment string `form:"private_comment" json:"private_comment"`
	// Message to include in email to applicant.
	// Will be included only if send_email is true.
	Message string `form:"message" json:"message"`
	// Send an email to the applicant informing
	// them that their sign-up has been rejected.
	SendEmail bool `form:"send_email" json:"send_email"`
}
