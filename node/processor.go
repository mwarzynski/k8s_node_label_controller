package node

import (
	v1 "k8s.io/api/core/v1"
)

type Processor interface {
	Name() string
	ProcessNode(node *v1.Node) error
}
