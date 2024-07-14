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
	"github.com/go-chi/cors"
	"github.com/go-chi/docgen"
	"github.com/go-chi/render"
)

var routes = flag.Bool("routes", false, "Generate router documentation")

func main() {
	r := chi.NewRouter()

	// Basic CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

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
		r.Get("/", GetAllTopics)
		r.Post("/", CreateTopic)
		r.Route("/{topicID}", func(r chi.Router) {
			r.Use(TopicCtx)
			r.Get("/", GetTopic)
			r.Put("/", UpdateTopic)
			r.Delete("/", DeleteTopic)
		})
	})

	if *routes {
		fmt.Println(docgen.MarkdownRoutesDoc(r, docgen.MarkdownOpts{
			ProjectPath: "github.com/go-chi/chi/v5",
			Intro:       "Welcome to the chi/_examples/rest generated docs.",
		}))
		return
	}

	http.ListenAndServe(":8083", r)
}

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

func GetTopic(w http.ResponseWriter, r *http.Request) {
	topic := r.Context().Value("topic").(*types.Topic)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(*topic)
}

func GetAllTopics(w http.ResponseWriter, r *http.Request) {
	topics, err := topic_service.GetAllTopics()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(topics)
	}
}

func CreateTopic(w http.ResponseWriter, r *http.Request) {
	// decode body
	var body types.Topic
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	// create topic
	createdTopic, err := topic_service.CreateTopic(&body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(*createdTopic)
}

func UpdateTopic(w http.ResponseWriter, r *http.Request) {
	// check if topic exists
	oldTopic := r.Context().Value("topic").(*types.Topic)

	// decode body
	var body types.Topic
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	// update topic
	newTopic, err := topic_service.UpdateTopic(oldTopic.Id, &body)

	// return updated topic
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(*newTopic)
}

func DeleteTopic(w http.ResponseWriter, r *http.Request) {
	topic := r.Context().Value("topic").(*types.Topic)

	err := topic_service.DeleteTopic(topic.Id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
