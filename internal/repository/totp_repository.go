package repository

import (
	"database/sql"
)

type TotpRepository interface {
	SaveSecret(userID string, secret string) error
	GetSecret(userID string) (string, error)
	EnableTOTP(userID string) error
}

type totpRepo struct {
	db *sql.DB
}

func NewTotpRepository(db *sql.DB) TotpRepository {
	return &totpRepo{db: db}
}

func (r *totpRepo) SaveSecret(userID string, secret string) error {
	_, err := r.db.Exec(`
		INSERT INTO user_totp(user_id, secret)
		VALUES ($1, $2)
		ON CONFLICT (user_id) DO UPDATE SET secret = EXCLUDED.secret
	`, userID, secret)
	return err
}

func (r *totpRepo) GetSecret(userID string) (string, error) {
	var secret string
	err := r.db.QueryRow(`SELECT secret FROM user_totp WHERE user_id = $1`, userID).Scan(&secret)
	if err != nil {
		return "", err
	}
	return secret, nil
}

func (r *totpRepo) EnableTOTP(userID string) error {
	_, err := r.db.Exec(`UPDATE users SET "enableTOTP" = TRUE WHERE id = $1`, userID)
	return err
}