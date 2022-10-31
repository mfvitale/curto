package main

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/mfvitale/curto/repository"
	"github.com/mfvitale/curto/services"
	"github.com/mfvitale/curto/services/config"
	"github.com/mfvitale/curto/services/core"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()
var rdb *redis.Client
var shortnerService services.ShortnerService

func init() {

	log.Out = os.Stdout
	log.SetLevel(logrus.InfoLevel)

	rdb = redis.NewClient(&redis.Options{
		Addr:     config.GetConfig().Redis.Endpoint,
		Username: config.GetConfig().Redis.Username,
		Password: config.GetConfig().Redis.Password,
		DB:       0,
	})

	redisRepo := repository.NewRedisUrlRepository(rdb)


	identifierService := core.NewSnowflakeGenerator(int64(getMachineId()), int64(config.GetConfig().App.DatacenterId))

	shortnerService = services.NewShortnerService(redisRepo, identifierService)

}

func getMachineId() int {

	var compRegEx = regexp.MustCompile(".*-([0-9]*)")
    match := compRegEx.FindStringSubmatch(config.GetConfig().App.PodName)
	id, _ := strconv.Atoi(match[1])
	return id
}

func Atoi(s string) {
	panic("unimplemented")
}

func main() {

	log.Infof("Server is running on port %s", config.GetConfig().App.Port)
	r := mux.NewRouter()
	r .HandleFunc("/", index)
	r .HandleFunc("/encode", encode)
	r .HandleFunc("/{hashValue}", decode)
	err := http.ListenAndServe(":"+config.GetConfig().App.Port, r)
	if err != nil {
		panic("Error while starting server")
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello! I'm curto, your URL shortner service!")
}

func encode(w http.ResponseWriter, r *http.Request) {

	log.Info("Pod "+ os.Getenv("POD_NAME")) //TODO remove
	url := r.URL.Query().Get("url")

	hashValue := shortnerService.Shorten(url)

	fmt.Fprintf(w, config.GetConfig().App.Domain+hashValue)
}

func decode(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
    hashValue, ok := vars["hashValue"]
	//TODO move this when hashvalue is not on redis
    if !ok {
        fmt.Fprintf(w, "URL not found")
		w.WriteHeader(http.StatusBadRequest)
    }

	originalUrl := shortnerService.Resolve(hashValue)

	http.Redirect(w,r, originalUrl, http.StatusSeeOther)
}
