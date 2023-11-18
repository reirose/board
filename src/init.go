package src

import (
	"database/sql"
	"github.com/go-chi/chi/v5"
)

var Router *chi.Mux
var Database *sql.DB
