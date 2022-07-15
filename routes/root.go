package routes

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/hookart/twitter-mentions/models"
	"github.com/spf13/viper"
)

func Serve() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	models.GetDBConnection()
	r.Get("/health", Healthcheck)
	r.Post("/verify", Verify)
	InitJWTKey()

	http.ListenAndServe(fmt.Sprintf(":%d", viper.GetInt("port")), r)
}
