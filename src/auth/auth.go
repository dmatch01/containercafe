package auth

import (
	"net/http"
	"log"
	"strings"
	"io/ioutil"
	"httphelper"  	// my httphelper
	"conf"  		// my conf package
)

//returns auth=true/false, compute node name, container/exec id, container id
func Auth(r *http.Request) (ok bool, node string, docker_id string, container string) {
	//parse r.RequestURI for container id or exec id
	uri := r.RequestURI
	id, id_type := get_id_and_type(uri)
	/*
	//1st: look for /containers/<id>/
	id := get_id_from_uri(uri, "/containers/")
	id_type:="Container"
	if id == "" {
		//2nd: look for /exec/<id>/
		id = get_id_from_uri(uri, "/exec/")
		if id == "" {
			id_type = "None"
			log.Printf("Auth: id not found in uri\n")
			log.Printf("Auth: id=%s, id_type=%s\n", id, id_type)
			//fail here
			//log.Printf("Auth result: ok=%t, node='%s'\n", ok, node)
			//return ok, node, docker_id, container
		}else {
			//id found in uri
			id_type = "Exec"
			log.Printf("Auth: id=%s, id_type=%s\n", id, id_type)
		}
	}
	*/
	//if type is Exec then get container id from redis cache to forward to getHost api
	var container_id string
	if id_type == "Exec" {
		container_id = conf.RedisGet(id)
	}else{
		container_id = id
	}

	//forward r header only without body to ccsapi auth endpoint, add X-Container-Id header
	var new_uri string
	if id_type == "None" {
		new_uri = "http://"+conf.GetCcsapiHost()+conf.GetCcsapiUri()+"getHost/NoneContainer"
	}else{
		new_uri = "http://"+conf.GetCcsapiHost()+conf.GetCcsapiUri()+"getHost/"+container_id
	}
	req, _ := http.NewRequest("GET", new_uri, nil)
	httphelper.CopyHeader(req.Header, r.Header)  //req.Header = r.Header
	req.URL.Host = conf.GetCcsapiHost()
	req.Header.Add(conf.GetCcsapiIdHeader(), container_id)
	req.Header.Add(conf.GetCcsapiIdTypeHeader(), "Container" /*id_type*/)
	client := &http.Client{
		CheckRedirect: nil,
	}
	log.Println("Auth: will invoke CCSAPI to authenticate...")
	resp, err := client.Do(req)
	if (err != nil) {
		log.Printf("Auth: Error in auth request... %v\n", err)

		log.Printf("Auth result: ok=%t, node='%s'\n", ok, node)
		return ok, node, docker_id, container
	}

	//get auth response status, and X-Compute-Node header
	if resp.StatusCode == 200 {
		ok = true

		//first check in header
		//node = httphelper.GetHeader(resp.Header, conf.GetCcsapiComputeNodeHeader())
		//second check for json response in body
		defer resp.Body.Close()
		body, e := ioutil.ReadAll(resp.Body)            //Default_redirect_host   //testing default
		if e != nil {
			log.Printf("error reading ccsapi response\n")
			log.Printf("Auth result: ok=%t, node='%s'\n", ok, node)
			return ok, node, docker_id, container
		}
		log.Printf("Auth: ccsapi raw response=%s\n", body)
		var resp GetHostResp
		err := parse_getHost_Response(body, &resp)
		if err != nil {
			log.Printf("error parsing ccsapi response\n")
			log.Printf("Auth result: ok=%t, node='%s'\n", ok, node)
			return ok, node, docker_id, container
		}
		node = resp.Host
		node = node + ":" + conf.GetDockerPort()
		container = resp.Container_id
		//container id needs nova- prefix
		//exec id does not need a prefix
		if id_type == "Container" {
			docker_id = "nova-"+container   //append nova to the id returned from getHost
		}
		if id_type == "Exec" {
			docker_id = id  //use the exec id that came in the original req
		}
		if id_type == "None" {
			docker_id = ""
		}

		if resp.Swarm {
			node = resp.Mgr_host    //Mgr_host = host:port
			//insert space_id in the header to be forwarded
			r.Header.Set("X-Auth-Token", resp.Space_id)
			if id_type == "Container" {
				docker_id = container
			}
			if id_type == "Exec" {
				docker_id = id  //use the exec id that came in the original req
			}
			if id_type == "None" {
				docker_id = ""
			}
			if !resp.Swarm_tls{
				conf.SetTlsOutboundOverride(true)
			}
		}
	}else {
		//status!=200
	}

	log.Printf("Auth result: ok=%t, node=%s, docker_id=%s, container=%s\n", ok, node, docker_id, container)
	return ok, node, docker_id, container
}

//Convert /v*/containers/id/*  to  /<docker_api_ver>/containers/<redirect_resource_id>/*
//Convert /v*/exec/id/*  to  /<docker_api_ver>/exec/<redirect_resource_id>/*
func RewriteURI(reqURI string, redirect_resource_id string) string{
	var redirectURI string
	sl := strings.Split(reqURI, "/")
	if redirect_resource_id == "" {
		//supports /v../containers/json  /v../build  /v../build?foo=bar
		//redirectURI = conf.GetDockerApiVer()+"/"+sl[2]+"/"+sl[3]
		redirectURI = conf.GetDockerApiVer()
		for i:=2; i < len(sl); i++ {
			redirectURI += "/" + sl[i]
		}
	}else {
		//redirectURI = conf.GetDockerApiVer()+"/"+sl[2]+"/"+redirect_resource_id+"/"+sl[4]
		redirectURI = conf.GetDockerApiVer()+"/"+sl[2]+"/"+redirect_resource_id
		for i:=4; i < len(sl); i++ {
			redirectURI += "/" + sl[i]
		}
		//what if there is ?foo=bar in last slice and last slice is resource_id e.g., DELETE /v/containers/123?foo=bar
		if len(sl) <= 4 {
			sl2 := strings.Split(sl[len(sl)-1], "?")
			if len(sl2) > 1{
				redirectURI += "?" + sl2[1]
			}
		}
	}
	log.Printf("RewriteURI: '%s' --> '%s'\n", reqURI, redirectURI)
	return redirectURI
}

func get_id_from_uri(uri string, pattern string) string{
	var id string
	slice1 := strings.Split(uri, pattern)
	//log.Printf("get_id_from_uri: pattern=%s, slice1=%v\n", pattern, slice1)
	if len(slice1) > 1 {
		slice2 := strings.Split(slice1[1], "/")
		id = slice2[0]
		slice3 := strings.Split(id, "?")
		id = slice3[0]
	}else{
		id=""
	}
	//log.Printf("get_id_from_uri: id=%s\n", id)
	return id
}

func get_id_and_type(uri string) (id string, id_type string){
	patterns := []string {
		"/containers/json",
		"/containers/create",
		"/images",
		"/build",
		"/version",
		"/_ping",
	}

	found := false
	for i:=0; i < len(patterns); i++ {
		//if uri contains patterns[i]
		if strings.Contains(uri, patterns[i]) {
			found = true
		}
	}
	if found {
		id=""
		id_type = "None"
		log.Printf("id=%s, id_type=%s\n", id, id_type)
		return id, id_type
	}

	//1st: look for /containers/<id>/
	id = get_id_from_uri(uri, "/containers/")
	id_type="Container"
	if id == "" {
		//2nd: look for /exec/<id>/
		id = get_id_from_uri(uri, "/exec/")
		if id == "" {
			id_type = "None"
		}else {
			//id found in uri
			id_type = "Exec"
		}
	}
	log.Printf("id=%s, id_type=%s\n", id, id_type)
	return id, id_type
}


