package account

import (
	"time"
)

type Account struct {
	ID        int
	Owner     string
	Balance   uint
	Currency  string
	CreatedAt time.Time
}
