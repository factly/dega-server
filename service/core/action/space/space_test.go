package space

import (
	"os"
	"testing"

	"github.com/factly/dega-server/config"
	"github.com/factly/dega-server/service/core/model"
	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	// os.Setenv("DSN", "postgres://postgres:postgres@localhost:5432/dega?sslmode=disable")
	err := godotenv.Load("../../../../.env")
	_ = err

	config.SetupDB()

	config.DB.CreateTable(&model.Space{})

	exitValue := m.Run()
	config.DB.DropTable((&model.Space{}))

	os.Exit(exitValue)
}
