package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/urfave/cli"
	"github.com/zgiber/js2svg"
)

func main() {
	app := &cli.App{
		Action: run,
		Name:   "generate", // TODO: give this a good name
		Usage:  "Generate SVG diagrams for JSON Schema objects",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "src",
				Usage: "The (absolute) path to the JSON document containing the Schema object(s). Can be uri with a file:// or http:// scheme.",
			},
			&cli.StringFlag{
				Name:  "path",
				Usage: "The path of the selected object within the JSON document. (eg.: 'components.schemas.myAwesomeSchema')",
			},
			&cli.StringFlag{
				Name:  "out",
				Usage: "The file path for the SVG file to be written. If a file with the same name already exists it will be overwritten. If empty, stdout is used.",
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func run(ctx *cli.Context) error {
	srcPath := ctx.String("src")
	u, err := url.ParseRequestURI(srcPath)
	if err != nil {
		return err
	}

	var src io.ReadCloser
	switch u.Scheme {
	case "http", "https":
		src, err = openRemoteSrc(u.String())
		if err != nil {
			return err
		}
	case "file":
		src, err = openFileSrc(u.Path)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown src scheme: '%s'", u.Scheme)
	}
	defer src.Close()

	d, err := js2svg.ParseToDiagram(src, ctx.String("path"))
	if err != nil {
		return err
	}

	dst := os.Stdout
	if len(ctx.String("out")) > 0 {
		dst, err = os.Create(ctx.String("out"))
		if err != nil {
			return err
		}
	}

	return d.Render(dst)
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
