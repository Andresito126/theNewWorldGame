package spawner

import (
	"context"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/Andresito126/theNewWorldGame/src/domain"
)

// mensaje que el spawner enviará a la ui
type NewNodeRequest struct {
	ResourceType string
	X            float64
	Y            float64
}

// la gorotune para crear recursos
func ResourceSpawnerLoop(
	ctx context.Context,
	wg *sync.WaitGroup,
	// canal de pura escritura
	newNodeChan chan<- NewNodeRequest, 
) {
	// avisa que murio
	defer wg.Done() 
	log.Println("spaw generando recursos")

	//tickers 
	treeTicker := time.NewTicker(20 * time.Second)
	scrapTicker := time.NewTicker(20 * time.Second)
	defer treeTicker.Stop()
	defer scrapTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("cerrandooooo")
			return

		case <-treeTicker.C:
			//arbol
			log.Println("spaw generando arbol")
			req := NewNodeRequest{
				ResourceType: domain.ResourceMutantTree,
				X: rand.Float64() * 640, 
				Y: rand.Float64() * 480,
			}
			// envia la peticion
			newNodeChan <- req 

		case <-scrapTicker.C:
			//charaarra
			log.Println("spaw generando chatarra")
			req := NewNodeRequest{
				ResourceType: domain.ResourceScrapPile,
				X: rand.Float64() * 640,
				Y: rand.Float64() * 480,
			}
			// envia la peticion
			newNodeChan <- req 
		}
	}
}