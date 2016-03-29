package main

import (
	"github.com/drone/routes"
	"log"
  "encoding/json"
	"net/http"
	"github.com/mkilling/goejdb"
	"labix.org/v2/mgo/bson"
	"os"
	"io/ioutil"
)

/*
*
All structurs to handle json
*
*/
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
		Id string `json:"_id,omitempty"`
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



type DB struct {
	jb *goejdb.Ejdb
	coll *goejdb.EjColl
}

func dbinit() (*goejdb.Ejdb,*goejdb.EjColl){
	jb, err1 := goejdb.Open("userdetails", goejdb.JBOWRITER | goejdb.JBOCREAT | goejdb.JBOTRUNC)
	if err1 != nil {
			os.Exit(1)
	}
	coll, _ := jb.CreateColl("details",nil)
return jb,coll
}



/*
*
Variables for interface{} asserrtion
*
*/
var Profiles = make(map[string]User)
var random map[string]interface{}
var music map[string]interface{}
var moviea =make(map[string]interface{})
var food map[string]interface{}
var travel map[string]interface{}
var flight map[string]interface{}
var seat map[string]interface{}

var tmp User
var db DB

/*
*
Main REST methods
*
*/


func GetProfile(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		email := params.Get(":email")
		var searchstr=`{"email" : "`+email+`"}`
		res, _ := db.coll.Find(searchstr)
		log.Println("Length= ")
		log.Println(len(res))

		for _, bs := range res {
				bson.Unmarshal(bs, &tmp)
				log.Println(tmp)
		}

		bson.Unmarshal(res[0], &tmp)
		a, err := json.Marshal(tmp)
		if err != nil {
			log.Println("error:", err)
		}
		if tmp.Email!="" {
			log.Println("prof in get")
			log.Println(tmp)
			w.Write([]byte(a))
		}
		if tmp.Email ==""{
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(""))
		}
}

func PutProfile(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		email := params.Get(":email")

		var prof User
		var searchstr=`{"email" : "`+email+`"}`
		res, _ := db.coll.Find(searchstr)
		bson.Unmarshal(res[0], &prof)
		//prof1:=prof
		if prof.Email !=""{
			body, err := ioutil.ReadAll(r.Body)
	 		if err != nil {
			 	panic(err)
	 		}
			json.Unmarshal([]byte(body), &random)
			for key, value := range random {

				if key=="favorite_sport"{
					prof.Favorite_sport=value.(string)
					//Profiles[email]=prof
				}
				if key=="zip"{
					prof.Zip=value.(string)
					//Profiles[email]=prof
				}
				if key=="country"{
					prof.Country=value.(string)
					//Profiles[email]=prof
				}
				if key=="profession"{
					prof.Profession=value.(string)
					//Profiles[email]=prof
				}
				if key=="favorite_color"{
					prof.Favorite_color=value.(string)
					//Profiles[email]=prof
				}
				if key=="is_smoking"{
					prof.Is_smoking=value.(string)
					//Profiles[email]=prof
				}
				if key=="music"{
					music = value.(map[string]interface{})
						for keya, valuea := range music {
							if keya=="spotify_user_id"{
								prof.Music.Spotify_user_id=valuea.(string)
								//Profiles[email]=prof
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
						//Profiles[email]=prof
				}
				if key=="food"{
					food = value.(map[string]interface{})
						for keya, valuea := range food {
							if keya=="drink_alcohol"{
								prof.Food.Drink_alcohol=valuea.(string)
								//Profiles[email]=prof
							}
							if keya=="type"{
								prof.Food.Type=valuea.(string)
							//	Profiles[email]=prof
							}
						}
						//Profiles[email]=prof
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
							//Profiles[email]=prof
						}
			}
			db.coll.Sync()
			db.coll.BeginTransaction()
			log.Println("In Put")
			deleted:=db.coll.RmBson(prof.Id)
			bsrecput, _ := bson.Marshal(prof)
			log.Println("deleted is ")
			log.Println(deleted)
			db.coll.SaveBson(bsrecput)
			// a,_:=json.Marshal(prof)
			// c,_:=json.Marshal(prof1)
			// log.Println(prof)
			// db.coll.Update(string(c),string(a))
			db.coll.CommitTransaction()
			db.coll.Sync()
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
	var searchstr=`{"email" : "`+email+`"}`
	res, _ := db.coll.Find(searchstr)
	bson.Unmarshal(res[0], &tmp)

	db.coll.Sync()
	db.coll.BeginTransaction()
	log.Println(goejdb.IsValidOidStr(tmp.Id))
	log.Println(bson.IsObjectIdHex(tmp.Id))
	db.coll.RmBson(tmp.Id)
	db.coll.CommitTransaction()
	db.coll.Sync()

	log.Println("in Del")
	log.Println(tmp)

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
	prof.Id=bson.ObjectId.Hex(bson.NewObjectId())
	bsrec, _ := bson.Marshal(prof)
	db.coll.Sync()
	db.coll.BeginTransaction()
	db.coll.SaveBson(bsrec)
	db.coll.CommitTransaction()
	db.coll.Sync()
	log.Printf("\nSaved "+prof.Email)

	rw.WriteHeader(http.StatusCreated)
	rw.Write([]byte(""))
}


//
// func PostProfileRpc(inp string) {
// 	decoder := json.NewDecoder(inp)
// 	var prof User
// 	err := decoder.Decode(&prof)
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	prof.Id=bson.ObjectId.Hex(bson.NewObjectId())
// 	bsrec, _ := bson.Marshal(prof)
// 	db.coll.Sync()
// 	db.coll.BeginTransaction()
// 	db.coll.SaveBson(bsrec)
// 	db.coll.CommitTransaction()
// 	db.coll.Sync()
// 	log.Printf("\nSaved "+prof.Email)
// }
//

func main() {
			db.jb,db.coll=dbinit()
			mux := routes.New()
			mux.Get("/profile/:email", GetProfile)
			mux.Del("/profile/:email", DelProfile)
			mux.Put("/profile/:email", PutProfile)
			mux.Post("/profile",PostProfile)
			http.Handle("/", mux)
		 	log.Println("Listening...")
		 	http.ListenAndServe(":3000", nil)
}
