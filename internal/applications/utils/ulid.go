package utils

import (
	"math/rand"
	"time"

	"github.com/oklog/ulid"
)

func GenULID(t time.Time) string {
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	return ulid.MustNew(ulid.Timestamp(t), entropy).String()
}
