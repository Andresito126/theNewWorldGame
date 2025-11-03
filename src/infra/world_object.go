package infra

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// es un objeto clickeable 
type ResourceNode struct {
	X float64
	Y float64
	Sprite *ebiten.Image
	ResourceType string 
	ID int 
	Health int 
}

// mi constructor 
func NewResourceNode(id int, x, y float64, sprite *ebiten.Image, resType string) *ResourceNode {
	return &ResourceNode{
		X: x,
		Y: y,
		Sprite: sprite,
		ResourceType: resType,
		ID: id,
		Health: 20, 
	}
}

// devuelve el rectangulo del sprite para detectar clics
func (r *ResourceNode) GetBounds() image.Rectangle {
	bounds := r.Sprite.Bounds()
	return image.Rect(
		int(r.X),
		int(r.Y),
		//ancho
		int(r.X)+bounds.Dx(),
		// alto
		int(r.Y)+bounds.Dy(), 
	)
}

// wasclicked revisa si las coordenadas  están dentro del sprite
func (r *ResourceNode) WasClicked(mx, my int) bool {
	bounds := r.GetBounds()
	return mx >= bounds.Min.X && mx < bounds.Max.X &&
		my >= bounds.Min.Y && my < bounds.Max.Y
}