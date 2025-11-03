package infra

import (
	"github.com/Andresito126/theNewWorldGame/src/domain"
	"math"
)

const (
	survivorSpeed = 2.0
	BaseX         = 320
	BaseY         = 240
)

// representa un sobreviviente
type VisualSurvivor struct {
	ID int
	// "IDLE", "MOVING_TO_RESOURCE", "GATHERING", "MOVING_TO_BASE"
	State      string
	X          float64
	Y          float64
	TargetX    float64
	TargetY    float64
	ActiveTask domain.Task
}
func NewVisualSurvivor(id int, x, y float64) *VisualSurvivor {
	return &VisualSurvivor{
		ID:    id,
		State: "IDLE",
		X:     BaseX,
		Y:     BaseY,
	}
}

// UpdatePosition es de solo el movimiento
// da true si llego
func (s *VisualSurvivor) UpdatePosition() bool {

	var targetX, targetY float64

	if s.State == "MOVING_TO_RESOURCE" {
		targetX = s.TargetX
		targetY = s.TargetY
	} else if s.State == "MOVING_TO_BASE" {
		// pos del refugio
		targetX = BaseX
		targetY = BaseY
	} else {
		return false
	}
	// logica de movemnet
	dx := targetX - s.X
	dy := targetY - s.Y
	distance := math.Sqrt(dx*dx + dy*dy)

	if distance < survivorSpeed {
		// llego
		s.X = targetX
		s.Y = targetY
		return true
	}

	// se anda moviendo
	angle := math.Atan2(dy, dx)
	s.X += math.Cos(angle) * survivorSpeed
	s.Y += math.Sin(angle) * survivorSpeed
	return false 
}
