package lightning

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const blinkGraphQLEndpoint = "https://api.blink.sv/graphql"

type InvoiceResult struct {
	PaymentRequest string
	PaymentHash    string
}

func CreateInvoice(
	apiKey string,
	walletID string,
	amountSats int64,
	memo string,
) (*InvoiceResult, error) {
	return createInvoice(
		http.DefaultClient,
		blinkGraphQLEndpoint,
		apiKey,
		walletID,
		amountSats,
		memo,
	)
}

func createInvoice(
	client *http.Client,
	endpoint string,
	apiKey string,
	walletID string,
	amountSats int64,
	memo string,
) (*InvoiceResult, error) {
	payload := map[string]any{
		"query": `
            mutation LnInvoiceCreate($input: LnInvoiceCreateInput!) {
                lnInvoiceCreate(input: $input) {
                    invoice {
                        paymentRequest
                        paymentHash
                    }
                    errors {
                        message
                    }
                }
            }
        `,
		"variables": map[string]any{
			"input": map[string]any{
				"walletId": walletID,
				"amount":   amountSats,
				"memo":     memo,
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		endpoint,
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-KEY", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("blink returned HTTP %s", resp.Status)
	}

	var result struct {
		Data struct {
			LnInvoiceCreate struct {
				Invoice struct {
					PaymentRequest string `json:"paymentRequest"`
					PaymentHash    string `json:"paymentHash"`
				} `json:"invoice"`
				Errors []struct {
					Message string `json:"message"`
				} `json:"errors"`
			} `json:"lnInvoiceCreate"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	invoice := result.Data.LnInvoiceCreate.Invoice

	if len(result.Data.LnInvoiceCreate.Errors) > 0 {
		return nil, fmt.Errorf(
			"blink error: %s",
			result.Data.LnInvoiceCreate.Errors[0].Message,
		)
	}

	if invoice.PaymentRequest == "" {
		return nil, fmt.Errorf("blink returned empty invoice")
	}

	return &InvoiceResult{
		PaymentRequest: invoice.PaymentRequest,
		PaymentHash:    invoice.PaymentHash,
	}, nil
}
