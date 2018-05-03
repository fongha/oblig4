package main

import (
	"fmt"
	"net/http"
	"log"
	"io/ioutil"
	"encoding/json"
	"html/template"
	"strings"
	"strconv"
)
var tilbakeMld string
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
	TilbakeMeldinger struct {
		TilbakeMelding1 TilbakeMelding
		TilbakeMeldingene []TilbakeMelding
	}
}
type TilbakeMelding struct {
	TilbakeMld string
}
func (f *L) addMld(melding TilbakeMelding) []TilbakeMelding{
	f.TilbakeMeldinger.TilbakeMeldingene = append(f.TilbakeMeldinger.TilbakeMeldingene, melding)
	return f.TilbakeMeldinger.TilbakeMeldingene
}
type L2 struct {
	List []struct{
		Main struct {
			Temp      float64 `json:"temp"`
		}
	}
}

func main(){
	http.HandleFunc("/", welcome)
	http.HandleFunc("/forecast", showForecast)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.ListenAndServe(":8080", nil)
}
var temperatur L2
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
	er2 := json.Unmarshal(jSonInfo, &temperatur)
	if err != nil {
		fmt.Println("error:", er2)
	}
	hentData()

	fmt.Printf("%+v, %s", temperatur, "hallo")

	tmpl, err := template.ParseFiles("forecast.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	hentData()
	if err := tmpl.Execute(w, weather); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func welcome (w http.ResponseWriter, r *http.Request) {
	//fmt.Println("method:", r.Method) //get request method
	if r.Method == "GET" {
		t, _ := template.ParseFiles("index.html")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
	}
}

func hentData() {
	//fmt.Println(len(temperatur.List))
	for i := 0; i < len(temperatur.List); i++{
		a := fmt.Sprintf("%+v", temperatur.List[i])
		b := strings.Replace(a, "{", "", -1)
		c := strings.Replace(b, "}", "", -1)
		d := strings.Replace(c, "Main:Temp:", "", -1)
		e, err := strconv.ParseFloat(d, 64)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(len(temperatur.List))
		fmt.Println(1)
		if e >= 0 {
			fmt.Println("Det er mer enn 0 grader")
			tilbakeMld = "Det er mer enn 0 grader"
		}else {
			tilbakeMld = "blaaaorgn"
		}
		Test:= temperatur.List
		pushData(tilbakeMld)
		//fmt.Println(tilbakeMld)
		fmt.Print(Test)
	}
	fmt.Println(2)
}

func pushData(mld string) {
	tilbakeMeldingen := L{}
	tilbakeMeldingen.TilbakeMeldinger.TilbakeMelding1 = TilbakeMelding{tilbakeMld}
	tilbakeMeldingen.TilbakeMeldinger.TilbakeMeldingene = []TilbakeMelding{}
}
