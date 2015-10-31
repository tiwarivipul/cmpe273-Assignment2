package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"net/http"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
    "github.com/julienschmidt/httprouter"
)


type JsonResult struct {
	
	Results []struct {
		
		Geometry     struct {
			   Location struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			   } `json:"location"`
			
		} `json:"geometry"`
		
	} `json:"results"`
}
	

type JsonName struct {
	Id         bson.ObjectId `json:"id" bson:"_id"`
	Name       string        `json:"name"`
	Address    string        `json:"address"`
	City       string        `json:"city"`
	State      string        `json:"state"`
	Zip    string            `json:"zip"`
	CoordinateField struct {
		Lat float64 `json:"lat"`
		Lng float64 `json:"lng"`
	} `json:"coordinate_field"`
}

func Getlocations(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {

	session, err := mgo.Dial("mongodb://Vipul01:1234@ds029804.mongolab.com:29804/cmpe273")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	c := session.DB("cmpe273").C("people")
	id := p.ByName("name")
	oid := bson.ObjectIdHex(id)
	var result JsonName
	c.FindId(oid).One(&result)
    if err != nil {
		log.Fatal(err)
	}


    fmt.Println("SeqNo:", result.Id.String())
	oid = bson.ObjectId(result.Id)

	b2, err := json.Marshal(result)
	if err != nil {
	}

	fmt.Fprintf(rw, string(b2))
	fmt.Println("Method : " + req.Method)
}

func Postlocations(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {

	var crud JsonName
	args := json.NewDecoder(req.Body)
	err := args.Decode(&crud)
	ApiCall := "http://maps.google.com/maps/api/geocode/json?address="  // Calling the API
	LocationCall := crud.Address + " " + crud.City + " " + crud.State   // Getting the location details
	LocationCall = strings.Replace(LocationCall, " ", "+", -1)
	FinalCall := "&sensor=false"
	Link := ApiCall + LocationCall + FinalCall
	fmt.Println("Published URL: " + Link)
	res, err := http.Get(Link)
	if err != nil {
		log.Fatal(err)
	}
	robots, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	var jsonresult JsonResult
	err = json.Unmarshal(robots, &jsonresult)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(jsonresult.Results[0].Geometry.Location.Lat)
	fmt.Println(jsonresult.Results[0].Geometry.Location.Lng)

	crud.Id = bson.NewObjectId()
	
    crud.CoordinateField.Lat = jsonresult.Results[0].Geometry.Location.Lat
	crud.CoordinateField.Lng = jsonresult.Results[0].Geometry.Location.Lng

	fmt.Println("Name " + crud.Name)
	fmt.Println("\nAddress: " + crud.Address)
	fmt.Println("\nCity:  " + crud.City)
	fmt.Println("\nState: " + crud.State)
	fmt.Println("\nLat and long : ")

	fmt.Println(crud.CoordinateField.Lat)
	fmt.Println(crud.CoordinateField.Lng)
	fmt.Println("\n Zip: " + crud.Zip)
		

	if err != nil {
	}

	session, err := mgo.Dial("mongodb://Vipul01:1234@ds029804.mongolab.com:29804/cmpe273")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	c := session.DB("cmpe273").C("people")

	err = c.Insert(crud)
	if err != nil {
		log.Fatal(err)
	}

	result := JsonName{}
	id := crud.Id.Hex()
	oid := bson.ObjectIdHex(id)
	c.FindId(oid).One(&result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Updated Name:", result.Name)
	fmt.Println("Address:", result.Address)
	fmt.Println("SeqNo:", result.Id.String())
	oid = bson.ObjectId(result.Id)

	b2, err := json.Marshal(result)
	if err != nil {
	}

	fmt.Fprintf(rw, string(b2))
}

func PutLocations(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
	var crud JsonName
	args := json.NewDecoder(req.Body)
	err := args.Decode(&crud)

	session, err := mgo.Dial("mongodb://Vipul01:1234@ds029804.mongolab.com:29804/cmpe273")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	c := session.DB("cmpe273").C("people")

	id := p.ByName("name")
	oid := bson.ObjectIdHex(id)
	var result JsonName
	c.FindId(oid).One(&result)
	
// Printing all the values......

	fmt.Println("Name:", result.Name)
	fmt.Println("Address:", result.Address)
	fmt.Println("City:", result.City)
	fmt.Println("State:", result.State)
	fmt.Println("Zip:", result.Zip)
	fmt.Println("latitude:", result.CoordinateField.Lat)
	fmt.Println("longitude:", result.CoordinateField.Lng)

	if crud.Name != "" {
		result.Name = crud.Name
	}
	if crud.Address != "" {
		result.Address = crud.Address
	}
	if crud.City != "" {
		result.City = crud.City
	}
	if crud.State != "" {
		result.State = crud.State
	}
	if crud.Zip != "" {
		result.Zip = crud.Zip
	}

	ApiCall := "http://maps.google.com/maps/api/geocode/json?address="
	LocationCall := result.Address + " " + result.City + " " + result.State
	LocationCall = strings.Replace(LocationCall, " ", "+", -1)

	FinalCall := "&sensor=false"
	Link := ApiCall + LocationCall + FinalCall
	fmt.Println("Published URL: " + Link)

	res, err := http.Get(Link)
	if err != nil {
		log.Fatal(err)
	}
	robots, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	var jsonresult JsonResult
	err = json.Unmarshal(robots, &jsonresult)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("New latitude longiteude")

	fmt.Println(jsonresult.Results[0].Geometry.Location.Lat)
	fmt.Println(jsonresult.Results[0].Geometry.Location.Lng)

	result.CoordinateField.Lat = jsonresult.Results[0].Geometry.Location.Lat
	result.CoordinateField.Lng = jsonresult.Results[0].Geometry.Location.Lng

	c.UpdateId(oid, result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\n\nNew Name:", result.Name)
	fmt.Println("Updated Address:", result.Address)
	fmt.Println("Updated City:", result.City)
	fmt.Println("Updated State:", result.State)
	fmt.Println("Updated Zip:", result.Zip)
	fmt.Println("Updated latitude:", result.CoordinateField.Lat)
	fmt.Println("Updated longitude:", result.CoordinateField.Lng)

	fmt.Println("SeqNo:", result.Id.String())
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