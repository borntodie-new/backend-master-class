package sqlc

import (
	"context"
)

type CreateUserTxParams struct {
	// CreateUserParams 基本就是创建用户的请求参数
	CreateUserParams
	// AfterCreate 向外暴露一个函数，当用户新建完成后，需要执行的操作由用户自定义
	AfterCreate func(user User) error
}

type CreateUserTxResult struct {
	User User
}

func (store *SQLStore) CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error) {
	var result CreateUserTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		result.User, err = q.CreateUser(ctx, arg.CreateUserParams)
		if err != nil {
			return err
		}
		return arg.AfterCreate(result.User)
	})
	return result, err
}
