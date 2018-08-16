package blockchain

import (
	"fmt"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
	"github.com/pkg/errors"

	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	packager "github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
)

type FabricSetup struct {
	ConfigFile      string
	OrgID           string
	OrdererID       string
	ChannelID       string
	ChaincodeID     string
	initialized     bool
	ChannelConfig   string
	ChaincodeGoPath string
	ChaincodePath   string
	OrgAdmin        string
	OrgName         string
	UserName        string
	client          *channel.Client
	admin           *resmgmt.Client
	sdk             *fabsdk.FabricSDK
	event           *event.Client
}

func (setup *FabricSetup) Initialize() error {

	if setup.initialized {
		return errors.New("SDK already initialized")
	}

	sdk, err := fabsdk.New(config.FromFile(setup.ConfigFile))
	if err != nil {
		return errors.WithMessage(err, "Failed to create SDK")
	}
	setup.sdk = sdk
	fmt.Println("---- SDK CREATED ----")

	resourceManagerClientContext := setup.sdk.Context(fabsdk.WithUser(setup.OrgAdmin), fabsdk.WithOrg(setup.OrgName))
	if err != nil {
		return errors.WithMessage(err, "Failed to load admin")
	}

	resMgmtClient, err := resmgmt.New(resourceManagerClientContext)
	if err != nil {
		return errors.WithMessage(err, "Failed to create channel management client w/ ablove admin credentials")
	}
	setup.admin = resMgmtClient
	fmt.Println("---- RESOURCE MANAGEMENT CLIENT CREATED ----")

	mspClient, err := mspclient.New(sdk.Context(), mspclient.WithOrg(setup.OrgName))
	if err != nil {
		return errors.WithMessage(err, "Failed to create MSP Client")
	}
	adminIdentity, err := mspClient.GetSigningIdentity(setup.OrgAdmin)
	if err != nil {
		return errors.WithMessage(err, "Failed to get admin signing identity")
	}

	req := resmgmt.SaveChannelRequest{
		ChannelID:         setup.ChannelID,
		ChannelConfigPath: setup.ChannelConfig,
		SigningIdentities: []msp.SigningIdentity{adminIdentity}}
	txID, err := setup.admin.SaveChannel(req, resmgmt.WithOrdererEndpoint(setup.OrdererID))
	if err != nil || txID.TransactionID == "" {
		return errors.WithMessage(err, "Failed to make admin join channel")
	}
	fmt.Println("---- ADMIN JOINED CHANNEL ----")

	fmt.Println("---- INITIALIZATION SUCCESSFUL ----")
	setup.initialized = true
	return nil
}

func (setup *FabricSetup) InstallAndInstantiateCC() error {

	ccPkg, err := packager.NewCCPackage(setup.ChaincodePath, setup.ChaincodeGoPath)
	if err != nil {
		return errors.WithMessage(err, "Failed to create Chaincode Package")
	}
	fmt.Println("---- CREATED CHAINCODE PACKAGE ----")

	installCCReq := resmgmt.InstallCCRequest{
		Name:    setup.ChaincodeID,
		Path:    setup.ChaincodePath,
		Version: "0",
		Package: ccPkg}
	_, err = setup.admin.InstallCC(installCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		return errors.WithMessage(err, "Failed to install chaincode")
	}
	fmt.Println("---- INSTALLED CHAINCODE ----")

	ccPolicy := cauthdsl.SignedByAnyMember([]string{"org1.example.com"})
	resp, err := setup.admin.InstantiateCC(setup.ChannelID, resmgmt.InstantiateCCRequest{
		Name:    setup.ChaincodeID,
		Path:    setup.ChaincodeGoPath,
		Version: "0",
		Args:    [][]byte{[]byte("init")},
		Policy:  ccPolicy})
	if err != nil || resp.TransactionID == "" {
		return errors.WithMessage(err, "Failed to instantiate chaincode")
	}
	fmt.Println("---- INSTANTIATED CHAINCODE ----")

	clientContext := setup.sdk.ChannelContext(setup.ChaincodeID, fabsdk.WithUser(setup.UserName))
	setup.client, err = channel.New(clientContext)
	if err != nil {
		return errors.WithMessage(err, "Failed to create new channel client")
	}
	fmt.Println("---- CHANNEL CLIENT CREATED ----")

	setup.event, err = event.New(clientContext)
	if err != nil {
		return errors.WithMessage(err, "Failed to create new event client")
	}
	fmt.Println("---- EVENT CLIENT CREATED ----")

	fmt.Println("---- CHAINCODE INSTALLED AND INSTANTIATED ----")
	return nil
}

func (setup *FabricSetup) CloseSDK() {
	setup.sdk.Close()
}
