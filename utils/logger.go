package utils


import (
    "io"
    "log"
    "os"
    "path/filepath"
)

var (
    InfoLogger  *log.Logger
    ErrorLogger *log.Logger
)

func InitLogger() {
    // Ensure the logs directory exists
    err := os.MkdirAll("logs", os.ModePerm)
    if err != nil {
        log.Fatalf("Failed to create logs directory: %v", err)
    }

    // Open log file (append if it exists)
    logFilePath := filepath.Join("logs", "app.log")
    file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        log.Fatalf("Failed to open log file: %v", err)
    }

    // Multi-writer to write both to file and console
    multiWriter := io.MultiWriter(file, os.Stdout)

    // Initialize loggers
    InfoLogger = log.New(multiWriter, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
    ErrorLogger = log.New(multiWriter, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}
