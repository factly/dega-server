package space

import (
	"os"
	"testing"

	"github.com/factly/dega-server/config"
	"github.com/factly/dega-server/service/core/model"
	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	err := godotenv.Load("../../../../.env")
	_ = err

	config.SetupDB()

	config.DB.CreateTable(&model.Space{})

	exitValue := m.Run()
	config.DB.DropTable((&model.Space{}))

	os.Exit(exitValue)
}
