package main

import (
	"errors"
	"sort"
	"strings"
)

var ErrInvalidKey error = errors.New("invalid string name. string cannot be empty")

type Workspace struct {
	vectors map[string]Vector2D
}

func NewWorkspace() *Workspace {
	return &Workspace{
		vectors: make(map[string]Vector2D),
	}
}

func (w *Workspace) Define(name string, v Vector2D) error {
	name = FixName(name)
	if name == "" {
		return ErrInvalidKey
	}
	w.vectors[name] = v
	return nil
}

func (w *Workspace) Get(n string) (Vector2D, bool) {
	n = FixName(n)
	vector, exists := w.vectors[n]
	return vector, exists
}

func (w *Workspace) List() []string {
	vecNames := []string{}
	for vector := range w.vectors {
		vecNames = append(vecNames, vector)
	}
	sort.Strings(vecNames)
	return vecNames
}

func (w *Workspace) Delete(n string) bool {
	n = FixName(n)
	_, exists := w.vectors[n]
	if exists {
		delete(w.vectors, n)
		return true
	}
	return false
}

func FixName(n string) string {
	return strings.ToLower(strings.TrimSpace(n))
}
