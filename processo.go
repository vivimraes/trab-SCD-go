
package main

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"
)

const (
	address = "localhost:8080"
	F       = 10
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run processo.go <ProcessID> <Repetitions>")
		return
	}

	processID := os.Args[1]
	repetitions := atoi(os.Args[2])

	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("Error connecting to coordinator:", err)
		return
	}
	defer conn.Close()

	for i := 0; i < repetitions; i++ {
		
		requestMsg := fmt.Sprintf("1|%s|%06d", processID, 0)
		conn.Write([]byte(requestMsg))

		
		buffer := make([]byte, F)
		_, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error receiving GRANT:", err)
			return
		}
		fmt.Println("Received:", string(buffer))

		
		writeToCriticalSection(processID)

		
		releaseMsg := fmt.Sprintf("3|%s|%06d", processID, 0)
		conn.Write([]byte(releaseMsg))

		
		time.Sleep(time.Duration(random(1, 5)) * time.Second)
	}
}

func writeToCriticalSection(processID string) {
	file, err := os.OpenFile("resultado.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error opening resultado.txt:", err)
		return
	}
	defer file.Close()

	timestamp := time.Now().Format("2006-01-02 15:04:05.000")
	_, err = file.WriteString(fmt.Sprintf("Process %s - %s\n", processID, timestamp))
	if err != nil {
		fmt.Println("Error writing to resultado.txt:", err)
	}
}

func random(min, max int) int {
	return min + rand.Intn(max-min)
}

func atoi(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}
