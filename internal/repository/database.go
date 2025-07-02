package repository

import (
	"ellas-corner/internal/db"
)

var database *db.Database

func SetDatabase(d *db.Database) {
	database = d
}
