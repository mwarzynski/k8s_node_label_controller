package node

import v1 "k8s.io/api/core/v1"

type Updater interface {
	Update(node *v1.Node) (*v1.Node, error)
}
