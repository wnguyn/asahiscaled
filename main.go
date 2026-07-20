package main

// read file /sys/class/power_supply/macsmc_battery
import (
	"fmt"
	"os"
	"strconv"
	"time"
	// tea "charm.land/bubbletea/v2"
)



const (
	NoFinger FingerStatus = foo
	FingerOnPad
	ActiveMeasure
)
const url = "/sys/class/power_supply/macsmc_battery"

type hapticInfo struct {
	hapticBuffer [64]int 
	tail          int
	enumTail int // meant to keep tabs w/ enum status
	calcValue int 
}

func main() {
	fmt.Println("WE ARE READING ON \n FILE: {}", url)

	hapticInfo := hapticInfo{
		hapticBuffer: []int,
		tail:          0,
	}

	for {
		value := handleFile()
		fillBuffer(&hapticInfo, value)
		sortEnum()
	}

}

func fillBuffer(a *hapticInfo, val int) {
	idx := a.tail
	a.hapticBuffer[idx] = val
	if val != -1 {
		a.tail++
	}
}

func handleFile() int {
	sensorVal, fileErr := os.ReadFile(url)
	// This is horrible holy shit
	if fileErr != nil {
		return -1
	}

	fileOutput, failInt := strconv.Atoi(string(sensorVal))
	if failInt != nil {
		return -1
	}
	return fileOutput
}
