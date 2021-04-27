package main

import (
	_ "github.com/corinz/apigateway/docs"
	"github.com/corinz/apigateway/internal/app"
)

func main() {
	a := app.NewApp()
	a.Startup()
}
