package main

import (
	"log"

	"github.com/Kong/go-pdk"
	"github.com/Kong/go-pdk/server"
)

var (
	Version  = "1.0.0"
	Priority = 1
)

func main() {
	server.StartServer(New, Version, Priority)
}

type Config struct {
}

func New() interface{} {
	return &Config{}
}

func (conf Config) Access(kong *pdk.PDK) {
	auth, err := kong.Request.GetHeader("Authorization")
	if err != nil {
		log.Println("Error GET Header: ", err)
		return
	}

	log.Println("Authorization Header", "Authorization", auth)
	kong.Log.Debug("KONG-DEBUG-WARD-PLUGIN: ", auth)
}
