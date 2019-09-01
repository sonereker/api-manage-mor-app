package main

import (
	_ "github.com/lib/pq"
	"github.com/sonereker/api-manage-mor-app/app"
	"github.com/sonereker/api-manage-mor-app/config"
)

func main() {
	c := config.GetConfig()

	a := &app.App{}
	a.Initialize(c)
	a.Run(c.App.Host)
}
