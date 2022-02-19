package memorystorage

import (
	"context"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/brisk84/home_work/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	st := New()

	id1 := uuid.New().String()
	id2 := uuid.New().String()
	id3 := uuid.New().String()

	ctx := context.Background()

	ev1 := storage.Event{
		ID:           id1,
		Title:        "First event",
		TimeStart:    time.Date(2022, 1, 1, 13, 0, 0, 0, time.Local),
		TimeEnd:      time.Date(2022, 1, 1, 14, 0, 0, 0, time.Local),
		Description:  "First event test",
		UserID:       uuid.New().String(),
		NotifyBefore: time.Date(2022, 1, 1, 13, 0, 0, 0, time.Local).Add(-15 * time.Minute),
	}
	ev2 := storage.Event{
		ID:           id1,
		Title:        "Second event",
		TimeStart:    time.Date(2022, 1, 1, 14, 0, 0, 0, time.Local),
		TimeEnd:      time.Date(2022, 1, 1, 15, 0, 0, 0, time.Local),
		Description:  "Second event test",
		UserID:       uuid.New().String(),
		NotifyBefore: time.Date(2022, 1, 1, 14, 0, 0, 0, time.Local).Add(-15 * time.Minute),
	}
	ev3 := storage.Event{
		ID:           id3,
		Title:        "Duplicate event",
		TimeStart:    time.Date(2022, 1, 1, 15, 0, 0, 0, time.Local),
		TimeEnd:      time.Date(2022, 1, 1, 16, 0, 0, 0, time.Local),
		Description:  "Third event test",
		UserID:       uuid.New().String(),
		NotifyBefore: time.Date(2022, 1, 1, 14, 0, 0, 0, time.Local).Add(-15 * time.Minute),
	}

	err := st.AddEvent(ctx, ev1)
	require.NoError(t, err)

	err = st.AddEvent(ctx, ev1)
	errDate := storage.ErrDateBusy
	errUUID := storage.ErrUUIDBusy
	require.ErrorAs(t, err, &errDate)
	require.ErrorAs(t, err, &errUUID)

	err = st.AddEvent(ctx, ev2)
	require.ErrorIs(t, err, storage.ErrUUIDBusy)

	ev2.ID = id2
	err = st.AddEvent(ctx, ev2)
	require.NoError(t, err)

	err = st.AddEvent(ctx, ev3)
	require.NoError(t, err)

	ev01, err := st.GetEvent(ctx, id1)
	require.NoError(t, err)
	require.Equal(t, ev1, ev01)

	ev02, err := st.GetEvent(ctx, id2)
	require.NoError(t, err)
	require.Equal(t, ev2, ev02)

	ev03, err := st.GetEvent(ctx, id3)
	require.NoError(t, err)
	require.Equal(t, ev3, ev03)

	ev01.Title = "First event (edited)"
	err = st.EditEvent(ctx, ev01)
	require.NoError(t, err)

	ev001, err := st.GetEvent(ctx, id1)
	require.NoError(t, err)
	require.NotEqual(t, ev1, ev001)

	evs1 := st.ListEvents(ctx)
	require.Equal(t, 3, len(evs1))

	err = st.DeleteEvent(ctx, id2)
	require.NoError(t, err)

	err = st.DeleteEvent(ctx, id2)
	require.ErrorIs(t, err, storage.ErrNotFound)

	ev2.Title = "Second event (edited)"
	err = st.EditEvent(ctx, ev2)
	require.ErrorIs(t, err, storage.ErrNotFound)

	ev3.TimeStart = ev1.TimeStart
	err = st.EditEvent(ctx, ev3)
	require.ErrorIs(t, err, storage.ErrDateBusy)

	evs2 := st.ListEvents(ctx)
	require.Equal(t, 2, len(evs2))
}

func TestStorageConcurency(t *testing.T) {
	st := New()

	wg := sync.WaitGroup{}
	threadsCount := 1000

	ctx := context.Background()

	wg.Add(threadsCount)
	for i := 0; i < threadsCount; i++ {
		go func(i int) {
			defer wg.Done()
			ev1 := storage.Event{
				ID:           uuid.New().String(),
				Title:        "Event #" + strconv.Itoa(i),
				TimeStart:    time.Now().AddDate(0, 0, i),
				TimeEnd:      time.Date(2022, 1, 1, 14, 0, 0, 0, time.Local),
				Description:  "Event test",
				UserID:       uuid.New().String(),
				NotifyBefore: time.Date(2022, 1, 1, 13, 0, 0, 0, time.Local).Add(-15 * time.Minute),
			}
			err := st.AddEvent(ctx, ev1)
			require.NoError(t, err)
		}(i)
	}
	wg.Wait()
	evs := st.ListEvents(ctx)
	require.Equal(t, threadsCount, len(evs))

	wg.Add(2 * threadsCount)
	for i := 0; i < threadsCount; i++ {
		go func(i int) {
			defer wg.Done()
			ev := evs[i]
			ev.Title += "[edited]"
			err := st.EditEvent(ctx, ev)
			require.NoError(t, err)
		}(i)
		go func(i int) {
			defer wg.Done()
			_, err := st.GetEvent(ctx, evs[i].ID)
			require.NoError(t, err)
		}(i)
	}
	wg.Wait()

	wg.Add(threadsCount)
	for i := 0; i < threadsCount; i++ {
		go func(i int) {
			defer wg.Done()
			ev, err := st.GetEvent(ctx, evs[i].ID)
			require.NoError(t, err)
			require.NotEqual(t, ev, evs[i])
		}(i)
	}
	wg.Wait()

	wg.Add(threadsCount)
	for i := 0; i < threadsCount; i++ {
		go func(i int) {
			defer wg.Done()
			st.DeleteEvent(ctx, evs[i].ID)
		}(i)
	}
	wg.Wait()
	evs = st.ListEvents(ctx)
	require.Equal(t, 0, len(evs))
}
