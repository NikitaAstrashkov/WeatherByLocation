package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ip2location/ip2location-go"
	"net/http"
	"os"
	"strings"
)

var (
	db                   *ip2location.DB
	weatherUnits         = "&units=metric"
	weatherApiKey        = "77699bd83349c01e6271ac022ba11ded"
	weatherByLocationReq = "https://api.openweathermap.org/data/2.5/weather?lat={lat}&lon={lon}&appid=" + weatherApiKey + weatherUnits
	weatherByCityReq     = "https://api.openweathermap.org/data/2.5/weather?q={city name}&appid=" + weatherApiKey + weatherUnits
	latitude             string
	longtitude           string
)

type Location struct {
	City string
}

type Response struct {
	Coord Coordinates
	Main  Weather
	Name  string
}

type Weather struct {
	Temp       float32
	Feels_like float32
	Temp_min   float32
	Temp_max   float32
	Pressure   int
	Humidity   int
}

type Coordinates struct {
	Lon float32
	Lat float32
}

var currentUserLoc Location
var apiResponse Response

func getWeather(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	var request string
	if currentUserLoc.City == "" {
		latitude, longtitude = getLocationByIp(r)
		request = strings.Replace(strings.Replace(weatherByLocationReq, "{lon}", longtitude, 1), "{lat}", latitude, 1)
	} else {
		request = strings.Replace(weatherByCityReq, "{city name}", currentUserLoc.City, 1)
	}
	response, err := http.Get(request)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	err = json.NewDecoder(response.Body).Decode(&apiResponse)

	if currentUserLoc.City == "" {
		currentUserLoc.City = apiResponse.Name //
	}

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	ourResp, _ := json.Marshal(apiResponse)

	_, err = w.Write(ourResp)
}

func getLocationByIp(r *http.Request) (string, string) {
	var lat, long string

	ip := r.RemoteAddr
	results, err := db.Get_all(ip)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	lat = fmt.Sprintf("%f", results.Latitude)
	long = fmt.Sprintf("%f", results.Longitude)
	return lat, long
}

func setLocation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err := json.NewDecoder(r.Body).Decode(&currentUserLoc)

	if err != nil {
		fmt.Print(err.Error())
		w.WriteHeader(http.StatusLocked)
		os.Exit(1)
	}

	w.WriteHeader(http.StatusCreated)
}

func main() {
	var err error
	db, err = ip2location.OpenDB("./IP2LOCATION-LITE-DB5.BIN")

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	router := mux.NewRouter()
	router.HandleFunc("/weather", getWeather).Methods(http.MethodGet)
	router.HandleFunc("/setLocation", setLocation).Methods(http.MethodPost)
	http.Handle("/", router)

	_ = http.ListenAndServe(":8181", nil)
}
