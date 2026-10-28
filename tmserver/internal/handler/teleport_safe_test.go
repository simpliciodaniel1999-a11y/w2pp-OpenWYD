package handler

import (
	"testing"

	"github.com/jeanluca/w2pp-openwyd/tmserver/internal/content"
	"github.com/jeanluca/w2pp-openwyd/tmserver/internal/route"
)

func teleportTestHeights(dim int) *content.Grid {
	return &content.Grid{Dim: dim, Data: make([]byte, dim*dim)}
}

func TestNearestTeleportCell(t *testing.T) {
	t.Run("célula válida mais próxima", func(t *testing.T) {
		h := teleportTestHeights(9)
		h.Data[4*9+4] = route.Blocked
		x, y, ok := nearestTeleportCell(h, 4, 4, 3, nil, func(x, y int) bool { return false })
		if !ok || x != 4 || y != 3 {
			t.Fatalf("recebido (%d,%d,%v), esperado (4,3,true)", x, y, ok)
		}
	})

	t.Run("célula ocupada é ignorada", func(t *testing.T) {
		h := teleportTestHeights(9)
		x, y, ok := nearestTeleportCell(h, 4, 4, 2, nil, func(x, y int) bool { return x == 4 && y == 4 })
		if !ok || x != 4 || y != 3 {
			t.Fatalf("recebido (%d,%d,%v), esperado (4,3,true)", x, y, ok)
		}
	})

	t.Run("nenhum espaço livre no raio", func(t *testing.T) {
		h := teleportTestHeights(5)
		x, y, ok := nearestTeleportCell(h, 2, 2, 1, nil, func(x, y int) bool { return true })
		if ok {
			t.Fatalf("recebido (%d,%d,true), esperado nenhuma célula", x, y)
		}
	})

	t.Run("limites do mapa", func(t *testing.T) {
		h := teleportTestHeights(5)
		x, y, ok := nearestTeleportCell(h, 0, 0, 2, nil, func(x, y int) bool { return false })
		if !ok || x != 0 || y != 0 {
			t.Fatalf("recebido (%d,%d,%v), esperado (0,0,true)", x, y, ok)
		}
	})

	t.Run("terreno bloqueado é ignorado", func(t *testing.T) {
		h := teleportTestHeights(5)
		for y := 1; y <= 3; y++ {
			for x := 1; x <= 3; x++ {
				h.Data[y*5+x] = byte(route.Blocked)
			}
		}
		x, y, ok := nearestTeleportCell(h, 2, 2, 2, nil, func(x, y int) bool { return false })
		if !ok || (x >= 1 && x <= 3 && y >= 1 && y <= 3) {
			t.Fatalf("recebido (%d,%d,%v), esperado célula livre fora da área bloqueada", x, y, ok)
		}
	})

	t.Run("destino permanece dentro da região permitida", func(t *testing.T) {
		h := teleportTestHeights(9)
		region := &teleportRegion{MinX: 4, MinY: 4, MaxX: 6, MaxY: 6}
		h.Data[4*9+4] = byte(route.Blocked)
		x, y, ok := nearestTeleportCell(h, 4, 4, 3, region, func(x, y int) bool { return false })
		if !ok || x < region.MinX || x > region.MaxX || y < region.MinY || y > region.MaxY {
			t.Fatalf("recebido (%d,%d,%v), esperado célula dentro da região", x, y, ok)
		}
	})

	t.Run("destino e região inválidos são rejeitados", func(t *testing.T) {
		h := teleportTestHeights(5)
		ocupada := func(x, y int) bool { return false }
		if _, _, ok := nearestTeleportCell(h, -1, 0, 2, nil, ocupada); ok {
			t.Fatal("destino fora do mapa foi aceito")
		}
		regiao := &teleportRegion{MinX: 3, MinY: 3, MaxX: 2, MaxY: 4}
		if _, _, ok := nearestTeleportCell(h, 2, 2, 2, regiao, ocupada); ok {
			t.Fatal("região inválida foi aceita")
		}
	})
}
