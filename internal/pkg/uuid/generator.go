package uuid

import "github.com/google/uuid"

func Next() uuid.UUID {
	res, err := uuid.NewV7()
	if err != nil {
		panic(err)
	}
	return res
}
