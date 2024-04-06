package store

type Store struct {
	driver Driver
}

func NewStore(driver Driver) *Store {
	return &Store{
		driver: driver,
	}
}
