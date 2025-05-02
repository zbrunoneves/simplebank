package api

import (
	"bytes"
	"database/sql"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/hibiken/asynq"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
)

func NewWebServer(db *sql.DB, tasksClient *asynq.Client) *echo.Echo {
	e := echo.New()
	e.Validator = &Validator{validate: validator.New()}

	e.HidePort = true
	e.HideBanner = true
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Recover())
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogMethod: true,
		LogStatus: true,
		LogError:  true,
		BeforeNextFunc: func(c echo.Context) {
			var buf bytes.Buffer
			tee := io.TeeReader(c.Request().Body, &buf)

			b, err := io.ReadAll(tee)
			if err != nil {
				return
			}
			s := strings.ReplaceAll(string(b), "\n", " ")
			re := regexp.MustCompile(`\s{2,}`)
			s = re.ReplaceAllString(s, " ")
			c.Set("body", strings.TrimSpace(s))

			c.Request().Body = io.NopCloser(bytes.NewReader(buf.Bytes()))
		},
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			l := log.Info()
			if v.Error != nil {
				l = log.Error()
			}

			body, ok := c.Get("body").(string)
			if !ok || len(body) > 0 {
				l = l.Str("body", body)
			}

			l.Str("uri", v.URI).
				Str("method", v.Method).
				Int("status", v.Status).
				Msg("request")

			return nil
		},
	}))

	routes(e, db, tasksClient)

	return e
}

type Validator struct {
	validate *validator.Validate
}

func (v *Validator) Validate(i any) error {
	if err := v.validate.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	return nil
}
