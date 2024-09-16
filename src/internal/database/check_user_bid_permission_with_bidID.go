package database

import (
	"context"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725725133-team-77957/zadanie-6105/src/internal/api/handlers"
)

// CheckUserBidPermissionWithBidID проверяет что пользователь - автор предложения или ответственный за организацию
func (db Database) CheckUserBidPermissionWithBidID(ctx context.Context, username string, bidId string) error {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM bid b
			LEFT JOIN employee e ON (
				(b.bid_author_type = 'User' AND b.bid_author_id = e.id) -- Если пользователь является автором предложения
				OR 
				(b.bid_author_type = 'Organization' AND EXISTS (  -- Если организация автор, проверить, является ли пользователь ответственным
					SELECT 1 
					FROM organization_responsible ore
					JOIN tender t ON ore.organization_id = t.organization_id
					WHERE ore.user_id = e.id
					AND t.tender_id = b.tender_id
				))
			)
			WHERE e.username = $1
			AND b.bid_id = $2
		);
    `

	var exists bool

	err := db.db.QueryRowContext(ctx, query, username, bidId).Scan(&exists)
	if err != nil {
		return handlers.ErrMsgNotPermission
	}

	return nil
}
