# node-label-controller

## Goal:
We want to use the *ContainerLinuxUpdateOperator*.

## Problem:
We have a cluster with a mix of Nodes. 1x Ubuntu, 1x CentOS and 1x `ContainerLinux`.
The *ContainerLinuxUpdateOperator* should only run on `ContainerLinux` nodes.
Our initial idea was to use a DaemonSet with a NodeSelector, but unfortunately the Kubernetes Node
Object has no label for the Operating system.

## Idea:
Have a controller which watches the Kubernetes nodes and sets a label to the Kubernetes Node object when the node uses ContainerLinux as operating system.

## Task
- Write a controller which:
    - Watches the Kubernetes Node objects
    - Check if any Node uses ContainerLinux
    - Attaches a label to the Node if it uses _ContainerLinux_ (`kubermatic.io/uses-container-linux: 'true'`)

- Deploy the *ContainerLinuxUpdateOperator* (but set the NodeSelector) using the manifests from the repository
- Validate that the agent of the *ContainerLinuxUpdateOperator* only gets deployed on nodes with the `kubermatic.io/uses-container-linux: 'true'` label.
- Write a Dockerfile for the controller
- Write a Kubernetes Deployment for the controller
- Write the RBAC manifests which are required for the controller

## Relevant information
- The operating system can be found inside the Kubernetes Node object
- The controller must be written in Go
- You will be given a Kubernetes cluster to test with. Each node in the cluster uses a different operating system.
- When you are done - upload your result to Github or send us the code via email

## Some good starting points:
- https://github.com/coreos/container-linux-update-operator
- https://github.com/kubernetes-sigs/kubebuilder
- https://github.com/kubernetes/client-go/tree/master/examples/workqueue
- https://github.com/kubernetes/sample-controller

