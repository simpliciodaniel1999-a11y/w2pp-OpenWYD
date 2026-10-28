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
	t.Run("nearest valid cell", func(t *testing.T) {
		h := teleportTestHeights(9)
		h.Data[4*9+4] = route.Blocked
		x, y, ok := nearestTeleportCell(h, 4, 4, 3, nil, func(x, y int) bool { return false })
		if !ok || x != 3 || y != 4 {
			t.Fatalf("got (%d,%d,%v), want (3,4,true)", x, y, ok)
		}
	})
	t.Run("occupied cell is skipped", func(t *testing.T) {
		h := teleportTestHeights(9)
		x, y, ok := nearestTeleportCell(h, 4, 4, 2, nil, func(x, y int) bool { return x == 4 && y == 4 })
		if !ok || x != 3 || y != 4 {
			t.Fatalf("got (%d,%d,%v), want (3,4,true)", x, y, ok)
		}
	})
	t.Run("no space in radius", func(t *testing.T) {
		h := teleportTestHeights(5)
		x, y, ok := nearestTeleportCell(h, 2, 2, 1, nil, func(x, y int) bool { return true })
		if ok { t.Fatalf("got (%d,%d,true), want no cell", x, y) }
	})
	t.Run("map boundaries", func(t *testing.T) {
		h := teleportTestHeights(5)
		x, y, ok := nearestTeleportCell(h, 0, 0, 2, nil, func(x, y int) bool { return false })
		if !ok || x != 0 || y != 0 {
			t.Fatalf("got (%d,%d,%v), want (0,0,true)", x, y, ok)
		}
	})
	t.Run("restricted region", func(t *testing.T) {
		h := teleportTestHeights(9)
		region := &teleportRegion{MinX: 4, MinY: 4, MaxX: 6, MaxY: 6}
		h.Data[4*9+4] = route.Blocked
		x, y, ok := nearestTeleportCell(h, 4, 4, 3, region, func(x, y int) bool { return false })
		if !ok || x < region.MinX || x > region.MaxX || y < region.MinY || y > region.MaxY {
			t.Fatalf("got (%d,%d,%v), want cell inside region", x, y, ok)
		}
	})
}
