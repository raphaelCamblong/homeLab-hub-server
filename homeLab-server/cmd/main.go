package main

import (
	"crypto/tls"
	"homelab.com/homelab-server/homeLab-server/app"
	"net/http"
)

func useInsecureHttpTLS() {
	http.DefaultClient.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
}

func main() {
	useInsecureHttpTLS()

	app.NewApp().Start()
}
