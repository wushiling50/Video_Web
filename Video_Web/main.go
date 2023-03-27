package main

import (
	"main/Video_Web/conf"
	"main/Video_Web/routes"
)

func main() {
	conf.Init()
	r := routes.NewRouter()
	r.Run(conf.HttpPort)
}
