package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"exp/count"
	"exp/reverse"

	"github.com/ServiceWeaver/weaver"
)

//go:generate ../../cmd/weaver/weaver generate

type Response struct {
	Message string `json:"message"`
	Value   string `json:"value"`
}

type Reverser struct {
	Reverser weaver.Ref[reverse.Reverser]
}

func (rv Reverser) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get query param name
	name := r.URL.Query().Get("name")

	// reverse name

	reversed, _ := rv.Reverser.Get().Reverse(r.Context(), name)

	// write response
	resp := Response{
		Message: "reverse",
		Value:   reversed,
	}
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

type Counter struct {
	Counter weaver.Ref[count.Counter]
}

func (c Counter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get query param name
	name := r.URL.Query().Get("name")

	// count name
	ctx := r.Context()
	counter := c.Counter.Get()
	count, _ := counter.Count(ctx, name)

	// write response
	resp := Response{
		Message: "count",
		Value:   strconv.Itoa(count),
	}
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
