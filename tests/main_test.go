package tests

import (
	"os"
	"testing"

	"github.com/ashishnagargoje0/backend/config"
	"github.com/ashishnagargoje0/backend/controllers"
)

func TestMain(m *testing.M) {
	// ✅ Connect to MongoDB before running any tests
	config.ConnectMongoDB()

	// ✅ Initialize all collections used in tests
	controllers.InitCartCollection()
	//	controllers.InitWishlistCollection()
	//	controllers.InitCompareCollection()
	controllers.InitAuthCollection()
	// controllers.InitProductCollection()
	controllers.InitOrderCollection()

	// ✅ Run tests
	os.Exit(m.Run())
}
