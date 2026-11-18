package main

import (
	"bytes"
	"database/sql"
	"satgate/internal/forms"
	"satgate/internal/invoices"
	"satgate/internal/messages"
	"satgate/internal/middleware"
	"satgate/internal/webhooks"

	// "embed"
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
	Forms         *forms.Service
	BlinkApiKey   string
	BlinkWalletID string
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
	formRepository := forms.NewRepository(db)
	formService := forms.NewService(formRepository)

	invoiceRepository := invoices.NewRepository(db)
	invoiceService := invoices.NewService(invoiceRepository)
	invoiceHandler := invoices.NewHandler(invoiceService)

	webhookRepository := webhooks.NewRepository(db)
	webhookService := webhooks.NewService(webhookRepository, invoiceService, webhookSecret)
	webhookHandler := webhooks.NewHandler(webhookService)

	state := &AppState{
		DB:            db,
		Forms:         formService,
		BlinkApiKey:   getenv("BLINK_API_KEY", ""),
		BlinkWalletID: getenv("BLINK_WALLET_ID", ""),
	}

	formHandler := forms.NewHandler(formService)

	messageRepository := messages.NewRepository(db)
	messageService := messages.NewService(messageRepository)
	messageHandler := messages.NewHandler(messageService)

	log.Printf(
		"Blink configured: api_key=%t wallet_id=%t",
		state.BlinkApiKey != "",
		state.BlinkWalletID != "",
	)

	appEnv := strings.ToLower(getenv("APP_ENV", "production"))
	jwtSecret := getenv("JWT_SECRET", "")
	allowMockTokens := appEnv == "development"

	if appEnv != "development" && jwtSecret == "" {
		log.Fatal("JWT_SECRET must be configured outside development")
	}

	r := gin.Default()
	r.Use(cors.Default())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	protected := r.Group("/api")
	protected.Use(middleware.Auth(jwtSecret, allowMockTokens))

	protected.GET("/forms", formHandler.List)
	protected.POST("/forms", formHandler.Create)
	protected.PUT("/forms/:id", formHandler.Update)
	protected.DELETE("/forms/:id", formHandler.Delete)

	protected.GET("/messages", messageHandler.List)

	r.POST("/api/invoices", state.createInvoice)
	r.GET("/api/invoices/:id", invoiceHandler.Get)

	if allowMockTokens {
		r.POST("/api/invoices/:id/mock-pay", invoiceHandler.MockPay)
	}

	r.POST("/api/webhooks/lightning", webhookHandler.Handle)

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

func (s *AppState) createInvoice(c *gin.Context) {
	log.Printf(
		"CREATE INVOICE: BlinkApiKey=%t BlinkWalletID=%t",
		s.BlinkApiKey != "",
		s.BlinkWalletID != "",
	)

	var req invoices.CreateInvoiceRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("CREATE INVOICE: form_id=%q", req.FormID)

	if err := validateMessagePayload(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	form, err := s.Forms.GetPublic(c.Request.Context(), req.FormID)
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

	c.JSON(http.StatusOK, invoices.InvoiceResponse{
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

func validateMessagePayload(req invoices.CreateInvoiceRequest) error {
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
			owner_id    TEXT,
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
	if err != nil {
		return err
	}

	// Upgrade existing databases without destroying existing data.
	rows, err := db.Query("PRAGMA table_info(forms)")
	if err != nil {
		return err
	}

	hasOwnerID := false

	for rows.Next() {
		var cid, notNull, primaryKey int
		var name, columnType string
		var defaultValue sql.NullString

		if err := rows.Scan(
			&cid,
			&name,
			&columnType,
			&notNull,
			&defaultValue,
			&primaryKey,
		); err != nil {
			rows.Close()
			return err
		}

		if name == "owner_id" {
			hasOwnerID = true
		}
	}

	if err := rows.Err(); err != nil {
		rows.Close()
		return err
	}
	rows.Close()

	if !hasOwnerID {
		if _, err := db.Exec("ALTER TABLE forms ADD COLUMN owner_id TEXT"); err != nil {
			return err
		}
	}

	// Preserve ownership of the existing development demo form.
	if _, err := db.Exec(
		"UPDATE forms SET owner_id = ? WHERE id = ? AND owner_id IS NULL",
		"demo-user",
		"demo-form",
	); err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_forms_owner_id
		ON forms(owner_id)
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
		"INSERT INTO forms (id, owner_id, name, domain, amount_sats, created_at) VALUES (?, ?, ?, ?, ?, ?)",
		"demo-form", "demo-user", "Sarah's Portfolio", "localhost", 5, time.Now().UTC(),
	)
	return err
}
