package vfs

// routes are helpers to manage writes to a database
type route struct {
	parts []*routePart
}

type routePart struct {
	parent *route
	nodes []*routeNode
	next *id
}

type routeNode struct {
	parent *routePart
	idx uint8
	id *id
	branch *branch
}

func (r *route) file() *id {
	if len(r.parts) <= 0 {
		return nil
	}
	return r.parts[len(r.parts)-1].next
}

func (p *routePart) add(dir *directory) {
}

func (p *routeNode) cacheKey() string {
	return ""
}
