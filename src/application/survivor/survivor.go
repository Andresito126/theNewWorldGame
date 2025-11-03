package survivor

import (
	"context"
	"log"
	"github.com/Andresito126/theNewWorldGame/src/domain"
	"sync"
	"time"
)

//  es la función que corre en su propia goroutine
//  la vida de un sobre
func SurvivorMainLoop(
	// para saber cuando debe morir, cuando se cierra
	ctx context.Context, 
	//  para avisar que ya murio
	wg *sync.WaitGroup,  
	id int,
	// canal de solo lectura de tareas
	jobs <-chan domain.Task, 
	// canal de solo escritura de resultados
	results chan<- domain.Result, 
	store *domain.Store, 
) {
	// es para que avise al waitgroup cuando la func muera
	defer wg.Done()
	log.Printf("superviviendte %d libre para trabjar", id)

	for {
		select {
		case <-ctx.Done():
			//  el context fue cancelado, o sea que se termina
			log.Printf("superv %d terminadooo", id)
			return 

		case task, ok := <-jobs:
			if !ok {
				// canal de jobs termina
				log.Printf("superv %d canal de tareas cerrado.", id)
				return 
			}

			// recibio la tarea
			log.Printf("superv %d tienes una nueva  tarea: %s trabajando por %s", id, task.Type, task.Duration)

			
			time.Sleep(task.Duration)

			// resultado
			botin := domain.Result{
				Resource:     task.Resource,
				Amount:       10, 
				TaskTargetID: task.TargetID, 
			}

			// con esto se regresa el resultado de regreso
			results <- botin
			log.Printf("superv %d tarea completadaa, entregando %d de %s", id, botin.Amount, botin.Resource)
		}
	}
}