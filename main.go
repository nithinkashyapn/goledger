package main

import (
	"fmt"
	"os"

	"./blockchain"
	"github.com/chainHero/heroes-service/web/controllers"
)

func main() {
	fSetup := blockchain.FabricSetup{
		OrdererID:       "orderer.example.com",
		ChannelID:       "channelone",
		ChannelConfig:   "./fixtures/channel-artifacts/channelone.tx",
		ChaincodeID:     "example-service",
		ChaincodeGoPath: os.Getenv("GOPATH"),
		ChaincodePath:   "github.com/hyperledger/fabric/chaincode/",
		OrgAdmin:        "Admin",
		OrgName:         "org1",
		ConfigFile:      "config.yaml",
		UserName:        "User1",
	}

	err := fSetup.Initialize()
	if err != nil {
		fmt.Println("Unable to initialize SDK %v\n", err)
		return
	}

	defer fSetup.CloseSDK()

	err = fSetup.InstallAndInstantiateCC()
	if err != nil {
		fmt.Printf("Unable to install and instantiate chaincode %v\n", err)
	}

	app := &controllers.Application{
		Fabric: &fSetup,
	}
	web.serve(app)
}
