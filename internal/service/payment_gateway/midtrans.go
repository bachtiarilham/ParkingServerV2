package payment_gateway

// import (
// 	"context"
// 	"fmt"
// 	"time"

// 	// Import library Midtrans Go official
// 	"github.com/veritrans/go-midtrans"
// )

// type MidtransConfig struct {
// 	ServerKey   string
// 	ClientKey   string
// 	Environment midtrans.EnvironmentType // Sandbox atau Production
// }

// type MidtransService struct {
// 	client *midtrans.Client
// }

// func NewMidtransService(cfg MidtransConfig) (*MidtransService, error) {
// 	client := new(midtrans.Client)
// 	client.New(cfg.ServerKey, cfg.Environment)

// 	return &MidtransService{
// 		client: client,
// 	}, nil
// }

// func (s *MidtransService) RequestPayment(ctx context.Context, details PaymentDetails) (QRString string, VAString string, BankName string, ExpiredAt string, err error) {
// 	req := &midtrans.ChargeReq{
// 		PaymentType: midtrans.SourceQris,
// 		TransactionDetails: midtrans.TransactionDetails{
// 			OrderID:  details.OrderID,
// 			GrossAmt: int64(details.Amount),
// 		},
// 		Items: &[]midtrans.ItemDetail{},
// 		// Expiry: &midtrans.ExpiryDetails{
// 		// 	StartTime: time.Now().Format("2006-01-02 15:04:05 -0700 MST"),
// 		// 	Unit:      "minutes",
// 		// 	Duration:  details.ExpiryDuration / 60, // Convert seconds to minutes
// 		// },
// 	}

// 	for _, item := range details.ItemDetails {
// 		*req.Items = append(*req.Items, midtrans.ItemDetail{
// 			ID:    item.ID,
// 			Name:  item.Name,
// 			Price: item.Price,
// 			Qty:   item.Quantity,
// 		})
// 	}

// 	if details.CustomerDetails != nil {
// 		req.CustomerDetail = &midtrans.CustDetail{
// 			FName: details.CustomerDetails.FName,
// 			LName: details.CustomerDetails.LName,
// 			Email: details.CustomerDetails.Email,
// 			Phone: details.CustomerDetails.Phone,
// 		}
// 	}

// 	resp, err := s.client.ChargeTransaction(req)
// 	if err != nil {
// 		return "", "", "", "", fmt.Errorf("midtrans charge error: %w", err)
// 	}

// 	// Ambil detail dari response
// 	qrString := ""
// 	vaString := ""
// 	bankName := ""
// 	expiredAt := ""

// 	if resp.PaymentType == "qris" {
// 		qrString = resp.Actions[0].URL // Atau QRString dari response jika tersedia
// 		// Midtrans biasanya tidak langsung memberikan QR string, tapi URL untuk generate QR
// 		// Kamu perlu ambil dari Actions atau Snap Token
// 		// Misalnya, jika Actions berisi URL untuk generate QR:
// 		// qrString = resp.Actions[0].URL
// 		// Untuk demo, kita asumsikan QRIS string langsung ada
// 		qrString = resp.QRString // Ganti dengan cara yang benar dari library Midtrans
// 	}

// 	// Midtrans QRIS tidak memiliki bank name
// 	// Jika VA, maka VAString dan BankName akan diisi
// 	if resp.PaymentType == "bank_transfer" {
// 		vaString = resp.VANumbers[0].VANumber
// 		bankName = resp.VANumbers[0].Bank
// 	}

// 	// Ambil expiry time jika tersedia
// 	if resp.Expiry != nil {
// 		expiredAt = resp.Expiry.DateTime
// 	} else {
// 		// Jika tidak ada expiry di response, buat sendiri berdasarkan duration
// 		expTime := time.Now().Add(time.Duration(details.ExpiryDuration) * time.Second)
// 		expiredAt = expTime.Format("2006-01-02T15:04:05Z07:00")
// 	}

// 	return qrString, vaString, bankName, expiredAt, nil
// }

// func (s *MidtransService) VerifyPayment(ctx context.Context, orderID string) (Status string, PaidAt *string, err error) {
// 	// Gunakan API Get Status dari Midtrans
// 	statusResp, err := s.client.CheckTransaction(orderID)
// 	if err != nil {
// 		return "", nil, fmt.Errorf("midtrans check transaction error: %w", err)
// 	}

// 	status := statusResp.TransactionStatus
// 	var paidAt *string
// 	if statusResp.SettlementTime != "" {
// 		paidAt = &statusResp.SettlementTime // Format ISO8601
// 	}

// 	return status, paidAt, nil
// }
