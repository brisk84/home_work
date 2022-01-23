package sqlstorage

// func TestStorage(t *testing.T) {
// 	require.True(t, true)
// 	return
// 	ConnStr := "postgresql://localhost/calendar?user=postgres&password=sqlSync24&sslmode=disable"
// 	st := New("postgres", ConnStr, 20)

// 	err := st.Connect(context.TODO())
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	defer st.Close(context.TODO())

// 	id1 := uuid.New().String()
// 	id2 := uuid.New().String()
// 	id3 := uuid.New().String()

// 	ev1 := storage.Event{
// 		ID:           id1,
// 		Title:        "First event",
// 		TimeStart:    time.Date(2022, 1, 1, 13, 0, 0, 0, time.Local),
// 		TimeEnd:      time.Date(2022, 1, 1, 14, 0, 0, 0, time.Local),
// 		Description:  "First event test",
// 		UserID:       uuid.New().String(),
// 		NotifyBefore: time.Date(2022, 1, 1, 13, 0, 0, 0, time.Local).Add(-15 * time.Minute),
// 	}
// 	ev2 := storage.Event{
// 		ID:           id1,
// 		Title:        "Second event",
// 		TimeStart:    time.Date(2022, 1, 1, 14, 0, 0, 0, time.Local),
// 		TimeEnd:      time.Date(2022, 1, 1, 15, 0, 0, 0, time.Local),
// 		Description:  "Second event test",
// 		UserID:       uuid.New().String(),
// 		NotifyBefore: time.Date(2022, 1, 1, 14, 0, 0, 0, time.Local).Add(-15 * time.Minute),
// 	}
// 	ev3 := storage.Event{
// 		ID:           id3,
// 		Title:        "Duplicate event",
// 		TimeStart:    time.Date(2022, 1, 1, 15, 0, 0, 0, time.Local),
// 		TimeEnd:      time.Date(2022, 1, 1, 16, 0, 0, 0, time.Local),
// 		Description:  "Third event test",
// 		UserID:       uuid.New().String(),
// 		NotifyBefore: time.Date(2022, 1, 1, 14, 0, 0, 0, time.Local).Add(-15 * time.Minute),
// 	}

// 	st.ClearCalendar()

// 	err = st.AddEvent(ev1)
// 	require.NoError(t, err)

// 	err = st.AddEvent(ev1)
// 	errDate := storage.ErrDateBusy
// 	errUUID := storage.ErrUUIDBusy
// 	require.ErrorAs(t, err, &errDate)
// 	require.ErrorAs(t, err, &errUUID)

// 	err = st.AddEvent(ev2)
// 	require.ErrorIs(t, err, storage.ErrUUIDBusy)

// 	ev2.ID = id2
// 	err = st.AddEvent(ev2)
// 	require.NoError(t, err)

// 	err = st.AddEvent(ev3)
// 	require.NoError(t, err)

// 	ev01, err := st.GetEvent(id1)
// 	require.NoError(t, err)
// 	require.Equal(t, ev1, ev01)

// 	ev02, err := st.GetEvent(id2)
// 	require.NoError(t, err)
// 	require.Equal(t, ev2, ev02)

// 	ev03, err := st.GetEvent(id3)
// 	require.NoError(t, err)
// 	require.Equal(t, ev3, ev03)

// 	ev01.Title = "First event (edited)"
// 	err = st.EditEvent(ev01)
// 	require.NoError(t, err)

// 	ev001, err := st.GetEvent(id1)
// 	require.NoError(t, err)
// 	require.NotEqual(t, ev1, ev001)

// 	evs1 := st.ListEvents()
// 	require.Equal(t, 3, len(evs1))

// 	err = st.DeleteEvent(id2)
// 	require.NoError(t, err)

// 	err = st.DeleteEvent(id2)
// 	require.ErrorIs(t, err, storage.ErrNotFound)

// 	ev2.Title = "Second event (edited)"
// 	err = st.EditEvent(ev2)
// 	require.ErrorIs(t, err, storage.ErrNotFound)

// 	ev3.TimeStart = ev1.TimeStart
// 	err = st.EditEvent(ev3)
// 	require.ErrorIs(t, err, storage.ErrDateBusy)

// 	evs2 := st.ListEvents()
// 	require.Equal(t, 2, len(evs2))

// 	st.ClearCalendar()
// }

// func TestStorageConcurency(t *testing.T) {
// 	require.True(t, true)
// 	return
// 	ConnStr := "postgresql://localhost/calendar?user=postgres&password=sqlSync24&sslmode=disable"
// 	st := New("postgres", ConnStr, 20)

// 	err := st.Connect(context.TODO())
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	defer st.Close(context.TODO())

// 	wg := sync.WaitGroup{}
// 	threadsCount := 100

// 	wg.Add(threadsCount)
// 	for i := 0; i < threadsCount; i++ {
// 		go func(i int) {
// 			defer wg.Done()
// 			ev1 := storage.Event{
// 				ID:           uuid.New().String(),
// 				Title:        "Event #" + strconv.Itoa(i),
// 				TimeStart:    time.Now().AddDate(0, 0, i),
// 				TimeEnd:      time.Date(2022, 1, 1, 14, 0, 0, 0, time.Local),
// 				Description:  "Event test",
// 				UserID:       uuid.New().String(),
// 				NotifyBefore: time.Date(2022, 1, 1, 13, 0, 0, 0, time.Local).Add(-15 * time.Minute),
// 			}
// 			err := st.AddEvent(ev1)
// 			require.NoError(t, err)
// 		}(i)
// 	}
// 	wg.Wait()
// 	evs := st.ListEvents()
// 	require.Equal(t, threadsCount, len(evs))

// 	wg.Add(2 * threadsCount)
// 	for i := 0; i < threadsCount; i++ {
// 		go func(i int) {
// 			defer wg.Done()
// 			ev := evs[i]
// 			ev.Title += "[edited]"
// 			err := st.EditEvent(ev)
// 			require.NoError(t, err)
// 		}(i)
// 		go func(i int) {
// 			defer wg.Done()
// 			_, err := st.GetEvent(evs[i].ID)
// 			require.NoError(t, err)
// 		}(i)
// 	}
// 	wg.Wait()

// 	wg.Add(threadsCount)
// 	for i := 0; i < threadsCount; i++ {
// 		go func(i int) {
// 			defer wg.Done()
// 			ev, err := st.GetEvent(evs[i].ID)
// 			require.NoError(t, err)
// 			require.NotEqual(t, ev, evs[i])
// 		}(i)
// 	}
// 	wg.Wait()

// 	wg.Add(threadsCount)
// 	for i := 0; i < threadsCount; i++ {
// 		go func(i int) {
// 			defer wg.Done()
// 			st.DeleteEvent(evs[i].ID)
// 		}(i)
// 	}
// 	wg.Wait()
// 	evs = st.ListEvents()
// 	require.Equal(t, 0, len(evs))
// }
