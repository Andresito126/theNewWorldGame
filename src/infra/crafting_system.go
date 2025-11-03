package infra

import (
	"log"

	"github.com/Andresito126/theNewWorldGame/src/application" 
	"github.com/Andresito126/theNewWorldGame/src/domain"
)

//logica de crafteo
type CraftingSystem struct {
	refugeRecipe map[string]int
	barrierRecipe map[string]int 
}

func NewCraftingSystem() *CraftingSystem {
	// crafteos
	refugeRecipe := map[string]int{
		"Wood":  50, 
		"Scrap": 30,
	}

	barrierRecipe := map[string]int{
		"Wood": 15, 
	}

return &CraftingSystem{
		refugeRecipe:  refugeRecipe,
		barrierRecipe: barrierRecipe,
	}
}

// la que se mandarpa a llamar 
func (cs *CraftingSystem) AttemptCraftRefuge(
	service *application.GameService,
	survivors [5]*VisualSurvivor,
	baseX, baseY float64,
) {
	log.Println("crafteando refugio")
	// cosnume los recursos
	canBuild := service.Store.ConsumeResources(cs.refugeRecipe)

	if canBuild {

		log.Println("crafteando refugioo")
		for _, s := range survivors {
			if s.State == "IDLE" {
				task, _ := domain.NewCraftingTask(domain.CraftTypeRefuge, int(baseX), int(baseY))
				
                s.State = "MOVING_TO_RESOURCE"
                s.TargetX = float64(task.TargetX)
                s.TargetY = float64(task.TargetY)
                s.ActiveTask = task
				break 
			}
		}
	} else {
		log.Println("recursos insuficientes")
	}
}

func (cs *CraftingSystem) AttemptCraftBarrier(
	service *application.GameService,
	survivors [5]*VisualSurvivor,
	baseX, baseY float64,
) {
	log.Println("crafteando barrera")

	canBuild := service.Store.ConsumeResources(cs.barrierRecipe)

	if canBuild {
		log.Println("crafteando barrera")
		for _, s := range survivors {
			if s.State == "IDLE" {
				task, _ := domain.NewCraftingTask(domain.CraftTypeBarrier, int(baseX), int(baseY))

                s.State = "MOVING_TO_RESOURCE"
                s.TargetX = float64(task.TargetX)
                s.TargetY = float64(task.TargetY)
                s.ActiveTask = task
				break
			}
		}
	} else {
		log.Println("recursos insuficientes")
	}
}