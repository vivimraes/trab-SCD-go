package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"sync"
)

const (
	n           = 4 // Número de processos a serem iniciados
	repetitions = 5 // Número de repetições para cada processo
)

func main() {
	var wg sync.WaitGroup

	for i := 1; i <= n; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// Monta o comando para iniciar o processo
			cmd := exec.Command("go", "run", "processo.go", strconv.Itoa(id), strconv.Itoa(repetitions))

			// Redireciona a saída do comando para o terminal
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			// Executa o comando
			err := cmd.Run()
			if err != nil {
				fmt.Printf("Erro ao iniciar o processo %d: %v\n", id, err)
			}
		}(i)
	}

	// Aguarda todos os processos terminarem
	wg.Wait()
	fmt.Println("Todos os processos foram finalizados.")
}
