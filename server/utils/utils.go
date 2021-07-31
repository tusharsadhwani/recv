package utils

import (
	"math"
	"math/rand"
	"time"

	. "github.com/tusharsadhwani/recv/constants"
)

var seed = rand.NewSource(time.Now().UnixNano())
var rng = rand.New(seed)

func GenerateRoomCode() RoomID {
	lowerLimit := int(math.Pow10(RoomCodeLength - 1))
	upperLimit := int(math.Pow10(RoomCodeLength))
	roomCode := lowerLimit + rng.Intn(upperLimit-lowerLimit)
	return RoomID(roomCode)
}
