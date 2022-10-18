// Coded by TD, SD & KD the 10/10/22
// Memory & CPU info
package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("Here are your memory informations :")
	read_mem_info()
	fmt.Println()
	fmt.Println("Here are your CPU informations :")
	fmt.Println("ID User Nice System Idle IOWait IRQ SoftIRQ Steal Guest GuestNice")
	get_cpu_stat()
}

// Function to open the memory file
func read_mem_info() {
	file, err := os.Open("/proc/meminfo")
	// Control that the file was open with no error or exiting
	if err != nil {
		fmt.Println("Could not open the memory file - Exiting - : ", err)
		os.Exit(1)
	}
	// Used to close the file at the end of main()
	defer file.Close()

	// Define a scanner to read the file line by line
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	// Scan all the line and retrieve our targets
	array := [13]string{"MemTotal", "MemFree", "MemAvailable", "Buffers", "Cached", "SwapCached", "Active", "Inactive", "SwapTotal", "SwapFree", "Dirty", "Writeback", "Shmem"}
	i := 0
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), array[i]) && i < len(array) {
			fmt.Println(scanner.Text())
			i++
			if i == len(array) {
				break
			}
		}
	}
}

// Function to retrieve the number of CPU cores
func get_cpu_number() string {
	output := exec.Command("nproc")
	out, err := output.CombinedOutput()
	if err != nil {
		fmt.Println("Could not read nproc command - Exiting - : ", err)
		os.Exit(1)
	}
	return string(out)
}

// Function to get all CPU stats
func get_cpu_stat() {
	file, err := os.Open("/proc/stat")
	if err != nil {
		fmt.Println("Could not open the /proc/stat file - Exiting - :", err)
		os.Exit(1)
	}
	defer file.Close()

	// Read every lines of the file and put it in rawBytes
	rawBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("Could not read the lines - Exiting - :", err)
		os.Exit(1)
	}

	// Splitting every lines by a \n
	lines := strings.Split(string(rawBytes), "\n")
	// Retrieving the number of cpu core, adding 1 to count the first line of the file : all CPU info
	cpu_nb, err := strconv.Atoi(strings.TrimSuffix(get_cpu_number(), "\n"))
	if err != nil {
		fmt.Println("Could no get cpu number", err)
		os.Exit(1)
	}

	// For every cpu, display informations
	for i := 0; i < cpu_nb+1; i++ {
		fmt.Println(lines[i])
	}
}
