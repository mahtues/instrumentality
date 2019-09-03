package handler

import (
	"fmt"
	"math"

	"encoding/json"
	"net/http"
)

type VertexId string

type Vertex struct{}

type EdgeId string

type Edge struct {
	Source VertexId `json:"source"`
	Target VertexId `json:"target"`
	Weight uint32   `json:"weight"`
}

type Graph struct {
	Vertices map[VertexId]Vertex `json:"vertices"`
	Edges    map[EdgeId]Edge     `json:"edges"`
}

func ShortestPathV2() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		g := Graph{}

		if err := json.NewDecoder(r.Body).Decode(&g); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		m := make(map[VertexId]map[VertexId]uint32)
		for id, _ := range g.Vertices {
			m[id] = make(map[VertexId]uint32)
		}

		for _, e := range g.Edges {
			if _, exists := g.Vertices[e.Source]; !exists {
				message := fmt.Sprintf("%s vertex not found", e.Source)
				http.Error(w, message, http.StatusBadRequest)
				return
			}

			if _, exists := g.Vertices[e.Target]; !exists {
				message := fmt.Sprintf("%s vertex not found", e.Target)
				http.Error(w, message, http.StatusBadRequest)
				return
			}

			m[e.Source][e.Target] = e.Weight
		}

		fmt.Fprint(w, shortestPath(m, VertexId("a"), VertexId("b")))
	})
}

func ShortestPath() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var g map[VertexId]map[VertexId]uint32

		json.NewDecoder(r.Body).Decode(&g)

		fmt.Fprint(w, shortestPath(g, VertexId("a"), VertexId("d")))
	})
}

func shortestPath(g map[VertexId]map[VertexId]uint32, s VertexId, t VertexId) uint32 {
	if s == t {
		return 0
	}

	var result uint32 = math.MaxUint32

	if _, exists := g[s]; !exists {
		return result
	}

	for u, w := range g[s] {
		result = min(result, w+shortestPath(g, u, t))
	}

	return result
}

func min(a uint32, b uint32) uint32 {
	if a < b {
		return a
	} else {
		return b
	}
}
