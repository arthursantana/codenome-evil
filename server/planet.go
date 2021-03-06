package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
)

type Planet struct {
	Id      int    `json:"id"`
	OwnerId int    `json:"ownerId"`
	Name    string `json:"name"`

	Position [2]float64 `json:"position"`
	R        int        `json:"r"`

	// resources
	Workers  int `json:"workers"`
	Cattle   int `json:"cattle"`
	Energy   int `json:"energy"`
	Obtanium int `json:"obtanium"`

	BusyWorkers int `json:"busyWorkers"`
	BusyCattle  int `json:"busyCattle"`
	BusyEnergy  int `json:"busyEnergy"`

	NotEnoughWorkers bool `json:"notEnoughWorkers"`
	NotEnoughCattle  bool `json:"notEnoughCattle"`
	NotEnoughEnergy  bool `json:"notEnoughEnergy"`
	NotEnoughFarms   bool `json:"notEnoughFarms"`

	Buildings [][]Building `json:"buildings"`

	DockSpace      int `json:"dockSpace"`
	UnitSpace      int `json:"unitSpace"`
	EnemyUnitSpace int `json:"enemyUnitSpace"`
}

var defaultNames = []string{"Big Rock", "Mercury", "Venus", "Earth", "Mars", "Jupiter", "Saturn", "Neptune", "Uranus", "Pluto", "Urectum", "Tessia", "Sur'Kesh", "Tuchanka", "Omega", "Palaven", "Rannoch", "3834 Zappafrank", "Omicron Persei 8", "Planet 9 from Outer Space"}

func (p *Planet) randomize() {
	index := rand.Intn(len(defaultNames))
	p.Name = defaultNames[index]
	defaultNames = append(defaultNames[:index], defaultNames[index+1:]...)

	width := 1300
	height := 750

	p.Position[0] = 50 + float64(rand.Intn(width-100))
	p.Position[1] = 50 + float64(rand.Intn(height-100))
	p.R = 4 // + rand.Intn(5)

	p.Workers = 35000
	p.Cattle = 60000
	p.Energy = 0
	p.Obtanium = 4000

	p.DockSpace = 10
	p.UnitSpace = 10
	p.EnemyUnitSpace = 10

	p.Buildings = make([][]Building, p.R)
	for i := range p.Buildings {
		p.Buildings[i] = make([]Building, p.R)

		for j := range p.Buildings[i] {
			p.Buildings[i][j] = Building{"", true, 0, 0, []string{}}
		}
	}

	p.Buildings[0][0].Type = "hq"

	p.OwnerId = -1
}

type PlanetList []Planet

var planets PlanetList

func (planets PlanetList) writeJSON(w io.WriteCloser) {
	fmt.Fprintf(w, "\"planets\": [")
	separator := ""
	for i := range planets {
		p := planets[i]
		pJson, err := json.Marshal(p)
		if err != nil {
			log.Println(err)
			return
		}

		fmt.Fprintf(w, "%v%v", separator, string(pJson))
		separator = ",\n"
	}
	fmt.Fprintf(w, "]")
}

func (p *Planet) combat() {
	homeSuccessfulShots := 0
	awaySuccessfulShots := 0

	for i := range units {
		if units[i].PlanetId == p.Id {
			if units[i].OwnerId == p.OwnerId {
				if units[i].hits() {
					homeSuccessfulShots++
				}
			} else {
				if units[i].hits() {
					awaySuccessfulShots++
				}
			}
		}
	}

	anybodyFromHome := false
	anybodyFromAway := false
	for i := range units {
		if units[i].PlanetId == p.Id {
			if units[i].OwnerId == p.OwnerId {
				if awaySuccessfulShots > 0 {
					units[i].PlanetId = -1
					units[i].OwnerId = -1

					awaySuccessfulShots--
					p.UnitSpace++
				} else {
					anybodyFromHome = true
				}
			} else {
				if homeSuccessfulShots > 0 {
					units[i].PlanetId = -1
					units[i].OwnerId = -1

					homeSuccessfulShots--
					p.EnemyUnitSpace++
				} else {
					anybodyFromAway = true
				}
			}
		}
	}

	if anybodyFromAway && !anybodyFromHome { // conquest!
		for i := range units {
			if units[i].PlanetId == p.Id {
				if units[i].OwnerId != p.OwnerId {
					p.OwnerId = units[i].OwnerId

					p.UnitSpace, p.EnemyUnitSpace = p.EnemyUnitSpace, p.UnitSpace // away and home team switch places

					for j := range ships { // steals present ships
						if ships[j].PlanetId == p.Id {
							ships[j].OwnerId = units[i].OwnerId
						}
					}

					break
				}
			}
		}
	}
}
