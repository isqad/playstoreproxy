package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"playstoreproxy/internal/handlers"
	"playstoreproxy/internal/log"
	"playstoreproxy/internal/utils"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog/hlog"
	"github.com/spf13/viper"

	flag "github.com/spf13/pflag"
)

func main() {
	logLevel := "info"

	flag.String("listen", "0.0.0.0:8085", "Listening address for http server")
	flag.Bool("debug", true, "Run application in debug env")
	flag.Parse()

	viper.BindPFlags(flag.CommandLine)

	if viper.GetBool("debug") {
		logLevel = "debug"
	}

	log.Init(logLevel, "playstoreproxy")

	r := chi.NewRouter()
	// Some basic middlewares
	r.Use(
		middleware.RealIP,
		hlog.NewHandler(log.Logger),
		hlog.RemoteAddrHandler("ip"),
		hlog.UserAgentHandler("user_agent"),
		hlog.RefererHandler("referer"),
		middleware.Recoverer,
	)

	r.Get("/playstore/check_version", handlers.NewPlayStoreHandler().ServeHTTP)

	quit := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	// Serve static assets
	// serves files from web/static dir
	staticDir, err := utils.StaticDir()
	if err != nil {
		log.Panicf("Failed to get static dir: %v", err)
	}
	r.Method("GET", utils.StaticPrefix+"*", http.StripPrefix(utils.StaticPrefix, http.FileServer(http.Dir(staticDir))))

	r.Get("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		if err := utils.ServeStaticFile(staticDir+"/favicon.ico", w); err != nil {
			log.Infof("Failed to serve static file: %v", err)
		}
	})

	r.Get("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		if err := utils.ServeStaticFile(staticDir+"/robots.txt", w); err != nil {
			log.Infof("Failed to serve static file: %v", err)
		}
	})

	// Handle 404
	r.NotFound(handlers.NewNotFoundHandler().ServeHTTP)

	signal.Notify(quit, os.Interrupt)

	// Configure the HTTP server
	server := &http.Server{
		Addr:              viper.GetString("listen"),
		Handler:           r,
		ReadHeaderTimeout: 1 * time.Second,
		WriteTimeout:      10 * time.Second,
	}

	// Handle shutdown
	server.RegisterOnShutdown(func() {
		close(done)
	})

	// Shutdown the HTTP server
	go func() {
		<-quit
		log.Warn("Server is going shutting down...")

		// Wait 30 seconds for close http connections
		waitIdleConnCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(waitIdleConnCtx); err != nil {
			log.Panicf("Cannot gracefully shutdown the server: %v\n", err)
		}
	}()

	// Start HTTP server
	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Panicf("Server has been closed immediatelly: %v\n", err)
	}

	<-done
	log.Warn("Server stopped...")
}
