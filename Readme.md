## Change the Race Course APP.
Need to add "data: raceContract.unlinked_binary" to the raceContracts.new()


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

Currently RaceCourse yaml takes 2 spec fields. Number of Replicas and Image Name. Need to add resource request to speed up the
boot time.

There is no automated way to initialize the contract on the quorum nodes. We need to manually exec into quorum nodes and run

geth account new

geth attach qdata/etherum/geth.ipc

raceAbi = ""
raceBin = "0x"
raceInterface = eth.contract(raceAbi)
raceHex = raceInterface.new({from: eth.accounts[0] , data: raceHex , gas: 1000000 })

"npm run start" takes a while for the racecourse app.

TODO:
Add a ready check for the racecourse deployment.
Modify the Spec field of race course app to add resource request.

## Currently Working:
once the RaceCourse app is listening on port 3000. We can do http://quorum-service.kl-quorum:22001 to connect to the blockchain.
