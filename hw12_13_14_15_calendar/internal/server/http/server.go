package internalhttp

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/brisk84/home_work/hw12_13_14_15_calendar/internal/storage"
)

type Server struct {
	addr       string
	httpServer http.Server
	logg       Logger
	app        Application
	ctx        context.Context
}

type Logger interface {
	Info(msg string)
	Error(msg string)
}

type Application interface {
	AddEvent(context.Context, storage.Event) error
	GetEvent(context.Context, string) (storage.Event, error)
	EditEvent(context.Context, storage.Event) error
	DeleteEvent(context.Context, string) error
	ListEvents(context.Context) ([]storage.Event, error)
}

func NewServer(logger Logger, app Application, addr string) *Server {
	server := &Server{
		addr: addr,
		logg: logger,
		app:  app,
	}
	return server
}

func (s *Server) AddEvent(writer http.ResponseWriter, request *http.Request) {
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		s.logg.Error(err.Error())
		return
	}
	ev := storage.Event{}
	err = json.Unmarshal(body, &ev)
	if err != nil {
		s.logg.Error("AddEvent:" + err.Error())
		return
	}
	s.logg.Info(ev.ID + ev.Title + ev.TimeStart.String() + ev.TimeEnd.String() + ev.UserID + ev.NotifyBefore.String() + ev.Description)

	err = s.app.AddEvent(s.ctx, ev)
	if err != nil {
		s.logg.Error("AddEvent:" + err.Error())
		writer.Write([]byte(err.Error()))
		return
	}
	writer.Write([]byte("ok"))
}

func (s *Server) GetEvent(writer http.ResponseWriter, request *http.Request) {
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		s.logg.Error(err.Error())
		return
	}
	ev := storage.Event{}
	err = json.Unmarshal(body, &ev)
	if err != nil {
		s.logg.Error("GetEvent:" + err.Error())
		return
	}
	s.logg.Info(ev.ID)

	ev, err = s.app.GetEvent(s.ctx, ev.ID)
	if err != nil {
		s.logg.Error("GetEvent:" + err.Error())
		writer.Write([]byte(err.Error()))
		return
	}
	data, err := json.Marshal(ev)
	if err != nil {
		s.logg.Error("GetEvent:" + err.Error())
		writer.Write([]byte(err.Error()))
		return
	}
	writer.Write(data)
}

func (s *Server) EditEvent(writer http.ResponseWriter, request *http.Request) {
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		s.logg.Error(err.Error())
		return
	}
	ev := storage.Event{}
	err = json.Unmarshal(body, &ev)
	if err != nil {
		s.logg.Error("EditEvent:" + err.Error())
		return
	}
	s.logg.Info(ev.ID)

	err = s.app.EditEvent(s.ctx, ev)
	if err != nil {
		s.logg.Error("EditEvent:" + err.Error())
		writer.Write([]byte(err.Error()))
		return
	}
	writer.Write([]byte("ok"))
}

func (s *Server) DeleteEvent(writer http.ResponseWriter, request *http.Request) {
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		s.logg.Error(err.Error())
		return
	}
	ev := storage.Event{}
	err = json.Unmarshal(body, &ev)
	if err != nil {
		s.logg.Error("GetEvent:" + err.Error())
		return
	}
	s.logg.Info(ev.ID)

	err = s.app.DeleteEvent(s.ctx, ev.ID)
	if err != nil {
		s.logg.Error("DeleteEvent:" + err.Error())
		writer.Write([]byte(err.Error()))
		return
	}
	writer.Write([]byte("ok"))
}

func (s *Server) ListEvents(writer http.ResponseWriter, request *http.Request) {
	events, err := s.app.ListEvents(s.ctx)
	if err != nil {
		s.logg.Error("ListEvents:" + err.Error())
		writer.Write([]byte(err.Error()))
		return
	}
	if events == nil {
		s.logg.Error("ListEvents: no events")
		writer.Write([]byte(err.Error()))
		return
	}
	data, err := json.Marshal(events)
	if err != nil {
		s.logg.Error("ListEvents:" + err.Error())
		writer.Write([]byte(err.Error()))
		return
	}
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
