package domain

import (
	"fmt"
	"time"
)
// recursos 
const (
	ResourceMutantTree = "MutantTree"
	ResourceScrapPile  = "ScrapPile"
	ResourceWaterPuddle = "WaterPuddle"
)

// para pasarselas a la tarea.
func NewTaskFromResource(resourceType string, targetID int, x, y int) (Task, error) {

	switch resourceType {
	case ResourceMutantTree:
		return Task{
			Type:     "GATHER",
			Resource: "Wood",
			Duration: 3 * time.Second,
			TargetX:  x,
			TargetY:  y,
			TargetID: targetID,
		}, nil

	case ResourceScrapPile:
		return Task{
			Type:     "GATHER",
			Resource: "Scrap",
			Duration: 5 * time.Second,
			TargetX:  x,
			TargetY:  y,
			TargetID: targetID,
		}, nil

	case ResourceWaterPuddle:
		return Task{
			Type:     "GATHER",
			Resource: "DirtyWater",
			Duration: 2 * time.Second,
			TargetX:  x,
			TargetY:  y,
			TargetID: targetID,
		}, nil

	default:
		return Task{}, fmt.Errorf("unknown resource type: %s", resourceType)
	}
}