package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/go-chi/chi"
	"github.com/urfave/cli"
	"github.com/zgiber/js2svg"
)

func main() {
	app := &cli.App{
		Action: run,
		Name:   "serve",
		Usage:  "Serve SVG diagrams generated from JSONSchema documents.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "port",
				Usage: "The port where the http server accepts incoming connections.",
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func run(ctx *cli.Context) error {
	js2svg.ExternalDivider = "."
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	s, err := newService()
	if err != nil {
		return err
	}
	r := chi.NewMux()

	r.Get("/{collection}/{object}", s.renderObject)
	r.Get("/{collection}", s.listObjects)
	r.Post("/collections", s.registerObjects)

	return http.ListenAndServe(":"+ctx.String("port"), r)
}

type service struct {
	l       sync.RWMutex
	schemas map[string]map[string]interface{}
}

func newService() (*service, error) {
	s := &service{
		l:       sync.RWMutex{},
		schemas: map[string]map[string]interface{}{},
	}

	return s, nil
}

func (s *service) registerObjects(w http.ResponseWriter, r *http.Request) {
	type requestObject struct {
		URL        string `json:"url,omitempty"`
		Collection string `json:"collection,omitempty"`
		SchemaPath string `json:"schema_path,omitempty"`
	}

	requestData := &requestObject{}
	err := json.NewDecoder(r.Body).Decode(requestData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	u, err := url.ParseRequestURI(requestData.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var src io.ReadCloser
	switch u.Scheme {
	case "http", "https":
		src, err = openRemoteSrc(u.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}
	// case "file":
	// 	src, err = openFileSrc(u.Path)
	// 	if err != nil {
	// 		http.Error(w, err.Error(), http.StatusInternalServerError)
	// 		return
	// 	}
	default:
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer src.Close()

	s.l.Lock()
	defer s.l.Unlock()
	schema, err := js2svg.ParseToMap(src, requestData.SchemaPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}

	s.schemas[requestData.Collection] = schema
}

func (s *service) renderObject(w http.ResponseWriter, r *http.Request) {
	// locking everything like this is slow of course however
	// it's not an issue (unless this tool becomes a platform .. lol)
	s.l.RLock()
	defer s.l.RUnlock()

	objectName := chi.URLParam(r, "object")
	collection := chi.URLParam(r, "collection")
	schema, exists := s.schemas[collection]
	if !exists {
		http.Error(w, "collection not found", http.StatusNotFound)
		return
	}

	m := js2svg.GetObject(schema, objectName)
	debug(m)

	d, err := js2svg.MakeDiagram(m, objectName) // path here is really just used for naming the root item

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = d.Render(w)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *service) listObjects(w http.ResponseWriter, r *http.Request) {
	s.l.RLock()
	defer s.l.RUnlock()

	collection := chi.URLParam(r, "collection")
	schema, exists := s.schemas[collection]
	if !exists {
		http.Error(w, "collection not found", http.StatusNotFound)
		return
	}

	for k, v := range schema {
		switch t := v.(type) {
		case map[string]interface{}:
			if t["type"] != nil && t["type"] == "object" {
				link := fmt.Sprintf(strings.TrimRight(r.RequestURI, "/") + "/" + k)
				fmt.Fprintf(w, `<a href="%[1]s">%[1]s</a><br/>`, link)
			}
		default:
		}
	}
}

func openFileSrc(path string) (io.ReadCloser, error) {
	return os.Open(path)
}

func openRemoteSrc(url string) (io.ReadCloser, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("remote retuned status %v", resp.StatusCode)
	}

	return resp.Body, nil
}

func debug(v interface{}) {
	b, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(b))
}
