# node-label-controller

Note: I described my **line of thinking** below the task description.

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


# My line of thinking

###  Kubernetes Controller

I have some intuition about what is the Kubernetes Controller (it watches the Kubernetes objects and may interract with them).
We need to code it, but what is it actually (based on documentation)?

**Controllers are control loops that watch the state of your cluster, then make or request changes where needed.**
(I briefly read: https://kubernetes.io/docs/concepts/architecture/controller/)

In our case, controller needs to 'watch te state' of Nodes inside the cluster and make requests to set/unset labels.

### NodeLabeler Controller

I copy pasted most of the 'watcher' code from `workqueue`. ContainerLinux labels management is described best by source
code: `node/container_linux/labeler.go` (also: see tests).

Thank you for a helpful 'Relevant information' section.

### ContainerLinux update operator

> Deploy the *ContainerLinuxUpdateOperator* (but set the NodeSelector) using the manifests from the repository.

I added (modified) manifests to the `k8s/container-linux-update-operator` folder.

> Validate that the agent of the *ContainerLinuxUpdateOperator* only gets deployed on nodes with the `kubermatic.io/uses-container-linux: 'true'` label.

```
(⎈|loodse:default) $ k describe pods container-linux-update-agent-sjn82 -n reboot-coordinator | grep Node
Node:           ip-172-31-11-14.eu-central-1.compute.internal/172.31.11.14
Node-Selectors:  kubermatic.io/uses-container-linux=true
(⎈|loodse:default) $ k describe nodes ip-172-31-11-14.eu-central-1.compute.internal | grep 'OS Image'
 OS Image:                   Container Linux by CoreOS 2247.6.0 (Rhyolite)
(⎈|loodse:default) $ k describe pods container-linux-update-operator-749954844-7hdpc -n reboot-coordinator | grep Node
Node:           ip-172-31-11-14.eu-central-1.compute.internal/172.31.11.14
Node-Selectors:  kubermatic.io/uses-container-linux=true
(⎈|loodse:default) $ k describe nodes ip-172-31-11-14.eu-central-1.compute.internal | grep 'OS Image'
 OS Image:                   Container Linux by CoreOS 2247.6.0 (Rhyolite)
```

### Makefile

`make`: fetches dependencies, runs tests and builds the controller into the `.bin/` folder. I assume you have correctly
set up the Golang. (Otherwise, to build the binary I would use the `golang` image from Dockerhub.)

`make docker_image`: creates a new Docker Image tagged as `container-linux-node-labeler:0.0.1`. Normally the tag would
be evaluated based on the git tags, but it's fine for the purpose of interview challenge.

### Running inside Kubernetes

Firstly, I pushed image to the Dockerhub (https://hub.docker.com/repository/docker/mwarzynski/container-linux-node-labeller). I did it manually, but it should be set up with a CD pipeline (which should watch the git tags and maybe push docker images accordingly).

At this point I didn't know how to configure controller in order to work inside Kubernetes cluster. However, there is some magic inside which falls back to the `rest.InClusterConfig()`. Anyway, it is as simple as running a binary without any arguments.

I decided to deploy my toys to a separate namespace `node-labels`.

To start somewhere, I created a deployment for the `mwarzynski/container-linux-node-labeller:0.0.1`. At the beginning tt just said: Kubernetes should deploy one instance of this image somewhere. However, running pods didn't work, because by default it had no permissions to manage the Nodes (at least fetch and modify their labels).

It brings us to `ClusterRole` (btw: `Role` doesn't allow to set permissions for cluster-wide resources as Nodes).
I added a new `ClusterRole` named `manage-nodes` which allowed to manage `nodes` with following
actions: `["get", "list", "watch", "update"]` (no idea how the workqueue works underneath; maybe not all of them are required).

Yeah, cool that we have a cluster-wide role, but actually how to bind the `ClusterRole` with our controller's deployment?

Well, we need to use `ClusterRoleBinding` with a defined `Subject` (Group, User or ServiceAccount). I did a few Google searches and landed at Kubernetes documentation page 'Managing Service Accounts' (https://kubernetes.io/docs/reference/access-authn-authz/service-accounts-admin/).
> User accounts are for humans. Service accounts are for processes, which run in pods.

So, we should use the `ServiceAccount`, shouldn't we? I hope that executable binary file isn't a human yet.
I created a `node-manager` `ServiceAccount` inside our `node-labels` namespace. And also the `ClusterRoleBinding` which
made the ServiceAccount so powerful that anything using this account might access (and change!) the Nodes.

Finally, we have to tell Kubernetes that controller's Deployment should use this 'powerful' ServiceAccount: `serviceAccountName: node-manager`.

PS. one-liner: `kubectl apply -f ./k8s`.
