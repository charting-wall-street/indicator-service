package web

import (
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"inca/internal/compute"
	"inca/internal/config"
	"inca/internal/definition"
	"inca/internal/indicator"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func GetIndicatorList(w http.ResponseWriter, r *http.Request) {
	response := &struct {
		Indicators []indicator.ListItem `json:"indicators"`
	}{Indicators: indicator.List()}
	sendResponse(w, r, response)
}

type IndicatorOptions struct {
	Name     string
	Def      *definition.Definition
	Block    int64
	Interval int64
	Params   []int
	Symbol   string
}

func requestOptions(w http.ResponseWriter, r *http.Request) (*IndicatorOptions, bool) {

	var err error
	var ok bool
	opts := new(IndicatorOptions)

	// get indicator name
	opts.Name = mux.Vars(r)["name"]

	// get symbol
	opts.Symbol = r.URL.Query().Get("symbol")

	// get block number
	opts.Block, err = strconv.ParseInt(r.URL.Query().Get("block"), 10, 64)
	if err != nil {
		http.Error(w, "bad block number", http.StatusBadRequest)
		return nil, false
	}

	// get time interval
	opts.Interval, err = strconv.ParseInt(r.URL.Query().Get("interval"), 10, 64)
	if err != nil {
		http.Error(w, "bad interval number", http.StatusBadRequest)
		return nil, false
	}

	if opts.Interval == 1 {
		http.Error(w, "interval cannot be 1", http.StatusBadRequest)
		return nil, false
	}

	// get parameter list
	paramListS := strings.Split(r.URL.Query().Get("params"), ",")
	opts.Params = make([]int, len(paramListS))
	for i, p := range paramListS {
		opts.Params[i], err = strconv.Atoi(p)
		if err != nil {
			http.Error(w, "invalid indicator parameter", http.StatusBadRequest)
			return nil, false
		}
	}

	// find definition
	opts.Def, ok = indicator.ByName(opts.Name)
	if !ok {
		http.Error(w, "indicator does not exist", http.StatusNotFound)
		return nil, false
	}

	return opts, true
}

func GetTransitionIndicator(w http.ResponseWriter, r *http.Request) {

	opts, ok := requestOptions(w, r)
	if !ok {
		return
	}

	resolutionS := r.URL.Query().Get("resolution")
	resolution, err := strconv.ParseInt(resolutionS, 10, 64)
	if err != nil {
		http.Error(w, "invalid resolution parameter", http.StatusBadRequest)
		return
	}

	comp := compute.NewCalculable(opts.Def, opts.Symbol, opts.Block*resolution/opts.Interval, opts.Interval, resolution)
	comp.UseCache = r.URL.Query().Get("cache") != "no-cache"

	res, err := comp.Transition(opts.Block, resolution, opts.Params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sendResponseIndicator(w, r, res)
}

func GetIndicator(w http.ResponseWriter, r *http.Request) {

	opts, ok := requestOptions(w, r)
	if !ok {
		return
	}

	comp := compute.NewCalculable(opts.Def, opts.Symbol, opts.Block, opts.Interval, opts.Interval)
	comp.UseCache = r.URL.Query().Get("cache") != "no-cache"

	res, err := comp.Compute(opts.Params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sendResponseIndicator(w, r, res)
}

func router() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/indicators", GetIndicatorList).Methods("GET")
	r.HandleFunc("/indicators/{name}", GetIndicator).Methods("GET")
	r.HandleFunc("/indicators/t/{name}", GetTransitionIndicator).Methods("GET")
	return r
}

func Start() {

	// Middleware and routes
	app := negroni.New(negroni.NewRecovery())
	app.UseHandler(router())

	// CORS
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization", "Cache-Control"})
	exposedOk := handlers.ExposedHeaders([]string{"Link"})
	originsOk := handlers.AllowedOrigins(strings.Split(config.ServiceConfig().AllowedOrigins(), ","))
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS", "DELETE"})
	handler := handlers.CORS(originsOk, headersOk, methodsOk, exposedOk)(app)

	// Setup server
	server := &http.Server{
		Addr:    ":" + config.ServiceConfig().Port(),
		Handler: handler,
	}

	fmt.Printf("listening on port %s\n", config.ServiceConfig().Port())

	log.Fatal(server.ListenAndServe())
}
