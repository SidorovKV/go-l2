package server

import (
	"dev11/model"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type EventSvc interface {
	CreateEvent(event model.Event) (model.Event, error)
	UpdateEvent(event model.Event) (model.Event, error)
	DeleteEvent(event model.Event) error
	EventsForDay(userId uint, day time.Time) ([]model.Event, error)
	EventsForWeek(userId uint, startDay time.Time) ([]model.Event, error)
	EventsForMonth(userId uint, month time.Month, year int) ([]model.Event, error)
}

type Server struct {
	EventSvc
}

func NewServer(repository EventSvc) *Server {
	return &Server{repository}
}

func (s *Server) Start(port uint16) {
	http.HandleFunc("/create_event", s.createEvent)
	http.HandleFunc("/update_event", s.updateEvent)
	http.HandleFunc("/delete_event", s.deleteEvent)
	http.HandleFunc("/events_for_day", s.eventsForDay)
	http.HandleFunc("/events_for_week", s.eventsForWeek)
	http.HandleFunc("/events_for_month", s.eventsForMonth)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
}

func (s *Server) createEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Println("Method not allowed on /create_event")
		sendError(http.StatusMethodNotAllowed, "Method not allowed", w)
		return
	}

	var event model.Event
	err := unmarshalEvent(r, &event)
	if err != nil {
		log.Println("Error unmarshalling request body on /create_event")
		sendError(http.StatusBadRequest, "Bad request body", w)
		return
	}
	if err = validateEvent(event); err != nil {
		log.Println(err)
		sendError(http.StatusBadRequest, err.Error(), w)
		return
	}

	createdEvent, err := s.CreateEvent(event)
	if err != nil {
		log.Println(err)
		sendError(http.StatusServiceUnavailable, "Internal server error", w)
		return
	}

	result := successResponse{createdEvent}
	resp, err := json.Marshal(result)
	if err != nil {
		log.Println(err)
		sendError(http.StatusInternalServerError, "Internal server error", w)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(resp)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Event created successfully")
}

func (s *Server) updateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Println("Method not allowed on /update_event")
		sendError(http.StatusMethodNotAllowed, "Method not allowed", w)
		return
	}

	var event model.Event
	err := unmarshalEvent(r, &event)
	if err != nil {
		log.Println("Error unmarshalling request body on /update_event")
		sendError(http.StatusBadRequest, "Bad request body", w)
		return
	}

	if err = validateEvent(event); err != nil {
		log.Println(err)
		sendError(http.StatusBadRequest, err.Error(), w)
		return
	}

	updatedEvent, err := s.UpdateEvent(event)
	if err != nil {
		log.Println(err)
		sendError(http.StatusServiceUnavailable, "Internal server error", w)
		return
	}

	result := successResponse{updatedEvent}
	resp, err := json.Marshal(result)
	if err != nil {
		log.Println(err)
		sendError(http.StatusInternalServerError, "Internal server error", w)
		return
	}

	_, err = w.Write(resp)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Event updated successfully")
}

func (s *Server) deleteEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Println("Method not allowed on /delete_event")
		sendError(http.StatusMethodNotAllowed, "Method not allowed", w)
		return
	}

	var event model.Event
	err := unmarshalEvent(r, &event)
	if err != nil {
		log.Println("Error unmarshalling request body on /update_event")
		sendError(http.StatusBadRequest, "Bad request body", w)
		return
	}

	err = s.DeleteEvent(event)
	if err != nil {
		log.Println(err)
		sendError(http.StatusServiceUnavailable, err.Error(), w)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("Event deleted"))
	log.Println("Event deleted successfully")
}

func (s *Server) eventsForDay(w http.ResponseWriter, r *http.Request) {
	var userId uint
	var date string
	if r.URL.Query().Has("user_id") {
		id, err := strconv.Atoi(r.URL.Query().Get("user_id"))
		if err != nil {
			log.Println("Error parsing user_id parameter")
			sendError(http.StatusBadRequest, "Bad request", w)
			return
		}
		if id < 0 {
			log.Println("user_id cannot be negative")
			sendError(http.StatusBadRequest, "Bad request", w)
			return
		}
		userId = uint(id)
	} else {
		log.Println("Missing user_id parameter")
		sendError(http.StatusBadRequest, "Bad request", w)
		return
	}

	if r.URL.Query().Has("date") {
		date = r.URL.Query().Get("date")
	} else {
		log.Println("Missing date parameter")
		sendError(http.StatusBadRequest, "Bad request", w)
		return
	}

	day, err := time.Parse("2006-01-02", date)
	if err != nil {
		log.Println("Error parsing date parameter")
		sendError(http.StatusBadRequest, "Bad request", w)
		return
	}

	events, err := s.EventsForDay(userId, day)
	if err != nil {
		log.Println(err)
		sendError(http.StatusServiceUnavailable, err.Error(), w)
		return
	}

	result := successResponse{events}
	resp, err := json.Marshal(result)
	if err != nil {
		log.Println(err)
		sendError(http.StatusInternalServerError, err.Error(), w)
		return
	}

	_, err = w.Write(resp)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Events for day retrieved successfully")
}

func (s *Server) eventsForWeek(w http.ResponseWriter, r *http.Request) {
	var userId uint
	var date string
	if r.URL.Query().Has("user_id") {
		id, err := strconv.Atoi(r.URL.Query().Get("user_id"))
		if err != nil {
			log.Println("Error parsing user_id parameter")
			sendError(http.StatusBadRequest, "Bad request", w)
			return
		}
		if id < 0 {
			log.Println("user_id cannot be negative")
			sendError(http.StatusBadRequest, "Bad request", w)
			return
		}
		userId = uint(id)
	} else {
		log.Println("Missing user_id parameter")
		sendError(http.StatusBadRequest, "Bad request", w)
		return
	}

	if r.URL.Query().Has("date") {
		date = r.URL.Query().Get("date")
	} else {
		log.Println("Missing date parameter")
		sendError(http.StatusBadRequest, "Bad request", w)
		return
	}

	day, err := time.Parse("2006-01-02", date)
	if err != nil {
		log.Println("Error parsing date parameter")
		sendError(http.StatusBadRequest, "Bad request", w)
		return
	}

	events, err := s.EventsForWeek(userId, day)
	if err != nil {
		log.Println(err)
		sendError(http.StatusServiceUnavailable, err.Error(), w)
		return
	}

	result := successResponse{events}
	resp, err := json.Marshal(result)
	if err != nil {
		log.Println(err)
		sendError(http.StatusInternalServerError, err.Error(), w)
		return
	}

	_, err = w.Write(resp)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Events for week retrieved successfully")
}

func (s *Server) eventsForMonth(w http.ResponseWriter, r *http.Request) {
	var userId uint
	var date string
	if r.URL.Query().Has("user_id") {
		id, err := strconv.Atoi(r.URL.Query().Get("user_id"))
		if err != nil {
			log.Println("Error parsing user_id parameter")
			sendError(http.StatusBadRequest, "Bad request", w)
			return
		}
		if id < 0 {
			log.Println("user_id cannot be negative")
			sendError(http.StatusBadRequest, "Bad request", w)
			return
		}
		userId = uint(id)
	} else {
		log.Println("Missing user_id parameter")
		sendError(http.StatusBadRequest, "Bad request", w)
		return
	}

	if r.URL.Query().Has("date") {
		date = r.URL.Query().Get("date")
	} else {
		log.Println("Missing date parameter")
		sendError(http.StatusBadRequest, "Bad request", w)
		return
	}

	day, err := time.Parse("2006-01-02", date)
	if err != nil {
		log.Println("Error parsing date parameter")
		sendError(http.StatusBadRequest, "Bad request", w)
		return
	}

	year, month, _ := day.Date()
	events, err := s.EventsForMonth(userId, month, year)
	if err != nil {
		log.Println(err)
		sendError(http.StatusServiceUnavailable, err.Error(), w)
		return
	}

	result := successResponse{events}
	resp, err := json.Marshal(result)
	if err != nil {
		log.Println(err)
		sendError(http.StatusInternalServerError, err.Error(), w)
		return
	}

	_, err = w.Write(resp)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Events for month retrieved successfully")
}

func sendError(status int, errorString string, w http.ResponseWriter) {
	w.WriteHeader(status)
	errResp := errorResponse{errorString}
	resp, _ := json.Marshal(errResp)
	_, err := w.Write(resp)
	if err != nil {
		log.Println(err)
	}
}

func validateEvent(event model.Event) error {
	y, m, d := event.Date.Date()
	yn, mn, dn := time.Now().Date()
	if y < yn {
		return fmt.Errorf("date cannot be in the past")
	}
	if y == yn && m < mn {
		return fmt.Errorf("date cannot be in the past")
	}
	if y == yn && m == mn && d < dn {
		return fmt.Errorf("date cannot be in the past")
	}

	name := strings.Trim(event.Name, " ")
	if name == "" {
		return fmt.Errorf("name of event cannot be empty")
	}
	return nil
}

func unmarshalEvent(r *http.Request, event *model.Event) error {
	format := "2006-01-02"
	if r.Form.Has("id") {
		id, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			return err
		}
		if id < 0 {
			return fmt.Errorf("id cannot be negative")
		}
		event.Id = uint(id)
	}
	if r.Form.Has("name") {
		event.Name = r.FormValue("name")
	}
	if r.Form.Has("date") {
		date, err := time.Parse(format, r.FormValue("date"))
		if err != nil {
			return err
		}
		event.Date = date
	}
	if r.Form.Has("description") {
		event.Description = r.FormValue("description")
	}
	if r.Form.Has("user_id") {
		id, err := strconv.Atoi(r.FormValue("user_id"))
		if err != nil {
			return err
		}
		if id < 0 {
			return fmt.Errorf("user_id cannot be negative")
		}
		event.CreatorId = uint(id)
	}
	return nil
}

type errorResponse struct {
	Error string `json:"error"`
}

type successResponse struct {
	Result any `json:"result"`
}
