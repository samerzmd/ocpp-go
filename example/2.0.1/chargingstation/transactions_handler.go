package main

import (
	"github.com/samerzmd/ocpp-go/ocpp"
	"github.com/samerzmd/ocpp-go/ocpp2.0.1/transactions"
	"github.com/samerzmd/ocpp-go/ocppj"
)

func (handler *ChargingStationHandler) OnGetTransactionStatus(request *transactions.GetTransactionStatusRequest) (response *transactions.GetTransactionStatusResponse, err error) {
	logDefault(request.GetFeatureName()).Warnf("Unsupported feature")
	return nil, ocpp.NewHandlerError(ocppj.NotSupported, "Not supported")
}
