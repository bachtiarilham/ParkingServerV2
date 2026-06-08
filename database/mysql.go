//koneksi DB

package database

import (
	"database/sql"
	"fmt"
	"net/url"
	"time"

	"modulegue/config"

	_ "github.com/go-sql-driver/mysql"
)

func OpenMySQL(cfg config.Config) (*sql.DB, error) {
	params := url.Values{}
	params.Set("parseTime", "true")
	params.Set("charset", "utf8mb4")
	params.Set("multiStatements", "true")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
		params.Encode(),
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}
