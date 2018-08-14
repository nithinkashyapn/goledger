# Goledger

A complete implementation of hyperledger in Go.

+ Create the project structure
    
    + blockchain - `helper` functions
    + chaincode - `chaincode` to interact with BC
    + fixtures - `network` files
    + web - `frontend`

+ Create the `fixtures` folder
```
mkdir fixtures && cd fixtures
```

+ Copy the `bin` folder respective to the computer
```

```

+ Create the `crypto-config.yaml` file
```
vi crypto-config.yaml
```

+ Generate certificates
```
./bin/cryptogen generate --config=./crypto-config.yaml
```

+ Create the `configtx.yaml` file
```
vi configtx.yaml
```

+ Create `channel-artifacts` folder
```
mkdir channel-artifacts
```

+ Create genesis block
```
./bin/configtxgen -profile OrdererGenesis -outputBlock ./channel-artifacts/genesis.block
```

+ Generate channel transactions
```
./bin/configtxgen -profile ChannelOne -outputCreateChannelTx ./channel-artifacts/channelone.tx -channelID channelone
```

+ Generate anchor peer
```
./bin/configtxgen -profile ChannelOne -outputAnchorPeersUpdate ./channel-artifacts/Org1MSPanchors_channelone.tx -channelID channelone -asOrg Org1MSP
```

+ Create `docker-compose.yaml` file
```
vi docker-compose.yaml
```

+ Start the network
```
docker-compose up -d
```

+ Verify the network
```
docker ps -a
```

> Yay!! You've setup the network.
> Continue with README_2.md