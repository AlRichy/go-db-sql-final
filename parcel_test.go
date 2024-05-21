package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	// randSource источник псевдо случайных чисел.
	// Для повышения уникальности в качестве seed
	// используется текущее время в unix формате (в виде числа)
	randSource = rand.NewSource(time.Now().UnixNano())
	// randRange использует randSource для генерации случайных чисел
	randRange = rand.New(randSource)
)

// getTestParcel возвращает тестовую посылку
func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

func ConnectDb() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "tracker.db")
	return db, err
}

// TestAddGetDelete проверяет добавление, получение и удаление посылки
func TestAddGetDelete(t *testing.T) {
	// prepare
	db, err := ConnectDb()
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	p, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotEmpty(t, p)

	parcel.Number = p

	gp, err := store.Get(p)
	require.NoError(t, err)
	require.Equal(t, parcel, gp)

	err = store.Delete(p)
	require.NoError(t, err)

	_, err = store.Get(p)
	require.ErrorIs(t, err, sql.ErrNoRows)
}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	// prepare
	db, err := ConnectDb()
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	p, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotEmpty(t, p)

	newAddress := "new test address"
	err = store.SetAddress(p, newAddress)
	require.NoError(t, err)

	gp, err := store.Get(p)
	require.NoError(t, err)
	require.Equal(t, newAddress, gp.Address)
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	// prepare
	db, err := ConnectDb()
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	p, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotEmpty(t, p)

	err = store.SetStatus(p, ParcelStatusSent)
	require.NoError(t, err)

	gp, err := store.Get(p)
	require.NoError(t, err)
	require.Equal(t, ParcelStatusSent, gp.Status)
}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
	// prepare
	db, err := ConnectDb()
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)

	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := map[int]Parcel{}

	client := randRange.Intn(10_000_000)
	for i := range parcels {
		parcels[i].Client = client
		id, err := store.Add(parcels[i])
		require.NoError(t, err)
		require.NotEmpty(t, id)
		parcels[i].Number = id
		parcelMap[id] = parcels[i]
	}

	storedParcels, err := store.GetByClient(client)
	require.NoError(t, err)
	assert.Equal(t, len(parcels), len(storedParcels))

	for _, parcel := range storedParcels {
		expectedParcel, ok := parcelMap[parcel.Number]
		require.True(t, ok)
		assert.Equal(t, expectedParcel, parcel)
	}
}
