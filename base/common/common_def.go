package common

const(
	MethodServiceCenterRegister = "ServiceCenter.Register"
	MethodServiceCenterDispatch = "ServiceCenter.Dispatch"
	MethodServiceNodeCall       = "ServiceNode.Call"
)

type ServiceCenterRegisterData struct {
	Name string `json:"name"`
	Addr string `json:"addr"`
}

type ServiceCenterDispatchData struct{
	Name string `json:"name"`
	Params string `json:"params"`
}