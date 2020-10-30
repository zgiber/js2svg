package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
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
				Name:  "src",
				Usage: "The uri pointing to the OpenAPI YAML document containing the Schema object(s). Can be file:// or http:// scheme.",
			},
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

type service struct {
	l       sync.RWMutex
	schema  map[string]interface{}
	objects []string
}

func newService(ctx *cli.Context) (*service, error) {
	srcPath := ctx.String("src")
	u, err := url.ParseRequestURI(srcPath)
	if err != nil {
		return nil, err
	}

	var src io.ReadCloser
	switch u.Scheme {
	case "http", "https":
		src, err = openRemoteSrc(u.String())
		if err != nil {
			return nil, err
		}
	case "file":
		src, err = openFileSrc(u.Path)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unknown src scheme: '%s'", u.Scheme)
	}
	defer src.Close()

	s := &service{
		l:       sync.RWMutex{},
		schema:  map[string]interface{}{},
		objects: []string{},
	}

	// contains all schemas
	s.schema, err = js2svg.ParseToMap(src, "components.schemas")
	if err != nil {
		return nil, err
	}

	log.Printf("parsed schema with %v items", len(s.schema))

	s.objects = listObjects(s.schema)
	return s, nil
}

func listObjects(m map[string]interface{}) []string {
	objects := []string{}
	for k, v := range m {
		if mm, ok := v.(map[string]interface{}); ok {
			if mm["type"] != nil {
				objects = append(objects, k)
			}
			childObjects := listObjects(mm)
			objects = append(objects, childObjects...)
		}
	}
	return objects
}

func (s *service) renderSVGObject(w http.ResponseWriter, r *http.Request) {
	s.l.RLock()
	defer s.l.RUnlock()

	path := chi.URLParam(r, "objectPath")
	m := js2svg.GetObject(s.schema, path)
	debug(m)
	d, err := js2svg.MakeDiagram(m, path) // path here is really just used for naming the root item
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := d.Render(w); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func debug(v interface{}) {
	b, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(b))
}

func (s *service) listObjects(w http.ResponseWriter, r *http.Request) {
	s.l.RLock()
	defer s.l.RUnlock()

	for k, v := range s.schema {
		switch t := v.(type) {
		case map[string]interface{}:
			if t["type"] != nil && t["type"] == "object" {
				link := fmt.Sprintf("http://%s/render/%s\n", r.Host, k)
				fmt.Fprintf(w, `<a href="%[1]s">%[1]s</a><br/>`, link)
			}
		default:
		}
	}
}

func run(ctx *cli.Context) error {
	js2svg.ExternalDivider = "."
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	s, err := newService(ctx)
	if err != nil {
		return err
	}
	r := chi.NewMux()

	r.Get("/render/{objectPath}", s.renderSVGObject)
	r.Get("/list", s.listObjects)

	return http.ListenAndServe(":"+ctx.String("port"), r)
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
