package work

import "sync"

// Интерфейс надо реализовать объектам, которые будут обрабатываться параллельно
type Worker interface {
	Task()
}

// Пул для выполнения
type Pool struct {
	// из этого канала будем брать задачи для обработки
	tasks chan Worker
	// для синхронизации работы
	wg sync.WaitGroup
}

// при создании пула передадим максимальное количество горутин
func New(maxGoroutines int) *Pool {
	pool := Pool{
		tasks: make(chan Worker), // канал, откуда брать задачи
	}
	// для ожидания завершения
	pool.wg.Add(maxGoroutines)
	for i := 0; i < maxGoroutines; i++ {
		// создадим горутины по указанному количеству maxGoroutines
		go func() {
			// забираем задачи из канала
			for w := range pool.tasks {
				// и выполняем
				w.Task()
			}
			// после закрытия канала нужно оповестить наш пул
			pool.wg.Done()
		}()
	}

	return &pool
}

// Передаем объект, который реализует интерфейс Worker
func (p *Pool) Run(w Worker) {
	// добавляем задачи в канал, из которого забирает работу пул
	p.tasks <- w
}

func (p *Pool) Shutdown() {
	// закроем канал с задачами
	close(p.tasks)
	// дождемся завершения работы уже запущенных задач
	p.wg.Wait()
}
