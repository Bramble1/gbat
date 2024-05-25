package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"
	"os"
	"os/exec"
)

const (
	BATTERY_LEVEL=68
	MESSAGE_DURATION="30000"
)


//get the battery percentage
func battery_status(path string) int{
	//read from battery status info from the file
	file := path + "/capacity"
	data, err := ioutil.ReadFile(file)
	if err != nil {
		os.Exit(-1)
	}

	//return the percentage removing any whitespace
	percentage := strings.TrimSpace(string(data))
	perc := convert_to_number(percentage)
	return perc
}

func PowInts(x, n int) int {
   if n == 0 { return 1 }
   if n == 1 { return x }
   y := PowInts(x, n/2)
   if n % 2 == 0 { return y*y }
   return x*y*y
}

func convert_to_number(buffer string) int{

	number := 0
	length := len(buffer)
	exp := length - 1
	for i:= 0; i<length; i++ {
		number += PowInts(10,exp) * int(buffer[i]-'0')
		exp-=1
	}

	return number

}

func execute_command(percentage int){
	message := fmt.Sprintf("Low Battery threshold: %d%%",percentage)
	cmd := exec.Command("notify-send", "-u", "critical", "-t", MESSAGE_DURATION,message)

	err := cmd.Run()

	if err != nil{
		fmt.Printf("Error executing command: %s",err)
		os.Exit(-1)
	}
}

func is_charging()bool{
	data,err := ioutil.ReadFile("/sys/class/power_supply/AC/online")
	if err != nil {
		fmt.Println("Error ",err)
		os.Exit(-1)
	}

	//convert data to string of bytes
	status := strings.TrimSpace(string(data))

	return status == "1"
}

func main(){
	battery_path := "/sys/class/power_supply/BAT0"
	
	for{
		if(is_charging()!=true){
			if battery_status(battery_path) <= BATTERY_LEVEL{
				execute_command(battery_status(battery_path))
			}
		}

		time.Sleep(30*time.Second)
	}
}
