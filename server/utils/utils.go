package utils

import (
	"math"
	"math/rand"
	"time"
)

const RoomCodeLength = 4

var seed = rand.NewSource(time.Now().UnixNano())
var rng = rand.New(seed)

func GenerateRoomCode() int {
	lowerLimit := int(math.Pow10(RoomCodeLength - 1))
	upperLimit := int(math.Pow10(RoomCodeLength))
	roomCode := lowerLimit + rng.Intn(upperLimit-lowerLimit)
	return roomCode
}
