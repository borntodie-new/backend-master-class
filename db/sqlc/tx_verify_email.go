package sqlc

import (
	"context"
	"database/sql"
)

type VerifyEmailTxParams struct {
	EmailId    int64
	SecretCode string
}

type VerifyEmailTxResult struct {
	User        User
	VerifyEmail VerifyEmail
}

func (store *SQLStore) VerifyEmailTx(ctx context.Context, arg VerifyEmailTxParams) (VerifyEmailTxResult, error) {
	var result VerifyEmailTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// update user's is_verify field
		result.VerifyEmail, err = q.UpdateVerifyEmail(ctx, UpdateVerifyEmailParams{
			ID:         arg.EmailId,
			SecretCode: arg.SecretCode,
		})
		if err != nil {
			return err
		}
		// update verify_email's is_used field
		result.User, err = q.UpdateUser(ctx, UpdateUserParams{
			IsEmailVerified: sql.NullBool{
				Bool:  true,
				Valid: true,
			},
			Username: result.VerifyEmail.Username,
		})
		if err != nil {
			return err
		}
		return err
	})
	return result, err
}
