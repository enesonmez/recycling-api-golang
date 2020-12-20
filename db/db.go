package db

// Composition(birleşim) inheritance yerine kullanılır.
type Db struct {
	dbHost     string
	dbPort     int
	dbUserName string
	dbPass     string
	dbName     string
}
