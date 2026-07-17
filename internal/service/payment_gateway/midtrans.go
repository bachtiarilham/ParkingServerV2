package payment_gateway

import (
	"bytes"
	"context"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type MidtransConfig struct {
	ServerKey               string
	ClientKey               string
	Environment             string
	BaseURL                 string
	Acquirer                string
	NotificationURL         string
	OverrideNotificationURL string
	HTTPClient              *http.Client
}

type MidtransService struct {
	cfg    MidtransConfig
	client *http.Client
}

func NewMidtransService(cfg MidtransConfig) (*MidtransService, error) {
	if cfg.ServerKey == "" {
		return nil, fmt.Errorf("midtrans server key is required")
	}
	if cfg.ClientKey == "" {
		return nil, fmt.Errorf("midtrans client key is required")
	}

	if cfg.Environment == "" {
		cfg.Environment = "sandbox"
	}
	if cfg.Acquirer == "" {
		cfg.Acquirer = "gopay"
	}
	if cfg.BaseURL == "" {
		cfg.BaseURL = defaultMidtransBaseURL(cfg.Environment)
	}

	client := cfg.HTTPClient
	if client == nil {
		client = &http.Client{Timeout: 30 * time.Second}
	}

	return &MidtransService{
		cfg:    cfg,
		client: client,
	}, nil
}

func (s *MidtransService) RequestPayment(ctx context.Context, details PaymentDetails) (string, string, string, string, error) {
	body := midtransChargeRequest{
		PaymentType: "qris",
		TransactionDetails: midtransTransactionDetails{
			OrderID:     details.OrderID,
			GrossAmount: details.Amount,
		},
		QRIS: midtransQrisRequest{
			Acquirer: s.cfg.Acquirer,
		},
	}

	if len(details.ItemDetails) > 0 {
		body.ItemDetails = make([]midtransItemDetail, 0, len(details.ItemDetails))
		for _, item := range details.ItemDetails {
			body.ItemDetails = append(body.ItemDetails, midtransItemDetail{
				ID:       item.ID,
				Name:     item.Name,
				Price:    item.Price,
				Quantity: item.Quantity,
			})
		}
	}

	if details.CustomerDetails != nil {
		body.CustomerDetails = &midtransCustomerDetails{
			FirstName: details.CustomerDetails.FName,
			LastName:  details.CustomerDetails.LName,
			Email:     details.CustomerDetails.Email,
			Phone:     details.CustomerDetails.Phone,
		}
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return "", "", "", "", fmt.Errorf("marshal midtrans charge request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.cfg.BaseURL+"/v2/charge", bytes.NewReader(payload))
	if err != nil {
		return "", "", "", "", fmt.Errorf("create midtrans charge request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", basicAuthHeader(s.cfg.ServerKey))
	if s.cfg.OverrideNotificationURL != "" {
		req.Header.Set("X-Override-Notification", s.cfg.OverrideNotificationURL)
	}
	if s.cfg.NotificationURL != "" {
		req.Header.Set("X-Append-Notification", s.cfg.NotificationURL)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return "", "", "", "", fmt.Errorf("midtrans charge request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", "", "", fmt.Errorf("read midtrans charge response: %w", err)
	}

	var chargeResp midtransChargeResponse
	if err := json.Unmarshal(respBody, &chargeResp); err != nil {
		return "", "", "", "", fmt.Errorf("unmarshal midtrans charge response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		if chargeResp.StatusMessage == "" {
			chargeResp.StatusMessage = resp.Status
		}
		return "", "", "", "", fmt.Errorf("midtrans charge failed: %s", chargeResp.StatusMessage)
	}

	qrString := chargeResp.QRString
	if qrString == "" {
		for _, action := range chargeResp.Actions {
			if strings.Contains(strings.ToLower(action.Name), "qr-code") && action.URL != "" {
				qrString = action.URL
				break
			}
		}
	}

	expiredAt := ""
	if chargeResp.ExpiryTime != "" {
		expiredAt = chargeResp.ExpiryTime
	}

	return qrString, "", "", expiredAt, nil
}

func (s *MidtransService) VerifyPayment(ctx context.Context, orderID string) (string, *string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.cfg.BaseURL+"/v2/"+orderID+"/status", nil)
	if err != nil {
		return "", nil, fmt.Errorf("create midtrans status request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", basicAuthHeader(s.cfg.ServerKey))

	resp, err := s.client.Do(req)
	if err != nil {
		return "", nil, fmt.Errorf("midtrans status request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil, fmt.Errorf("read midtrans status response: %w", err)
	}

	var statusResp midtransStatusResponse
	if err := json.Unmarshal(respBody, &statusResp); err != nil {
		return "", nil, fmt.Errorf("unmarshal midtrans status response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		if statusResp.StatusMessage == "" {
			statusResp.StatusMessage = resp.Status
		}
		return "", nil, fmt.Errorf("midtrans status failed: %s", statusResp.StatusMessage)
	}

	paidAt := statusResp.SettlementTime
	if paidAt == "" && (statusResp.TransactionStatus == "capture" || statusResp.TransactionStatus == "settlement") {
		paidAt = statusResp.TransactionTime
	}

	if statusResp.TransactionStatus == "" {
		return "", nil, nil
	}

	if paidAt == "" {
		return statusResp.TransactionStatus, nil, nil
	}

	return statusResp.TransactionStatus, &paidAt, nil
}

func (s *MidtransService) VerifySignature(orderID, statusCode, grossAmount, signature string) bool {
	sum := sha512.Sum512([]byte(orderID + statusCode + grossAmount + s.cfg.ServerKey))
	return fmt.Sprintf("%x", sum) == strings.ToLower(signature)
}

func defaultMidtransBaseURL(environment string) string {
	if strings.EqualFold(environment, "production") {
		return "https://api.midtrans.com"
	}
	return "https://api.sandbox.midtrans.com"
}

func basicAuthHeader(serverKey string) string {
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(serverKey+":"))
}

type midtransChargeRequest struct {
	PaymentType        string                     `json:"payment_type"`
	TransactionDetails midtransTransactionDetails `json:"transaction_details"`
	ItemDetails        []midtransItemDetail       `json:"item_details,omitempty"`
	CustomerDetails    *midtransCustomerDetails   `json:"customer_details,omitempty"`
	QRIS               midtransQrisRequest        `json:"qris"`
}

type midtransTransactionDetails struct {
	OrderID     string `json:"order_id"`
	GrossAmount int64  `json:"gross_amount"`
}

type midtransItemDetail struct {
	ID       string `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Price    int64  `json:"price,omitempty"`
	Quantity int32  `json:"quantity,omitempty"`
}

type midtransCustomerDetails struct {
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Email     string `json:"email,omitempty"`
	Phone     string `json:"phone,omitempty"`
}

type midtransQrisRequest struct {
	Acquirer string `json:"acquirer"`
}

type midtransChargeResponse struct {
	StatusCode        string           `json:"status_code"`
	StatusMessage     string           `json:"status_message"`
	TransactionID     string           `json:"transaction_id"`
	OrderID           string           `json:"order_id"`
	GrossAmount       string           `json:"gross_amount"`
	TransactionStatus string           `json:"transaction_status"`
	PaymentType       string           `json:"payment_type"`
	TransactionTime   string           `json:"transaction_time"`
	FraudStatus       string           `json:"fraud_status"`
	Acquirer          string           `json:"acquirer"`
	QRString          string           `json:"qr_string"`
	ExpiryTime        string           `json:"expiry_time"`
	Actions           []midtransAction `json:"actions"`
}

type midtransAction struct {
	Name   string `json:"name"`
	Method string `json:"method"`
	URL    string `json:"url"`
}

type midtransStatusResponse struct {
	StatusCode        string `json:"status_code"`
	StatusMessage     string `json:"status_message"`
	TransactionID     string `json:"transaction_id"`
	OrderID           string `json:"order_id"`
	GrossAmount       string `json:"gross_amount"`
	TransactionStatus string `json:"transaction_status"`
	FraudStatus       string `json:"fraud_status"`
	PaymentType       string `json:"payment_type"`
	TransactionTime   string `json:"transaction_time"`
	SettlementTime    string `json:"settlement_time"`
}
