package controllers

import (
	"net/http"
	"../../blockchain"
)

type Application struct {
	Fabric *blockchain.FabricSetup
}

func renderTemplate(w http.ResponseWriter, r *http.Request, templateName string, data interface{}){
	
}