package main

import (
	"sync/atomic"

	"github.com/faisalalhatimi/http-server/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
	jwtSecret      string
	polkAPIKEY     string
}
