package main

import (
	"container/list"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

type FingerStatus int

const (
	NoFinger FingerStatus = iota
	FingerOnPad
	ActiveMeasure
	Err
)

func (s FingerStatus) String() string {
	switch s {
	case NoFinger:
		return "No Finger"
	case FingerOnPad:
		return "Finger On Pad"
	case ActiveMeasure:
		return "Pressing"
	case Err:
		return "Error"
	}
	return "wtf"
}

const sensorPath = "/sys/class/power_supply/macsmc_battery"

type hapticInfo struct {
	last_value *list.List
}

func readSensor() int {
	raw, err := os.ReadFile(sensorPath)
	if err != nil {
		return -1
	}
	val, err := strconv.Atoi(strings.TrimSpace(string(raw)))
	if err != nil {
		return -1
	}
	return val
}

func fillBuffer(a *hapticInfo, value int) {
	a.last_value.PushFront(value)
	if a.last_value.Len() > 64 {
		a.last_value.Remove(a.last_value.Back())
	}
}

func collectValues(a *hapticInfo) []float64 {
	vals := make([]float64, 0, a.last_value.Len())
	for e := a.last_value.Front(); e != nil; e = e.Next() {
		vals = append(vals, float64(e.Value.(int)))
	}
	return vals
}

func sortEnum(a *hapticInfo) FingerStatus {
	vals := collectValues(a)
	if len(vals) < 4 {
		return NoFinger
	}
	recent := avg(vals[:min(8, len(vals))])
	baselineStart := min(32, len(vals)-1)
	baseline := avg(vals[baselineStart:])
	delta := math.Abs(recent - baseline)
	recentVariance := variance(vals[:min(16, len(vals))])

	const threshold = 50.0
	switch {
	case delta > threshold && recentVariance < 200:
		return ActiveMeasure
	case delta > threshold/2:
		return FingerOnPad
	default:
		return NoFinger
	}
}

func avg(vals []float64) float64 {
	if len(vals) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range vals {
		sum += v
	}
	return sum / float64(len(vals))
}

func variance(vals []float64) float64 {
	if len(vals) < 2 {
		return 0
	}
	m := avg(vals)
	sumSq := 0.0
	for _, v := range vals {
		d := v - m
		sumSq += d * d
	}
	return sumSq / float64(len(vals))
}

func main() {
	info := hapticInfo{
		last_value: list.New(),
	}

	fmt.Println("reading sensor...")

	for {
		val := readSensor()
		fillBuffer(&info, val)
		status := sortEnum(&info)
		handleFunction(&info, status)
		time.Sleep(80 * time.Millisecond)
	}
}
