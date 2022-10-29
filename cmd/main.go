package main

import (
	"fmt"
	"net/http"
	log "github.com/sirupsen/logrus"

	"github.com/go-redis/redis/v8"
	"github.com/mfvitale/curto/services"
	"github.com/mfvitale/curto/services/config"
)

var rdb *redis.Client

func init() {

	rdb = redis.NewClient(&redis.Options{
		Addr:     config.GetConfig().Redis.Endpoint,
		Username: config.GetConfig().Redis.Username,
		Password: config.GetConfig().Redis.Password,
		DB:       0,
	})

}
func main() {

	log.Info("Server running on port "+ config.GetConfig().App.Port)

	http.HandleFunc("/", index)
	http.HandleFunc("/encode", encode)
	http.ListenAndServe(":"+config.GetConfig().App.Port, nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello! I'm curto, your URL shortner service!")
}

func encode(w http.ResponseWriter, r *http.Request) {

	url := r.URL.Query().Get("url")

	shortnerService := services.NewShortnerService(rdb)
	hashValue := shortnerService.Encode(url)

	fmt.Fprintf(w, hashValue)
	//http.Redirect(w,r, "https://curto-url-shortner-mfvitale.cloud.okteto.net/"+hashValue, http.StatusSeeOther)
}
