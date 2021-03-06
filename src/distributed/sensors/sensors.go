package main

import (
	"bytes"
	"encoding/gob"
	"flag"
	"go_rabbitmq/src/distributed/dto"
	"go_rabbitmq/src/distributed/qutils"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
)

var name = flag.String("name", "sensor", "name of the sensor")
var freq = flag.Uint("freq", 5, "update frequency in cycles/sec")
var max = flag.Float64("max", 5., "maximum value for generated readings")
var min = flag.Float64("min", 1., "minimum value for generated readings")
var stepSize = flag.Float64("step", 0.1, "maximum allowable change per measurement")

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

var value = r.Float64()*(*max-*min) + *min
var nom = (*max-*min)/2 + *min

var url = "amqp://" + os.Getenv("RABBITMQ_USER") + ":" + os.Getenv("RABBITMQ_PASSWORD") + "@" + os.Getenv("RABBITMQ_HOST") + ":" + os.Getenv("RABBITMQ_PORT")

func init() {
	// load environment variables
	dotenvErr := godotenv.Load()
	if dotenvErr != nil {
		log.Fatal(dotenvErr)
	}
}

func main() {
	flag.Parse()

	conn, ch := qutils.GetChannel(url)
	defer conn.Close()
	defer ch.Close()

	dataQueue := qutils.GetQueue(*name, ch)

	// time between each signal
	dur, _ := time.ParseDuration(strconv.Itoa(1000/int(*freq)) + "ms")

	signal := time.Tick(dur)

	// encode data to the buffer for transmission
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)

	for range signal {
		calcValue()
		reading := dto.SensorMessage{
			Name:      *name,
			Value:     value,
			Timestamp: time.Now(),
		}
		buf.Reset()         // reset the buffer - remove previous data
		enc.Encode(reading) // encode the message

		// capture client message sent to server
		msg := amqp.Publishing{
			Body: buf.Bytes(),
		}

		ch.Publish("", // default exchange
			dataQueue.Name,
			false,
			false,
			msg)
		log.Printf("Reading sent. Value: %v\n", value)
	}
}

func calcValue() {
	var maxStep, minStep float64

	if value < nom {
		maxStep = *stepSize
		minStep = -1 * *stepSize * (value - *min) / (nom - *min)
	} else {
		maxStep = *stepSize * (*max - value) / (*max - nom)
		minStep = -1 * *stepSize
	}

	value += r.Float64()*(maxStep-minStep) + minStep
}
