package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/dath-251-thuanle/file-sharing-web-backend2/internal/domain"
	"github.com/dath-251-thuanle/file-sharing-web-backend2/pkg/utils"
)

type SharedRepository interface {
	ShareFileWithUsers(ctx context.Context, fileID string, userIDs []string) error
	GetUsersSharedWith(ctx context.Context, fileID string) (*domain.Shared, error)
}

type sharedRepository struct {
	db *sql.DB
}

func NewSharedRepository(db *sql.DB) SharedRepository {
	return &sharedRepository{db: db}
}

func (r *sharedRepository) ShareFileWithUsers(ctx context.Context, fileID string, userIDs []string) error {
	if len(userIDs) == 0 {
		return nil
	}

	// 1. Chuẩn bị query
	var valuePlaceholders []string
	var args []interface{}

	i := 1
	for _, userID := range userIDs {
		// (user_id, file_id)
		valuePlaceholders = append(valuePlaceholders, fmt.Sprintf("($%d, $%d)", i, i+1))
		args = append(args, userID, fileID)
		i += 2
	}

	query := fmt.Sprintf(`
		INSERT INTO shared (user_id, file_id) 
		VALUES %s 
		ON CONFLICT (user_id, file_id) DO NOTHING
	`, strings.Join(valuePlaceholders, ", "))

	// 2. Thực thi query
	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute share file query: %w", err)
	}

	return nil
}

func (r *sharedRepository) GetUsersSharedWith(ctx context.Context, fileID string) (*domain.Shared, error) {
	// SELECT * FROM shared_with WHERE file_id = $1

	query := `
		SELECT user_id FROM shared WHERE file_id = $1
	`

	share := domain.Shared{
		FileId:  fileID,
		UserIds: make([]string, 0, 10),
	}

	rows, err := r.db.QueryContext(ctx, query, fileID)
	if err != nil {
		log.Println(err)
		return nil, utils.NewError("Failed to query for shared list", utils.ErrCodeInternal)
	}

	for rows.Next() {
		var userid_tmp string

		if err := rows.Scan(&userid_tmp); err != nil {
			log.Println(err)
			return nil, utils.NewError("Error when scanning query result.", utils.ErrCodeInternal)
		}

		share.UserIds = append(share.UserIds, userid_tmp)
	}

	return &share, nil
}
