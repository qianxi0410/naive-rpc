package selector

type Selector interface {
	Select(service string) (*Node, error)
	Update(node *Node, err error) error
}

// represent a host
type Node struct {
	Network string
	Address string
}
