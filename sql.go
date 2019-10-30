package apollo

import (
	"database/sql"
	"fmt"
)

type DatabaseConnection struct {
	Server     string `json:"Server"`
	Port       int    `json:"Port"`
	UserId     string `json:"UserId"`
	Password   string `json:"Password"`
	Database   string `json:"Database"`
	DriverName string `json:"DriverName"`
	conn       *sql.DB
}

func (dc *DatabaseConnection) PostgresConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dc.Server, dc.Port, dc.UserId, dc.Password, dc.Database)
}

func (dc *DatabaseConnection) GetDbConnection() (*sql.DB, error) {
	if dc.conn != nil {
		return dc.conn, nil
	}

	conn, err := sql.Open(dc.DriverName, dc.PostgresConnectionString())
	if err != nil {
		return nil, fmt.Errorf("unable to open db connection: %w", err)
	}
	dc.conn = conn

	return conn, nil
}
