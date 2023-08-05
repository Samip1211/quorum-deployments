## Build Operator Image.
cd operator && make docker-build

## Build Artifacts for RaceCourse Operator Deployment.
cd operator && make artifacts

## Delete Previous Kind Cluster
kind delete cluster

## Bring up Quorum Cluster. This should also load all the images to K8s cluster.

sh local-reg.sh

## Bring up Boot Node.
k apply -f bootnode.yaml


## Bring up quorum Nodes.
k apply -f quorum-nodes.yaml


## Apply RaceCourse Manifest.
k apply -f racecourse-operator.yaml

## Apply RaceCourse Yaml.
k apply -f racecourse.yaml
