package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"text/template"
)

type Kiosk struct {
	Name string `json:"name"`
	Url  string `json:"url"`
	Ip   string `json:"ip"`
}

type Kiosks struct {
	Kiosks []Kiosk `json:"kiosks"`
}

// function that takes in kiosk and url and prints them out
func update(configFile string, newkiosk string, url string) {
	fmt.Println("updating")
	fmt.Println("kiosk:", newkiosk)
	fmt.Println("url:", url)

	//read kiosk_config.json into variable
	fmt.Println("configFile:", configFile)

	// read the kiosk_config.json file
	kioskConfig, err := os.Open(configFile)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("successfully opened kiosk_config.json")
	defer kioskConfig.Close()

	byteValue, _ := ioutil.ReadAll(kioskConfig)

	var kiosks Kiosks

	json.Unmarshal(byteValue, &kiosks)

	// update the url for newkiosk
	for i := 0; i < len(kiosks.Kiosks); i++ {
		if kiosks.Kiosks[i].Name == newkiosk {
			kiosks.Kiosks[i].Url = url
		}
	}

	// convert it back to json and write to file
	json, _ := json.MarshalIndent(kiosks, "", "  ")
	fmt.Println(string(json))
	_ = ioutil.WriteFile(configFile, json, 0644)

}

func getKioskNames(configFile string) []string {
	//read kiosk_config.json into variable
	// fmt.Println("configFile:", configFile)

	// read the kiosk_config.json file
	kioskConfig, err := os.Open(configFile)
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println("successfully opened kiosk_config.json")
	defer kioskConfig.Close()

	byteValue, _ := ioutil.ReadAll(kioskConfig)

	var kiosks Kiosks

	json.Unmarshal(byteValue, &kiosks)

	// update the url for newkiosk
	var kioskNames []string
	for i := 0; i < len(kiosks.Kiosks); i++ {
		kioskNames = append(kioskNames, kiosks.Kiosks[i].Name)
	}

	return kioskNames
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	configFile := os.Getenv("KIOSK_CONFIG")

	if r.Method == "GET" {
		fmt.Println(r)
		t, _ := template.ParseFiles("./websource/update.html")

		update(configFile, r.FormValue("kiosk"), r.FormValue("url"))

		t.Execute(w, nil)

	} else {
		http.Redirect(w, r, "/", 301)
	}
}

func getOptions(kioskNames []string) string {
	options := ""
	for i := 0; i < len(kioskNames); i++ {
		//append to options
		options = options + "<option value=\"" + kioskNames[i] + "\">" + kioskNames[i] + "</option>"
	}

	return options
}

type indexData struct {
	PageTitle string
	Options   string
}

func viewHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		t, _ := template.ParseFiles("./websource/index.html")

		kioskNames := getKioskNames(os.Getenv("KIOSK_CONFIG"))
		data := indexData{
			PageTitle: "Kiosk Controller",
			Options:   getOptions(kioskNames),
		}
		t.Execute(w, data)
	} else {
		r.ParseForm()
	}
}

func kiosksHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("displaying kiosks")
	http.ServeFile(w, r, os.Getenv("KIOSK_CONFIG"))
}

func main() {

	// TODO: Check for KIOSK_CONFIG environment variable
	// If its not set it to kiosk_config.json you get a nil pointer error

	http.HandleFunc("/", viewHandler)
	http.HandleFunc("/update", updateHandler)
	http.HandleFunc("/kiosks", kiosksHandler)
	fmt.Println("Starting server at port 8090")
	// http.ListenAndServe(":8090", nil)
	err := http.ListenAndServe(":8090", nil) // setting listening port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
