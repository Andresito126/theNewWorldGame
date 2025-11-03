package main

import (
	"context"
	"log"
	"sync"

	"github.com/Andresito126/theNewWorldGame/src/application" 
	"github.com/Andresito126/theNewWorldGame/src/infra"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	// prepara al context y wait group para aseguraar que se cancelen las goroutines 
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// basicamente el contador de goroutines activas
	var wg sync.WaitGroup 


	//se crea el game service y se lanzan las goroutines
	log.Println("incianod el game service, inciando los  workers")
	gameService := application.NewGameService(ctx, &wg)

	// le pasamos el servicio que ya está corriendo
	gameUI := infra.NewGame(gameService)

	ebiten.SetWindowSize(1080, 720)
	ebiten.SetWindowTitle("The new world")

	log.Println("icniandoooooooo")
	if err := ebiten.RunGame(gameUI); err != nil {
		log.Fatal(err)
	}

	log.Println("enviando señal de cierree")
	cancel()

	wg.Wait() 
	log.Println("se apagaron todos los workers")
}