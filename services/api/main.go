package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"

	// "embed"
	"encoding/hex"
	"encoding/json"
	"fmt"

	// "io/fs"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)



type AppState struct {
    DB            *sql.DB
    WebhookSecret string
    BlinkApiKey   string
    BlinkWalletID string
}

type Form struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Domain     string    `json:"domain"`
	AmountSats int64     `json:"amount_sats"`
	CreatedAt  time.Time `json:"created_at"`
}

type Message struct {
	ID          string     `json:"id"`
	FormID      string     `json:"form_id"`
	SenderName  string     `json:"sender_name"`
	SenderEmail string     `json:"sender_email"`
	Body        string     `json:"body"`
	Status      string     `json:"status"`
	AmountSats  int64      `json:"amount_sats"`
	CreatedAt   time.Time  `json:"created_at"`
	PaidAt      *time.Time `json:"paid_at"`
}

type CreateFormRequest struct {
	Name       string `json:"name"`
	Domain     string `json:"domain"`
	AmountSats int64  `json:"amount_sats"`
}

type UpdateFormRequest struct {
	Name       string `json:"name"`
	Domain     string `json:"domain"`
	AmountSats int64  `json:"amount_sats"`
}

type CreateInvoiceRequest struct {
	FormID      string `json:"form_id"`
	SenderName  string `json:"sender_name"`
	SenderEmail string `json:"sender_email"`
	Body        string `json:"body"`
}

type InvoiceResponse struct {
	InvoiceID      string    `json:"invoice_id"`
	MessageID      string    `json:"message_id"`
	PaymentRequest string    `json:"payment_request"`
	AmountSats     int64     `json:"amount_sats"`
	Status         string    `json:"status"`
	ExpiresAt      time.Time `json:"expires_at"`
}

type InvoiceStatusResponse struct {
	InvoiceID string     `json:"invoice_id"`
	MessageID string     `json:"message_id"`
	Status    string     `json:"status"`
	PaidAt    *time.Time `json:"paid_at"`
}

type LightningWebhook struct {
	InvoiceID   string  `json:"invoice_id"`
	PaymentHash *string `json:"payment_hash"`
	Preimage    *string `json:"preimage"`
	Status      string  `json:"status"`
}

type blinkInvoiceResult struct {
    PaymentRequest string
    PaymentHash    string
}

func main() {
	dbURL := getenv("DATABASE_URL", "satgate.db")
	webhookSecret := getenv("LIGHTNING_WEBHOOK_SECRET", "satgate-dev-secret")
	port := getenv("PORT", "8080")

	db, err := sql.Open("sqlite3", dbURL)
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	defer db.Close()
	db.SetMaxOpenConns(5)

	if err := runMigrations(db); err != nil {
		log.Fatalf("migration failed: %v", err)
	}
	if err := seedDefaultForm(db); err != nil {
		log.Fatalf("seed failed: %v", err)
	}
state := &AppState{
    DB:            db,
    WebhookSecret: webhookSecret,
    BlinkApiKey:   getenv("BLINK_API_KEY", ""),
    BlinkWalletID: getenv("BLINK_WALLET_ID", ""),
}

log.Printf(
    "Blink configured: api_key=%t wallet_id=%t",
    state.BlinkApiKey != "",
    state.BlinkWalletID != "",
)
	r := gin.Default()
	r.Use(cors.Default())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.GET("/api/forms", state.listForms)
	r.POST("/api/forms", state.createForm)
	r.PUT("/api/forms/:id", state.updateForm)
	r.DELETE("/api/forms/:id", state.deleteForm)

	r.GET("/api/messages", state.listMessages)

	r.POST("/api/invoices", state.createInvoice)
	r.GET("/api/invoices/:id", state.getInvoice)
	r.POST("/api/invoices/:id/mock-pay", state.mockPayInvoice)

	r.POST("/api/webhooks/lightning", state.lightningWebhook)

	// sub, err := fs.Sub(staticFiles, "dist")
	// if err != nil {
	// 	log.Fatalf("failed to open embedded dist: %v", err)
	// }
	// fileServer := http.FileServer(http.FS(sub))
	// r.NoRoute(func(c *gin.Context) {
	// 	path := c.Request.URL.Path
	// 	if _, err := fs.Stat(sub, strings.TrimPrefix(path, "/")); err != nil {
	// 		c.Request.URL.Path = "/"
	// 	}
	// 	fileServer.ServeHTTP(c.Writer, c.Request)
	// })

	log.Printf("SatGate API listening on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func (s *AppState) listForms(c *gin.Context) {
	rows, err := s.DB.Query(
		"SELECT id, name, domain, amount_sats, created_at FROM forms ORDER BY created_at DESC",
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	forms := []Form{}
	for rows.Next() {
		var f Form
		if err := rows.Scan(&f.ID, &f.Name, &f.Domain, &f.AmountSats, &f.CreatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		forms = append(forms, f)
	}
	c.JSON(http.StatusOK, forms)
}

func (s *AppState) createForm(c *gin.Context) {
	var req CreateFormRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Domain = strings.TrimSpace(req.Domain)

	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "form name is required"})
		return
	}
	if req.Domain == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "allowed domain is required"})
		return
	}
	if req.AmountSats < 1 || req.AmountSats > 10_000 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "amount_sats must be between 1 and 10000"})
		return
	}

	form := Form{
		ID:         uuid.New().String(),
		Name:       req.Name,
		Domain:     req.Domain,
		AmountSats: req.AmountSats,
		CreatedAt:  time.Now().UTC(),
	}

	_, err := s.DB.Exec(
		"INSERT INTO forms (id, name, domain, amount_sats, created_at) VALUES (?, ?, ?, ?, ?)",
		form.ID, form.Name, form.Domain, form.AmountSats, form.CreatedAt,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, form)
}

func (s *AppState) updateForm(c *gin.Context) {
	id := c.Param("id")

	var req UpdateFormRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Domain = strings.TrimSpace(req.Domain)

	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "form name is required"})
		return
	}
	if req.Domain == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "allowed domain is required"})
		return
	}
	if req.AmountSats < 1 || req.AmountSats > 10_000 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "amount_sats must be between 1 and 10000"})
		return
	}

	res, err := s.DB.Exec(
		"UPDATE forms SET name = ?, domain = ?, amount_sats = ? WHERE id = ?",
		req.Name, req.Domain, req.AmountSats, id,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	form, err := s.findForm(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, form)
}

func (s *AppState) deleteForm(c *gin.Context) {
	id := c.Param("id")

	tx, err := s.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer tx.Rollback()

	queries := []string{
		"DELETE FROM webhook_events WHERE invoice_id IN (SELECT id FROM invoices WHERE message_id IN (SELECT id FROM messages WHERE form_id = ?))",
		"DELETE FROM invoices WHERE message_id IN (SELECT id FROM messages WHERE form_id = ?)",
		"DELETE FROM messages WHERE form_id = ?",
	}
	for _, q := range queries {
		if _, err := tx.Exec(q, id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	res, err := tx.Exec("DELETE FROM forms WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	tx.Commit()
	c.Status(http.StatusNoContent)
}

func (s *AppState) listMessages(c *gin.Context) {
	rows, err := s.DB.Query(`
		SELECT m.id, m.form_id, m.sender_name, m.sender_email, m.body, m.status,
		       COALESCE(i.amount_sats, 0), m.created_at, m.paid_at
		FROM messages m
		LEFT JOIN invoices i ON i.message_id = m.id
		WHERE m.status = 'paid'
		ORDER BY m.paid_at DESC, m.created_at DESC
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	messages := []Message{}
	for rows.Next() {
		var m Message
		if err := rows.Scan(
			&m.ID, &m.FormID, &m.SenderName, &m.SenderEmail,
			&m.Body, &m.Status, &m.AmountSats, &m.CreatedAt, &m.PaidAt,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		messages = append(messages, m)
	}
	c.JSON(http.StatusOK, messages)
}
func (s *AppState) createInvoice(c *gin.Context) {
   log.Printf(
        "CREATE INVOICE: BlinkApiKey=%t BlinkWalletID=%t",
        s.BlinkApiKey != "",
        s.BlinkWalletID != "",
    )
	
	var req CreateInvoiceRequest

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

	    log.Printf("CREATE INVOICE: form_id=%q", req.FormID)


    if err := validateMessagePayload(req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    form, err := s.findForm(req.FormID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "form not found"})
        return
    }

    messageID := uuid.New().String()
    now := time.Now().UTC()

    memo := fmt.Sprintf(
        "SatGate message from %s",
        strings.TrimSpace(req.SenderName),
    )

println(s.BlinkWalletID)

    blink, err := createBlinkInvoice(
        s.BlinkApiKey,
        s.BlinkWalletID,
        form.AmountSats,
        memo,
    )
    if err != nil {
        log.Printf("Blink invoice creation failed: %v", err)

        c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
        })
        return
    }

    invoiceID := blink.PaymentHash
    paymentRequest := blink.PaymentRequest
    paymentHash := blink.PaymentHash
    expiresAt := now.Add(24 * time.Hour)

    tx, err := s.DB.Begin()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer tx.Rollback()

    _, err = tx.Exec(`
        INSERT INTO messages (
            id,
            form_id,
            sender_name,
            sender_email,
            body,
            status,
            payment_hash,
            created_at,
            paid_at
        )
        VALUES (?, ?, ?, ?, ?, 'pending', ?, ?, NULL)
    `,
        messageID,
        req.FormID,
        strings.TrimSpace(req.SenderName),
        strings.TrimSpace(req.SenderEmail),
        strings.TrimSpace(req.Body),
        paymentHash,
        now,
    )

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    _, err = tx.Exec(`
        INSERT INTO invoices (
            id,
            message_id,
            amount_sats,
            payment_request,
            payment_hash,
            status,
            preimage,
            created_at,
            paid_at,
            expires_at
        )
        VALUES (?, ?, ?, ?, ?, 'pending', NULL, ?, NULL, ?)
    `,
        invoiceID,
        messageID,
        form.AmountSats,
        paymentRequest,
        paymentHash,
        now,
        expiresAt,
    )

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    if err := tx.Commit(); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    log.Printf(
        "Blink invoice created: invoice_id=%s amount=%d payment_hash=%s",
        invoiceID,
        form.AmountSats,
        paymentHash,
    )

    c.JSON(http.StatusOK, InvoiceResponse{
        InvoiceID:      invoiceID,
        MessageID:      messageID,
        PaymentRequest: paymentRequest,
        AmountSats:     form.AmountSats,
        Status:         "pending",
        ExpiresAt:      expiresAt,
    })
}
func createBlinkInvoice(
    apiKey string,
    walletID string,
    amountSats int64,
    memo string,
) (*blinkInvoiceResult, error) {

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
        "https://api.blink.sv/graphql",
        bytes.NewReader(body),
    )
    if err != nil {
        return nil, err
    }

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("X-API-KEY", apiKey)

    resp, err := http.DefaultClient.Do(req)
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

    return &blinkInvoiceResult{
        PaymentRequest: invoice.PaymentRequest,
        PaymentHash:    invoice.PaymentHash,
    }, nil
}

func (s *AppState) getInvoice(c *gin.Context) {
	id := c.Param("id")
	resp, err := s.fetchInvoiceStatus(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (s *AppState) mockPayInvoice(c *gin.Context) {
	id := c.Param("id")
	preimage := "mock-preimage"
	if err := s.markInvoicePaid(id, &preimage); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	resp, err := s.fetchInvoiceStatus(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (s *AppState) lightningWebhook(c *gin.Context) {
	body, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not read body"})
		return
	}

	if err := verifySignature(c.Request.Header, string(body), s.WebhookSecret); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "signature verification failed"})
		return
	}

	var payload LightningWebhook
	if err := json.Unmarshal(body, &payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid webhook body"})
		return
	}

	if payload.Status != "paid" && payload.Status != "settled" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "webhook status is not paid"})
		return
	}

	eventID := c.GetHeader("x-satgate-event-id")
	if eventID == "" {
		eventID = uuid.New().String()
	}

	res, err := s.DB.Exec(
		"INSERT OR IGNORE INTO webhook_events (id, invoice_id, payload, created_at) VALUES (?, ?, ?, ?)",
		eventID, payload.InvoiceID, string(body), time.Now().UTC(),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	n, _ := res.RowsAffected()
	if n == 0 {
		resp, err := s.fetchInvoiceStatus(payload.InvoiceID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusOK, resp)
		return
	}

	if payload.PaymentHash != nil {
		var existingHash string
		err := s.DB.QueryRow("SELECT payment_hash FROM invoices WHERE id = ?", payload.InvoiceID).
			Scan(&existingHash)
		if err == nil && existingHash != *payload.PaymentHash {
			c.JSON(http.StatusBadRequest, gin.H{"error": "payment hash does not match invoice"})
			return
		}
	}

	if err := s.markInvoicePaid(payload.InvoiceID, payload.Preimage); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	resp, err := s.fetchInvoiceStatus(payload.InvoiceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (s *AppState) findForm(id string) (*Form, error) {
	    log.Printf("FIND FORM: looking for id=%q", id)

	var f Form
	err := s.DB.QueryRow(
		"SELECT id, name, domain, amount_sats, created_at FROM forms WHERE id = ?", id,
	).Scan(&f.ID, &f.Name, &f.Domain, &f.AmountSats, &f.CreatedAt)
	if err != nil {
		        log.Printf("FIND FORM: query failed: %v", err)

		return nil, err
	}

	    log.Printf("FIND FORM: found %+v", f)

		 log.Printf("FIND FORM: found id=%q name=%q amount=%d",
        f.ID,
        f.Name,
        f.AmountSats,
    )

	return &f, nil
}

func (s *AppState) fetchInvoiceStatus(id string) (*InvoiceStatusResponse, error) {
	var resp InvoiceStatusResponse
	err := s.DB.QueryRow(
		"SELECT id, message_id, status, paid_at FROM invoices WHERE id = ?", id,
	).Scan(&resp.InvoiceID, &resp.MessageID, &resp.Status, &resp.PaidAt)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (s *AppState) markInvoicePaid(invoiceID string, preimage *string) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	paidAt := time.Now().UTC()

	res, err := tx.Exec(`
		UPDATE invoices SET status = 'paid', paid_at = ?, preimage = COALESCE(?, preimage)
		WHERE id = ? AND status != 'paid'`,
		paidAt, preimage, invoiceID,
	)
	if err != nil {
		return err
	}

	var messageID string
	err = tx.QueryRow("SELECT message_id FROM invoices WHERE id = ?", invoiceID).Scan(&messageID)
	if err != nil {
		return fmt.Errorf("invoice not found")
	}

	n, _ := res.RowsAffected()
	if n > 0 {
		_, err = tx.Exec(
			"UPDATE messages SET status = 'paid', paid_at = ? WHERE id = ?",
			paidAt, messageID,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func validateMessagePayload(req CreateInvoiceRequest) error {
	if strings.TrimSpace(req.FormID) == "" {
		return fmt.Errorf("form_id is required")
	}
	if strings.TrimSpace(req.SenderName) == "" {
		return fmt.Errorf("sender_name is required")
	}
	if !strings.Contains(req.SenderEmail, "@") {
		return fmt.Errorf("a valid sender_email is required")
	}
	if len(strings.TrimSpace(req.Body)) < 10 {
		return fmt.Errorf("message body must be at least 10 characters")
	}
	return nil
}

func buildDemoInvoice(invoiceID string, amountSats int64) string {
	return fmt.Sprintf("satgate-demo-invoice:%s:%d", invoiceID, amountSats)
}

func verifySignature(headers http.Header, body, secret string) error {
	sig := headers.Get("x-satgate-signature")
	if sig == "" {
		return fmt.Errorf("missing signature")
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(body))
	expected := hex.EncodeToString(mac.Sum(nil))

	if expected != sig {
		return fmt.Errorf("signature mismatch")
	}
	return nil
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func runMigrations(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS forms (
			id          TEXT PRIMARY KEY,
			name        TEXT NOT NULL,
			domain      TEXT NOT NULL,
			amount_sats INTEGER NOT NULL,
			created_at  DATETIME NOT NULL
		);

		CREATE TABLE IF NOT EXISTS messages (
			id           TEXT PRIMARY KEY,
			form_id      TEXT NOT NULL,
			sender_name  TEXT NOT NULL,
			sender_email TEXT NOT NULL,
			body         TEXT NOT NULL,
			status       TEXT NOT NULL DEFAULT 'pending',
			payment_hash TEXT,
			created_at   DATETIME NOT NULL,
			paid_at      DATETIME
		);

		CREATE TABLE IF NOT EXISTS invoices (
			id              TEXT PRIMARY KEY,
			message_id      TEXT NOT NULL,
			amount_sats     INTEGER NOT NULL,
			payment_request TEXT NOT NULL,
			payment_hash    TEXT,
			status          TEXT NOT NULL DEFAULT 'pending',
			preimage        TEXT,
			created_at      DATETIME NOT NULL,
			paid_at         DATETIME,
			expires_at      DATETIME NOT NULL
		);

		CREATE TABLE IF NOT EXISTS webhook_events (
			id         TEXT PRIMARY KEY,
			invoice_id TEXT NOT NULL,
			payload    TEXT NOT NULL,
			created_at DATETIME NOT NULL
		);
	`)
	return err
}

func seedDefaultForm(db *sql.DB) error {
	var count int
	db.QueryRow("SELECT COUNT(*) FROM forms").Scan(&count)
	if count > 0 {
		return nil
	}
	_, err := db.Exec(
		"INSERT INTO forms (id, name, domain, amount_sats, created_at) VALUES (?, ?, ?, ?, ?)",
		"demo-form", "Sarah's Portfolio", "localhost", 5, time.Now().UTC(),
	)
	return err
}
