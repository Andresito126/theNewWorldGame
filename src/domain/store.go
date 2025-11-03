package domain

import "sync"

// se  guarda los recursos recolectados el almacen
// es la memoria compartida que debe ser protegida
type Store struct {
	mutex    sync.Mutex
	resources map[string]int
}

func NewStore() *Store {
	return &Store{
		resources: make(map[string]int),
	}
}

// la forma segura de añadir recursos
func (s *Store) AddResource(resource string, amount int) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.resources[resource] += amount
}

// es la forma seguro de leer los recursos
// devuelve una copia para que la ui la lea sin rollos
func (s *Store) GetResources() map[string]int {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// creo una copia para evitar race conditions al leer
	copy := make(map[string]int)
	for k, v := range s.resources {
		copy[k] = v
	}
	return copy
}

//comstricciones

func (s *Store) ConsumeResources(recipe map[string]int) bool {
	// bloquea el store
	s.mutex.Lock()
	// quita el bloqueo al final
	defer s.mutex.Unlock()

	// checa si hay suficiente para craftear
	for resourceName, requiredAmount := range recipe {
		if s.resources[resourceName] < requiredAmount {
			return false 
		}
	}

	for resourceName, requiredAmount := range recipe {
		s.resources[resourceName] -= requiredAmount
	}

	return true 
}