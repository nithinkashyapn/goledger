# Goledger

A complete implementation of hyperledger in Go.

+ Create the `chaincode` folder in root directory
```
mkdir chaincode && cd chaincode
```

+ Write the chaincode in the `main.go` file
```
nano main.go
```

+ Create the `blockchain` directory at the root
```
mkdir blockchain && cd blockchain
```

+ Create the `setup.go` file - This handled dynamic channel creation and CC install and instantiate
```
nano setup.go
```

+ Create the helper files - This handles data input and query w/ the blockchain
```
nano query.go 
nano invoke.go
```
