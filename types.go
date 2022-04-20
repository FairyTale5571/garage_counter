package main

import (
	"database/sql"
	"time"
)

var db *sql.DB
var cdb *sql.DB

type ServerPerformance struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time

	ServerType string
	Fps        uint
	Cps        uint
	Players    uint
	Objects    uint
}
