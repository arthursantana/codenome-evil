package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

var (
	players PlayerList
	planets PlanetList
	ships   ShipList
	units   UnitList

	lastDataUpdate int64
)

func main() {
	var (
		nPlanets         = flag.Int("planets", 10, "number of planets")
		dataUpdatePeriod = flag.Int("dataUpdatePeriod", 100, "number of milliseconds between data updates")
		tickPeriod       = flag.Int("tickPeriod", 1000, "number of milliseconds between ticks (has to be multiple of dataUpdatePeriod)")
	)

	flag.Parse()

	// VALIDATE NUMBERS HERE
	// END

	if *nPlanets > 20 {
		*nPlanets = 20
	}

	lastDataUpdate = 0
	planets = make([]Planet, *nPlanets)
	players = make([]Player, 0)
	ships = make([]Ship, 0)

	// GENERATE RANDOM STUFF
	rand.Seed(time.Now().UTC().UnixNano())
	for i := range planets {
		planets[i].randomize()
		planets[i].Id = i
	}

	// SERVE
	fs := http.FileServer(http.Dir("static"))
	fsHandler := http.StripPrefix("/static/", fs)
	http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		fsHandler.ServeHTTP(w, r)
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/static/index.htm", http.StatusMovedPermanently)
	})

	http.HandleFunc("/data/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "{")
		planets.writeJSON(w)
		fmt.Fprintf(w, ",")
		players.writeJSON(w)
		fmt.Fprintf(w, ",")
		ships.writeJSON(w)
		fmt.Fprintf(w, ",")
		units.writeJSON(w)
		fmt.Fprintf(w, "}")
	})

	http.HandleFunc("/ws/", handleWebsocket)

	go func() { // tick
		dataUpdatesPerTick := *tickPeriod / *dataUpdatePeriod

		dataUpdatesSinceLastTick := 0
		for {
			time.Sleep(time.Duration(*dataUpdatePeriod) * time.Millisecond)

			if dataUpdatesSinceLastTick >= dataUpdatesPerTick {
				dataUpdatesSinceLastTick = 0
				tick()
			} else {
				dataUpdatesSinceLastTick++
			}

			lastDataUpdate = time.Now().UTC().UnixNano()
		}
	}()

	http.ListenAndServe(":8081", nil)
}
