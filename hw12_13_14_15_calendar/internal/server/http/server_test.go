package internalhttp

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brisk84/home_work/hw12_13_14_15_calendar/internal/app"
	"github.com/brisk84/home_work/hw12_13_14_15_calendar/internal/logger"
	"github.com/brisk84/home_work/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/brisk84/home_work/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/stretchr/testify/require"
)

func TestServer(t *testing.T) {
	var stor storage.Calendar
	logg := logger.New("stdout", "INFO")
	stor = memorystorage.New()
	calendar := app.New(logg, &stor)
	srv := NewServer(logg, calendar, "")

	ev1 := `{
		"id": "1",
		"title": "Event01",
		"time_start": "2022-01-23T17:45:00Z",
		"time_end": "2022-01-23T18:00:00Z",
		"description": "Event01 test",
		"user_id": "1",
		"notify_before": "2022-01-23T17:30:00Z"
		}`

	req := httptest.NewRequest(http.MethodPost, "/AddEvent", bytes.NewBuffer([]byte(ev1)))
	w := httptest.NewRecorder()
	srv.AddEvent(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)
	require.Equal(t, "ok", string(data))

	evId := `{ "id": "1" }`
	req = httptest.NewRequest(http.MethodPost, "/GetEvent", bytes.NewBuffer([]byte(evId)))
	w = httptest.NewRecorder()
	srv.GetEvent(w, req)
	res = w.Result()
	defer res.Body.Close()
	data, err = ioutil.ReadAll(res.Body)
	require.NoError(t, err)
	getEv := storage.Event{}
	err = json.Unmarshal(data, &getEv)
	require.NoError(t, err)
	myEv := storage.Event{}
	err = json.Unmarshal([]byte(ev1), &myEv)
	require.NoError(t, err)
	require.Equal(t, myEv, getEv)

	myEv.Title = myEv.Title + "[edited]"
	edEv, err := json.Marshal(myEv)
	require.NoError(t, err)
	req = httptest.NewRequest(http.MethodPost, "/EditEvent", bytes.NewBuffer(edEv))
	w = httptest.NewRecorder()
	srv.EditEvent(w, req)
	res = w.Result()
	defer res.Body.Close()
	data, err = ioutil.ReadAll(res.Body)
	require.NoError(t, err)
	require.Equal(t, "ok", string(data))

	req = httptest.NewRequest(http.MethodPost, "/GetEvent", bytes.NewBuffer([]byte(evId)))
	w = httptest.NewRecorder()
	srv.GetEvent(w, req)
	res = w.Result()
	defer res.Body.Close()
	data, err = ioutil.ReadAll(res.Body)
	require.NoError(t, err)
	err = json.Unmarshal(data, &getEv)
	require.NoError(t, err)
	require.Equal(t, myEv, getEv)

	ev2 := `{
		"id": "2",
		"title": "Event02",
		"time_start": "2022-01-23T17:45:00Z",
		"time_end": "2022-01-23T18:00:00Z",
		"description": "Event02 test",
		"user_id": "1",
		"notify_before": "2022-01-23T17:30:00Z"
		}`
	req = httptest.NewRequest(http.MethodPost, "/AddEvent", bytes.NewBuffer([]byte(ev2)))
	w = httptest.NewRecorder()
	srv.AddEvent(w, req)
	res = w.Result()
	defer res.Body.Close()
	data, err = ioutil.ReadAll(res.Body)
	require.NoError(t, err)
	require.Contains(t, string(data), "date is busy")

	ev3 := `{
		"id": "3",
		"title": "Event03",
		"time_start": "2022-01-24T17:45:00Z",
		"time_end": "2022-01-24T18:00:00Z",
		"description": "Event03 test",
		"user_id": "1",
		"notify_before": "2022-01-24T17:30:00Z"
		}`
	req = httptest.NewRequest(http.MethodPost, "/AddEvent", bytes.NewBuffer([]byte(ev3)))
	w = httptest.NewRecorder()
	srv.AddEvent(w, req)
	res = w.Result()
	defer res.Body.Close()
	data, err = ioutil.ReadAll(res.Body)
	require.NoError(t, err)
	require.Equal(t, "ok", string(data))

	req = httptest.NewRequest(http.MethodPost, "/ListEvents", nil)
	w = httptest.NewRecorder()
	srv.ListEvents(w, req)
	res = w.Result()
	defer res.Body.Close()
	data, err = ioutil.ReadAll(res.Body)
	require.NoError(t, err)

	ev3s := storage.Event{}
	json.Unmarshal([]byte(ev3), &ev3s)
	evs := []storage.Event{}
	evs = append(evs, myEv)
	evs = append(evs, ev3s)
	evsr := []storage.Event{}
	err = json.Unmarshal(data, &evsr)
	require.NoError(t, err)
	require.Equal(t, evs, evsr)

	req = httptest.NewRequest(http.MethodPost, "/DeleteEvent", bytes.NewBuffer([]byte(evId)))
	w = httptest.NewRecorder()
	srv.DeleteEvent(w, req)
	res = w.Result()
	defer res.Body.Close()
	data, err = ioutil.ReadAll(res.Body)
	require.NoError(t, err)
	require.Equal(t, "ok", string(data))

	req = httptest.NewRequest(http.MethodPost, "/ListEvents", nil)
	w = httptest.NewRecorder()
	srv.ListEvents(w, req)
	res = w.Result()
	defer res.Body.Close()
	data, err = ioutil.ReadAll(res.Body)
	require.NoError(t, err)

	evs2 := []storage.Event{}
	evs2 = append(evs2, ev3s)
	err = json.Unmarshal(data, &evsr)
	require.NoError(t, err)
	require.Equal(t, evs2, evsr)

}
