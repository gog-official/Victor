package workspace

import (
	"errors"
	"sort"
	"strings"

	"github.com/gog-official/victor/internal/vector"
)

const MaxVectors = 5

var (
	ErrInvalidKey    error = errors.New("invalid string name. string cannot be empty")
	ErrWorkspaceFull error = errors.New("workspace can hold at most 5 vectors")
)

type Workspace struct {
	vectors map[string]vector.Vector2D
}

func NewWorkspace() *Workspace {
	return &Workspace{
		vectors: make(map[string]vector.Vector2D),
	}
}

func (w *Workspace) Define(name string, v vector.Vector2D) error {
	name = FixName(name)
	if name == "" {
		return ErrInvalidKey
	}
	if _, exists := w.vectors[name]; !exists && len(w.vectors) >= MaxVectors {
		return ErrWorkspaceFull
	}
	w.vectors[name] = v
	return nil
}

func (w *Workspace) Get(n string) (vector.Vector2D, bool) {
	n = FixName(n)
	v, exists := w.vectors[n]
	return v, exists
}

func (w *Workspace) List() []string {
	vecNames := []string{}
	for vecName := range w.vectors {
		vecNames = append(vecNames, vecName)
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
