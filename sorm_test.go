package sorm

import (
	"testing"
)

func TestOpen(t *testing.T) {
	t.Parallel()

	var db, err = Open(
		"postgres",
		"user=postgres password=GYAaTx57jJput dbname=postgres port=55432 sslmode=disable",
	)
	if err != nil || db.GetRawDB().Ping() != nil {
		panic(err)
	}

	if err := db.Close(); err != nil {
		panic(err)
	}
}
