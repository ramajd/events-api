package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/ramajd/events-api/handlers"
	"github.com/ramajd/events-api/objects"
	"github.com/ramajd/events-api/store"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	router    *mux.Router
	flushAll  func(t *testing.T)
	createOne func(t *testing.T, name string) *objects.Event
	getOne    func(t *testing.T, id string, wantErr bool) *objects.Event
)

func TestMain(t *testing.M) {
	log.Println("Registering")

	conn := "postgres://user:password@localhost:5432/db?sslmode=disable"
	if c := os.Getenv("DB_CONN"); c != "" {
		conn = c
	}

	router = mux.NewRouter().PathPrefix("/api/v1/").Subrouter()
	st := store.NewPostgresEventStore(conn)
	hnd := handlers.NewEventHandler(st)
	RegisterAllRoutes(router, hnd)

	flushAll = func(t *testing.T) {
		db, err := gorm.Open(postgres.Open(conn), nil)
		if err != nil {
			t.Fatal(err)
		}
		db.Delete(&objects.Event{}, "1=1")
	}

	createOne = func(t *testing.T, name string) *objects.Event {
		evt := &objects.Event{
			Name:        name,
			Description: "Description of " + name,
			Website:     "https://" + name + ".com",
			Slot: &objects.TimeSlot{
				StartTime: time.Now().UTC(),
				EndTime:   time.Now().UTC().Add(time.Hour),
			},
		}
		if err := st.Create(context.TODO(), &objects.CreateRequest{Event: evt}); err != nil {
			t.Fatal(err)
		}
		return evt
	}
	getOne = func(t *testing.T, id string, wantErr bool) *objects.Event {
		evt, err := st.Get(context.TODO(), &objects.GetRequest{ID: id})
		if err != nil && wantErr {
			t.Fatal(err)
		}
		return evt
	}

	log.Println("Starting")
	os.Exit(t.Run())

}

func Do(req *http.Request) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func TestUnknownEndpoints(t *testing.T) {
	tests := []struct {
		name  string
		setup func(t *testing.T) *http.Request
	}{
		{
			name: "root",
			setup: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodGet, "/", nil)
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
		},
		{
			name: "api-root",
			setup: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodGet, "/api/v1", nil)
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
		},
		{
			name: "random",
			setup: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodGet, "/random", nil)
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := Do(tt.setup(t))
			_ = assert.Equal(t, http.StatusNotFound, w.Code) &&
				assert.Equal(t, "404 page not found\n", string(w.Body.Bytes()))
		})
	}
}

func TestListEndpoint(t *testing.T) {

	flushAll(t)

	tests := []struct {
		name    string
		code    int
		setup   func(t *testing.T) *http.Request
		listLen int
	}{
		{
			name: "Zero",
			setup: func(t *testing.T) *http.Request {
				flushAll(t)
				req, err := http.NewRequest(http.MethodGet, "/api/v1/events", nil)
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			code:    http.StatusOK,
			listLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := Do(tt.setup(t))
			got := &objects.EventResponseWrapper{}
			assert.Equal(t, tt.code, w.Code)
			assert.Nil(t, json.Unmarshal(w.Body.Bytes(), got))
			assert.Equal(t, len(got.Events), tt.listLen)
		})
	}

}
