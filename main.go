package main

import (
	"fmt"
	"net/http"
	"log"
	"io/ioutil"
	"encoding/json"
	"html/template"
	"time"
)

type L struct {
	Cod     string  `json:"cod"`
	Message float64 `json:"message"`
	Cnt     int     `json:"cnt"`
	List    []struct {
		Dt   int `json:"dt"`
		Main struct {
			Temp      float64 `json:"temp"`
			TempMin   float64 `json:"temp_min"`
			TempMax   float64 `json:"temp_max"`
			Pressure  float64 `json:"pressure"`
			SeaLevel  float64 `json:"sea_level"`
			GrndLevel float64 `json:"grnd_level"`
			Humidity  int     `json:"humidity"`
			TempKf    float64 `json:"temp_kf"`
			Celsius	   float64
			Comment		string
			Time 		string
		} `json:"main"`
		Weather []struct {
			ID          int    `json:"id"`
			Main        string `json:"main"`
			Description string `json:"description"`
			Icon        string `json:"icon"`
		} `json:"weather"`
		Clouds struct {
			All int `json:"all"`
		} `json:"clouds"`
		Wind struct {
			Speed float64 `json:"speed"`
			Deg   float64 `json:"deg"`
		} `json:"wind"`
		Rain struct {
			ThreeH float64 `json:"3h"`
		} `json:"rain"`
		Sys struct {
			Pod string `json:"pod"`
		} `json:"sys"`
		DtTxt string `json:"dt_txt"`
	} `json:"list"`
	City struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Coord struct {
			Lat float64 `json:"lat"`
			Lon float64 `json:"lon"`
		} `json:"coord"`
		Country    string `json:"country"`
		Population int    `json:"population"`
	} `json:"city"`
}

func main(){
	http.HandleFunc("/", welcome)
	http.HandleFunc("/forecast", showForecast)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.ListenAndServe(":8080", nil)
}


var weather L

func showForecast(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	page, err := http.Get("http://api.openweathermap.org/data/2.5/forecast?q="+name+"%2Cno&units=imperial&appid=0f4ee0e05eebd5458c5e59798b05a962")
	if err != nil {
		log.Fatal(err)
	}
	jSonInfo, err := ioutil.ReadAll(page.Body)
	page.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	er := json.Unmarshal(jSonInfo, &weather)
	if err != nil {
		fmt.Println("error:", er)
	}
	weather.convert()
	weather.dateTime()


	tmpl, err := template.ParseFiles("forecast.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, weather); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func welcome (w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("index.html")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
	}
}

//temp := weather.List[0].Main.Temp

func toCelsius(temp float64) float64{
	c := (temp-32)/1.8000
	return c
}

func  (w *L) convert(){
	for i := 0; i < len(w.List); i++{
		c := toCelsius(w.List[i].Main.Temp)
		w.List[i].Main.Celsius = c

		/**
		mindre enn 0: Welcome to Norway
		0-5 : You will need a thick jacket
		5-15: Bring a jacket
		15-20: It's summer time!
		20-40: Don't get burned
		50 -> : Good luck Chuck
		 */

		if c < 0{
			w.List[i].Main.Comment = "Welcome to Norway"
		} else if c <= 5 {
			w.List[i].Main.Comment = "You will need a thick jacket."
		}
	}
}

func (w *L) dateTime() {
	var oldDay, newDay string;
	for i := 0; i < len(w.List); i++ {
		layout := "2006-01-02 15:04:05"
		t, _ := time.Parse(layout, w.List[i].DtTxt)
		newDay = t.Format("Monday")
		if newDay != oldDay {
			oldDay = newDay
			w.List[i].DtTxt = t.Format("Monday 2. Jan 15:04")
		} else {
			w.List[i].DtTxt = t.Format("15:04")
		}
	}
}