package main

import (
	"database/sql"
	"satgate/internal/database"
	"satgate/internal/forms"
	"satgate/internal/invoices"
	"satgate/internal/lightning"
	"satgate/internal/messages"
	"satgate/internal/middleware"
	"satgate/internal/webhooks"

	// "embed"
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

	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("migration failed: %v", err)
	}
	if err := database.SeedDefaultForm(db); err != nil {
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

	blink, err := lightning.CreateInvoice(
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
