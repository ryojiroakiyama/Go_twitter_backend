package dao

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

// Interface of configureation
type DBConfig interface {
	FormatDSN() string
}

// Prepare sqlx.DB
func initDb(config DBConfig) (*sqlx.DB, error) {
	driverName := "mysql"
	db, err := sqlx.Open(driverName, config.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("sqlx.Open failed: %w", err)
	}

	// 以下, mysqlのサーバが立ち上がる前に接続確認をしてしまうためコメントアウト
	// err = db.Ping()
	// if err != nil {
	// 	return nil, fmt.Errorf("sqlx.Ping failed: %w", err)
	// }

	// 必要に応じて以下の設定
	//db.SetConnMaxLifetime(0)
	//pool.SetMaxIdleConns(3)
	//pool.SetMaxOpenConns(3)

	return db, nil
}
