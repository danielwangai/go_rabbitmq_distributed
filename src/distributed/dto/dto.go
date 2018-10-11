package dto

// holds definitions for any type sent through the message broker

import (
	"encoding/gob"
	"time"
)

type SensorMessage struct {
	Name      string
	Value     float64
	Timestamp time.Time
}

func init() {
	gob.Register(SensorMessage{})
}
