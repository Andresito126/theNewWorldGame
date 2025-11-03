package domain

import (
	"fmt"
	"time"
)

const (
	CraftTypeRefuge  = "CRAFT_REFUGE"
	CraftTypeBarrier = "CRAFT_BARRIER"
)

func NewCraftingTask(craftType string, baseX, baseY int) (Task, error) {
	switch craftType {	

	case CraftTypeRefuge:

		return Task{
			Type: "BUILD",
			Resource: "Refuge",
			Duration: 10 * time.Second,
			TargetX:  baseX,
			TargetY:  baseY,
		}, nil

	case CraftTypeBarrier:
		return Task{
			Type: "BUILD",
			Resource: "Barrier", 
			Duration: 4 * time.Second, 
			TargetX:  baseX,
			TargetY:  baseY,
		}, nil

	default:
		return Task{}, fmt.Errorf("tipo de crafteo desconocido: %s", craftType)
	}
}