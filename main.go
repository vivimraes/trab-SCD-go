package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"sync"
)

const (
	n           = 3 
	repetitions = 5 
)

func main() {
	var wg sync.WaitGroup

	for i := 1; i <= n; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			
			cmd := exec.Command("go", "run", "processo.go", strconv.Itoa(id), strconv.Itoa(repetitions))

			
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			
			err := cmd.Run()
			if err != nil {
				fmt.Printf("Erro ao iniciar o processo %d: %v\n", id, err)
			}
		}(i)
	}

	
	wg.Wait()
	fmt.Println("Todos os processos foram finalizados.")
}
