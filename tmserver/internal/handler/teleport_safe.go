package handler

import (
	"github.com/jeanluca/w2pp-openwyd/tmserver/internal/content"
	"github.com/jeanluca/w2pp-openwyd/tmserver/internal/route"
)

// teleportRegion representa um retângulo inclusivo opcional que mantém o
// destino alternativo dentro do mapa ou da plataforma que autorizou o teleporte.
type teleportRegion struct {
	MinX, MinY, MaxX, MaxY int
}

// nearestTeleportCell procura a célula livre e não bloqueada mais próxima,
// dentro de maxDistance, usando distância euclidiana. Empates são resolvidos
// de forma determinística: de cima para baixo e, depois, da esquerda para a direita.
// Com region nil, qualquer célula dentro dos limites do mapa pode ser considerada.
func nearestTeleportCell(
	heights *content.Grid,
	targetX, targetY, maxDistance int,
	region *teleportRegion,
	occupied func(x, y int) bool,
) (int, int, bool) {
	if heights == nil || heights.Dim <= 0 || len(heights.Data) != heights.Dim*heights.Dim ||
		maxDistance < 0 || occupied == nil {
		return 0, 0, false
	}
	dim := heights.Dim
	// Um raio maior que a dimensão do mapa não encontra células adicionais.
	// Limitá-lo também evita estouro nos cálculos de coordenadas e distância.
	if maxDistance > dim {
		maxDistance = dim
	}
	if targetX < 0 || targetY < 0 || targetX >= dim || targetY >= dim {
		return 0, 0, false
	}
	if region != nil && (region.MinX < 0 || region.MinY < 0 ||
		region.MaxX >= dim || region.MaxY >= dim ||
		region.MinX > region.MaxX || region.MinY > region.MaxY ||
		targetX < region.MinX || targetX > region.MaxX ||
		targetY < region.MinY || targetY > region.MaxY) {
		return 0, 0, false
	}

	bestX, bestY := 0, 0
	bestDistance := int(^uint(0) >> 1)
	found := false
	minX, maxX := targetX-maxDistance, targetX+maxDistance
	minY, maxY := targetY-maxDistance, targetY+maxDistance
	if minX < 0 {
		minX = 0
	}
	if minY < 0 {
		minY = 0
	}
	if maxX >= dim {
		maxX = dim - 1
	}
	if maxY >= dim {
		maxY = dim - 1
	}
	if region != nil {
		if minX < region.MinX {
			minX = region.MinX
		}
		if minY < region.MinY {
			minY = region.MinY
		}
		if maxX > region.MaxX {
			maxX = region.MaxX
		}
		if maxY > region.MaxY {
			maxY = region.MaxY
		}
	}
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			dx, dy := x-targetX, y-targetY
			distance := dx*dx + dy*dy
			if distance > maxDistance*maxDistance || distance > bestDistance {
				continue
			}
			// BASE_ApplyAttribute grava route.Blocked (127) na grade de alturas
			// para representar células bloqueadas. Os valores são bytes com sinal
			// na lógica legada.
			if int(int8(heights.At(x, y))) == route.Blocked || occupied(x, y) {
				continue
			}
			if !found || distance < bestDistance {
				bestX, bestY, bestDistance, found = x, y, distance, true
			}
		}
	}
	return bestX, bestY, found
}
