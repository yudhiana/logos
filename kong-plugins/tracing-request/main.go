package main

import (
	"log"
	"time"

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

	kong.Response.SetHeader("X-Request-ID", time.Now().Format(time.RFC3339))

	log.Println("Authorization Header", "Authorization", auth)
	kong.Log.Debug("KONG-DEBUG-WARD-PLUGIN: ", auth)
}


func (conf Config) Log(kong *pdk.PDK) {
	kong.Log.Debug("KONG-DEBUG-WARD-PLUGIN-LOG: ", "Hello World!")
}