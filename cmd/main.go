package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

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
var analyticsService *services.AnalyticsService
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
	analyticsRepo := repository.NewRedisAnalyticsRepository(rdb)

	identifierService := core.NewDefaultSnowflakeGenerator(int64(getMachineId()), int64(appConfigurationService.GetConfig().App.DatacenterId))

	shortnerService = services.NewShortnerService(redisRepo, identifierService)
	analyticsService = services.NewAnalyticsService(analyticsRepo)

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
	r.HandleFunc("/", index)
	r.HandleFunc("/shorten", shorten)
	r.HandleFunc("/analytics/{hashValue}", getAnalytics).Methods("GET")
	r.HandleFunc("/analytics/{hashValue}/clicks", getClickHistory).Methods("GET")
	r.HandleFunc("/{hashValue}", decode)
	err := http.ListenAndServe(":"+appConfigurationService.GetConfig().App.Port, r)
	if err != nil {
		panic("Error while starting server")
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	contents, err := os.ReadFile("index.txt")
	if err != nil {
	   log.Errorf("Error while reading index.txt: %s", err.Error())
	}
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

	// Record click analytics asynchronously
	go func() {
		ip := getClientIP(r)
		userAgent := r.UserAgent()
		referer := r.Referer()
		if err := analyticsService.RecordClick(hashValue, ip, userAgent, referer); err != nil {
			log.Errorf("Error recording click analytics: %s", err.Error())
		}
	}()

	http.Redirect(w, r, originalUrl, http.StatusSeeOther)
}

// getClientIP extracts the client IP address from the request
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header first (for proxies/load balancers)
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		// Take the first IP in the list
		ips := strings.Split(forwarded, ",")
		return strings.TrimSpace(ips[0])
	}

	// Check X-Real-IP header
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// Fall back to RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

// getAnalytics returns analytics summary for a short URL
func getAnalytics(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hashValue := vars["hashValue"]

	analytics, err := analyticsService.GetAnalytics(hashValue)
	if err != nil {
		log.Errorf("Error getting analytics: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, INTERNAL_ERROR_MESSAGE)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(analytics); err != nil {
		log.Errorf("Error encoding analytics response: %s", err.Error())
	}
}

// getClickHistory returns detailed click history for a short URL
func getClickHistory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hashValue := vars["hashValue"]

	// Parse limit from query params, default to 100
	limitStr := r.URL.Query().Get("limit")
	limit := int64(100)
	if limitStr != "" {
		if parsedLimit, err := strconv.ParseInt(limitStr, 10, 64); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	analytics, err := analyticsService.GetAnalyticsWithHistory(hashValue, limit)
	if err != nil {
		log.Errorf("Error getting click history: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, INTERNAL_ERROR_MESSAGE)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(analytics); err != nil {
		log.Errorf("Error encoding click history response: %s", err.Error())
	}
}
