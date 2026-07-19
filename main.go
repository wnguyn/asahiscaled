package main
// read file /sys/class/power_supply/macsmc_battery
import (
	"fmt"
	"os"
	"strconv"
//	tea "charm.land/bubbletea/v2"
)

const url = "/sys/class/power_supply/macsmc_battery";

type haptic_info struct {
	haptic_buffer [64]int
}

func main() {
	fmt.Println("WE ARE READING ON \n FILE: {}", url)
	for {
		value := handle_file()
		fmt.Println("", value)
	}

}



func handle_file() int {
	sensor_val, file_err := os.ReadFile(url)
	// This is horrible
	if file_err != nil {
		return -1
	}
	
	file_output, fail_int := strconv.Atoi(string(sensor_val));
	if fail_int != nil {
		return -1
	}
	return file_output;
}
