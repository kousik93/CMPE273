package main

import (
	"github.com/drone/routes"
	"log"
  "encoding/json"
	"net/http"
	"io/ioutil"


)
type Food struct{
	Type string `json:"type"`
	Drink_alchohol string `json:"drink_alchohol"`
}

type Music struct{
	Spotify_user_id string `json:"spotify_user_id"`

}

type Movie struct{
	Tv_shows []string `json:"tv_shows"`
	Movies []string `json:"movies"`
}

type Flight struct{
	Seat string `json:"seat"`
}

type Travel struct{
		Flight struct{
		Seat string `json:"seat"`
	}`json:"flight"`
}

type User struct {
    Email string `json:"email"`
    Zip string `json:"zip"`
    Country string `json:"country"`
    Profession string `json:"profession"`
    Favorite_color string `json:"favorite_color"`
    Is_smoking string `json:"is_smoking"`
    Favorite_sport string `json:"favorite_sport"`

		Food struct{
			Type string `json:"type"`
			Drink_alcohol string `json:"drink_alcohol"`
		}`json:"food"`

		Music struct{
			Spotify_user_id string `json:"spotify_user_id"`

		}`json:"music"`

		Movie struct{
			Tv_shows []string `json:"tv_shows"`
			Movies []string `json:"movies"`
		}`json:"movie"`

	   Travel struct{
		 Flight struct{
				Seat string `json:"seat"`
			}`json:"flight"`

		}`json:"travel"`
 }

//Variables for interface{} asserrtion
var Profiles = make(map[string]User)
var random map[string]interface{}
var music map[string]interface{}
var moviea =make(map[string]interface{})
var food map[string]interface{}
var travel map[string]interface{}
var flight map[string]interface{}
var seat map[string]interface{}

func GetProfile(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		email := params.Get(":email")
		prof:=Profiles[email]
		a, err := json.Marshal(prof)
		if err != nil {
			log.Println("error:", err)
		}
		if prof.Email!="" {
			w.Write([]byte(a))
		}
		if prof.Email ==""{
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(""))
		}
}

func PutProfile(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		email := params.Get(":email")
		prof:=Profiles[email]
		if prof.Email !=""{
			body, err := ioutil.ReadAll(r.Body)
	 		if err != nil {
			 	panic(err)
	 		}
			json.Unmarshal([]byte(body), &random)
			for key, value := range random {

				if key=="favorite_sport"{
					prof.Favorite_sport=value.(string)
					Profiles[email]=prof
				}
				if key=="zip"{
					prof.Zip=value.(string)
					Profiles[email]=prof
				}
				if key=="country"{
					prof.Country=value.(string)
					Profiles[email]=prof
				}
				if key=="profession"{
					prof.Profession=value.(string)
					Profiles[email]=prof
				}
				if key=="favorite_color"{
					prof.Favorite_color=value.(string)
					Profiles[email]=prof
				}
				if key=="is_smoking"{
					prof.Is_smoking=value.(string)
					Profiles[email]=prof
				}
				if key=="music"{
					music = value.(map[string]interface{})
						for keya, valuea := range music {
							if keya=="spotify_user_id"{
								prof.Music.Spotify_user_id=valuea.(string)
								Profiles[email]=prof
							}
						}
					}
				if key=="movie"{
					moviea = value.(map[string]interface{})
					for keya, valuea := range moviea {
							switch vv := valuea.(type) {
								case []interface{}:
										if keya=="movies"{
											prof.Movie.Movies=nil
											for i, u := range vv {
													if i==0{}
													 prof.Movie.Movies = append(prof.Movie.Movies, u.(string))
											}
										}
										if keya=="tv_shows"{
											prof.Movie.Tv_shows=nil
											for i, u := range vv {
													if i==0{}
													prof.Movie.Tv_shows = append(prof.Movie.Tv_shows, u.(string))
											}
										}
								}
						}
						Profiles[email]=prof
				}
				if key=="food"{
					food = value.(map[string]interface{})
						for keya, valuea := range food {
							if keya=="drink_alcohol"{
								prof.Food.Drink_alcohol=valuea.(string)
								Profiles[email]=prof
							}
							if keya=="type"{
								prof.Food.Type=valuea.(string)
								Profiles[email]=prof
							}
						}
						Profiles[email]=prof
					}
					if key=="travel"{
						travel = value.(map[string]interface{})
							for keya, valuea := range travel {
								if keya==""{}
								flight=valuea.(map[string]interface{})
								for keyaa, valueaa := range travel {
									if keyaa==""{}
									seat=valueaa.(map[string]interface{})
									for keyaaa,valueaaa:=range seat{
										if keyaaa==""{}
										prof.Travel.Flight.Seat=valueaaa.(string)
									}
								}
							}
							Profiles[email]=prof
						}
			}
			w.WriteHeader(http.StatusNoContent)
			w.Write([]byte(""))
		}
		if prof.Email==""{
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(""))
		}
}

func DelProfile(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	email := params.Get(":email")
	delete(Profiles, email)
	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte(""))
}

func PostProfile(rw http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var prof User
	err := decoder.Decode(&prof)
	if err != nil {
		panic(err)
	}
	Profiles[prof.Email]=prof
	rw.WriteHeader(http.StatusCreated)
	rw.Write([]byte(""))
}

func main() {
			mux := routes.New()
			mux.Get("/profile/:email", GetProfile)
			mux.Del("/profile/:email", DelProfile)
			mux.Put("/profile/:email", PutProfile)
			mux.Post("/profile",PostProfile)
			http.Handle("/", mux)
		 	log.Println("Listening...")
		 	http.ListenAndServe(":3000", nil)
}
