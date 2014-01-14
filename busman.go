package main

import (
    "fmt"
    "net/http"
    "io/ioutil"
    "encoding/json"
    "time"
)


type Stop struct {
    V int                   `json:"__v,int"`
    Id string               `json:"_id"`
    City string             `json:"city"`
    Name string             `json:"name"`
    Sort_id int             `json:"sort_id,int"`
    Destinations [] Destination    `json:"destinations"`
}

type Destination struct {
    V int                   `json:"__v,int"`
    Id string               `json:"_id"`
    City string             `json:"city"`
    Name string            `json:"name"`
    Sort_id int             `json:"sort_id, int"`
    Destinations [] string    `json:"destinations,[]string"`
}


type Time struct {
    FromId string
    From string
    ToId string
    To string
    Date string
    //Arrival string
    Route int
    Busstop string
    Hash string
    Id string `json:"_id"`
    V int `json:"__v"`
}

func main() {

    busStopsURL := "http://rutebuss.no/stops"
    
    resp, err := http.Get(busStopsURL) 
    if err != nil {
        fmt.Print("Bus stops could not be downlaoded...", err)
    }
    
    body, err := ioutil.ReadAll(resp.Body)
    
    if err != nil {
        fmt.Println("Could not read body ", err)
    }


    stops := make([]Stop, 100)
    
    err = json.Unmarshal(body, &stops)
    
    if err != nil {
        fmt.Println(string(body))
        fmt.Println("Could not unmarshal body... ", err,  string(body[5]))
    }


    fmt.Println("From", stops[0].Name,"to", stops[1].Name)

    from := stops[0].Id
    to := stops[1].Id

    t := time.Now()
    timeString := t.UTC().Format(time.RFC3339)
    

    travelURL :=
        "http://rutebuss.no/departure?from="+from+"&to="+to+"&date="+timeString
    
    resp, err = http.Get(travelURL) 
    if err != nil {
        fmt.Print("Bus times could not be downlaoded...", err)
    }
    
    body, err = ioutil.ReadAll(resp.Body)
    
    if err != nil {
        fmt.Println("Could not read body ", err)
    }


    times := make([]Time, 10)
    
    err = json.Unmarshal(body, &times)
    
    if err != nil {
        fmt.Println(string(body))
        fmt.Println("Could not unmarshal body... ", err,  string(body[5]))
    }



    for _, t := range times {
        departure, err := time.Parse(time.RFC3339, t.Date)
        if err != nil {
            fmt.Println("Parsing of date went horrible... ",err)
        }
        untilDeparture := departure.Sub(time.Now()) 


        fmt.Println(untilDeparture.String())
    }

}


