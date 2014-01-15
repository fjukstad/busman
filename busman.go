package main

import (
    "fmt"
    "net/http"
    "io/ioutil"
    "encoding/json"
    "time"
    "flag" 
)


type Stop struct {
    V int                   `json:"__v,int"`
    Id string               `json:"_id"`
    City string             
    Name string             
    Sort_id int             
    Destinations [] Destination    
}

type Destination struct {
    V int                   `json:"__v,int"`
    Id string               `json:"_id"`
    City string            
    Name string            
    Sort_id int             
    Destinations [] string    
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

    fromString := flag.String("from", "UiT",
                        "The bus stop you are traveling from")
    toString := flag.String("to", "Sentrum",
                        "The bus stop you are traveling to") 

    flag.Parse() 
                        
    busStopsURL := "http://rutebuss.no/stops"
    
    resp, err := http.Get(busStopsURL) 
    if err != nil {
        fmt.Print("Bus stops could not be downlaoded...", err)
        return
    }
    
    body, err := ioutil.ReadAll(resp.Body)
    
    if err != nil {
        fmt.Println("Could not read body ", err)
        return
    }


    stops := make([]Stop, 100)
    
    err = json.Unmarshal(body, &stops)
    
    if err != nil {
        fmt.Print("Bus route ",*fromString,"-",*toString,"does not exist")
        return
    }

    
    fmt.Println("From", *fromString ,"to", *toString)

    var from, to string

    for _, stop := range stops { 
        if stop.Name == *fromString {
            from = stop.Id
        }
        if stop.Name == *toString {
            to = stop.Id
        } 
    } 

    if from == "" && to == "" {
        fmt.Println("Bus route ",from,"-",to,"does not exist")
        return
    }

    t := time.Now()
    timeString := t.UTC().Format(time.RFC3339)
    

    travelURL :=
        "http://rutebuss.no/departure?from="+from+"&to="+to+"&date="+timeString
    
    resp, err = http.Get(travelURL) 
    if err != nil {
        fmt.Print("Bus times could not be downlaoded...", err)
        return
    }
    
    body, err = ioutil.ReadAll(resp.Body)
    
    if err != nil {
        fmt.Println("Could not read body ", err)
        return
    }


    times := make([]Time, 10)
    
    err = json.Unmarshal(body, &times)
    
    if err != nil {
        fmt.Println("Bus route ",*fromString,"-",*toString,"does not exist")
        return
    }

    if len(times) == 0 {
        fmt.Println("Bus route ",*fromString,"-",*toString,"does not exist")
    }

    for _, t := range times {
        departure, err := time.Parse(time.RFC3339, t.Date)
        if err != nil {
            fmt.Println("Parsing of date went horrible... ",err)
        }
        untilDeparture := departure.Sub(time.Now()) 

        hour := time.Hour

        departure = departure.Add(hour)
        
        tm := departure.Format(time.Kitchen)
        fmt.Println("Bus",t.Route,"leaves at ", tm, " in ",untilDeparture.String(),
                    "from", t.Busstop)
    }

}


