package models

// constants represents string values that we use for request
type constants struct {
	WeatherApiKey        string
	WeatherByLocationReq string
	WeatherByCityReq     string
}

// response used for structuring fields from json response
type Response struct {
	Coordinates coordinates
	Main        weather
	Name        string
}

// weather is part of Response type
type weather struct {
	Temp       float32
	Feels_like float32
	Temp_min   float32
	Temp_max   float32
	Pressure   int
	Humidity   int
}

// coordinates is part of Response type
type coordinates struct {
	Lon float32
	Lat float32
}

// Misc containing data required in WeatherApp
type misc struct {
	Coords coordinates
	Consts constants
	City   string
}

var Misc = misc{
	Consts: constants{
		"77699bd83349c01e6271ac022ba11ded",
		"https://api.openweathermap.org/data/2.5/weather?lat={lat}&lon={lon}&appid=" + "77699bd83349c01e6271ac022ba11ded" + "&units=metric",
		"https://api.openweathermap.org/data/2.5/weather?q={city name}&appid=" + "77699bd83349c01e6271ac022ba11ded" + "&units=metric",
	},
}
