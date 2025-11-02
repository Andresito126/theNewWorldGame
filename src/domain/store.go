package domain

import "sync"

// Store guarda los recursos recolectados el almacen
// es la memoria compartida que debe ser protegida
type Store struct {
	mutex    sync.Mutex
	resources map[string]int
}

// NewStore es el constructor que prepara el store
func NewStore() *Store {
	return &Store{
		resources: make(map[string]int),
	}
}

// AddResource es la forma segura de añadir recursos
func (s *Store) AddResource(resource string, amount int) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.resources[resource] += amount
}

// GetResources es la forma seguro de leer los recursos
// Devuelve una copia para que la ui la lea sin rollos
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