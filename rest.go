package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/julienschmidt/httprouter"
)

type Person struct {
	Name  string
	Phone string
}
type MyJsonResult struct {
	Results []struct {
		AddressComponents []struct {
			LongName  string   `json:"long_name"`
			ShortName string   `json:"short_name"`
			Types     []string `json:"types"`
		} `json:"address_components"`
		FormattedAddress string `json:"formatted_address"`
		Geometry         struct {
			Location struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"location"`
			LocationType string `json:"location_type"`
			Viewport     struct {
				Northeast struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"northeast"`
				Southwest struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"southwest"`
			} `json:"viewport"`
		} `json:"geometry"`
		PlaceID string   `json:"place_id"`
		Types   []string `json:"types"`
	} `json:"results"`
	Status string `json:"status"`
}

type MyJsonName struct {
	Id         bson.ObjectId `json:"id" bson:"_id"`
	Name       string        `json:"name"`
	Address    string        `json:"address"`
	City       string        `json:"city"`
	State      string        `json:"state"`
	Zip        string        `json:"zip"`
	Coordinate struct {
		Lat float64 `json:"lat"`
		Lng float64 `json:"lng"`
	} `json:"coordinate"`
}

func Getlocations(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {

	session, err := mgo.Dial("mongodb://Vipul01:1234@ds029804.mongolab.com:29804/cmpe273")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("cmpe273").C("people")
	id := p.ByName("name")
	oid := bson.ObjectIdHex(id)
	var result MyJsonName
	c.FindId(oid).One(&result)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Searched Name:", result.Name)
	fmt.Println("Searched Address:", result.Address)
	fmt.Println("Searched City:", result.City)
	fmt.Println("Searched State:", result.State)
	fmt.Println("Searched Zip:", result.Zip)
	fmt.Println("Searched latitude:", result.Coordinate.Lat)
	fmt.Println("Searched longitude:", result.Coordinate.Lng)

	fmt.Println("Id2:", result.Id.String())
	oid = bson.ObjectId(result.Id)

	b2, err := json.Marshal(result)
	if err != nil {
	}

	fmt.Fprintf(rw, string(b2))
	fmt.Println("Method Name: " + req.Method)
}

func Postlocations(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {

	var myjson3 MyJsonName
	s3 := json.NewDecoder(req.Body)
	err := s3.Decode(&myjson3)
	StartQuery := "http://maps.google.com/maps/api/geocode/json?address="
	WhereQuery := myjson3.Address + " " + myjson3.City + " " + myjson3.State
	WhereQuery = strings.Replace(WhereQuery, " ", "+", -1)
	EndQuery := "&sensor=false"
	Url1 := StartQuery + WhereQuery + EndQuery
	fmt.Println("Published URL: " + Url1)
	res, err := http.Get(Url1)
	if err != nil {
		log.Fatal(err)
	}
	robots, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	var myjsonresult1 MyJsonResult
	err = json.Unmarshal(robots, &myjsonresult1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(myjsonresult1.Results[0].Geometry.Location.Lat)
	fmt.Println(myjsonresult1.Results[0].Geometry.Location.Lng)

	myjson3.Id = bson.NewObjectId()
	fmt.Println("Check1")
	fmt.Println(string(myjson3.Id))
	fmt.Println(myjson3.Id.Hex())
	fmt.Println(myjson3.Id.String())
	fmt.Println(myjson3.Id.Pid())

	myjson3.Coordinate.Lat = myjsonresult1.Results[0].Geometry.Location.Lat
	myjson3.Coordinate.Lng = myjsonresult1.Results[0].Geometry.Location.Lng

	fmt.Println("Name " + myjson3.Name)
	fmt.Println("\nAddress: " + myjson3.Address)
	fmt.Println("\nCity:  " + myjson3.City)
	fmt.Println("\nState: " + myjson3.State)
	fmt.Println("\nLat and long : ")

	fmt.Println(myjson3.Coordinate.Lat)
	fmt.Println(myjson3.Coordinate.Lng)

	if err != nil {
	}

	session, err := mgo.Dial("mongodb://Vipul01:1234@ds029804.mongolab.com:29804/cmpe273")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("cmpe273").C("people")

	err = c.Insert(myjson3)
	if err != nil {
		log.Fatal(err)
	}

	result := MyJsonName{}
	fmt.Println()
	id := myjson3.Id.Hex()
	oid := bson.ObjectIdHex(id)
	c.FindId(oid).One(&result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("New Name:", result.Name)
	fmt.Println("Address:", result.Address)
	fmt.Println("Id2:", result.Id.String())
	oid = bson.ObjectId(result.Id)

	b2, err := json.Marshal(result)
	if err != nil {
	}

	fmt.Fprintf(rw, string(b2))
}

func PutLocations(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
	var myjson3 MyJsonName
	s3 := json.NewDecoder(req.Body)
	err := s3.Decode(&myjson3)

	session, err := mgo.Dial("mongodb://Vipul01:1234@ds029804.mongolab.com:29804/cmpe273")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("cmpe273").C("people")

	id := p.ByName("name")
	oid := bson.ObjectIdHex(id)
	var result MyJsonName
	c.FindId(oid).One(&result)
	fmt.Println("Old Name:", result.Name)
	fmt.Println("Old Address:", result.Address)
	fmt.Println("Old City:", result.City)
	fmt.Println("Old State:", result.State)
	fmt.Println("Old Zip:", result.Zip)
	fmt.Println("Old latitude:", result.Coordinate.Lat)
	fmt.Println("Old longitude:", result.Coordinate.Lng)

	if myjson3.Name != "" {
		result.Name = myjson3.Name
	}
	if myjson3.Address != "" {
		result.Address = myjson3.Address
	}
	if myjson3.City != "" {
		result.City = myjson3.City
	}
	if myjson3.State != "" {
		result.State = myjson3.State
	}
	if myjson3.Zip != "" {
		result.Zip = myjson3.Zip
	}

	StartQuery := "http://maps.google.com/maps/api/geocode/json?address="
	WhereQuery := result.Address + " " + result.City + " " + result.State
	WhereQuery = strings.Replace(WhereQuery, " ", "+", -1)

	EndQuery := "&sensor=false"
	Url1 := StartQuery + WhereQuery + EndQuery
	fmt.Println("Published URL: " + Url1)

	res, err := http.Get(Url1)
	if err != nil {
		log.Fatal(err)
	}
	robots, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	var myjsonresult1 MyJsonResult
	err = json.Unmarshal(robots, &myjsonresult1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("New latitude longiteude")

	fmt.Println(myjsonresult1.Results[0].Geometry.Location.Lat)
	fmt.Println(myjsonresult1.Results[0].Geometry.Location.Lng)

	result.Coordinate.Lat = myjsonresult1.Results[0].Geometry.Location.Lat
	result.Coordinate.Lng = myjsonresult1.Results[0].Geometry.Location.Lng

	c.UpdateId(oid, result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\n\nNew Name:", result.Name)
	fmt.Println("New Address:", result.Address)
	fmt.Println("New City:", result.City)
	fmt.Println("New State:", result.State)
	fmt.Println("New Zip:", result.Zip)
	fmt.Println("New latitude:", result.Coordinate.Lat)
	fmt.Println("New longitude:", result.Coordinate.Lng)

	fmt.Println("Id2:", result.Id.String())
	oid = bson.ObjectId(result.Id)

	b2, err := json.Marshal(result)
	if err != nil {
	}

	fmt.Fprintf(rw, string(b2))
	fmt.Println("Method Name: " + req.Method)
}

func DeleteLocations(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {

	fmt.Fprintf(rw, "Deleting Id, %s!\n", p.ByName("name"))

	session, err := mgo.Dial("mongodb://Vipul01:1234@ds029804.mongolab.com:29804/cmpe273")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("cmpe273").C("people")
	id := p.ByName("name")
	oid := bson.ObjectIdHex(id)
	c.RemoveId(oid)
	fmt.Fprintf(rw, "Deleted, %s!\n", p.ByName("name"))
}

func main() {

	mux := httprouter.New()
	mux.GET("/locations/:name", Getlocations)
	mux.POST("/locations", Postlocations)
	mux.PUT("/locations/:name", PutLocations)
	mux.DELETE("/locations/:name", DeleteLocations)
	server := http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: mux,
	}
	server.ListenAndServe()
}