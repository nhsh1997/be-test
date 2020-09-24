package stores

type MainStore struct {
}

// NewMainStore creates a new instance of MainStore, contains whole common functions
// for a service

func NewMainStore() *MainStore {
	return &MainStore{}
}
