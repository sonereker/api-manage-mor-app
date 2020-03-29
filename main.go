package main

import (
	_ "github.com/lib/pq"
	"github.com/sonereker/kule-app-api/app"
	"github.com/sonereker/kule-app-api/config"
)

func main() {
	c := config.GetConfig()

	a := &app.App{}
	a.Initialize(c)
	a.Run(c.App.Host)
}
