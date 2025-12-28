package main

/*
#include <stdlib.h>

typedef enum {
	MATERIAL_UNKNOWN = 0,
	MATERIAL_PYRO,
	MATERIAL_HYDRO,
	MATERIAL_AERO,
	MATERIAL_NATURO,
	MATERIAL_GEO,
	MATERIAL_AETHER,
	MATERIAL_BLOOD,
	MATERIAL_PALEO,
	MATERIAL_MELO,
	MATERIAL_ALLEGRO,
	MATERIAL_FAITH
} Material;

typedef struct {
	Material material;
	int value;
} MaterialCostView;

typedef struct {
	int card_type;
	char* tag;
	char* uuid;
	int cost_count;
	MaterialCostView* costs;
	int atk;
	int hp;
	int mp;
} CardView;

typedef struct {
	char* id;
	int x;
	int y;
	int hp;
	int mp;
	int atk;
	char* owner_id;
	int used_action_cost;
	int max_action_cost;
} PieceView;

typedef struct {
	int width;
	int height;
	int turn;
	int side;
	int piece_count;
	PieceView* pieces;
} BoardView;
*/
import "C"

import (
	"unsafe"

	"frontage/pkg/engine/impl/card"
	"frontage/pkg/engine/model"
)

func materialCostsToC(m model.Materials) (*C.MaterialCostView, C.int) {
	if len(m) == 0 {
		return nil, 0
	}
	count := len(m)
	itemSize := unsafe.Sizeof(C.MaterialCostView{})
	raw := C.malloc(C.size_t(count) * C.size_t(itemSize))
	slice := (*[1 << 30]C.MaterialCostView)(raw)[:count:count]
	i := 0
	enumMap := materialEnumMap()
	for k, v := range m {
		if enumVal, ok := enumMap[k]; ok {
			slice[i].material = enumVal
		} else {
			slice[i].material = C.MATERIAL_UNKNOWN
		}
		slice[i].value = C.int(v)
		i++
	}
	return (*C.MaterialCostView)(raw), C.int(count)
}

func materialEnumMap() map[model.Material]C.Material {
	materials := model.EnumerateMaterial()
	result := make(map[model.Material]C.Material, len(materials))
	for i, m := range materials {
		result[m] = C.Material(i + 1)
	}
	return result
}

func FillCardView(dst *C.CardView, src model.Card) C.int {
	if dst == nil || src == nil {
		return -1
	}
	dst.card_type = C.int(src.Type())
	dst.tag = C.CString(string(src.LocalizeTag()))
	dst.uuid = C.CString(src.Id().String())
	dst.atk = 0
	dst.hp = 0
	dst.mp = 0
	if p, ok := src.(card.Piece); ok {
		dst.atk = C.int(p.ATK())
		dst.hp = C.int(p.HP())
		dst.mp = C.int(p.MP())
	}
	dst.costs, dst.cost_count = materialCostsToC(src.PlayCost())
	return 0
}

func FillPieceView(dst *C.PieceView, src model.Piece) C.int {
	if dst == nil || src == nil {
		return -1
	}
	dst.id = C.CString(src.Id().String())
	dst.x = C.int(src.Position().X)
	dst.y = C.int(src.Position().Y)
	dst.hp = C.int(src.HP())
	dst.mp = C.int(src.MP())
	dst.atk = C.int(src.ATK())
	dst.owner_id = nil
	if owner := src.Owner(); owner != nil {
		dst.owner_id = C.CString(owner.ID().String())
	}
	dst.used_action_cost = 0
	dst.max_action_cost = 0
	if mp, ok := src.(model.MutablePiece); ok {
		dst.used_action_cost = C.int(mp.UsedActionCost())
		dst.max_action_cost = C.int(mp.MaxActionCost())
	}
	return 0
}

func FillBoardView(dst *C.BoardView, src *model.Board) C.int {
	if dst == nil || src == nil {
		return -1
	}
	size := src.Size()
	dst.width = C.int(size.Width)
	dst.height = C.int(size.Height)
	dst.turn = C.int(src.Turn())
	dst.side = C.int(src.Phase())

	entities := src.Entities()
	count := 0
	for x := range entities {
		for y := range entities[x] {
			if entities[x][y] != nil {
				count++
			}
		}
	}
	dst.piece_count = C.int(count)
	if count == 0 {
		dst.pieces = nil
		return 0
	}
	itemSize := unsafe.Sizeof(C.PieceView{})
	raw := C.malloc(C.size_t(count) * C.size_t(itemSize))
	slice := (*[1 << 30]C.PieceView)(raw)[:count:count]
	idx := 0
	for x := range entities {
		for y := range entities[x] {
			piece := entities[x][y]
			if piece == nil {
				continue
			}
			_ = FillPieceView(&slice[idx], piece)
			idx++
		}
	}
	dst.pieces = (*C.PieceView)(raw)
	return 0
}

//export FreeCardView
func FreeCardView(view *C.CardView) {
	if view == nil {
		return
	}
	if view.tag != nil {
		C.free(unsafe.Pointer(view.tag))
	}
	if view.uuid != nil {
		C.free(unsafe.Pointer(view.uuid))
	}
	if view.costs != nil {
		count := int(view.cost_count)
		costs := (*[1 << 30]C.MaterialCostView)(unsafe.Pointer(view.costs))[:count:count]
		_ = costs
		C.free(unsafe.Pointer(view.costs))
	}
}

//export FreeBoardView
func FreeBoardView(view *C.BoardView) {
	if view == nil {
		return
	}
	if view.pieces != nil {
		count := int(view.piece_count)
		pieces := (*[1 << 30]C.PieceView)(unsafe.Pointer(view.pieces))[:count:count]
		for i := range pieces {
			if pieces[i].id != nil {
				C.free(unsafe.Pointer(pieces[i].id))
			}
			if pieces[i].owner_id != nil {
				C.free(unsafe.Pointer(pieces[i].owner_id))
			}
		}
		C.free(unsafe.Pointer(view.pieces))
	}
}
