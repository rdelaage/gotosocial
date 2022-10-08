/*
   GoToSocial
   Copyright (C) 2021-2022 GoToSocial Authors admin@gotosocial.org

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

package bundb

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/superseriousbusiness/gotosocial/internal/db"
	"github.com/superseriousbusiness/gotosocial/internal/gtsmodel"
	"github.com/uptrace/bun"
)

type relationshipDB struct {
	conn *DBConn
}

func (r *relationshipDB) newBlockQ(block *gtsmodel.Block) *bun.SelectQuery {
	return r.conn.
		NewSelect().
		Model(block).
		Relation("Account").
		Relation("TargetAccount")
}

func (r *relationshipDB) newFollowQ(follow interface{}) *bun.SelectQuery {
	return r.conn.
		NewSelect().
		Model(follow).
		Relation("Account").
		Relation("TargetAccount")
}

func (r *relationshipDB) IsBlocked(ctx context.Context, account1 string, account2 string, eitherDirection bool) (bool, db.Error) {
	q := r.conn.
		NewSelect().
		TableExpr("? AS ?", bun.Ident("blocks"), bun.Ident("block")).
		Column("block.id")

	if eitherDirection {
		q = q.
			WhereGroup(" OR ", func(inner *bun.SelectQuery) *bun.SelectQuery {
				return inner.
					Where("? = ?", bun.Ident("block.account_id"), account1).
					Where("? = ?", bun.Ident("block.target_account_id"), account2)
			}).
			WhereGroup(" OR ", func(inner *bun.SelectQuery) *bun.SelectQuery {
				return inner.
					Where("? = ?", bun.Ident("block.account_id"), account2).
					Where("? = ?", bun.Ident("block.target_account_id"), account1)
			})
	} else {
		q = q.
			Where("? = ?", bun.Ident("block.account_id"), account1).
			Where("? = ?", bun.Ident("block.target_account_id"), account2)
	}

	return r.conn.Exists(ctx, q)
}

func (r *relationshipDB) GetBlock(ctx context.Context, account1 string, account2 string) (*gtsmodel.Block, db.Error) {
	block := &gtsmodel.Block{}

	q := r.newBlockQ(block).
		Where("? = ?", bun.Ident("block.account_id"), account1).
		Where("? = ?", bun.Ident("block.target_account_id"), account2)

	if err := q.Scan(ctx); err != nil {
		return nil, r.conn.ProcessError(err)
	}
	return block, nil
}

func (r *relationshipDB) GetRelationship(ctx context.Context, requestingAccount string, targetAccount string) (*gtsmodel.Relationship, db.Error) {
	rel := &gtsmodel.Relationship{
		ID: targetAccount,
	}

	// check if the requesting account follows the target account
	follow := &gtsmodel.Follow{}
	if err := r.conn.
		NewSelect().
		Model(follow).
		Column("follow.show_reblogs", "follow.notify").
		Where("? = ?", bun.Ident("follow.account_id"), requestingAccount).
		Where("? = ?", bun.Ident("follow.target_account_id"), targetAccount).
		Limit(1).
		Scan(ctx); err != nil {
		if err := r.conn.ProcessError(err); err != db.ErrNoEntries {
			return nil, fmt.Errorf("GetRelationship: error fetching follow: %s", err)
		}
		// no follow exists so these are all false
		rel.Following = false
		rel.ShowingReblogs = false
		rel.Notifying = false
	} else {
		// follow exists so we can fill these fields out...
		rel.Following = true
		rel.ShowingReblogs = *follow.ShowReblogs
		rel.Notifying = *follow.Notify
	}

	// check if the target account follows the requesting account
	followedByQ := r.conn.
		NewSelect().
		TableExpr("? AS ?", bun.Ident("follows"), bun.Ident("follow")).
		Column("follow.id").
		Where("? = ?", bun.Ident("follow.account_id"), targetAccount).
		Where("? = ?", bun.Ident("follow.target_account_id"), requestingAccount)
	followedBy, err := r.conn.Exists(ctx, followedByQ)
	if err != nil {
		return nil, fmt.Errorf("GetRelationship: error checking followedBy: %s", err)
	}
	rel.FollowedBy = followedBy

	// check if there's a pending following request from requesting account to target account
	requestedQ := r.conn.
		NewSelect().
		TableExpr("? AS ?", bun.Ident("follow_requests"), bun.Ident("follow_request")).
		Column("follow_request.id").
		Where("? = ?", bun.Ident("follow_request.account_id"), requestingAccount).
		Where("? = ?", bun.Ident("follow_request.target_account_id"), targetAccount)
	requested, err := r.conn.Exists(ctx, requestedQ)
	if err != nil {
		return nil, fmt.Errorf("GetRelationship: error checking requested: %s", err)
	}
	rel.Requested = requested

	// check if the requesting account is blocking the target account
	blockingQ := r.conn.
		NewSelect().
		TableExpr("? AS ?", bun.Ident("blocks"), bun.Ident("block")).
		Column("block.id").
		Where("? = ?", bun.Ident("block.account_id"), requestingAccount).
		Where("? = ?", bun.Ident("block.target_account_id"), targetAccount)
	blocking, err := r.conn.Exists(ctx, blockingQ)
	if err != nil {
		return nil, fmt.Errorf("GetRelationship: error checking blocking: %s", err)
	}
	rel.Blocking = blocking

	// check if the requesting account is blocked by the target account
	blockedByQ := r.conn.
		NewSelect().
		TableExpr("? AS ?", bun.Ident("blocks"), bun.Ident("block")).
		Column("block.id").
		Where("? = ?", bun.Ident("block.account_id"), targetAccount).
		Where("? = ?", bun.Ident("block.target_account_id"), requestingAccount)
	blockedBy, err := r.conn.Exists(ctx, blockedByQ)
	if err != nil {
		return nil, fmt.Errorf("GetRelationship: error checking blockedBy: %s", err)
	}
	rel.BlockedBy = blockedBy

	return rel, nil
}

func (r *relationshipDB) IsFollowing(ctx context.Context, sourceAccount *gtsmodel.Account, targetAccount *gtsmodel.Account) (bool, db.Error) {
	if sourceAccount == nil || targetAccount == nil {
		return false, nil
	}

	q := r.conn.
		NewSelect().
		TableExpr("? AS ?", bun.Ident("follows"), bun.Ident("follow")).
		Column("follow.id").
		Where("? = ?", bun.Ident("follow.account_id"), sourceAccount.ID).
		Where("? = ?", bun.Ident("follow.target_account_id"), targetAccount.ID)

	return r.conn.Exists(ctx, q)
}

func (r *relationshipDB) IsFollowRequested(ctx context.Context, sourceAccount *gtsmodel.Account, targetAccount *gtsmodel.Account) (bool, db.Error) {
	if sourceAccount == nil || targetAccount == nil {
		return false, nil
	}

	q := r.conn.
		NewSelect().
		TableExpr("? AS ?", bun.Ident("follow_requests"), bun.Ident("follow_request")).
		Column("follow_request.id").
		Where("? = ?", bun.Ident("follow_request.account_id"), sourceAccount.ID).
		Where("? = ?", bun.Ident("follow_request.target_account_id"), targetAccount.ID)

	return r.conn.Exists(ctx, q)
}

func (r *relationshipDB) IsMutualFollowing(ctx context.Context, account1 *gtsmodel.Account, account2 *gtsmodel.Account) (bool, db.Error) {
	if account1 == nil || account2 == nil {
		return false, nil
	}

	// make sure account 1 follows account 2
	f1, err := r.IsFollowing(ctx, account1, account2)
	if err != nil {
		return false, err
	}

	// make sure account 2 follows account 1
	f2, err := r.IsFollowing(ctx, account2, account1)
	if err != nil {
		return false, err
	}

	return f1 && f2, nil
}

func (r *relationshipDB) AcceptFollowRequest(ctx context.Context, originAccountID string, targetAccountID string) (*gtsmodel.Follow, db.Error) {
	var follow *gtsmodel.Follow

	if err := r.conn.RunInTx(ctx, func(tx bun.Tx) error {
		// get original follow request
		followRequest := &gtsmodel.FollowRequest{}
		if err := tx.
			NewSelect().
			Model(followRequest).
			Where("? = ?", bun.Ident("follow_request.account_id"), originAccountID).
			Where("? = ?", bun.Ident("follow_request.target_account_id"), targetAccountID).
			Scan(ctx); err != nil {
			return err
		}

		// create a new follow to 'replace' the request with
		follow = &gtsmodel.Follow{
			ID:              followRequest.ID,
			AccountID:       originAccountID,
			TargetAccountID: targetAccountID,
			URI:             followRequest.URI,
		}

		// if the follow already exists, just update the URI -- we don't need to do anything else
		if _, err := tx.
			NewInsert().
			Model(follow).
			On("CONFLICT (?,?) DO UPDATE set ? = ?", bun.Ident("account_id"), bun.Ident("target_account_id"), bun.Ident("uri"), follow.URI).
			Exec(ctx); err != nil {
			return err
		}

		// now remove the follow request
		if _, err := tx.
			NewDelete().
			TableExpr("? AS ?", bun.Ident("follow_requests"), bun.Ident("follow_request")).
			Where("? = ?", bun.Ident("follow_request.id"), followRequest.ID).
			Exec(ctx); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, r.conn.ProcessError(err)
	}

	// return the new follow
	return follow, nil
}

func (r *relationshipDB) RejectFollowRequest(ctx context.Context, originAccountID string, targetAccountID string) (*gtsmodel.FollowRequest, db.Error) {
	followRequest := &gtsmodel.FollowRequest{}

	if err := r.conn.RunInTx(ctx, func(tx bun.Tx) error {
		// get original follow request
		if err := tx.
			NewSelect().
			Model(followRequest).
			Where("? = ?", bun.Ident("follow_request.account_id"), originAccountID).
			Where("? = ?", bun.Ident("follow_request.target_account_id"), targetAccountID).
			Scan(ctx); err != nil {
			return err
		}

		// now delete it from the database by ID
		if _, err := tx.
			NewDelete().
			TableExpr("? AS ?", bun.Ident("follow_requests"), bun.Ident("follow_request")).
			Where("? = ?", bun.Ident("follow_request.id"), followRequest.ID).
			Exec(ctx); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, r.conn.ProcessError(err)
	}

	// return the deleted follow request
	return followRequest, nil
}

func (r *relationshipDB) GetAccountFollowRequests(ctx context.Context, accountID string) ([]*gtsmodel.FollowRequest, db.Error) {
	followRequests := []*gtsmodel.FollowRequest{}

	q := r.newFollowQ(&followRequests).
		Where("? = ?", bun.Ident("follow_request.target_account_id"), accountID).
		Order("follow_request.updated_at DESC")

	if err := q.Scan(ctx); err != nil {
		return nil, r.conn.ProcessError(err)
	}

	return followRequests, nil
}

func (r *relationshipDB) GetAccountFollows(ctx context.Context, accountID string) ([]*gtsmodel.Follow, db.Error) {
	follows := []*gtsmodel.Follow{}

	q := r.newFollowQ(&follows).
		Where("? = ?", bun.Ident("follow.account_id"), accountID).
		Order("follow.updated_at DESC")

	if err := q.Scan(ctx); err != nil {
		return nil, r.conn.ProcessError(err)
	}

	return follows, nil
}

func (r *relationshipDB) CountAccountFollows(ctx context.Context, accountID string, localOnly bool) (int, db.Error) {
	q := r.conn.
		NewSelect().
		TableExpr("? AS ?", bun.Ident("follows"), bun.Ident("follow"))

	if localOnly {
		q = q.
			Join("JOIN ? AS ? ON ? = ?", bun.Ident("accounts"), bun.Ident("account"), bun.Ident("follow.target_account_id"), bun.Ident("account.id")).
			Where("? = ?", bun.Ident("follow.account_id"), accountID).
			Where("? IS NULL", bun.Ident("account.domain"))
	} else {
		q = q.Where("? = ?", bun.Ident("follow.account_id"), accountID)
	}

	return q.Count(ctx)
}

func (r *relationshipDB) GetAccountFollowedBy(ctx context.Context, accountID string, localOnly bool) ([]*gtsmodel.Follow, db.Error) {
	follows := []*gtsmodel.Follow{}

	q := r.conn.
		NewSelect().
		Model(&follows).
		Order("follow.updated_at DESC")

	if localOnly {
		q = q.
			Join("JOIN ? AS ? ON ? = ?", bun.Ident("accounts"), bun.Ident("account"), bun.Ident("follow.account_id"), bun.Ident("account.id")).
			Where("? = ?", bun.Ident("follow.target_account_id"), accountID).
			Where("? IS NULL", bun.Ident("account.domain"))
	} else {
		q = q.Where("? = ?", bun.Ident("follow.target_account_id"), accountID)
	}

	err := q.Scan(ctx)
	if err != nil && err != sql.ErrNoRows {
		return nil, r.conn.ProcessError(err)
	}
	return follows, nil
}

func (r *relationshipDB) CountAccountFollowedBy(ctx context.Context, accountID string, localOnly bool) (int, db.Error) {
	q := r.conn.
		NewSelect().
		TableExpr("? AS ?", bun.Ident("follows"), bun.Ident("follow"))

	if localOnly {
		q = q.
			Join("JOIN ? AS ? ON ? = ?", bun.Ident("accounts"), bun.Ident("account"), bun.Ident("follow.account_id"), bun.Ident("account.id")).
			Where("? = ?", bun.Ident("follow.target_account_id"), accountID).
			Where("? IS NULL", bun.Ident("account.domain"))
	} else {
		q = q.Where("? = ?", bun.Ident("follow.target_account_id"), accountID)
	}

	return q.Count(ctx)
}
