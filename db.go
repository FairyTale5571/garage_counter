package main

import (
	"database/sql"
	"fmt"
	"garage_counter/logger"
	sq "github.com/Masterminds/squirrel"
	_ "github.com/go-sql-driver/mysql"
)

// Global variable in bd package
var db *sql.DB

// ConnectDatabase connect to database
func ConnectDatabase() (*sql.DB, error) {

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		dbCon.User,
		dbCon.Password,
		dbCon.Ip,
		dbCon.Port,
		dbCon.Database))
	if err != nil {
		logger.PrintLog("Database not connected\n")
		return nil, err
	}
	db.SetMaxOpenConns(10)

	statement, err := db.Prepare("SELECT VERSION()")
	if err != nil {
		logger.PrintLog(err.Error())
		return nil, err
	}
	rows, err := statement.Query() // execute our select statement
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	for rows.Next() {
		var title string
		rows.Scan(&title)
		logger.PrintLog("Database version: %v\n", title)
	}
	dbConnected = true

	return db, nil
}

func countVeh(id string,array []string) string {
	stmt := sq.Select("COUNT(*)").From("vehicles").Where(sq.Eq{"pid":id},sq.NotEq{"classname": array}).Where(sq.Eq{"deleted_at":nil}).Where(sq.NotEq{"classname":array})
	sql, args, _ := stmt.ToSql()
	rows, err := db.Query(sql,args...)
	if err != nil {
		logger.PrintLog("error: %v",err.Error())
		return "0"
	}
	defer rows.Close()

	var ret string
	for rows.Next() {
		if err := rows.Scan(&ret); err != nil {
			logger.PrintLog("error: %v",err.Error())
		}
		return ret
	}
	return "0"
}
