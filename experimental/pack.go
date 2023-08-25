package main

import (
	"fmt"
	"image"
	"sort"
)

type Packable interface {
	Len() int
	Size(n int) (width, height int)
	Place(n, x, y int)
}

func Pack(p Packable) (width, height int) {
	numBlocks := p.Len()

	if numBlocks == 0 {
		return 0, 0
	}

	w, h := p.Size(0)
	root := &node{
		x:      0,
		y:      0,
		width:  w,
		height: h,
	}

	p.Place(0, 0, 0)

	for i := 0; i < numBlocks; i++ {
		w, h = p.Size(i)

		node := root.find(w, h)
		if node != nil {
			node = node.split(w, h)

			// Update block in-place
			p.Place(i, node.x, node.y)

		} else {
			newRoot, grown := root.grow(w, h)
			if newRoot == nil {
				return -1, -1
			}

			// Update block in-place
			p.Place(i, grown.x, grown.y)

			root = newRoot
		}
	}

	return root.width, root.height
}

type node struct {
	x, y, width, height int
	right, down         *node
}

func (n *node) find(width, height int) *node {
	if n.right != nil || n.down != nil {
		right := n.right.find(width, height)
		if right != nil {
			return right
		}
		return n.down.find(width, height)
	} else if width <= n.width && height <= n.height {
		return n
	}
	return nil
}

func (n *node) split(width, height int) *node {
	n.down = &node{
		x:      n.x,
		y:      n.y + height,
		width:  n.width,
		height: n.height - height,
	}

	n.right = &node{
		x:      n.x + width,
		y:      n.y,
		width:  n.width - width,
		height: height,
	}

	return n
}

func (n *node) grow(width, height int) (root, grown *node) {
	canGrowDown := width <= n.width
	canGrowRight := height <= n.height

	// attempt to keep square-ish by growing right when height is much greater than width
	shouldGrowRight := canGrowRight && (n.height >= (n.width + width))

	// attempt to keep square-ish by growing down when width is much greater than height
	shouldGrowDown := canGrowDown && (n.width >= (n.height + height))

	if shouldGrowRight {
		return n.growRight(width, height)
	} else if shouldGrowDown {
		return n.growDown(width, height)
	} else if canGrowRight {
		return n.growRight(width, height)
	} else if canGrowDown {
		return n.growDown(width, height)
	}

	// need to ensure sensible root starting size to avoid this happening
	return nil, nil
}

func (n *node) growRight(width, height int) (root, grown *node) {
	newRoot := &node{
		x:      0,
		y:      0,
		width:  n.width + width,
		height: n.height,
		down:   n,
		right: &node{
			x:      n.width,
			y:      0,
			width:  width,
			height: n.height,
		},
	}

	node := newRoot.find(width, height)
	if node != nil {
		return newRoot, node.split(width, height)
	}
	return nil, nil
}

func (n *node) growDown(width, height int) (root, grown *node) {
	newRoot := &node{
		x:      0,
		y:      0,
		width:  n.width,
		height: n.height + height,
		down: &node{
			x:      0,
			y:      n.height,
			width:  n.width,
			height: height,
		},
		right: n,
	}

	node := newRoot.find(width, height)
	if node != nil {
		return newRoot, node.split(width, height)
	}
	return nil, nil
}

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

func (b *Block) Add(item interface{}, size image.Point) {
	b.items = append(b.items, item)
	b.sizes = append(b.sizes, size)
	b.place = append(b.place, image.Point{})
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

func main() {
	b := NewBlock()
	b.Add("c", image.Point{X: 5, Y: 3})   // 2
	b.Add("b", image.Point{X: 7, Y: 7})   // 1
	b.Add("e", image.Point{X: 2, Y: 1})   // 4
	b.Add("a", image.Point{X: 10, Y: 10}) // 0
	b.Add("d", image.Point{X: 2, Y: 14})  // 3

	for i := 0; i < b.Len(); i++ {
		fmt.Println(b.items[i].(string), " ", b.sizes[i], " ", b.place[i])
	}

	fmt.Println("")
	b.sort()
	fmt.Println("")

	for i := 0; i < b.Len(); i++ {
		fmt.Println(b.items[i].(string), " ", b.sizes[i], " ", b.place[i])
	}
}
