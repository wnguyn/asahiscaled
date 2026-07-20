package main

import "fmt"

func handleFunction(a *hapticInfo, status FingerStatus) {
	switch status {
	case NoFinger:
		fmt.Println("No finger detected")
	case FingerOnPad:
		fmt.Println("Finger on pad")
	case ActiveMeasure:
		fmt.Println("Active measurement — finger pressing")
	case Err:
		fmt.Println("Error state")
	}
}
