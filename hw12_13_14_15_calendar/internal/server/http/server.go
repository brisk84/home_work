package internalhttp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/brisk84/home_work/hw12_13_14_15_calendar/internal/app"
	"github.com/brisk84/home_work/hw12_13_14_15_calendar/internal/storage"
)

type Server struct {
	addr       string
	httpServer http.Server
	logg       Logger
	appl       *app.App
	ctx        context.Context
}

type Logger interface {
	Info(msg string)
	Error(msg string)
}

func NewServer(logger Logger, appl *app.App, addr string) *Server {
	server := &Server{
		addr: addr,
		logg: logger,
		appl: appl,
	}
	return server
}

func (s *Server) AddEvent(writer http.ResponseWriter, request *http.Request) {
	body, err := io.ReadAll(request.Body)
	if err != nil {
		s.logg.Error(err.Error())
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte(err.Error()))
		return
	}
	ev := storage.Event{}
	err = json.Unmarshal(body, &ev)
	if err != nil {
		s.logg.Error("AddEvent:" + err.Error())
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte(err.Error()))
		return
	}
	s.logg.Info(ev.ID + ev.Title + ev.TimeStart.String() + ev.TimeEnd.String() + ev.UserID +
		ev.NotifyBefore.String() + ev.Description)

	err = s.appl.AddEvent(s.ctx, ev)
	if err != nil {
		s.logg.Error("AddEvent:" + err.Error())
		writer.WriteHeader(http.StatusConflict)
		writer.Write([]byte(err.Error()))
		return
	}
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("ok"))
}

func (s *Server) GetEvent(writer http.ResponseWriter, request *http.Request) {
	body, err := io.ReadAll(request.Body)
	if err != nil {
		s.logg.Error(err.Error())
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte(err.Error()))
		return
	}
	ev := storage.Event{}
	err = json.Unmarshal(body, &ev)
	if err != nil {
		s.logg.Error("GetEvent:" + err.Error())
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte(err.Error()))
		return
	}
	s.logg.Info(ev.ID)

	ev, err = s.appl.GetEvent(s.ctx, ev.ID)
	if err != nil {
		s.logg.Error("GetEvent:" + err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(err.Error()))
		return
	}
	data, err := json.Marshal(ev)
	if err != nil {
		s.logg.Error("GetEvent:" + err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(err.Error()))
		return
	}
	writer.WriteHeader(http.StatusOK)
	writer.Write(data)
}

func (s *Server) EditEvent(writer http.ResponseWriter, request *http.Request) {
	body, err := io.ReadAll(request.Body)
	if err != nil {
		s.logg.Error(err.Error())
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte(err.Error()))
		return
	}
	ev := storage.Event{}
	err = json.Unmarshal(body, &ev)
	if err != nil {
		s.logg.Error("EditEvent:" + err.Error())
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte(err.Error()))
		return
	}
	s.logg.Info(ev.ID)

	err = s.appl.EditEvent(s.ctx, ev)
	if err != nil {
		s.logg.Error("EditEvent:" + err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(err.Error()))
		return
	}
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("ok"))
}

func (s *Server) DeleteEvent(writer http.ResponseWriter, request *http.Request) {
	body, err := io.ReadAll(request.Body)
	if err != nil {
		s.logg.Error(err.Error())
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte(err.Error()))
		return
	}
	ev := storage.Event{}
	err = json.Unmarshal(body, &ev)
	if err != nil {
		s.logg.Error("GetEvent:" + err.Error())
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte(err.Error()))
		return
	}
	s.logg.Info(ev.ID)

	err = s.appl.DeleteEvent(s.ctx, ev.ID)
	if err != nil {
		s.logg.Error("DeleteEvent:" + err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(err.Error()))
		return
	}
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("ok"))
}

func (s *Server) ListEvents(writer http.ResponseWriter, request *http.Request) {
	events, err := s.appl.ListEvents(s.ctx)
	if err != nil {
		s.logg.Error("ListEvents:" + err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(err.Error()))
		return
	}
	if events == nil {
		s.logg.Error("ListEvents: no events")
		writer.WriteHeader(http.StatusNotFound)
		writer.Write([]byte("no events"))
		return
	}
	data, err := json.Marshal(events)
	if err != nil {
		s.logg.Error("ListEvents:" + err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(err.Error()))
		return
	}
	writer.WriteHeader(http.StatusOK)
	writer.Write(data)
}

func (s *Server) Start(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", func(w http.ResponseWriter, req *http.Request) {
		ip, _, _ := net.SplitHostPort(req.RemoteAddr)
		ret := fmt.Sprintf("%s [%s] %s %s %s", ip, time.Now(), req.Method, req.URL.Path, req.UserAgent())
		fmt.Fprintf(w, "<p>%s</p>", ret)
	})
	mux.HandleFunc("/AddEvent", s.AddEvent)
	mux.HandleFunc("/GetEvent", s.GetEvent)
	mux.HandleFunc("/EditEvent", s.EditEvent)
	mux.HandleFunc("/DeleteEvent", s.DeleteEvent)
	mux.HandleFunc("/ListEvents", s.ListEvents)

	s.httpServer = http.Server{
		Addr:    s.addr,
		Handler: loggingMiddleware(mux, s.logg),
	}
	s.ctx = ctx

	s.logg.Info("Start http server")
	err := s.httpServer.ListenAndServe()
	<-ctx.Done()
	return err
}

func (s *Server) Stop(ctx context.Context) error {
	s.logg.Info("Stop http server")
	err := s.httpServer.Shutdown(ctx)
	return err
}
