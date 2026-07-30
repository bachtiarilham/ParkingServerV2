package midtrans

import (
	"context"
	"crypto/sha512"
	"encoding/hex"
	"fmt"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type MidtransClient struct {
	serverKey string
	snapCli   snap.Client
}

func NewMidtransClient(serverKey string, isProduction bool) *MidtransClient {
	env := midtrans.Sandbox
	if isProduction {
		env = midtrans.Production
	}

	var snapCli snap.Client
	snapCli.New(serverKey, env)

	return &MidtransClient{
		serverKey: serverKey,
		snapCli:   snapCli,
	}
}

func (c *MidtransClient) CreateSnapToken(ctx context.Context, req *snap.Request) (*snap.Response, error) {
	snapResp, err := c.snapCli.CreateTransaction(req)
	if err != nil {
		return nil, fmt.Errorf("midtrans SDK snap request: %w", err)
	}

	return snapResp, nil
}

// VerifySignature checks the Midtrans notification signature using SHA-512.
// Formula: SHA512(order_id + status_code + gross_amount + server_key)
func (c *MidtransClient) VerifySignature(orderID, statusCode, grossAmount, signature string) bool {
	raw := orderID + statusCode + grossAmount + c.serverKey
	hash := sha512.Sum512([]byte(raw))
	expected := hex.EncodeToString(hash[:])
	return expected == signature
}
