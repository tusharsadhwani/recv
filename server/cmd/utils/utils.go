package utils

import (
	"math"
	"math/rand"
	"time"
)

const LENGTH_OF_ROOM_CODE = 4

func GenerateRoomCode() int {
	rand.Seed(time.Hour.Nanoseconds())
	roomCode := 0
	for i := LENGTH_OF_ROOM_CODE - 1; i >= 0; i-- {
		roomCode += rand.Intn(10) * int(math.Pow(10, float64(i)))
	}
	return roomCode
}
