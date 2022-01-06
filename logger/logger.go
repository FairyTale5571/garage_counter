package logger

import (
	"log"
	"os"
	"path/filepath"
)

var(
	tempDir string
	logFile *os.File
	loger *log.Logger
)

func init() {
	SetupLogger()
	PrintLog("Logger started")
}

func PrintLog(format string, v ...interface{}) {
	log.Printf(format, v...)
	loger.Printf(format,v...)
}

func SetupLogger() *log.Logger {
	dir, _ := os.Getwd()
	tmpfn := filepath.Join(dir, "counter.log")
	logFile, err := os.OpenFile(tmpfn, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	loger = log.New(logFile, "", log.LstdFlags)
	return loger
}

func TeardownLogger() {
	if logFile != nil {
		logFile.Close()
	}
}