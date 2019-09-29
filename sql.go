package apollo

import "fmt"

type DatabaseConnection struct {
	Server     string
	Port       int
	UserId     string
	Password   string
	Database   string
	DriverName string
}

func (dc *DatabaseConnection) PostgresConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dc.Server, dc.Port, dc.UserId, dc.Password, dc.Database)
}