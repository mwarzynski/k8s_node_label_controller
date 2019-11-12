package container_linux_test

import (
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	containerLinux "github.com/mwarzynski/loodse-k8s-node-label-controller/node/container_linux"
)

func TestLabeler(t *testing.T) {
	tests := []struct {
		name           string
		nodeOS         string
		nodeLabels     map[string]string
		expectedUpdate bool
		expectedLabels map[string]string
	}{
		{
			name:           "set correct label for Container Linux node when there are no labels",
			nodeOS:         "Container Linux by CoreOS 2247.6.0 (Rhyolite)",
			expectedUpdate: true,
			expectedLabels: map[string]string{
				containerLinux.LabelUsesContainerLinuxKey: containerLinux.LabelUsesContainerLinuxValue,
			},
		},
		{
			name:   "Container Linux node has proper labels, do not update",
			nodeOS: "Container Linux by CoreOS 2247.6.0 (Rhyolite)",
			nodeLabels: map[string]string{
				containerLinux.LabelUsesContainerLinuxKey: containerLinux.LabelUsesContainerLinuxValue,
			},
			expectedUpdate: false,
			expectedLabels: map[string]string{
				containerLinux.LabelUsesContainerLinuxKey: containerLinux.LabelUsesContainerLinuxValue,
			},
		},
		{
			name:   "unset label if it's not a Container Linux node",
			nodeOS: "Ubuntu 19.04 LTS",
			nodeLabels: map[string]string{
				containerLinux.LabelUsesContainerLinuxKey: containerLinux.LabelUsesContainerLinuxValue,
			},
			expectedUpdate: true,
			expectedLabels: map[string]string{},
		},
		{
			name:           "not a Container Linux node (without container linux labels), do not update",
			nodeOS:         "Ubuntu 19.04 LTS",
			expectedUpdate: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var node *v1.Node
			// Mock Node Updater to check if update action was executed.
			nodeWasUpdated := false
			updater := &mockUpdater{
				UpdateFunc: func(updatedNode *v1.Node) (*v1.Node, error) {
					nodeWasUpdated = true
					node = updatedNode
					return updatedNode, nil
				},
			}

			// Create (Container Linux) Labeler component.
			labeler := containerLinux.NewLabeler(updater)
			if labeler.Name() == "" {
				t.Errorf("invalid name")
			}

			// Prepare and Process the Node.
			node = prepareNode(test.nodeLabels, test.nodeOS)
			err := labeler.ProcessNode(node)
			if err != nil {
				t.Fatal(err)
			}

			// Validate if all actions went fine.
			if nodeWasUpdated != test.expectedUpdate {
				t.Errorf("updater was updated=%v, but should be updated=%v", nodeWasUpdated, test.expectedUpdate)
			}
			if len(test.expectedLabels) != len(node.Labels) {
				t.Error("labels are different")
			}
			for k, v := range node.Labels {
				cv := test.expectedLabels[k]
				if v != cv {
					t.Errorf("%s: values are different: %q, %q", k, v, cv)
				}
			}
		})
	}
}

func prepareNode(labels map[string]string, OSImage string) *v1.Node {
	if labels == nil {
		labels = make(map[string]string)
	}
	node := &v1.Node{
		Status: v1.NodeStatus{
			NodeInfo: v1.NodeSystemInfo{
				OSImage: OSImage,
			},
		},
		ObjectMeta: metav1.ObjectMeta{
			Labels: labels,
		},
	}
	return node
}

type mockUpdater struct {
	UpdateFunc func(node *v1.Node) (*v1.Node, error)
}

func (mu *mockUpdater) Update(node *v1.Node) (*v1.Node, error) {
	if mu.UpdateFunc != nil {
		return mu.UpdateFunc(node)
	}
	return node, nil
}
