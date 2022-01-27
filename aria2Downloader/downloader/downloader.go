package downloader

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type MsgToAria2 struct {
	Jsonrpc string        `json:"jsonrpc"` // added for case insensitive to match aria2 jsonrpc format (lower case)
	Id      string        `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

type RetGidMsg struct {
	Jsonrpc string
	Id      string
	Result  string
}

type RetVersionMsg struct {
	Jsonrpc string
	Id      string
	Result  VersionResult
}
type VersionResult struct {
	EnabledFeatures []string
	Version         string
}

// for aria selective tellstatus
type RetStatusMsg struct {
	Jsonrpc string       `json:"jsonrpc"`
	Id      string       `json:"id"`
	Result  ResultParams `json:"resultParams"`
}
type ResultParams struct {
	Gid             string `json:"gid"`
	CompletedLength string `json:"completedLength"`
	TotalLength     string `json:"totalLength"`
	Status          string `json:"status"`
	Connections     string `json:"connections"`
}

var retGidJson RetGidMsg
var retVersionJson RetVersionMsg
var retStatusJson RetStatusMsg

//Add uris list to downloads queue and start download, using aria2 api.

func Start(urisList []string) (gid string) {

	urisListArr := []interface{}{urisList}

	jsonAddurisMsg := MsgToAria2{
		Jsonrpc: "2.0",
		Id:      "qwer",
		Method:  "aria2.addUri",
		Params:  urisListArr,
	}

	//call aria2 server
	callAria2Server(&jsonAddurisMsg, &retGidJson)

	gid = retGidJson.Result

	return gid
}

// remove from queue the download defined by its gid, using aria2 api

func Remove(gid string) (removedGid string) {

	gidArr := []interface{}{gid}

	jsonRemoveMsg := MsgToAria2{
		Jsonrpc: "2.0",
		Id:      "qwer",
		Method:  "aria2.remove",
		Params:  gidArr,
	}

	//call aria2 server
	callAria2Server(&jsonRemoveMsg, &retGidJson)

	removedGid = retGidJson.Result

	return removedGid

}
func Pause(gid string) (pausedGid string) {

	gidArr := []interface{}{gid}

	jsonPauseMsg := MsgToAria2{
		Jsonrpc: "2.0",
		Id:      "qwer",
		Method:  "aria2.pause",
		Params:  gidArr,
	}

	//call aria2 server
	callAria2Server(&jsonPauseMsg, &retGidJson)

	pausedGid = retGidJson.Result

	return pausedGid

}

func Resume(gid string) (resumedGid string) {

	gidArr := []interface{}{gid}
	//gidArr:= []string{gid,}

	jsonResumeMsg := MsgToAria2{
		Jsonrpc: "2.0",
		Id:      "qwer",
		Method:  "aria2.unpause",
		Params:  gidArr,
	}

	//call aria2 server
	callAria2Server(&jsonResumeMsg, &retGidJson)

	resumedGid = retGidJson.Result

	return resumedGid
}

//Get downloads selective status, by using aria2 api.
//Return status of many parameters - details in aria2 api
//To do - update prams status list
func DownloadStatus(gid string) (downloadStatus RetStatusMsg) {

	json_tellStatus_msg := MsgToAria2{ //Params1Msg
		Jsonrpc: "2.0",
		Id:      "qwer",
		Method:  "aria2.tellStatus",
		Params:  []interface{}{gid, []string{"gid", "completedLength", "totalLength", "status", "connections"}},
	}

	//call aria2 server
	callAria2Server(&json_tellStatus_msg, &retStatusJson)

	return retStatusJson
}

//	get aria2 version, by using aria2 api

func getAria2Version() (version string) {

	jsonGetVersionMsg := MsgToAria2{
		Jsonrpc: "2.0",
		Id:      "qwer",
		Method:  "aria2.getVersion",
		Params:  []interface{}{""},
	}
	//call aria2 server
	callAria2Server(&jsonGetVersionMsg, &retVersionJson)

	version = retVersionJson.Result.Version

	return version
}

//call aria2 server with json rpc message

func callAria2Server(smsg interface{}, retmsg interface{}) (body []byte) {

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

// For future use
func Authenticate(gid, user, passwd string) (res string) {

	//authenticate user before download, by using aria2 api. Return OK on success

	securityArr := []interface{}{gid, map[string]string{"http-user": user, " http-passwd": passwd}}

	json_changeOption_msg := MsgToAria2{
		Jsonrpc: "2.0",
		Id:      "qwer",
		Method:  "aria2.changeOption",
		Params:  securityArr,
	}

	//call aria2 server
	callAria2Server(&json_changeOption_msg, &retGidJson)

	//ret_gid_json Ret_gid_msg
	res = retGidJson.Result

	return res
}
