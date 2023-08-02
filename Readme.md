## Delete Previous Kind Cluster
kind delete cluster

## Bring up Quorum Cluster.
sh local-reg.sh

## Bring up Boot Node.
k apply -f bootnode.yaml


## Bring up quorum Nodes.
k apply -f quorum-nodes.yaml