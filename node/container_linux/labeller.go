package container_linux

import (
	"fmt"
	"strings"

	v1 "k8s.io/api/core/v1"

	"github.com/mwarzynski/loodse-k8s-node-label-controller/node"
)

const (
	LabelUsesContainerLinuxKey   = "kubermatic.io/uses-container-linux"
	LabelUsesContainerLinuxValue = "true"
)

type Labeller struct {
	updater node.Updater
}

func NewLabeller(updater node.Updater) *Labeller {
	return &Labeller{
		updater: updater,
	}
}

func (l *Labeller) Name() string {
	return "container-linux-node-labeller"
}

func (l *Labeller) ProcessNode(node *v1.Node) error {
	if node == nil {
		return nil
	}

	isContainerLinuxNode := strings.Contains(node.Status.NodeInfo.OSImage, "Container Linux by CoreOS")
	labelValue, labelFound := node.Labels[LabelUsesContainerLinuxKey]
	if (isContainerLinuxNode && labelValue == LabelUsesContainerLinuxValue) ||
		(!isContainerLinuxNode && !labelFound) {
		return nil
	}

	// Copy Node as not to operate on the cached one (and pollute the data).
	node = node.DeepCopy()
	if isContainerLinuxNode {
		node.Labels[LabelUsesContainerLinuxKey] = LabelUsesContainerLinuxValue
	} else {
		delete(node.Labels, LabelUsesContainerLinuxKey)
	}

	_, err := l.updater.Update(node)
	if err != nil {
		return fmt.Errorf("couldn't update node %q: %w", node.GetName(), err)
	}

	return nil
}
