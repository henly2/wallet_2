package common

import (
	"reflect"
)

const(
	MethodServiceCenterRegister = "ServiceCenter.Register"	// 服务向服务中心注册请求，对内
	MethodServiceCenterDispatch = "ServiceCenter.Dispatch"	// 客户向服务中心发送请求，对外
	MethodServiceCenterPingpong = "ServiceCenter.Pingpong"	// 服务向服务中心发送心跳，对内
	MethodServiceNodeCall       = "ServiceNode.Call"		// 服务中心向服务发送请求，对内
)

// 注册信息
type ServiceCenterRegisterData struct {
	Name string `json:"name"`			// service node name
	Addr string `json:"addr"`			// service node ip address
	Apis []string `json:"apis""`  		// service node apis
}

// 注册API
func (rd *ServiceCenterRegisterData)RegisterApi(api interface{})  {
	t := reflect.TypeOf(api)
	v := reflect.ValueOf(api)

	tName := reflect.Indirect(v).Type().Name()
	for m := 0; m < t.NumMethod(); m++ {
		method := t.Method(m)
		mName := method.Name

		rd.Apis = append(rd.Apis, tName+"."+mName)
	}
}

// 请求信息，作为rpc请求的params数据
// json like: {"api":"Arith.Add", "argv":"", "id":1}
type ServiceCenterDispatchData struct{
	Api  string `json:"api"`  // like "xxx.xxx"
	Argv string `json:"argv"` // json string
	Id   int    `json:"id"`     // id
}

// 应答信息，作为rpc应答的result数据
// json like: {"api":"Arith.Add", "id":1, "err":0, "errmsg":"", "ack":""}
type ServiceCenterDispatchAckData struct{
	Api    string `json:"api"`    // like "xxx.xxx"
	Id     int    `json:"id"`       // id
	Err    int    `json:"err"`      // like 0
	ErrMsg string `json:"errmsg"`   // string
	Ack    string `json:"ack"` // json string
}