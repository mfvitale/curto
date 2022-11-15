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
var appConfigurationService config.AppConfigurationService

const INTERNAL_ERROR_MESSAGE = "Ops! Seems something is not working. Plaese try again"

func init() {

	log.Out = os.Stdout
	log.SetLevel(logrus.InfoLevel)

	appConfigurationService = config.NewAppConfigurationService("config.yml", log)

	rdb = redis.NewClient(&redis.Options{
		Addr:     appConfigurationService.GetConfig().Redis.Endpoint,
		Username: appConfigurationService.GetConfig().Redis.Username,
		Password: appConfigurationService.GetConfig().Redis.Password,
		DB:       0,
	})

	redisRepo := repository.NewRedisUrlRepository(rdb)


	identifierService := core.NewSnowflakeGenerator(int64(getMachineId()), int64(appConfigurationService.GetConfig().App.DatacenterId))

	shortnerService = services.NewShortnerService(redisRepo, identifierService)

}

func getMachineId() int {

	if appConfigurationService.GetConfig().App.MachineId != -1 {
		return appConfigurationService.GetConfig().App.MachineId
	}

	var compRegEx = regexp.MustCompile(".*-([0-9]*)")
    match := compRegEx.FindStringSubmatch(appConfigurationService.GetConfig().App.PodName)
	id, _ := strconv.Atoi(match[1])
	return id
}

func Atoi(s string) {
	panic("unimplemented")
}

func main() {

	log.Infof("Server[%d] is running on port %s", getMachineId(), appConfigurationService.GetConfig().App.Port)
	r := mux.NewRouter()
	r .HandleFunc("/", index)
	r .HandleFunc("/shorten", shorten)
	r .HandleFunc("/{hashValue}", decode)
	err := http.ListenAndServe(":"+appConfigurationService.GetConfig().App.Port, r)
	if err != nil {
		panic("Error while starting server")
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	contents, _ := os.ReadFile("index.txt")

	fmt.Fprint(w, string(contents))
}

func shorten(w http.ResponseWriter, r *http.Request) {

	url := r.URL.Query().Get("url")

	hashValue, err := shortnerService.Shorten(url)
	if _, ok := err.(services.ShortenOperationError); ok {
		w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprint(w, INTERNAL_ERROR_MESSAGE)
		return
    }
	fmt.Fprintf(w, appConfigurationService.GetConfig().App.Domain+hashValue)
}

func decode(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
    hashValue := vars["hashValue"]

	originalUrl, err := shortnerService.Resolve(hashValue)
	if serr, ok := err.(repository.ShortURLNotFoundError); ok {
		w.WriteHeader(http.StatusNotFound)
        fmt.Fprint(w, serr.Error())
		return
    }

	http.Redirect(w,r, originalUrl, http.StatusSeeOther)
}
