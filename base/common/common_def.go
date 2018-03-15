package common

const(
	MethodServiceCenterRegister = "ServiceCenter.Register"
	MethodServiceCenterDispatch = "ServiceCenter.Dispatch"
	MethodServiceNodeCall       = "ServiceNode.Call"
)

type ServiceCenterRegisterData struct {
	Name string `json:"name"`
	Addr string `json:"addr"`
	Apis []string `json:"apis""`
}

type ServiceCenterDispatchData struct{
	Api string `json:"api"`
	Params string `json:"params"`
}