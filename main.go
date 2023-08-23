package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ServiceWeaver/weaver"
)

//go:embed index.html
var indexHtml string

func main() {
	if err := weaver.Run(context.Background(), run); err != nil {
		panic(err)
	}
}

type app struct {
	weaver.Implements[weaver.Main]
	searcher weaver.Ref[Searcher]
	lis      weaver.Listener `weaver:"emojis"`
}

func run(ctx context.Context, a *app) error {
	a.Logger(ctx).Info("emojis listener available.", "addr", a.lis)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		fmt.Fprint(w, indexHtml)
	})

	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")
		emojis, err := a.searcher.Get().Search(r.Context(), query)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		bytes, err := json.Marshal(emojis)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintln(w, string(bytes))
	})
	return http.Serve(a.lis, nil)
}
