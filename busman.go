package main

import (
        "encoding/json"
        "flag"
        "fmt"
        "io/ioutil"
        "net/http"
        "strings"
        "time"
        "os"
)

type Stop struct {
        V            int     `json:"__v,int"`
        Id           string  `json:"_id"`
        City         string
        Name         string
        Sort_id      int
        Destinations []Destination
}

type Destination struct {
        V            int     `json:"__v,int"`
        Id           string  `json:"_id"`
        City         string
        Name         string
        Sort_id      int
        Destinations []string
}

type Time struct {
        FromId  string
        From    string
        ToId    string
        To      string
        Date    string
        //Arrival string
        Route   int
        Busstop string
        Hash    string
        Id      string  `json:"_id"`
        V       int     `json:"__v"`
}

func main() {


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
                fmt.Print("Error contacting rutebuss.no")
                return
        }


        fromString := flag.String("from", "UiT",
                "The bus stop you are traveling from")
        toString := flag.String("to", "Sentrum",
                "The bus stop you are traveling to")

        invert := flag.Bool("i", false,
            "Switch the from and to values with each other")

        specificTime := flag.String("time", "now",
                "Specific time of departure. Formatted as hour:minute") 

        flag.Usage = func(){
            fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
            flag.PrintDefaults()
            fmt.Fprintf(os.Stderr, "\nAvailible stops: \n")//+availibleStops)
            for i, _ := range stops {
                if i%2 == 0 {
                    if i == len(stops) - 1 {
                        fmt.Fprintln(os.Stderr, stops[i].Name)
                    } else {
                        fmt.Fprintf(os.Stderr, "%-20s%s\n", stops[i].Name, stops[i+1].Name)
                    }
                }
            }
        }

        flag.Parse()
        if *invert {
            tmp := fromString
            fromString = toString
            toString = tmp
        }

        fmt.Println("From", *fromString, "to", *toString)


        var from, to string

        for _, stop := range stops {
                if strings.ToLower(stop.Name) == strings.ToLower(*fromString) {
                        from = stop.Id
                }
                if strings.ToLower(stop.Name) == strings.ToLower(*toString) {
                        to = stop.Id
                }
        }

        if from == "" && to == "" {
                fmt.Println("Bus route ", from, "-", to, "does not exist")
                return
        }
        
        var t time.Time
        var timeString string

        if *specificTime != "now" {
            const layout = "15:04"
            
            location, _ := time.LoadLocation("Local")
            
            t, err = time.ParseInLocation(layout, *specificTime, location)

            if err != nil {
                fmt.Println("Error: Could not parse time:", *specificTime) 
                fmt.Println(t)
                //flag.Usage()
                return
            }
            
            t_now := time.Now() 
            y,m,d := t_now.Date() 

            // date() gave day and month +1 :( 
            d = d - 1
            mnth := int(m) - 1
            
            // since t object is 0000-00-00Thour:minute we need to add today's
            // date and everything. 
            t = t.AddDate(y,mnth,d)

            timeString = t.UTC().Format(time.RFC3339)
        } else {
            t = time.Now()
            timeString = t.UTC().Format(time.RFC3339)
        }

        travelURL :=
                "http://rutebuss.no/departure?from=" + from + "&to=" + to + "&date=" + timeString
        
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
                fmt.Println("Bus route ", *fromString, "-", *toString, "does not exist")
                return
        }

        if len(times) == 0 {
                fmt.Println("Bus route ", *fromString, "-", *toString, "does not exist")
        }

        for _, departure_time := range times {

                departure, err := time.Parse(time.RFC3339, departure_time.Date)

                if err != nil {
                        fmt.Println("Parsing of date went horrible... ", err)
                }

                untilDeparture := departure.Sub(time.Now())
                untilDString := untilDeparture.String()
                untilDString = strings.Split(untilDString, ".")[0]
                untilDString = untilDString + "s"

                hour := time.Hour
                departure = departure.Add(hour)

                tm := departure.Format(time.Kitchen)
                fmt.Println("Bus", departure_time.Route, "leaves at", tm, "in",
                        untilDString, "from", departure_time.Busstop)
        }

}
