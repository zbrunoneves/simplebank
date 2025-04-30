package database

import (
	"context"
	"database/sql"
	"errors"
	"sync"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"

	_ "github.com/go-sql-driver/mysql"
)

func TestStore(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Store Suite")
}

var _ = Describe("Store", func() {
	var (
		ctx context.Context

		db        *sql.DB
		container *mysql.MySQLContainer
	)

	BeforeEach(func() {
		ctx = context.Background()

		var err error
		container, err = mysql.Run(ctx,
			"mysql:9.3.0",
			mysql.WithDatabase("bank"),
			mysql.WithUsername("admin"),
			mysql.WithPassword("123"),
			mysql.WithScripts("./schema.sql"),
		)
		if err != nil {
			log.Fatal().Msgf("failed to start mysql container: %s", err)
		}

		conn, err := container.ConnectionString(ctx)
		if err != nil {
			log.Fatal().Msgf("failed to get mysql connection string")
		}

		db, err = sql.Open("mysql", conn)
		if err != nil {
			log.Fatal().Msgf("failed to open db: %v", err)
		}
		if err := db.Ping(); err != nil {
			log.Fatal().Msgf("failed to ping db: %v", err)
		}
	})

	AfterEach(func() {
		db.Close()

		if err := testcontainers.TerminateContainer(container); err != nil {
			log.Fatal().Msgf("failed to terminate mysql container: %s", err)
		}
	})

	Describe("ExecTx", func() {
		It("should execute N transactions successfully", func() {

			_, err := db.ExecContext(context.Background(), "INSERT INTO account (owner, balance) VALUES (?, ?)", "John Doe", 1000)
			Expect(err).NotTo(HaveOccurred())

			wg := &sync.WaitGroup{}
			wg.Add(10)
			for i := 0; i < 10; i++ {
				go func() {
					_ = ExecTx(context.Background(), db, func(tx *sql.Tx) error {
						var curr int
						err = tx.QueryRow("SELECT balance FROM account WHERE owner = ? FOR UPDATE", "John Doe").Scan(&curr)
						if err != nil {
							return err
						}

						_, err := tx.Exec("UPDATE account set balance = ? WHERE owner = ?", curr-10, "John Doe")
						return err
					})

					wg.Done()
				}()
			}
			wg.Wait()

			var curr int
			err = db.QueryRow("SELECT balance FROM account WHERE owner = ?", "John Doe").Scan(&curr)
			Expect(err).NotTo(HaveOccurred())
			Expect(curr).To(Equal(900))
		})

		It("should abort half transactions on error", func() {
			_, err := db.ExecContext(context.Background(), "INSERT INTO account (owner, balance) VALUES (?, ?)", "John Doe", 1000)
			Expect(err).NotTo(HaveOccurred())

			wg := &sync.WaitGroup{}
			wg.Add(10)
			for i := 0; i < 10; i++ {
				go func() {
					_ = ExecTx(context.Background(), db, func(tx *sql.Tx) error {
						var curr int
						err = tx.QueryRow("SELECT balance FROM account WHERE owner = ? FOR UPDATE", "John Doe").Scan(&curr)
						if err != nil {
							return err
						}

						_, err := tx.Exec("UPDATE account set balance = ? WHERE owner = ?", curr-10, "John Doe")
						if err != nil {
							return err
						}

						if i%2 == 0 {
							return errors.New("forced error")
						}

						return nil
					})

					wg.Done()
				}()
			}
			wg.Wait()

			var curr int
			err = db.QueryRow("SELECT balance FROM account WHERE owner = ?", "John Doe").Scan(&curr)
			Expect(err).NotTo(HaveOccurred())
			Expect(curr).To(Equal(950))
		})
	})
})
