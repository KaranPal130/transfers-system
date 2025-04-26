package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/KaranPal130/transfers-system/internal/models"
	repository "github.com/KaranPal130/transfers-system/internal/repositories"
	service "github.com/KaranPal130/transfers-system/internal/services"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	accountService     *service.AccountService
	transactionService *service.TransactionService
}

func NewHandler(accountService *service.AccountService, transactionService *service.TransactionService) *Handler {
	return &Handler{
		accountService:     accountService,
		transactionService: transactionService,
	}
}

// CreateAccount handles account creation requests
// @Summary Create account
// @Description Create a new account
// @Tags accounts
// @Accept json
// @Produce json
// @Param account body models.AccountCreateRequest true "Account create request"
// @Success 201 {object} models.Account
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /accounts [post]
func (h *Handler) CreateAccount(c *gin.Context) {
	var req models.AccountCreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err := h.accountService.CreateAccount(c.Request.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidInitialBalance):
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid initial balance"})
		case errors.Is(err, service.ErrAccountAlreadyExists):
			c.JSON(http.StatusConflict, gin.H{"error": "Account already exists"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	account, err := h.accountService.GetAccount(c.Request.Context(), req.AccountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch created account"})
		return
	}

	c.JSON(http.StatusCreated, account)
}

// GetAccount handles account retrieval requests
// @Summary Get account
// @Description Get account by ID
// @Tags accounts
// @Produce json
// @Param account_id path int true "Account ID"
// @Success 200 {object} models.Account
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /accounts/{account_id} [get]
func (h *Handler) GetAccount(c *gin.Context) {
	accountIDStr := c.Param("account_id")

	accountID, err := strconv.ParseInt(accountIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	account, err := h.accountService.GetAccount(c.Request.Context(), accountID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrAccountNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, account)
}

// CreateTransaction handles transaction creation requests
// @Summary Create transaction
// @Description Create a new transaction
// @Tags transactions
// @Accept json
// @Produce json
// @Param transaction body models.TransactionRequest true "Transaction request"
// @Success 201 {object} models.TransactionRequest
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /transactions [post]
func (h *Handler) CreateTransaction(c *gin.Context) {
	var req models.TransactionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err := h.transactionService.CreateTransaction(c.Request.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidAmount):
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid amount"})
		case errors.Is(err, service.ErrInsufficientBalance):
			c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient balance"})
		case errors.Is(err, service.ErrSameSourceAndDest):
			c.JSON(http.StatusBadRequest, gin.H{"error": "Source and destination accounts must be different"})
		case errors.Is(err, repository.ErrAccountNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	c.JSON(http.StatusCreated, req)
}
