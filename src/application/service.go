package application

import (
	"context"
	"github.com/Andresito126/theNewWorldGame/src/application/survivor"
	"github.com/Andresito126/theNewWorldGame/src/domain"
	"sync"
)

// es el que conecta todo
type GameService struct {
	Store       *domain.Store
	jobsChan    chan domain.Task
	resultsChan chan domain.Result
}

// aca es donde por asi decirlo, donde se ensambla toda la logica de mi concu
func NewGameService(ctx context.Context, wg *sync.WaitGroup) *GameService {
	store := domain.NewStore()

	// canales, el de tareas y los resultados
	jobs := make(chan domain.Task, 100)   
	results := make(chan domain.Result, 100)

	// lanzo a las goroutines, de los supervivientes 
	const numSurvivors = 5
	for i := 0; i < numSurvivors; i++ {
		// para decir al wait group  que se aguantes por una goroutine 
		wg.Add(1) 
		go survivor.SurvivorMainLoop(ctx, wg, i, jobs, results, store)
	}

	// me da el service listo
	return &GameService{
		Store:       store,
		jobsChan:    jobs,
		resultsChan: results,
	}
}

//esta es la forma en que la ui añade una nueva tarea
func (s *GameService) AddTask(t domain.Task) {
	s.jobsChan <- t
}

// canal para solo leer
// para que pueda ver los resultados sin poder enviar nada
func (s *GameService) GetResultsChannel() <-chan domain.Result {
	return s.resultsChan
}