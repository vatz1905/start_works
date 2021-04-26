package aria2Downloader

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Params_dl_msg struct {
	Jsonrpc string        `json:"jsonrpc"` // added for case insensitive to match aria2 jsonrpc format (lower case)
	Id      string        `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

type Ret_gid_msg struct {
	Jsonrpc string
	Id      string
	Result  string
}

type Ret_version_msg struct {
	Jsonrpc string
	Id      string
	Result  Ver_result
}
type Ver_result struct {
	EnabledFeatures []string
	Version         string
}

// for aria selective tellstatus
type Ret_status_msg struct {
	Jsonrpc string      `json:"jsonrpc"`
	Id      string      `json:"id"`
	Result  Result_type `json:"result"`
}
type Result_type struct {
	Gid             string `json:"gid"`
	CompletedLength string `json:"completedLength"`
	TotalLength     string `json:"totalLength"`
	Status          string `json:"status"`
	Connections     string `json:"connections"`
}

var ret_gid_json Ret_gid_msg
var ret_ver_json Ret_version_msg
var ret_status_json Ret_status_msg

// aria2 server cli command: 'aria2c --enable-rpc --rpc-listen-port=6800'

//Add uris list to downloads queue, for downloading, using aria2 api.

func Download(uris_list []string) (gid string) { // , dir string
	uris_list_arr := []interface{}{uris_list}
	json_adduris_Msg := Params_dl_msg{
		Jsonrpc: "2.0",
		Id:      "qwer",
		Method:  "aria2.addUri",
		Params:  uris_list_arr,
	}
	//call aria2 server
	body := call_aria_server(&json_adduris_Msg, &ret_gid_json)

	fmt.Printf("body= %s\n", body)                  //for debug
	fmt.Printf("result= %s\n", ret_gid_json.Result) //for debug
	gid = ret_gid_json.Result
	return gid
}

// remove from queue the download defined by its gid, using aria2 api
func Remove_dl(gid string) (rm_gid string) {

	gid_arr := []interface{}{gid}

	json_remove_msg := Params_dl_msg{
		Jsonrpc: "2.0",
		Id:      "qwer",
		Method:  "aria2.remove",
		Params:  gid_arr,
	}
	//call aria2 server
	body := call_aria_server(&json_remove_msg, &ret_gid_json)

	fmt.Printf("body= %s\n", body)                  //for debug
	fmt.Printf("result= %s\n", ret_gid_json.Result) //for debug
	rm_gid = ret_gid_json.Result
	return rm_gid

}
func Pause_dl(gid string) (paused_gid string) {

	gid_arr := []interface{}{gid}

	json_pause_msg := Params_dl_msg{
		Jsonrpc: "2.0",
		Id:      "qwer",
		Method:  "aria2.pause",
		Params:  gid_arr,
	}
	//call aria2 server
	body := call_aria_server(&json_pause_msg, &ret_gid_json)

	fmt.Printf("body= %s\n", body)                  //for debug
	fmt.Printf("result= %s\n", ret_gid_json.Result) //for debug
	paused_gid = ret_gid_json.Result
	return paused_gid

}

func Unpause_dl(gid string) (unpaused_gid string) { //gid []string

	gidArr := []interface{}{gid}
	//gidArr:= []string{gid,}

	jsonUnpauseMsg := Params_dl_msg{ //ByGidMsg{
		Jsonrpc: "2.0",
		Id:      "qwer",
		Method:  "aria2.unpause",
		Params:  gidArr,
	}
	//call aria2 server
	body := call_aria_server(&jsonUnpauseMsg, &ret_gid_json)

	fmt.Printf("body= %s\n", body)                  //for debug
	fmt.Printf("result= %s\n", ret_gid_json.Result) //for debug
	unpaused_gid = ret_gid_json.Result
	return unpaused_gid
}

//get downloads selective status, by using aria2 api. Return status of many parameters - details in aria2 api
// to do - update prams status list
func GetStatus_dl(gid string) (dl_status Ret_status_msg) {
	fmt.Println("gid= ", gid)             //for debug
	json_tellStatus_msg := Params_dl_msg{ //Params1Msg
		Jsonrpc: "2.0",
		Id:      "qwer",
		Method:  "aria2.tellStatus",
		Params:  []interface{}{gid, []string{"gid", "completedLength", "totalLength", "status", "connections"}},
	}
	fmt.Printf("params= %s\n ", json_tellStatus_msg) //for debug
	//call aria2 server
	body := call_aria_server(&json_tellStatus_msg, &ret_status_json)

	fmt.Printf("body= %s\n", body)                      //for debug
	fmt.Printf("result= %+v\n", ret_status_json.Result) //for debug
	//status:=retStatusJson
	return ret_status_json
}

func GetUris_list(gid string) {
	//get downloads status, by using aria2 api. Return status of many parameters - details in aria2 api
	//data=json.dumps({"jsonrpc":"2.0", "id":"qwer", "method":"aria2.getUris", "params": [gid] })
}

//
func Authenticate(gid, user, passwd string) (res string) {
	//authenticate user before download, by using aria2 api. Return OK on success
	//data=json.dumps({..., "params": [gid, {'http-user':user,'http-passwd':passwd}] })  #'http-auth-challenge':'True',
	fmt.Println("gid= ", gid) //for debug

	sec_arr := []interface{}{gid, map[string]string{"http-user": user, " http-passwd": passwd}}

	json_changeOption_msg := Params_dl_msg{
		Jsonrpc: "2.0",
		Id:      "qwer",
		Method:  "aria2.changeOption",
		Params:  sec_arr,
	}
	fmt.Printf("params= %s\n ", json_changeOption_msg) //for debug
	//call aria2 server
	body := call_aria_server(&json_changeOption_msg, &ret_gid_json)

	//ret_gid_json Ret_gid_msg
	fmt.Printf("body= %s\n", body)                   //for debug
	fmt.Printf("result= %+v\n", ret_gid_json.Result) //for debug
	res = ret_gid_json.Result
	return res

}

//	get aria2 version, by using aria2 api

func get_aria2_ver() (version string) {

	json_getVersion_msg := Params_dl_msg{ //NoParamsMsg{
		Jsonrpc: "2.0",
		Id:      "qwer",
		Method:  "aria2.getVersion",
		Params:  []interface{}{""},
	}
	//call aria2 server
	body := call_aria_server(&json_getVersion_msg, &ret_ver_json)

	fmt.Printf("body= %s\n", string(body))        //for debug
	fmt.Printf("retVerJson= %+v\n", ret_ver_json) //for debug
	version = ret_ver_json.Result.Version

	return version
}

//help function: call aria2 server with json rpc message

func call_aria_server(smsg interface{}, retmsg interface{}) (body []byte) {
	//prepare smsg to json format
	b, err := json.Marshal(smsg)
	if err != nil {
		log.Fatal(err)
	}
	reader := strings.NewReader(string(b))
	//call aria2 server
	res, err := http.Post("http://localhost:6800/jsonrpc", "json", reader)
	if err != nil {
		log.Fatal(err)
	}
	body, err = ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	//unmarshal response
	err = json.Unmarshal(body, retmsg)
	if err != nil {
		log.Fatal(err)
	}

	return body
}
