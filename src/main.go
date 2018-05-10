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
	List    []struct {
		Dt   int `json:"dt"`
		Main struct {
			Temp      float64 `json:"temp"`
			Celsius	   float64
			Comment		string
			Time 		string
		} `json:"main"`
		Weather []struct {
			Main        string `json:"main"`
			Description string `json:"description"`
			Icon        string `json:"icon"`
		} `json:"weather"`
		Wind struct {
			Speed float64 `json:"speed"`
		} `json:"wind"`
		DtTxt string `json:"dt_txt"`
	} `json:"list"`
	City struct {
		Name  string `json:"name"`
		Country    string `json:"country"`
	} `json:"city"`
}

func main(){
	http.HandleFunc("/", Welcome)
	http.HandleFunc("/forecast", showForecast)
	http.HandleFunc("/error", errorPage)
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
	s := string(weather.Cod)
	if s != "200" {
		http.Redirect(w, r, "/error", http.StatusSeeOther)
	} else {

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
}

func errorPage (w http.ResponseWriter, r *http.Request) {
	e, _ := template.ParseFiles("error.html")
	e.Execute (w, nil)
}


func Welcome (w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("index.html")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
	}
}


func toCelsius(temp float64) float64{
	c := (temp-32)/1.8000
	return c
}

func  (w *L) convert(){
	var weatherMain string

	for i := 0; i < len(w.List); i++{

		c := toCelsius(w.List[i].Main.Temp)
		w.List[i].Main.Celsius = c
		weatherMain = w.List[i].Weather[0].Main

		if weatherMain == "Snow" &&  c < 0 {
			w.List[i].Main.Comment = "Better bring some warm clothes!"
		} else if ( weatherMain == "Clouds" || weatherMain  == "Clear" ) &&  c < 0 {
			w.List[i].Main.Comment = "It's chilly outside, but the weather is nice!"
		} else if weatherMain == "Rain" && c < 0 {
			w.List[i].Main.Comment = "It's cold and it's raining. Bring something warm and waterproof."
		} else if weatherMain == "Snow" && c >= 0 && c < 5 {
			w.List[i].Main.Comment = "It could be slippery outside. Heads up!"
		} else if ( weatherMain == "Clear" || weatherMain == "Clouds" ) && c >= 0 && c < 5 {
			w.List[i].Main.Comment = "Wear a thick jacket."
		} else if weatherMain == "Rain" && c >= 0 && c < 5 {
			w.List[i].Main.Comment = "Bring a raincoat and wear something warm under."
		} else if weatherMain == "Snow" && c >= 5 && c < 10 {
			w.List[i].Main.Comment = "It'll be slushy outside. Bring some waterproof protection for you feet."
		} else if ( weatherMain == "Clear" || weatherMain == "Clouds" ) && c >= 5 && c < 10 {
			w.List[i].Main.Comment = "Bring a jacket of your choosing."
		} else if weatherMain == "Rain" && c >= 5 && c < 10 {
			w.List[i].Main.Comment = "You should bring a raincoat or an umbrella today."
		} else if ( weatherMain == "Clear" || weatherMain == "Clouds" ) && c >= 10 && c < 15 {
			w.List[i].Main.Comment = "You could wear a thin jacket or a warm sweater."
		} else if weatherMain == "Rain" && c >= 10 && c < 15 {
			w.List[i].Main.Comment = "Wear a raincoat or bring an umbrella."
		} else if ( weatherMain == "Clear" || weatherMain == "Clouds" ) && c >= 15 && c < 20 {
			w.List[i].Main.Comment = "It is warm outside today! Wear something casual."
		} else if weatherMain == "Rain" && c >= 15 && c < 20 {
			w.List[i].Main.Comment = "Bummer! It's warm but it is still raining. Bring an umbrella."
		} else if ( weatherMain == "Clear" || weatherMain == "Clouds" ) && c >= 20 && c < 25 {
			w.List[i].Main.Comment = "The weather is great! Perhaps you want to spend the day outside?"
		} else if weatherMain == "Rain" && c >= 20 && c < 25 {
			w.List[i].Main.Comment = "How unfortunate that it is raining on a hot day like this."
		} else if ( weatherMain == "Clear" || weatherMain == "Clouds" ) &&  c < 30 {
			w.List[i].Main.Comment = "This temperature is surreal! A swim would be nice on a day like this."
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
