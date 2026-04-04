package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ GroupMembersModel = (*customGroupMembersModel)(nil)

type (
	// GroupMembersModel is an interface to be customized, add more methods here,
	// and implement the added methods in customGroupMembersModel.
	GroupMembersModel interface {
		groupMembersModel
		withSession(session sqlx.Session) GroupMembersModel
		FindOneByGroupIdUid(ctx context.Context, groupId, uid int64) (*GroupMembers, error)
		FindByGroupId(ctx context.Context, groupId int64) ([]*GroupMembers, error)
		FindByUid(ctx context.Context, uid int64) ([]*GroupMembers, error)
	}

	customGroupMembersModel struct {
		*defaultGroupMembersModel
	}
)

// NewGroupMembersModel returns a model for the database table.
func NewGroupMembersModel(conn sqlx.SqlConn) GroupMembersModel {
	return &customGroupMembersModel{
		defaultGroupMembersModel: newGroupMembersModel(conn),
	}
}

func (m *customGroupMembersModel) withSession(session sqlx.Session) GroupMembersModel {
	return NewGroupMembersModel(sqlx.NewSqlConnFromSession(session))
}

// FindByGroupId 查群组所有成员
func (m *defaultGroupMembersModel) FindByGroupId(ctx context.Context, groupId int64) ([]*GroupMembers, error) {
	var members []*GroupMembers
	query := fmt.Sprintf("SELECT %s FROM %s WHERE group_id=? ORDER BY role ASC, joined_at ASC",
		groupMembersRows, m.table)
	err := m.conn.QueryRowsCtx(ctx, &members, query, groupId)
	return members, err
}

// FindByUid 查用户加入的所有群
func (m *defaultGroupMembersModel) FindByUid(ctx context.Context, uid int64) ([]*GroupMembers, error) {
	var members []*GroupMembers
	query := fmt.Sprintf("SELECT %s FROM %s WHERE uid=?",
		groupMembersRows, m.table)
	err := m.conn.QueryRowsCtx(ctx, &members, query, uid)
	return members, err
}
