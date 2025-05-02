package api

import (
	"database/sql"
	"net/http"

	"simplebank/api/handlers"

	"github.com/hibiken/asynq"
	"github.com/labstack/echo/v4"
)

func routes(e *echo.Echo, db *sql.DB, tasksClient *asynq.Client) *echo.Echo {
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "WORKING")
	})

	e.POST("/v1/accounts", func(c echo.Context) error {
		h := handlers.NewAccountHandler(db, tasksClient)
		return h.CreateAccount(c)
	})

	e.GET("/v1/accounts", func(c echo.Context) error {
		h := handlers.NewAccountHandler(db, nil)
		return h.ListAccounts(c)
	})

	return e
}
