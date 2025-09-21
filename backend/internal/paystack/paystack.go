package paystack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/flexGURU/flower-haven/backend/internal/services"
	"github.com/flexGURU/flower-haven/backend/pkg"
)

var _ services.IPayStack = (*Paystack)(nil)

type Paystack struct {
	SecretKey   string
	CallbackURL string
	BaseURL     string
}

func NewPaystack(secretKey string, callbackUrl string) services.IPayStack {
	return &Paystack{
		SecretKey:   secretKey,
		BaseURL:     "https://api.paystack.co",
		CallbackURL: callbackUrl,
	}
}

func (ps Paystack) InitializePayment(email string, amount int64) (string, string, error) {
	payload := map[string]string{
		"email":        email,
		"amount":       fmt.Sprintf("%d", amount),
		"callback_url": ps.CallbackURL,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", "", pkg.Errorf(pkg.INTERNAL_ERROR, "failed to marshal payload: %s", err.Error())
	}

	req, err := http.NewRequest(http.MethodPost, ps.BaseURL+"/transaction/initialize", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", "", pkg.Errorf(pkg.INTERNAL_ERROR, "failed to create request: %s", err.Error())
	}

	req.Header.Set("Authorization", "Bearer "+ps.SecretKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", pkg.Errorf(pkg.INTERNAL_ERROR, "failed to send request: %s", err.Error())
	}
	defer resp.Body.Close()

	var result struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
		Data    struct {
			AuthorizationURL string `json:"authorization_url"`
			Reference        string `json:"reference"`
			AccessCode       string `json:"access_code"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", "", pkg.Errorf(pkg.INTERNAL_ERROR, "failed to decode response: %s", err.Error())
	}

	if !result.Status {
		return "", "", pkg.Errorf(pkg.INTERNAL_ERROR, "failed to initialize payment: %s", result.Message)
	}

	return result.Data.AccessCode, result.Data.Reference, nil
}

func (ps Paystack) VerifyPayment(reference string, amount int64) (string, error) {
	req, err := http.NewRequest(http.MethodGet, ps.BaseURL+"/transaction/verify/"+reference, nil)
	if err != nil {
		return "", pkg.Errorf(pkg.INTERNAL_ERROR, "failed to create request: %s", err.Error())
	}

	req.Header.Set("Authorization", "Bearer "+ps.SecretKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", pkg.Errorf(pkg.INTERNAL_ERROR, "failed to send request: %s", err.Error())
	}
	defer resp.Body.Close()

	var result struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
		Data    struct {
			Status string `json:"status"`
			Amount int64  `json:"amount"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", pkg.Errorf(pkg.INTERNAL_ERROR, "failed to decode response: %s", err.Error())
	}

	if !result.Status {
		return "", pkg.Errorf(pkg.NOT_FOUND_ERROR, "failed to verify payment: %s", result.Message)
	}

	if result.Data.Amount != amount {
		return "", pkg.Errorf(pkg.INVALID_ERROR, "amount mismatch: expected %d, got %d", amount, result.Data.Amount)
	}

	return result.Data.Status, nil
}
