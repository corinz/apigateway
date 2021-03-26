package main

import (
	"github.com/corinz/apigateway/internal/app"
)

func main() {
	a := app.NewApp()
	a.Startup(":8080")
}
