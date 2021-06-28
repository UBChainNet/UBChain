package types

type RpcTransferBody struct {
	Contract  string        `json:"contract"`
	Receivers []RpcReceiver `json:"receivers"`
}

type RpcReceiver struct {
	Address string `json:"address"`
	Amount  uint64 `json:"amount"`
}
