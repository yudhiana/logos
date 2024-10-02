package main

import (
	"log"
	"log/slog"

	"github.com/Kong/go-pdk"
	"github.com/Kong/go-pdk/server"
)

const (
	Version  = "0.2"
	Priority = 1
)

func main() {
	server.StartServer(New, Version, Priority)
}

type Tracing struct {
}

func New() interface{} {
	return &Tracing{}
}

func (t Tracing) Access(kong *pdk.PDK) {
	auth, err := kong.Request.GetHeader("Authorization")
	if err != nil {
		log.Println("Error GET Header: ", err)
		return
	}

	slog.Info("Authorization Header", "Authorization", auth)
}
