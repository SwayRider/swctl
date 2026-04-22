package internal

import "github.com/joho/godotenv"

func init() {
	_ = godotenv.Load("swctl.conf")
}
