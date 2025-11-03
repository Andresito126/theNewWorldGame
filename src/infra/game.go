package infra
	
import (
	"fmt"
	"image/color"
	"log"
	"math/rand"

	"github.com/Andresito126/theNewWorldGame/src/application"
	"github.com/Andresito126/theNewWorldGame/src/domain"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

)

// aca implementa la interfaz de ebiten y conecta todooo
type Game struct {
	service *application.GameService 
	Survivors [5]*VisualSurvivor
	survivorSprite *ebiten.Image

	treeSprite *ebiten.Image
	scrapSprite *ebiten.Image
	Resources []*ResourceNode

	background *ebiten.Image
}

// el constructor de la ui
func NewGame(svc *application.GameService) *Game {

	//  carga el mapa
	bgImg, _, err := ebitenutil.NewImageFromFile("assets/images/mapa4.png") 
	if err != nil {
		log.Fatalf("Error al cargar mapa.png: %v", err)
	}
	//  crea a los supervivienes
	survivors := [5]*VisualSurvivor{}
	for i := 0; i < 5; i++ {

    // constructor llamado 
		survivors[i] = NewVisualSurvivor(i, BaseX, BaseY) 
	}
	img, _, err := ebitenutil.NewImageFromFile("assets/images/superviviente3.jpg")
	if err != nil {
		log.Fatalf("Error al cargar el sprite: %v", err)
	}

	// recursos del munno 
	treeImg, _, err := ebitenutil.NewImageFromFile("assets/images/three.png")
	if err != nil { log.Fatalf("Error al cargar %v", err) }

	scrapImg, _, err := ebitenutil.NewImageFromFile("assets/images/scrap.png")
	if err != nil { log.Fatalf("Error al cargar: %v", err) }

	// crea los nodos de recursos 
	resources := make([]*ResourceNode, 0)
	var nodeCounter int = 0 

	//arboles 
	for i := 0; i < 5; i++ {
		x := rand.Float64() * 640 
		y := rand.Float64() * 480 
		resources = append(resources, NewResourceNode(nodeCounter, x, y, treeImg, domain.ResourceMutantTree))
		nodeCounter++ 
	}

	//chatarras 
		for i := 0; i < 5; i++ {
		x := rand.Float64() * 640
		y := rand.Float64() * 480
		resources = append(resources, NewResourceNode(nodeCounter, x, y, scrapImg, domain.ResourceScrapPile))
		nodeCounter++ 
	}
	
	return &Game{
		service:   svc,
		Survivors: survivors,
		survivorSprite: img,
		treeSprite:     treeImg,
        scrapSprite:    scrapImg,
        Resources:      resources,
		background:     bgImg, 
	}
}
func (g *Game) Update() error {

	// detecta los clicks
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition() 
		for _, res := range g.Resources {
			if res.WasClicked(mx, my) {
				
                // si si se crea la tarea usando la fabrica, res id 
				task, err := domain.NewTaskFromResource(res.ResourceType, res.ID, int(res.X), int(res.Y))
				if err != nil { 
                    log.Printf("Error al crear tarea: %v", err)
                    continue 
                }

				// se busca uno libre para la tarea
				for _, s := range g.Survivors {
					if s.State == "IDLE" {
                        log.Printf("Asignando tarea %s id %d a superv %d", task.Type, task.TargetID, s.ID)
						s.State = "MOVING_TO_RESOURCE"
						// el objetivo es el recurso
						s.TargetX = res.X 
						s.TargetY = res.Y
						s.ActiveTask = task
						break 
					}
				}

				break 
			}
		}
	}


	// movimiento y estado
	for _, s := range g.Survivors {
		llegueAlDestino := s.UpdatePosition()

		if llegueAlDestino {
			if s.State == "MOVING_TO_RESOURCE" {
				log.Printf(" %d llego al recurso,  ahora trabajando", s.ID)
				s.State = "GATHERING"
				// envia la tarea a la goroutine
				g.service.AddTask(s.ActiveTask) 

			} else if s.State == "MOVING_TO_BASE" {
				log.Printf(" %d llego a la base.", s.ID)
				s.State = "IDLE" 
			}
		}
	}

	// resultados 
	select {
	case result := <-g.service.GetResultsChannel():
        // guarda el botin
		g.service.Store.AddResource(result.Resource, result.Amount)
		
        // busca al sobreviviente que termino
		for _, s := range g.Survivors {
			if s.State == "GATHERING" && s.ActiveTask.TargetID == result.TaskTargetID { 
				log.Printf("%d moviendo a la base", s.ID)
				s.State = "MOVING_TO_BASE"
				break
			}
		}

        // nodos y funcionam
		// Lista temporal para los nodos vivos
        newResourceList := make([]*ResourceNode, 0)

        for _, res := range g.Resources {
            if res.ID == result.TaskTargetID {
                res.Health -= 10 
                log.Printf("Recurso %d dañado, vida restante: %d", res.ID, res.Health)
                if res.Health > 0 {
                    newResourceList = append(newResourceList, res)
                } else {
                    log.Printf("Recurso %d agotado", res.ID)
                }
            } else {
                // no es el nodo afectado, entonces se queda
                newResourceList = append(newResourceList, res)
            }
        }
        // reemplaza lista
        g.Resources = newResourceList

	default: 
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 30, G: 30, B: 30, A: 255})

	//mapa
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(g.background, op) 
	
	// dibuja los recursos
	for _, res := range g.Resources {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(res.X, res.Y)
		screen.DrawImage(res.Sprite, op)
	}

	// dibuja a los supervivientes
	for _, s := range g.Survivors {
		op := &ebiten.DrawImageOptions{}
		// mueve el sprite a la pos del superv
		op.GeoM.Translate(s.X, s.Y) 
		screen.DrawImage(g.survivorSprite, op)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d: %s", s.ID, s.State), int(s.X)+10, int(s.Y))
	}

	// invetario 
	resources := g.service.Store.GetResources()
	debugMsg := "THE NEW WORLD\n"
	debugMsg += " ALMACÉN (Store) \n"
	for resourceName, amount := range resources {
		debugMsg += fmt.Sprintf("%s: %d\n", resourceName, amount)
	}

	ebitenutil.DebugPrint(screen, debugMsg) 
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, 480
}

