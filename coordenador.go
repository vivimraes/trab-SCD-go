package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	address   = "localhost:8080"
	F         = 10 // Tamanho fixo da mensagem
	separator = "|"
	logFile   = "coordinator_log.txt"
)

type Coordinator struct {
	queue        []string
	mutex        sync.Mutex
	logs         []string
	processCount map[int]int
}

func (c *Coordinator) handleConnection(conn net.Conn) {
	defer conn.Close()
	buffer := make([]byte, F)
	for {
		n, err := conn.Read(buffer)
		if err != nil || n == 0 {
			break
		}
		message := string(buffer)
		fmt.Println("Received:", message)

		// Processamento da mensagem (REQUEST, GRANT, RELEASE)
		c.processMessage(message, conn)
	}
}

func (c *Coordinator) processMessage(message string, conn net.Conn) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Parse e log da mensagem
	fields := parseMessage(message)
	timestamp := time.Now().Format(time.RFC3339)
	c.logs = append(c.logs, fmt.Sprintf("%s - Received: %s", timestamp, message))

	// Gravar no arquivo de log
	c.logToFile(fmt.Sprintf("%s - Received: %s", timestamp, message))

	switch fields[0] {
	case "1": // REQUEST
		c.queue = append(c.queue, message)
		c.sendGrant(fields[1], conn)
	case "3": // RELEASE
		// Atualizar a fila e contar
		c.processCount[atoi(fields[1])]++
		c.queue = c.queue[1:]
		fmt.Printf("Process %s released the critical section.\n", fields[1])
	}
}

func (c *Coordinator) sendGrant(processID string, conn net.Conn) {
	grantMsg := fmt.Sprintf("2|%s|%06d", processID, 0)
	conn.Write([]byte(grantMsg))

	timestamp := time.Now().Format(time.RFC3339)
	c.logs = append(c.logs, fmt.Sprintf("%s - Sent: %s", timestamp, grantMsg))
	c.logToFile(fmt.Sprintf("%s - Sent: %s", timestamp, grantMsg))

	fmt.Printf("Granted access to process %s.\n", processID)
}

func (c *Coordinator) logToFile(logEntry string) {
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error opening log file:", err)
		return
	}
	defer file.Close()

	if _, err := file.WriteString(logEntry + "\n"); err != nil {
		fmt.Println("Error writing to log file:", err)
	}
}

func (c *Coordinator) handleTerminalInput() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		command := scanner.Text()
		c.executeCommand(command)
	}
}

func (c *Coordinator) executeCommand(command string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	switch command {
	case "1": // Imprimir a fila de pedidos atual
		fmt.Println("Fila de Pedidos:", c.queue)
	case "2": // Imprimir quantas vezes cada processo foi atendido
		fmt.Println("Contagem de Processos:", c.processCount)
	case "3": // Encerrar a execução
		c.shutdown()
	}
}

func (c *Coordinator) shutdown() {
	// Limpar os logs e resultado.txt
	os.Remove("resultado.txt")
	os.Remove(logFile)
	fmt.Println("Logs removidos e execução encerrada.")
	os.Exit(0)
}

func parseMessage(msg string) []string {
	return strings.Split(msg, separator)
}

func atoi(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func main() {
	coordinator := Coordinator{
		processCount: make(map[int]int),
	}

	go func() {
		ln, err := net.Listen("tcp", address)
		if err != nil {
			fmt.Println("Error starting server:", err)
			os.Exit(1)
		}
		defer ln.Close()

		fmt.Println("Coordinator started...")

		for {
			conn, err := ln.Accept()
			if err != nil {
				fmt.Println("Error accepting connection:", err)
				continue
			}
			go coordinator.handleConnection(conn)
		}
	}()

	// Iniciar a thread da interface de terminal
	coordinator.handleTerminalInput()
}
