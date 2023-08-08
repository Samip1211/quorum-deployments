## Change the Race Course APP.
The current race course app, throws error while deploying contract, when a new connection is made to blockchain.
To fix this update the following struct of raceContact 
```json
raceContract.new({from: web3.eth.accounts[0], gasPrice: 0, gas: 5000000})
```
To
```json
raceContract.new({from: web3.eth.accounts[0], gasPrice: 0, gas: 5000000, data: raceContract.unlinked_binary})
```

Then build the image using Dockerfile provided in the race-course-app-artifacts. 
Note: Copy the Dockerfile to the racecourse directory.

## Build Operator Image.
From this directory build the Operator Image.
```bash
cd operator && make docker-build
```


## Build Artifacts for RaceCourse Operator Deployment.
From this directory build the Operator Artifacts. This generates "racecourse-operator.yaml". This deploys the 
Racecourse CRD, some service accounts, some CL roles and CL role binding for the racecourse operator.

This operator listens on the "Racecourse CR".
```bash
cd operator && make artifacts
```


## Delete Previous Kind Cluster
Clean up previous Kind Deployments.
```bash
kind delete cluster
```


## Bring up Quorum Cluster. This should also load all the images to K8s cluster.
Before running this make sure to update "/Users/samip/Documents/kl/quorum-tools/examples/qdata_1/" path in kind/image.yaml
to the quorum-tools/ of your local machine. This mounts the Host's Local Directory to the Kind's Cluster Node. Using that 
we can mount this path to subsequent containers in the Kind Cluster.

```bash
sh local-reg.sh
```


## Bring up quorum Nodes.
Once the Kind Cluster is ready, deploy the following YAML. This is similar to running examples/docker-compose.yaml. In this yaml, the changes have been made so that quorum nodes can connect to boot-node. That is done using service.

```bash
k apply -f quorum-nodes.yaml
```

## Deploy Contract Manually.
Once quorum nodes are up, deploy the contract manually.

## Apply RaceCourse Manifest.
Once the quorum nodes are up and running, Deploy the following. This brings up race course operator. This is responsible for
bringing up racecourse application.

```bash
k apply -f racecourse-operator.yaml
```


## Apply RaceCourse Yaml.
Apply the following Yaml. This brings up the racecourse application.

```bash
k apply -f racecourse.yaml -n operator-system
```

Once yaml is applied check it's status. You can do that by

```bash
k get racecourses.kaleido.kaleido.com -o jsonpath="{.items[0].status}"
```

Once deploymentStatus is true, port-forward deployment or svc on 3000, to connect to racecourse app.
```bash
k port-forward -n operator-system svc/racecourse 3000
```

NOTE: It may take a while currently as cpu and mem requests/limits are not defiened on racecoruse app. However we can connect to it once the log mentions, Server started at port 3000. (I would like to add ready check, to make sure this process is easy.)

## Notes

1. Currently RaceCourse yaml takes 2 spec fields. Number of Replicas and Image Name. Need to add resource request to speed up the
boot time.

2. There is no automated way to initialize the contract on the quorum nodes. We need to manually exec into quorum nodes and run

## Deploying contract manulally.
geth account new

geth attach qdata/etherum/geth.ipc

raceAbi = ""
raceBin = "0x"
raceInterface = eth.contract(raceAbi)
raceHex = raceInterface.new({from: eth.accounts[0] , data: raceHex , gas: 1000000 })

"npm run start" takes a while for the racecourse app.

## TODO:
Add a ready check for the racecourse deployment.
Modify the Spec field of race course app to add resource request.

## Currently Working:
Once the RaceCourse app is listening on port 3000. We can do http://quorum-service.kl-quorum:22001 to connect to the blockchain.


## Issues being faced.
When we click "connect", it indefinetly waits. This is because, the contract is being deployed, however the transaction is not
getting mined. It works fine when doing docker-compose up. However even after trying various combinations and all, the transaction is stuck in pending state. I've verfied that there's no issue in connection, it fails to get mined when there's only 1 quorum node too. I even tried to see if there's CPU throttling going on, or CPU request is less. Yet CPU usage metrics looked fine. So not sure why it's not gettting mined.