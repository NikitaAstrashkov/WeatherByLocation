package main

import (
	"WeatherByLocation/models"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ip2location/ip2location-go"
	"net/http"
	"os"
	"strings"
)

func getWeather(w http.ResponseWriter, r *http.Request) {
	var err error
	var apiResponse models.Response
	var request string
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if models.Misc.City == "" {
		models.Misc.Coords.Lat, models.Misc.Coords.Lon = getLocationByIp(r)
		request = strings.Replace(models.Misc.Consts.WeatherByLocationReq, "{lon}", fmt.Sprintf("%f", models.Misc.Coords.Lon), 1)
		request = strings.Replace(request, "{lat}", fmt.Sprintf("%f", models.Misc.Coords.Lat), 1)
	} else {
		request = strings.Replace(models.Misc.Consts.WeatherByCityReq, "{city name}", models.Misc.City, 1)
	}

	response, err := http.Get(request)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	err = json.NewDecoder(response.Body).Decode(&apiResponse)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	models.Misc.City = apiResponse.Name
	ourResp, _ := json.Marshal(apiResponse)

	_, err = w.Write(ourResp)
}

func getLocationByIp(r *http.Request) (float32, float32) {
	db, err := ip2location.OpenDB("./IP2LOCATION-LITE-DB5.BIN")
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	ip := r.RemoteAddr
	results, err := db.Get_all(ip)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	return results.Latitude, results.Longitude
}

func setLocation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err := json.NewDecoder(r.Body).Decode(&models.Misc)

	if err != nil {
		fmt.Print(err.Error())
		w.WriteHeader(http.StatusLocked)
		os.Exit(1)
	}

	w.WriteHeader(http.StatusCreated)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/weather", getWeather).Methods(http.MethodGet)
	router.HandleFunc("/setLocation", setLocation).Methods(http.MethodPost)
	http.Handle("/", router)

	_ = http.ListenAndServe(":8181", nil)
}
