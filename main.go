package main

import (
	"encoding/csv"
	"os"
	"io"
	"strconv"
	"log"
	"time"
	"fmt"
)

type Graph struct {
	RawEdges    []Edge
	Vertices    map[uint64]bool
	VertexEdges map[uint64]map[uint64]float64
	Undirected  bool
	NegEdges    bool
}

type Edge struct {
	From   uint64
	To  uint64
	Weight float64
}

var VacationGuys = []uint64{}

func (gr *Graph) IsBipartite(origin uint64) bool {
	colours := map[uint64]bool{origin: false}
	queue := []uint64{origin}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		for v := range gr.VertexEdges[current] {
			if _, visited := colours[v]; !visited {
				colours[v] = !colours[current]
				queue = append(queue, v)
			} else {
				if colours[v] == colours[current] {
					return false
				}
			}
		}
	}
	return true
}

func cleanAndReset(elements []uint64) []uint64 {
	encountered := map[uint64]bool{}
	result := []uint64{}
	for v := range elements {
		if encountered[elements[v]] == true {} else {
			encountered[elements[v]] = true
			result = append(result, elements[v])
		}
	}
	return result
}

func readData(file string) ([]Edge, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	csvr := csv.NewReader(f)
	couples := []Edge{}
	for {
		row, err := csvr.Read()
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return couples, err
		}
		p := &Edge{}
		p.Weight = 0
		if p.From, err = strconv.ParseUint(row[0], 10,64); err != nil {
			return nil, err
		}
		if p.To, err = strconv.ParseUint(row[1], 10,64); err != nil {
			return nil, err
		}
		log.Println(p.From)
		log.Println(p.To)
		couples = append(couples,*p)
	}
}

func GetGraph(edges []Edge, undirected bool) (ug *Graph) {
	var weight float64
	ug = &Graph{
		RawEdges:    edges,
		Vertices:    make(map[uint64]bool),
		VertexEdges: make(map[uint64]map[uint64]float64),
		Undirected:  undirected,
		NegEdges:    false,
	}
	for _, edge := range edges {
		weight = edge.Weight
		if weight < 0 {
			ug.NegEdges = true
		}
		ug.Vertices[edge.From] = true
		ug.Vertices[edge.To] = true
		if _, ok := ug.VertexEdges[edge.From]; ok {
			ug.VertexEdges[edge.From][edge.To] = weight
		} else {
			ug.VertexEdges[edge.From] = map[uint64]float64{edge.To: weight}
		}
		if undirected {
			if _, ok := ug.VertexEdges[edge.To]; ok {
				ug.VertexEdges[edge.To][edge.From] = weight
			} else {
				ug.VertexEdges[edge.To] = map[uint64]float64{edge.From: weight}
			}
		}
	}
	return
}

func (gr *Graph) NewReversedGraph() (rev *Graph) {
	rev = &Graph{
		Vertices:    gr.Vertices,
		VertexEdges: make(map[uint64]map[uint64]float64),
		Undirected:  false,
	}

	for v, e := range gr.VertexEdges {
		for d, w := range e {
			if _, ok := rev.VertexEdges[d]; ok {
				rev.VertexEdges[d][v] = w
			} else {
				rev.VertexEdges[d] = map[uint64]float64{v: w}
			}
		}
	}
	return
}


func upTime() {
	for {
		now := time.Now()
		now.Second()
	}
}

func main() {
	go upTime()

	mapOfCouples, err := readData("couplesData.csv")
	if err != nil{
		log.Println("error de lectura", err)
	}
	Graph := GetGraph(mapOfCouples,false)
	if Graph.IsBipartite(0){
		log.Println("Parejas en forma bipartite. Ok!")
		for stockholm, LondonMatches := range Graph.VertexEdges {
			log.Println("************INIT ITERATION**************")
			currentVacationCandidate := stockholm
			currentMajorVertices := len(LondonMatches)
			for londonUser := range LondonMatches {
				log.Println("# equipos: ",len(Graph.NewReversedGraph().VertexEdges[londonUser]))
				log.Println("EvaluatedPartner: ",londonUser)
				if len(Graph.NewReversedGraph().VertexEdges[londonUser]) > currentMajorVertices{
					currentMajorVertices = len(Graph.NewReversedGraph().VertexEdges[londonUser])
					currentVacationCandidate = londonUser
				}
				log.Println("currentVacationCandidate: ",currentVacationCandidate)
			}
			VacationGuys = append(VacationGuys,currentVacationCandidate)
		}
	}else{
		log.Println("Las parejas no son correctas.")
	}
	log.Println("")
	log.Println(":::::::::::::::::::::::: VacationPeopleResult ::::::::::::::::::::::::")
	log.Println("")
	log.Println(cleanAndReset(VacationGuys))

	log.Println(":::::::::::::::::::::::::::::::::")
	log.Println("Execution Done, press Ctrl+X to exit")
	select {}
	fmt.Scanln()
}


