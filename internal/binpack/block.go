package binpack

import (
	"image"
	"sort"
)

// Block let's us binpack anything really, given some set of items and some idea of
// what size (x,y) means for it.
type Block struct {
	items []interface{}
	sizes []image.Point
	place []image.Point
}

func NewBlock() *Block {
	return &Block{
		items: make([]interface{}, 0),
		sizes: make([]image.Point, 0),
		place: make([]image.Point, 0),
	}
}

func (b *Block) Add(item interface{}, x, y int) {
	b.items = append(b.items, item)
	b.sizes = append(b.sizes, image.Point{x, y})
	b.place = append(b.place, image.Point{})
}

func (b *Block) Get(n int) (interface{}, int, int) {
	if n < 0 || n >= len(b.place) {
		return nil, -1, -1
	}
	return b.items[n], b.place[n].X, b.place[n].Y
}

func (b *Block) sort() {
	indexes := make([]int, len(b.sizes))
	for i := range indexes {
		indexes[i] = i
	}

	sort.Slice(indexes, func(i, j int) bool {
		a := b.sizes[indexes[i]]
		b := b.sizes[indexes[j]]
		if a.X > b.X { // we want big -> small
			return true
		} else if a.X == b.X {
			return a.Y > b.Y
		} else {
			return false
		}
	})

	// TODO: swap in-place to save memory
	newItems := make([]interface{}, len(b.items))
	newSizes := make([]image.Point, len(b.sizes))
	newPlace := make([]image.Point, len(b.place))

	for dst, src := range indexes {
		newItems[dst] = b.items[src]
		newSizes[dst] = b.sizes[src]
		newPlace[dst] = b.place[src]
	}

	b.items = newItems
	b.sizes = newSizes
	b.place = newPlace
}

func (b *Block) Len() int {
	return len(b.items)
}

func (b *Block) Size(n int) (width, height int) {
	if n < 0 || n >= len(b.sizes) {
		return -1, -1
	}
	return b.sizes[n].X, b.sizes[n].Y
}

func (b *Block) Place(n, x, y int) {
	if n < 0 || n >= len(b.place) {
		return
	}
	b.place[n] = image.Point{x, y}
}
