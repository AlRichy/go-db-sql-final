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
	db, err := // настройте подключение к БД
	store := NewParcelStore(db)
	parcel := getTestParcel()
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()

	// add
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора

	store := NewParcelStore(db)
	parcel := getTestParcel()
	// get
	// получите только что добавленную посылку, убедитесь в отсутствии ошибки
	// проверьте, что значения всех полей в полученном объекте совпадают со значениями полей в переменной parcel

	p, err := store.Add(parcel)
	parcel.Number = p
	require.NoError(t, err)
	require.NotEmpty(t, p)
	// delete
	// удалите добавленную посылку, убедитесь в отсутствии ошибки
	// проверьте, что посылку больше нельзя получить из БД
	gp, err := store.Get(p)
	require.NoError(t, err)
	require.Equal(t, parcel, gp)
	err = store.Delete(p)
	require.NoError(t, err)

	_, err = store.Get(p)
	require.ErrorIs(t, sql.ErrNoRows, err)
}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	// prepare
	db, err := // настройте подключение к БД
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()

	// add
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
	store := NewParcelStore(db)
	parcel := getTestParcel()

	p, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotEmpty(t, p)
	// set address
	// обновите адрес, убедитесь в отсутствии ошибки
	newAddress := "new test address"

	err = store.SetAddress(p, newAddress)
	require.NoError(t, err)
	// check
	// получите добавленную посылку и убедитесь, что адрес обновился
	gp, err := store.Get(p)
	require.NoError(t, err)
	require.Equal(t, newAddress, gp.Address)
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	// prepare
	db, err := // настройте подключение к БД

	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()
	// add
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
	store := NewParcelStore(db)
	parcel := getTestParcel()

	// set status
	// обновите статус, убедитесь в отсутствии ошибки

	p, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotEmpty(t, p)
	// check
	// получите добавленную посылку и убедитесь, что статус обновился
	err = store.SetStatus(p, ParcelStatusSent)
	require.NoError(t, err)

	gp, err := store.Get(p)
	require.NoError(t, err)
	require.Equal(t, ParcelStatusSent, gp.Status)
}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
	// prepare
	db, err := // настройте подключение к БД
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)

	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := map[int]Parcel{}
	// задаём всем посылкам один и тот же идентификатор клиента
	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client

	// add
	for i := 0; i < len(parcels); i++ {
		id, err := // добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
		id, err := store.Add(parcels[i])

		require.NoError(t, err)
		require.NotEmpty(t, id)
		// обновляем идентификатор добавленной у посылки
		id, err := store.Add(parcels[i])
		require.NoError(t, err)
		require.NotEmpty(t, id)
		parcels[i].Number = id

		// сохраняем добавленную посылку в структуру map, чтобы её можно было легко достать по идентификатору посылки
		parcelMap[id] = parcels[i]
	}

	// get by client
	storedParcels, err := // получите список посылок по идентификатору клиента, сохранённого в переменной client
	// получите список посылок по идентификатору клиента, сохранённого в переменной client
	// убедитесь в отсутствии ошибки
	// убедитесь, что количество полученных посылок совпадает с количеством добавленных
	var storedParcels []Parcel
	storedParcels, err = store.GetByClient(client)
	require.NoError(t, err)
	assert.Equal(t, len(parcels), len(storedParcels))

	// check
	for _, parcel := range storedParcels {
		// в parcelMap лежат добавленные посылки, ключ - идентификатор посылки, значение - сама посылка
		// убедитесь, что все посылки из storedParcels есть в parcelMap
		// убедитесь, что значения полей полученных посылок заполнены верно
		_, ok := parcelMap[parcel.Number]
		require.True(t, ok)
		assert.Equal(t, parcelMap[parcel.Number], parcel)
	}
}