package app

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"github.com/smolathon/internal/config"
	"github.com/smolathon/internal/database"
	"github.com/smolathon/internal/serve"
	"github.com/smolathon/pkg/client/mongodb"
	"github.com/smolathon/pkg/logging"
)

type App struct {
	cfg        *config.Config
	logger     *logging.Logger
	router     *httprouter.Router
	httpServer *http.Server
	//ctx        context.Context
	//pgClient   *pgxpool.Pool
}

func NewApp(config *config.Config, logger *logging.Logger, ctx context.Context) (App, error) {

	logger.Println("router initializing")
	router := httprouter.New()
	//--------------
	cfgMongo := config.MongoDB
	mongoDBClient, err := mongodb.NewCLient(context.Background(), cfgMongo.Host, cfgMongo.Port, cfgMongo.Username, cfgMongo.Password,
		cfgMongo.Database, cfgMongo.AuthDB)
	if err != nil {
		panic(err)
	}
	logger.Infof("%v", (*mongoDBClient).Client().ListDatabases)
	collections := []string{cfgMongo.MasterColletion, cfgMongo.CardsCollecion}
	storage := database.NewStorage(mongoDBClient, collections, logger)

	//-------------------
	ms := database.NewMasterStorage(storage)
	cs := database.NewCardStorage(storage)

	//fmt.Printf("%v -------------", collections)
	handler := serve.NewHandler(logger, ms, cs)
	handler.Register(router)
	return App{
		cfg:    config,
		logger: logger,
		router: router,
		// pgClient: pgClient,
	}, nil
}

func (a *App) Run() {
	a.startHTTP()
}

func (a *App) startHTTP() {
	a.logger.Info("start HTTP")

	var listener net.Listener

	if a.cfg.Listen.Type == config.LISTEN_TYPE_SOCK {
		appDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			a.logger.Fatal(err)
		}
		socketPath := path.Join(appDir, a.cfg.Listen.SocketFile)
		a.logger.Infof("socket path: %s", socketPath)

		a.logger.Info("create and listen unix socket")
		listener, err = net.Listen("unix", socketPath)
		if err != nil {
			a.logger.Fatal(err)
		}
	} else {
		a.logger.Infof("bind app to host %s and port %s,", a.cfg.Listen.BindIP, a.cfg.Listen.Port)
		var err error
		listener, err = net.Listen("tcp", fmt.Sprintf("%s:%s", a.cfg.Listen.BindIP, a.cfg.Listen.Port))
		if err != nil {
			a.logger.Fatal(err)
		}
	}

	c := cors.New(cors.Options{
		AllowedMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodPut, http.MethodOptions, http.MethodDelete},
		AllowedOrigins:     []string{"*"},
		AllowCredentials:   true,
		AllowedHeaders:     []string{"Location", "Charset", "Access-Control-Allow-Origin", "Content-Type", "content-type", "Origin"},
		OptionsPassthrough: true,
		ExposedHeaders:     []string{"Location", "Authorization", "Content-Disposition"},
		//Enable Debugging for testing, consider disabling in production
		Debug: false,
	})
	handler := c.Handler(a.router)

	a.httpServer = &http.Server{
		Handler:      handler,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	a.logger.Println("app completely initialized and started")

	if err := a.httpServer.Serve(listener); err != nil {
		switch {
		case errors.Is(err, http.ErrServerClosed):
			a.logger.Warn("server shutdown")
		default:
			a.logger.Fatal(err)
		}
	}

}
