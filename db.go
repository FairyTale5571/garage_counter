package main

import (
	"database/sql"
	"fmt"
	"garage_counter/logger"
	"github.com/ClickHouse/clickhouse-go/v2"
	sq "github.com/Masterminds/squirrel"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

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

	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(60 * time.Minute)

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

func ConnectClickHouse() (*sql.DB, error) {
	fmt.Printf("connecting to clickhouse...\n")
	const dsn = "tcp://ts.rimasrp.life:8123?database=rimas&username=default&password=1111&read_timeout=10&write_timeout=20"

	conn := clickhouse.OpenDB(&clickhouse.Options{
		Addr: []string{"ts.rimasrp.life:8123"},
		Auth: clickhouse.Auth{
			Database: "rimas",
			Username: "default",
			Password: "1111",
		},
		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
		DialTimeout: 5 * time.Second,
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
		Debug: true,
	})
	conn.SetMaxIdleConns(5)
	conn.SetMaxOpenConns(10)
	conn.SetConnMaxLifetime(time.Hour)

	statement, err := conn.Prepare("SELECT VERSION()")
	if err != nil {
		logger.PrintLog(err.Error())
		return nil, err
	}
	rows, err := statement.Query()
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	for rows.Next() {
		var title string
		rows.Scan(&title)
		logger.PrintLog("Clickhouse Database version: %v\n", title)
	}

	return conn, nil
}

func insertCpuC(server, cps, fps, players, objects string) string {

	return "0"
}

func insertCpu(server, cps, fps, players, objects string) string {
	stmt, args, _ := sq.Insert("server_performance").Columns("server_type", "cps", "fps", "players", "objects", "insert_time").
		Values(server, cps, fps, players, objects, time.Now()).ToSql()
	row, err := db.Query(stmt, args...)
	defer row.Close()
	if err != nil {
		logger.PrintLog("error: insert cpu: %s\n", err.Error())
		return "fail"
	}
	return "0"
}

func countVeh(id string, array []string) string {
	stmt := sq.
		Select("COUNT(*)").
		From("vehicles").
		Where(sq.Eq{"pid": id}, sq.NotEq{"classname": array}).
		Where(sq.Eq{"deleted_at": nil}).
		Where(sq.NotEq{"classname": array})

	sql, args, _ := stmt.ToSql()
	rows, err := db.Query(sql, args...)
	defer rows.Close()

	if err != nil {
		logger.PrintLog("error: countVeh: %v", err.Error())
		return "0"
	}

	var ret string
	for rows.Next() {
		if err := rows.Scan(&ret); err != nil {
			logger.PrintLog("error: countVeh: %v", err.Error())
		}
		return ret
	}
	return "0"
}

func deleteOldVehicles(vehicles []string) string {
	stmt, args, _ := sq.
		Update("vehicles").
		Set("deleted_at", time.Now()).
		Set("why_deleted", "Не использовалась 90 дней").
		Where(sq.NotEq{"classname": vehicles}).
		Where(sq.Eq{"deleted_at": nil}).
		Where(sq.Expr("last_used < DATE_SUB(NOW(), INTERVAL ? DAY)", 90)).
		ToSql()

	rows, err := db.Exec(stmt, args...)
	if err != nil {
		logger.PrintLog("error: deleteOldVehicles: %s\n", err.Error())
		return "fail"
	}
	affected, err := rows.RowsAffected()
	if err != nil {
		logger.PrintLog("cant get affected rows: %s", err.Error())
		return "fail"
	}
	logger.PrintLog("delete old cars: affected rows %d\n", affected)
	return "nice"
}
