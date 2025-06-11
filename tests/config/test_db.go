package tests

import (
	"os"
	"testing"

	"github.com/ashishnagargoje0/shetiseva-backend/config"
)

func TestDBConnection(t *testing.T) {
	os.Setenv("DATABASE_URL", "host=localhost user=postgres password=1016 dbname=shetiseva port=5432 sslmode=disable")
	config.ConnectDB()
	if config.DB == nil {
		t.Fatal("DB connection failed")
	}
}
