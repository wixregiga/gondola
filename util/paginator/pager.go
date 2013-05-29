package paginator

const (
	CURRENT   = 1 << 0
	DISABLED  = 1 << 1
	SEPARATOR = 1 << 2
	PREVIOUS  = 1 << 3
	NEXT      = 1 << 4
)

type Pager interface {
	Root() *Node
	Href(base string, page int) string
	Node(n *Node, page int, flags int) *Node
}
