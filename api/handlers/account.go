package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"simplebank/service/account"
	"simplebank/tasks"

	"github.com/hibiken/asynq"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type AccountHandler struct {
	AccountRepository *account.Repository
	TaskClient        *asynq.Client
}

func NewAccountHandler(db *sql.DB, taskClient *asynq.Client) *AccountHandler {
	return &AccountHandler{
		AccountRepository: account.NewRepository(db),
		TaskClient:        taskClient,
	}
}

type CreateAccountRequestBody struct {
	Owner    string `json:"owner" validate:"required"`
	Currency string `json:"currency" validate:"required,oneof=USD EUR BRL"`
}

func (h *AccountHandler) CreateAccount(c echo.Context) error {
	var reqBody CreateAccountRequestBody
	if err := c.Bind(&reqBody); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if err := c.Validate(reqBody); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	id, err := h.AccountRepository.CreateAccount(account.CreateAccountParams{
		Owner:    reqBody.Owner,
		Currency: reqBody.Currency,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	email := strings.ToLower(strings.ReplaceAll(reqBody.Owner, " ", "-")) + "@gmail.com"
	task, err := tasks.NewEmailDeliveryTask(email)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("could not create email task: %w", err))
	}

	info, err := h.TaskClient.Enqueue(task,
		asynq.Queue(tasks.QueueCritical),
		asynq.ProcessIn(5*time.Second),
		asynq.MaxRetry(3),
		asynq.Timeout(10*time.Second),
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("could not schedule email task: %w", err))
	}

	log.Info().Msgf("enqueued task: id=%s queue=%s", info.ID, info.Queue)

	return c.JSON(http.StatusOK, map[string]any{
		"id": id,
	})
}

type AccountResponseBody struct {
	ID        int       `json:"id"`
	Owner     string    `json:"owner"`
	Currency  string    `json:"currency"`
	Balance   uint      `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

func (h *AccountHandler) ListAccounts(c echo.Context) error {
	sizeStr := c.QueryParam("size")
	if sizeStr == "" {
		sizeStr = "10"
	}
	size, err := strconv.Atoi(sizeStr)
	if err != nil || size < 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid limit")
	}
	size = min(size, 20)

	pageStr := c.QueryParam("page")
	if pageStr == "" {
		pageStr = "1"
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid page")
	}

	accounts, err := h.AccountRepository.ListAccounts(account.ListAccountsParams{
		Limit:  size,
		Offset: (page - 1) * size,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	res := make([]AccountResponseBody, len(accounts))
	for i, a := range accounts {
		res[i] = AccountResponseBody{
			ID:        a.ID,
			Owner:     a.Owner,
			Currency:  a.Currency,
			Balance:   a.Balance,
			CreatedAt: a.CreatedAt,
		}
	}

	return c.JSON(http.StatusOK, res)

}
