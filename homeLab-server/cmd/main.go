package main

import (
	"crypto/tls"
	"net/http"

	"github.com/sirupsen/logrus"
	"homelab.com/homelab-server/homeLab-server/init/config"
	processor "homelab.com/homelab-server/homeLab-server/internal/core"
	"homelab.com/homelab-server/homeLab-server/pkg/tools"
)

func useInsecureHttpTLS() {
	http.DefaultClient.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
}

// @title task-mgmt-api
// @version 1.0
// @description This is task-mgmt api
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email balajichandrasekar17@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	useInsecureHttpTLS()
	defer handlePanic()
	config.Get().FlattenConfigQuick()
	tools.UpdateDocs(config.Get())

	processor := processor.NewProcessor()
	processor.Start()
}

func handlePanic() {
	if r := recover(); r != nil {
		logrus.Error("⚠️ Panic intercepted", r)
	}
}
