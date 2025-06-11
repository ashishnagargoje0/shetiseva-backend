package tests

import (
    "testing"

    "github.com/ashishnagargoje0/shetiseva-backend/config"
)

func TestConnectTestDB(t *testing.T) {
    // Initialize test DB connection
    config.ConnectTestDB()

    if config.TestDB == nil {
        t.Fatal("❌ TestDB connection is nil")
    }

    sqlDB, err := config.TestDB.DB()
    if err != nil {
        t.Fatalf("❌ Failed to retrieve *sql.DB from GORM: %v", err)
    }

    // Ensure connection can be pinged
    if err := sqlDB.Ping(); err != nil {
        t.Fatalf("❌ TestDB Ping failed: %v", err)
    }

    // Cleanup: close DB connection after test
    defer func() {
        if err := sqlDB.Close(); err != nil {
            t.Errorf("⚠️ Error closing DB: %v", err)
        }
    }()

    t.Log("✅ Successfully connected and pinged TestDB")
}
