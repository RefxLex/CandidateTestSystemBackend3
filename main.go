package main

import (
	"CandidateTestSystemBackend3/service/topic_service"
	"CandidateTestSystemBackend3/types"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/docgen"
	"github.com/go-chi/render"
)

var routes = flag.Bool("routes", false, "Generate router documentation")

func main() {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("root."))
	})

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	r.Get("/panic", func(w http.ResponseWriter, r *http.Request) {
		panic("test")
	})

	r.Route("/api/topic", func(r chi.Router) {
		r.Route("/{topicID}", func(r chi.Router) {
			r.Use(TopicCtx)
			r.Get("/", GetTopic)
			//r.Put("/", UpdateArticle)
			//r.Delete("/", DeleteArticle)
		})
	})

	if *routes {
		// fmt.Println(docgen.JSONRoutesDoc(r))
		fmt.Println(docgen.MarkdownRoutesDoc(r, docgen.MarkdownOpts{
			ProjectPath: "github.com/go-chi/chi/v5",
			Intro:       "Welcome to the chi/_examples/rest generated docs.",
		}))
		return
	}

	http.ListenAndServe(":8083", r)
}

// TopicCtx middleware is used to load an Topic object from
// the URL parameters passed through as the request. In case
// the Topic could not be found, we stop here and return a 404.
func TopicCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var topic *types.Topic
		var err error
		topicID := chi.URLParam(r, "topicID")

		topicIDConv, err := strconv.Atoi(topicID)
		if err != nil {
			panic(err)
		}

		if topicID != "" {
			topic, err = topic_service.GetTopicById(topicIDConv)
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(err.Error()))
			return
		}
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(err.Error()))
			return
		}

		ctx := context.WithValue(r.Context(), "topic", topic)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetTopic returns the specific Topic. You'll notice it just
// fetches the Topic right off the context, as its understood that
// if we made it this far, the Topic must be on the context. In case
// its not due to a bug, then it will panic, and our Recoverer will save us.
func GetTopic(w http.ResponseWriter, r *http.Request) {
	// Assume if we've reach this far, we can access the topic
	// context because this handler is a child of the TopicCtx
	// middleware. The worst case, the recoverer middleware will save us.
	topic := r.Context().Value("topic").(*types.Topic)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(*topic)
}

// UpdateArticle updates an existing Article in our persistent store.
// func UpdateArticle(w http.ResponseWriter, r *http.Request) {
// 	article := r.Context().Value("article").(*Article)

// 	data := &ArticleRequest{Article: article}
// 	if err := render.Bind(r, data); err != nil {
// 		render.Render(w, r, ErrInvalidRequest(err))
// 		return
// 	}
// 	article = data.Article
// 	dbUpdateArticle(article.ID, article)

// 	render.Render(w, r, NewArticleResponse(article))
// }

// DeleteArticle removes an existing Article from our persistent store.
// func DeleteArticle(w http.ResponseWriter, r *http.Request) {
// 	var err error

// 	// Assume if we've reach this far, we can access the article
// 	// context because this handler is a child of the ArticleCtx
// 	// middleware. The worst case, the recoverer middleware will save us.
// 	article := r.Context().Value("article").(*Article)

// 	article, err = dbRemoveArticle(article.ID)
// 	if err != nil {
// 		render.Render(w, r, ErrInvalidRequest(err))
// 		return
// 	}

// 	render.Render(w, r, NewArticleResponse(article))
// }
