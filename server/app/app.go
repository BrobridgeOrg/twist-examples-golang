package app

import (
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/soheilhy/cmux"
	"github.com/spf13/viper"

	"twistserver/app/datastore"
)

type App struct {
	connectionListener cmux.CMux
	httpServer         *HTTPServer
}

func CreateApp() *App {

	// expose port
	port := strconv.Itoa(viper.GetInt("host.port"))

	a := &App{
		httpServer: NewHTTPServer(":" + port),
	}

	return a
}
func (a *App) Init() error {
	datastore.CreateAccount("fred", 5000)
	datastore.CreateAccount("armani", 1000)
	// Initializing connection listener
	port := strconv.Itoa(viper.GetInt("host.port"))
	err := a.CreateConnectionListener(":" + port)
	if err != nil {
		return err
	}

	// Initialize HTTP server
	err = a.httpServer.Init(a)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) Run() error {

	// HTTP
	go func() {
		err := a.httpServer.Serve()
		if err != nil {
			log.Error(err)
		}
	}()

	err := a.Serve()
	if err != nil {
		return err
	}

	return nil
}
