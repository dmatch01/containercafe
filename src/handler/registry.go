package handler

import (
	"net/http"
	"io"
	"io/ioutil"

	"conf"
	"httphelper"
)

func InjectRegAuthHeader(r *http.Request) {
	tok := conf.GetRegAuthToken()
	r.Header.Set("X-Registry-Auth", tok)
}

func GetRegistryApiHost() (host string){
	// get service host from Consul
	service := "registry-api"
	host = conf.GetServiceHost(service)
	if (host == ""){
		Log.Printf("Failed to get Registry API host  service=%s", service)
	}
	return
}

// call internal registry api server to get image metadata
func invoke_reg_inspect(w http.ResponseWriter, r *http.Request, img string, namespace string, req_id string){
	host := GetRegistryApiHost()
	if (host == ""){
		ErrorHandlerWithMsg(w, r, 500, "Failed to get Registry API host")
		return
	}

	//Call service endpoint
	url := "http://"+host+"/v1/imageJson?imageName="+img
	Log.Printf("Will call Registry... url=%s req_id=%s", url, req_id)
	client := &http.Client{}
	req, _ := http.NewRequest("GET",url, nil)
	req.Header.Add("Accept", "application/json")
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	body,_:=ioutil.ReadAll(resp.Body)
	Log.Printf("Registry status code: %d", resp.StatusCode)
	Log.Printf("Registry response: %s", httphelper.PrettyJson(body) )

	//send response back to client
	w.WriteHeader(resp.StatusCode)
	io.WriteString(w, string(body))

	return
}

//implement image list by invoking search api of Containers registry
//return json to docker cli
//DockerHandler will print req processing completion message and exit right after this method
func invoke_reg_list(w http.ResponseWriter, r *http.Request, namespace string, req_id string){
	//Not recommended approach. These urls go to registry.ng.bluemix.net
	//ns_url := "http://" + conf.GetRegLocation() + "/v1/namespaces/" + namespace
	//lib_url := "http://" + conf.GetRegLocation() + "/v1/namespaces/library"

	// Recommended approach using internal api exposed by registry microservice
	// get service host from Consul
	host := GetRegistryApiHost()
	if (host == ""){
		ErrorHandlerWithMsg(w, r, 500, "Failed to get Registry API host")
		return
	}

	//Call service endpoint
	url := "http://" + host + "/v1/imageList/" + namespace
	Log.Printf("Will call Registry... url=%s req_id=%s", url, req_id)
	client := &http.Client{}
	req, _ := http.NewRequest("GET",url, nil)
	req.Header.Add("Accept", "application/json")
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	body,_:=ioutil.ReadAll(resp.Body)
	Log.Printf("Registry status code: %d", resp.StatusCode)
	Log.Printf("Registry response: %s", httphelper.PrettyJson(body))

	//send response back to client
	w.WriteHeader(resp.StatusCode)
	io.WriteString(w, string(body))
	return
}

func invoke_reg_rmi(w http.ResponseWriter, r *http.Request, namespace string, req_id string){
	//TODO
}
