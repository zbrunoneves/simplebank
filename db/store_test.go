package db

import (
	"context"
	"database/sql"
	"testing"

	"simplebank/db/sqlc"

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
		db  *sql.DB
		dbc *mysql.MySQLContainer
	)

	BeforeEach(func() {
		ctx := context.Background()

		var err error
		dbc, err = mysql.Run(ctx,
			"mysql:8.0.36",
			mysql.WithDatabase("bank"),
			mysql.WithUsername("admin"),
			mysql.WithPassword("123"),
			mysql.WithScripts("./schema.sql"),
		)
		if err != nil {
			log.Fatal().Msgf("failed to start mysql container: %s", err)
		}

		conn, err := dbc.ConnectionString(ctx)
		if err != nil {
			log.Fatal().Msgf("failed to get mysql connection string")
		}

		db, err = sql.Open("mysql", conn)
		if err != nil {
			log.Fatal().Msgf("failed to open database: %v", err)
		}
		if err := db.Ping(); err != nil {
			log.Fatal().Msgf("failed to ping database: %v", err)
		}
	})

	AfterEach(func() {
		db.Close()

		if err := testcontainers.TerminateContainer(dbc); err != nil {
			log.Printf("failed to terminate mysql container: %s", err)
		}
	})

	Describe("New", func() {
		It("should create a new store instance", func() {
			store := New(db)

			var result int
			err := store.db.QueryRow("SELECT 1").Scan(&result)

			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal(1))
		})
	})

	Describe("Account", func() {
		It("should create an account", func() {
			store := New(db)

			params := sqlc.CreateAccountParams{
				Owner:    "John Doe",
				Balance:  100.0,
				Currency: "$",
			}

			res, err := store.CreateAccount(context.Background(), params)

			Expect(err).NotTo(HaveOccurred())
			Expect(res.LastInsertId()).To(BeEquivalentTo(1))
			Expect(res.RowsAffected()).To(BeEquivalentTo(1))
		})
	})
})
