package processor

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	_ "homelab.com/homelab-server/homeLab-server/docs"
	"homelab.com/homelab-server/homeLab-server/infrastructure"
	"homelab.com/homelab-server/homeLab-server/infrastructure/cron"
	"homelab.com/homelab-server/homeLab-server/infrastructure/database"
	"homelab.com/homelab-server/homeLab-server/infrastructure/router/routes"
	"homelab.com/homelab-server/homeLab-server/infrastructure/streaming"
	"homelab.com/homelab-server/homeLab-server/init/config"
	"homelab.com/homelab-server/homeLab-server/internal/repositories"
)

type (
	Processor struct {
		Infra        *infrastructure.Infrastructure
		Repositories *repositories.Repositories
	}
)

func NewProcessor() *Processor {
	return &Processor{
		Infra: &infrastructure.Infrastructure{
			Router:              *initRouter(),
			Db:                  *initDB(),
			ExternalHttpService: initExternalHttpService(),
			StreamHub:           streaming.NewStreamHub(),
			Cron:                cron.NewCron(),
		},
	}
}

func (Processor *Processor) LoadConfig() {
	dotEnvConfig := config.LoadDotEnv(".env")
	if err := config.Get().FlattenConfig(*dotEnvConfig); err != nil {
		logrus.Error("❌ Failed to flatten config: %w", err)
		panic(err)
	}
}

func InjectDB(db database.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("mainDb", db)
		c.Next()
	}
}

func (Processor *Processor) BindContext() {
	Processor.Infra.Router.Get().Use(InjectDB(Processor.Infra.Db))
}

func (Processor *Processor) Start() {
	logrus.Info("Loading repositories")
	Processor.Repositories = repositories.NewRepositories(Processor.Infra)

	logrus.Info("Binding context")
	Processor.BindContext()

	logrus.Info("Binding api routes")
	routes.UserRoutes(Processor.Infra, Processor.Repositories.User)
	routes.ServiceRoutes(Processor.Infra, Processor.Repositories.Service)
	routes.PipelinesRoutes(Processor.Infra, Processor.Repositories.Pipeline)
	routes.DocsRoutes(Processor.Infra)
	routes.ClusterRoutes(Processor.Infra, Processor.Repositories.Cluster)
	routes.NotificationRoutes(Processor.Infra, Processor.Repositories.Notification)

	Processor.Infra.Router.Start()
}
