package apollo

import "github.com/gorilla/mux"

type DefaultApp struct {
	DbSettings     *DatabaseConnection
	RabbitSettings *RabbitConnection
	Router         *mux.Router
}
