package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-martini/martini"
	"github.com/wanelo/image-server/core"
	fetcher "github.com/wanelo/image-server/fetcher/http"
	"github.com/wanelo/image-server/parser"
)

func genericImageHandler(params martini.Params, r *http.Request, w http.ResponseWriter, sc *core.ServerConfiguration) {
	ic, err := parser.NameToConfiguration(sc, params["filename"])
	if err != nil {
		errorHandler(err, w, r, http.StatusNotFound, ic)
	}
	qs := r.URL.Query()

	ic.ServerConfiguration = sc
	ic.Namespace = params["namespace"]
	ic.ID = params["id1"] + params["id2"] + params["id3"]
	ic.Source = qs.Get("source")
	imageHandler(ic, w, r)

}

func multiImageHandler(params martini.Params, r *http.Request, w http.ResponseWriter, sc *core.ServerConfiguration) {
	qs := r.URL.Query()

	go func() {
		outputs := strings.Split(qs.Get("outputs"), ",")
		fmt.Println("multiImageHandler")
		fmt.Println(outputs)
		for _, filename := range outputs {
			fmt.Println(filename)
			ic, err := parser.NameToConfiguration(sc, filename)
			if err != nil {
				continue
			}
			qs := r.URL.Query()

			ic.ServerConfiguration = sc
			ic.Namespace = params["namespace"]
			ic.ID = params["id1"] + params["id2"] + params["id3"]
			ic.Source = qs.Get("source")

			allowed, _ := allowedImage(ic)
			if !allowed {
				return
			}

			_, err = sc.Adapters.Processor.CreateImage(ic)
			if err != nil {
				return
			}

			go func() {
				sc.Events.ImageProcessed <- ic
			}()
		}
	}()

	ic := &core.ImageConfiguration{
		ServerConfiguration: sc,
		Namespace:           params["namespace"],
		ID:                  params["id1"] + params["id2"] + params["id3"],
		Source:              qs.Get("source"),
	}
	err := fetcher.FetchOriginal(ic)
	if err != nil {
		errorHandler(err, w, r, http.StatusNotFound, ic)
		return
	}
}

func imageHandler(ic *core.ImageConfiguration, w http.ResponseWriter, r *http.Request) {
	allowed, err := allowedImage(ic)
	if !allowed {
		errorHandler(err, w, r, http.StatusNotAcceptable, ic)
		return
	}

	sc := ic.ServerConfiguration
	resizedPath, err := sc.Adapters.Processor.CreateImage(ic)
	if err != nil {
		errorHandler(err, w, r, http.StatusNotFound, ic)
		return
	}
	http.ServeFile(w, r, resizedPath)
}

func errorHandler(err error, w http.ResponseWriter, r *http.Request, status int, ic *core.ImageConfiguration) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		fmt.Fprint(w, "404 image not available. ", err)
	}

	go func() {
		sc := ic.ServerConfiguration
		sc.Events.ImageProcessedWithErrors <- ic
	}()
}

func allowedImage(ic *core.ImageConfiguration) (bool, error) {
	sc := ic.ServerConfiguration
	// verify maximum width
	if ic.Width > sc.MaximumWidth {
		err := fmt.Errorf("maximum width is: %v\n", sc.MaximumWidth)
		return false, err
	}

	// verify image format
	for _, ext := range sc.WhitelistedExtensions {
		if ext == ic.Format {
			return true, nil
		}
	}
	return false, fmt.Errorf("format not allowed %s", ic.Format)
}
