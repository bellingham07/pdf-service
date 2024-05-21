package model

import (
	"context"
	"gorm.io/gorm"
)

var _ UserModel = (*customUserModel)(nil)

type (
	// UserModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserModel.
	UserModel interface {
		userModel
		customUserLogicModel
	}

	customUserModel struct {
		*defaultUserModel
	}

	customUserLogicModel interface {
		FindByName(ctx context.Context, phone string) (*User, error)
	}
)

func (m *defaultUserModel) FindByName(ctx context.Context, phone string) (*User, error) {
	var resp User
	err := m.conn.WithContext(ctx).Model(&User{}).Where("phone = ?", phone).Take(&resp).Error
	switch err {
	case nil:
		return &resp, nil
	case gorm.ErrRecordNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// NewUserModel returns a model for the database table.
func NewUserModel(conn *gorm.DB) UserModel {
	return &customUserModel{
		defaultUserModel: newUserModel(conn),
	}
}

func (m *defaultUserModel) customCacheKeys(data *User) []string {
	if data == nil {
		return []string{}
	}
	return []string{}
}
