package Handlers

import (
	"encoding/json"
	"groupie-tracker/Error"
	"html/template"
	"net/http"
	"strings"
)

type MainData struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

type LocationData struct {
	Locations []string `json:"locations"`
	Dates     string   `json:"dates"`
}

type DatesData struct {
	Dates []string `json:"dates"`
}

type RelationsData struct {
	DatesLocations map[string][]string `json:"datesLocations"`
}

func fetchData(url string, target interface{}) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return err
	}

	err = json.NewDecoder(response.Body).Decode(&target)
	if err != nil {
		return err
	}

	return nil
}

func HandleStarDetailsPage(w http.ResponseWriter, r *http.Request) {
	// Fetch the main data
	id := strings.TrimPrefix(r.URL.Path, "/stardetails/")
	mainDataUrl := "https://groupietrackers.herokuapp.com/api/artists/" + id
	var mainData MainData
	if err := fetchData(mainDataUrl, &mainData); err != nil {
		Error.RenderErrorPage(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// Fetch the location data
	var locationData LocationData
	if err := fetchData(mainData.Locations, &locationData); err != nil {
		Error.RenderErrorPage(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// Fetch the dates data
	var datesData DatesData
	if err := fetchData(locationData.Dates, &datesData); err != nil {
		Error.RenderErrorPage(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// Fetch the relations data
	var relationsData RelationsData
	if err := fetchData(mainData.Relations, &relationsData); err != nil {
		Error.RenderErrorPage(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	// Combine the location and dates data based on relations
	locationDates := make(map[string][]string)
	for _, location := range locationData.Locations {
		locationDates[location] = relationsData.DatesLocations[location]
	}

	// Construct the final object with all the data
	finalData := struct {
		ID           int                 `json:"id"`
		Image        string              `json:"image"`
		Name         string              `json:"name"`
		Members      []string            `json:"members"`
		CreationDate int                 `json:"creationDate"`
		FirstAlbum   string              `json:"firstAlbum"`
		Locations    []string            `json:"locations"`
		Dates        map[string][]string `json:"dates"`
	}{
		ID:           mainData.ID,
		Image:        mainData.Image,
		Name:         mainData.Name,
		Members:      mainData.Members,
		CreationDate: mainData.CreationDate,
		FirstAlbum:   mainData.FirstAlbum,
		Locations:    locationData.Locations,
		Dates:        locationDates,
	}
	tmpl, err := template.ParseFiles("WebPages/StarDetails.html")
	if err != nil {
		Error.RenderErrorPage(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	err = tmpl.Execute(w, finalData)
	if err != nil {
		Error.RenderErrorPage(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
}
