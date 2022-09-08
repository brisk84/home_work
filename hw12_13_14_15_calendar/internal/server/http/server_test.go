package internalhttp

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brisk84/home_work/hw12_13_14_15_calendar/internal/app"
	"github.com/brisk84/home_work/hw12_13_14_15_calendar/internal/logger"
	"github.com/brisk84/home_work/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/brisk84/home_work/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/stretchr/testify/require"
)

var (
	calendar *app.App
	logg     *logger.Logger
	srv      *Server
	stor     storage.Calendar
	evID     string
	ev1      string
	ev1ed    string
)

func AddEvent(t *testing.T) {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/AddEvent", bytes.NewBuffer([]byte(ev1)))
	w := httptest.NewRecorder()
	srv.AddEvent(w, req)
	res := w.Result()
	require.Equal(t, http.StatusOK, res.StatusCode)
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	require.Equal(t, "ok", string(data))
}

func GetEvent(t *testing.T) {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/GetEvent", bytes.NewBuffer([]byte(evID)))
	w := httptest.NewRecorder()
	srv.GetEvent(w, req)
	res := w.Result()
	require.Equal(t, http.StatusOK, res.StatusCode)
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	getEv := storage.Event{}
	err = json.Unmarshal(data, &getEv)
	require.NoError(t, err)
	myEv := storage.Event{}
	err = json.Unmarshal([]byte(ev1), &myEv)
	require.NoError(t, err)
	require.Equal(t, myEv, getEv)
}

func EditEvent(t *testing.T) {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/EditEvent", bytes.NewBuffer([]byte(ev1ed)))
	w := httptest.NewRecorder()
	srv.EditEvent(w, req)
	res := w.Result()
	require.Equal(t, http.StatusOK, res.StatusCode)
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	require.Equal(t, "ok", string(data))
}

func GetEvent2(t *testing.T) {
	t.Helper()
	getEv := storage.Event{}
	evEd := storage.Event{}
	req := httptest.NewRequest(http.MethodPost, "/GetEvent", bytes.NewBuffer([]byte(evID)))
	w := httptest.NewRecorder()
	srv.GetEvent(w, req)
	res := w.Result()
	require.Equal(t, http.StatusOK, res.StatusCode)
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	err = json.Unmarshal(data, &getEv)
	require.NoError(t, err)
	err = json.Unmarshal([]byte(ev1ed), &evEd)
	require.NoError(t, err)
	require.Equal(t, evEd, getEv)
}

func TestServer(t *testing.T) {
	evID = `{ "id": "1" }`

	ev1 = `{
		"id": "1",
		"title": "Event01",
		"timeStart": "2022-01-23T17:45:00Z",
		"timeEnd": "2022-01-23T18:00:00Z",
		"description": "Event01 test",
		"userId": "1",
		"notifyBefore": "2022-01-23T17:30:00Z"
		}`

	ev1ed = `{
			"id": "1",
			"title": "Event01[edited]",
			"timeStart": "2022-01-23T17:45:00Z",
			"timeEnd": "2022-01-23T18:00:00Z",
			"description": "Event01 test",
			"userId": "1",
			"notifyBefore": "2022-01-23T17:30:00Z"
		}`

	logg = logger.New("stdout", "INFO")
	stor = memorystorage.New()
	calendar = app.New(logg, stor)
	srv = NewServer(logg, calendar, "")

	AddEvent(t)
	GetEvent(t)
	EditEvent(t)

	evEd := storage.Event{}
	err := json.Unmarshal([]byte(ev1ed), &evEd)
	require.NoError(t, err)

	GetEvent2(t)

	ev2 := `{
		"id": "2",
		"title": "Event02",
		"timeStart": "2022-01-23T17:45:00Z",
		"timeEnd": "2022-01-23T18:00:00Z",
		"description": "Event02 test",
		"userId": "1",
		"notifyBefore": "2022-01-23T17:30:00Z"
		}`
	req := httptest.NewRequest(http.MethodPost, "/AddEvent", bytes.NewBuffer([]byte(ev2)))
	w := httptest.NewRecorder()
	srv.AddEvent(w, req)
	res := w.Result()
	require.Equal(t, http.StatusConflict, res.StatusCode)
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	require.Contains(t, string(data), "date is busy")

	ev3 := `{
		"id": "3",
		"title": "Event03",
		"timeStart": "2022-01-24T17:45:00Z",
		"timeEnd": "2022-01-24T18:00:00Z",
		"description": "Event03 test",
		"userId": "1",
		"notifyBefore": "2022-01-24T17:30:00Z"
		}`
	req = httptest.NewRequest(http.MethodPost, "/AddEvent", bytes.NewBuffer([]byte(ev3)))
	w = httptest.NewRecorder()
	srv.AddEvent(w, req)
	res = w.Result()
	require.Equal(t, http.StatusOK, res.StatusCode)
	defer res.Body.Close()
	data, err = io.ReadAll(res.Body)
	require.NoError(t, err)
	require.Equal(t, "ok", string(data))

	req = httptest.NewRequest(http.MethodPost, "/ListEvents", nil)
	w = httptest.NewRecorder()
	srv.ListEvents(w, req)
	res = w.Result()
	require.Equal(t, http.StatusOK, res.StatusCode)
	defer res.Body.Close()
	data, err = io.ReadAll(res.Body)
	require.NoError(t, err)

	ev3s := storage.Event{}
	json.Unmarshal([]byte(ev3), &ev3s)
	evs := []storage.Event{}
	evs = append(evs, evEd)
	evs = append(evs, ev3s)
	evsr := []storage.Event{}
	err = json.Unmarshal(data, &evsr)
	require.NoError(t, err)
	require.Equal(t, evs, evsr)

	req = httptest.NewRequest(http.MethodPost, "/DeleteEvent", bytes.NewBuffer([]byte(evID)))
	w = httptest.NewRecorder()
	srv.DeleteEvent(w, req)
	res = w.Result()
	require.Equal(t, http.StatusOK, res.StatusCode)
	defer res.Body.Close()
	data, err = io.ReadAll(res.Body)
	require.NoError(t, err)
	require.Equal(t, "ok", string(data))

	req = httptest.NewRequest(http.MethodPost, "/ListEvents", nil)
	w = httptest.NewRecorder()
	srv.ListEvents(w, req)
	res = w.Result()
	require.Equal(t, http.StatusOK, res.StatusCode)
	defer res.Body.Close()
	data, err = io.ReadAll(res.Body)
	require.NoError(t, err)

	evs2 := []storage.Event{}
	evs2 = append(evs2, ev3s)
	err = json.Unmarshal(data, &evsr)
	require.NoError(t, err)
	require.Equal(t, evs2, evsr)
}
