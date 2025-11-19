package service

import (
	"encoding/base64"
	"fmt"

	"github.com/dath-251-thuanle/file-sharing-web-backend2/internal/repository"
	"github.com/pquerna/otp/totp"
	"github.com/skip2/go-qrcode"
)

type TOTPSetupResponse struct {
	Secret string `json:"secret"`
	QRCode string `json:"qrCode"`
}

type TotpService interface {
	SetupTOTP(userID string) (*TOTPSetupResponse, error)
	VerifyTOTP(userID string, code string) (bool, error)
}

type totpService struct {
	Repo   repository.TotpRepository
	Issuer string
}

func NewTotpService(repo repository.TotpRepository, issuer string) TotpService {
	return &totpService{
		Repo:   repo,
		Issuer: issuer,
	}
}

func (s *totpService) SetupTOTP(userID string) (*TOTPSetupResponse, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      s.Issuer,
		AccountName: fmt.Sprintf("user-%s", userID),
	})
	if err != nil {
		return nil, err
	}

	secret := key.Secret()
	otpURL := key.URL()

	if err := s.Repo.SaveSecret(userID, secret); err != nil {
		return nil, err
	}

	png, err := qrcode.Encode(otpURL, qrcode.Medium, 256)
	if err != nil {
		return nil, err
	}
	qrBase64 := "data:image/png;base64," + base64.StdEncoding.EncodeToString(png)

	return &TOTPSetupResponse{
		Secret: secret,
		QRCode: qrBase64,
	}, nil
}

func (s *totpService) VerifyTOTP(userID string, code string) (bool, error) {
	secret, err := s.Repo.GetSecret(userID)
	if err != nil {
		return false, err
	}

	valid := totp.Validate(code, secret)

	if valid {
		if err := s.Repo.EnableTOTP(userID); err != nil {
			return true, fmt.Errorf("verified but failed to enable status: %v", err)
		}
	}

	return valid, nil
}