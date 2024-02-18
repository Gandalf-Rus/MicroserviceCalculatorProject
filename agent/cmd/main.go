package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"
)

func main() {
	fmt.Println("Hi, I'm agent, I'm here")

	// Создаем канал для принятия сигналов
	stop := make(chan os.Signal, 1)
	// Регистрируем этот канал на получение определенных сигналов
	signal.Notify(stop, os.Interrupt)

	// Ждем сигнал остановки или завершения времени сна
	select {
	case <-time.After(time.Minute):
		fmt.Println("Time is up!")
	case <-stop:
		fmt.Println("Received interrupt signal")
	}

	fmt.Println("I pass away :(")
}
