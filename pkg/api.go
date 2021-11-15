package pkg

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type config struct {
	port int
	env  string
}

type application struct {
	config config
	logger *log.Logger
	pulse  *Pulse
}

func HttpAPI(p *Pulse, port int, env string) {
	cfg := config{
		port: port,
		env:  env,
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	app := &application{
		config: cfg,
		logger: logger,
		pulse:  p,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/status/all", app.statusHandler)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      mux,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Printf("starting %s server on %s", cfg.env, srv.Addr)
	err := srv.ListenAndServe()
	logger.Fatal(err)

}

func (app *application) statusHandler(w http.ResponseWriter, r *http.Request) {
	app.pulse.mutex.RLock()
	defer app.pulse.mutex.RUnlock()
	fmt.Fprintln(w, app.pulse.nodeMap)
}
