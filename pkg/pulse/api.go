package pulse

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// Configuration for HTTP server
type config struct {
	port int
	env  string
}

// Application has the Pulse node embedded to query information
type application struct {
	config config
	logger *log.Logger
	pulse  *Pulse
}

// REST API provides status update on the nodes being tracked
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
	mux.HandleFunc("/status/:id", app.statusSingleHandler)

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

func (app *application) statusSingleHandler(w http.ResponseWriter, r *http.Request) {
	app.pulse.mutex.RLock()
	defer app.pulse.mutex.RUnlock()
	fmt.Fprintln(w, app.pulse.nodeMap)
}
