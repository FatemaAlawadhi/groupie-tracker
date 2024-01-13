package Handlers

import (
	"encoding/json"
	"groupie-tracker/Error"
	"html/template"
	"net/http"
)

type Artist struct {
	ID      int      `json:"id"`
	Image   string   `json:"image"`
	Name    string   `json:"name"`
	Members []string `json:"members"`
}

func HandleStarsPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		response, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
		if err != nil {
			Error.RenderErrorPage(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		defer response.Body.Close()

		// Parse the JSON response
		var artists []Artist
		err = json.NewDecoder(response.Body).Decode(&artists)
		if err != nil {
			Error.RenderErrorPage(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		var artistsFiltered []Artist
		for _, artist := range artists {
			if artist.ID != 21 {
				artistsFiltered = append(artistsFiltered, artist)
			}
		}
		
		//Parse and Excute
		tmpl, err := template.ParseFiles("WebPages/Stars/Stars.html", "WebPages/Stars/Card.html")
		if err != nil {
			Error.RenderErrorPage(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		err = tmpl.Execute(w, artistsFiltered)
		if err != nil {
			Error.RenderErrorPage(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
	} else {
		Error.RenderErrorPage(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}
}
