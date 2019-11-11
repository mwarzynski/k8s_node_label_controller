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

type Labeler struct {
	updater node.Updater
}

func NewLabeler(updater node.Updater) *Labeler {
	return &Labeler{
		updater: updater,
	}
}

func (l *Labeler) Name() string {
	return "container-linux-node-labeler"
}

func (l *Labeler) ProcessNode(node *v1.Node) error {
	if node == nil {
		return nil
	}

	isContainerLinuxNode := strings.Contains(node.Status.NodeInfo.OSImage, "Container Linux by CoreOS")
	labelValue, labelFound := node.Labels[LabelUsesContainerLinuxKey]
	if (isContainerLinuxNode && labelValue == LabelUsesContainerLinuxValue) ||
		(!isContainerLinuxNode && !labelFound) {
		return nil
	}

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
