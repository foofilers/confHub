package main

import (
	"github.com/foofilers/confHub/server"
	"github.com/foofilers/confHub/conf"
)

func main() {
	conf.InitConfFromFile("confHub",".")
	server.Start("0.0.0.0:8080")
}

