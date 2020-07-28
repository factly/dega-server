package space

import (
	"os"
	"testing"

	"github.com/factly/dega-server/config"
	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	err := godotenv.Load("../../../../.env")
	_ = err

	config.SetupDB()

	exitValue := m.Run()

	os.Exit(exitValue)
}
