package main

import (
	"os"
	"testing"
	"net/http"
	"io/ioutil"
	"log"
	"encoding/json"
	"fmt"

	"net"
)

//funksjon som ser om filer finnes
func fileCheck(fName string) bool {
	if _, err := os.Stat(fName); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

//test index.html filen. For å teste andre filer: endre filnavnet.
func TestIndex(t *testing.T) {
	index := fileCheck("templates/index.html") //fil i samme dir som kildekoden.
	if index != true {
		t.Error("File not found.\n")
	}else{
		fmt.Print("File OK.\n")
	}
}

//test om URL'en er valid.
func TestUrl(t *testing.T) {

	type L struct {
		Cod string `json:"cod"` //COD er statuskoden.
	}
	var sL L
	by := "Kristiansand"
	page, err1 := http.Get("http://api.openweathermap.org/data/2.5/forecast?q=" + by + "%2Cno&units=imperial&appid=0f4ee0e05eebd5458c5e59798b05a962")
	if (err1 != nil ) { //sjekker http pakkens egen errorhandling.
		t.Error("Unvalid URL.\n")
	}else{
		fmt.Print("Valid URL.\n")
	}

	jSonInfo, err := ioutil.ReadAll(page.Body)
	page.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	er := json.Unmarshal(jSonInfo, &sL)
	if err != nil {
		fmt.Println("error:", er)
	}
	sLi := string(sL.Cod) //sjekker at statuskoden etter request er 200 (OK).
	if (sLi != "200" || err1 != nil ) {
		t.Error("Status message: " + sLi)
	}else{
		fmt.Print("Status message: 200 OK.\n")
	}

}

func TestCity(t *testing.T) {

	type L struct {
		City struct {
			Name string `json:"name"` //navnet på byen JSON-dataen omhandler.
		} `json:"city"`
	}
	var sL L
	by := "Kristiansand"
	page, err1 := http.Get("http://api.openweathermap.org/data/2.5/forecast?q=" + by + "%2Cno&units=imperial&appid=0f4ee0e05eebd5458c5e59798b05a962")
	if (err1 != nil ) {
		t.Error("unvalid URL.\n") //http-pakkens egen errorhandling.
	}

	jSonInfo, err := ioutil.ReadAll(page.Body)
	page.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	er := json.Unmarshal(jSonInfo, &sL)
	if err != nil {
		fmt.Println("error:", er)
	}

	byNavn := string(sL.City.Name)
	if (byNavn != by) { //sammenligner input-byen med bynavnet i JSON-dataen.
		t.Error("Expected " + by + " got: " + byNavn)
	}else{
		fmt.Print("Got the correct city.\n")
	}

}

//TESTER AT PORTEN :8080 ER ÅPEN
func TestConn(t *testing.T) {
	message := "Testing port...\n"
	port := ":8080"
	go func() {
		conn, err := net.Dial("tcp", port) //Egen tråd som starter forbindelse med :8080
		if err != nil {										  //ikke skriv feil port her da kjører tråden 4ever...
			t.Fatal(err)
		}
		defer conn.Close()

		if _, err := fmt.Fprintf(conn, message); err != nil { //skriver messagen "message" til porten :8080
			t.Fatal(err)
		}
	}()

	l, err := net.Listen("tcp", port) //lytter på :8080
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			return
		}
		defer conn.Close()

		buf, err := ioutil.ReadAll(conn)
		if err != nil {
			t.Fatal(err)
		}

		fmt.Println(string(buf[:]))
		if msg := string(buf[:]); msg != message {
			t.Fatalf("Unexpected message:\nGot:\t\t%s\nExpected:\t%s\n", msg, message)
		}else{
			fmt.Print("THE GATES ARE OPEN! (atleast " + port +")\n")
		}
		return // Done
	}

}