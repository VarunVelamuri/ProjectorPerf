package profile

import (
	"bytes"
	"encoding/json"
	"github.com/couchbase/indexing/secondary/logging"
	"net/http"
)

func StartProfile(name, username, password, projAddr string) {
	if true {
		client := &http.Client{}
		address := "http://" + projAddr + "/settings"
		jbody := make(map[string]interface{})
		jbody["projector.cpuProfFname"] = name + "_cpu.pprof"
		pbody, err := json.Marshal(jbody)
		if err != nil {
			panic(err)
		}
		req, _ := http.NewRequest("POST", address, bytes.NewBuffer(pbody))
		req.SetBasicAuth(username, password)
		req.Header.Add("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
			logging.Fatalf(address)
			logging.Fatalf("%v", req)
			logging.Fatalf("%v", resp)
			logging.Fatalf("cpuProfName failed for test: %v\n", name)
		}
	}
	if true {
		client := &http.Client{}
		address := "http://" + projAddr + "/settings"
		jbody := make(map[string]interface{})
		jbody["projector.memProfFname"] = name + "_mem_start.pprof"
		pbody, err := json.Marshal(jbody)
		if err != nil {
			panic(err)
		}
		req, _ := http.NewRequest("POST", address, bytes.NewBuffer(pbody))
		req.SetBasicAuth(username, password)
		req.Header.Add("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
			logging.Fatalf(address)
			logging.Fatalf("%v", req)
			logging.Fatalf("%v", resp)
			logging.Fatalf("memProfName failed for test: %v\n", name)
		}
	}
	// Start the CPU profile
	if true {
		client := &http.Client{}
		address := "http://" + projAddr + "/settings"
		jbody := make(map[string]interface{})
		jbody["projector.cpuProfile"] = true
		pbody, err := json.Marshal(jbody)
		if err != nil {
			panic(err)
		}
		req, _ := http.NewRequest("POST", address, bytes.NewBuffer(pbody))
		req.SetBasicAuth(username, password)
		req.Header.Add("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
			logging.Fatalf(address)
			logging.Fatalf("%v", req)
			logging.Fatalf("%v", resp)
			logging.Fatalf("cpuProfile failed for test: %v\n", name)
		}
	}
	if true {
		client := &http.Client{}
		address := "http://" + projAddr + "/settings"
		jbody := make(map[string]interface{})
		jbody["projector.memProfile"] = true
		pbody, err := json.Marshal(jbody)
		if err != nil {
			panic(err)
		}
		req, _ := http.NewRequest("POST", address, bytes.NewBuffer(pbody))
		req.SetBasicAuth(username, password)
		req.Header.Add("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
			logging.Fatalf(address)
			logging.Fatalf("%v", req)
			logging.Fatalf("%v", resp)
			logging.Fatalf("memProfile failed for test: %v\n", name)
		}
	}
	if true {
		client := &http.Client{}
		address := "http://" + projAddr + "/settings"
		jbody := make(map[string]interface{})
		jbody["projector.memProfile"] = false
		pbody, err := json.Marshal(jbody)
		if err != nil {
			panic(err)
		}
		req, _ := http.NewRequest("POST", address, bytes.NewBuffer(pbody))
		req.SetBasicAuth(username, password)
		req.Header.Add("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
			logging.Fatalf(address)
			logging.Fatalf("%v", req)
			logging.Fatalf("%v", resp)
			logging.Fatalf("memProfile failed for test: %v\n", name)
		}
	}
}
func StopProfile(name, username, password, projAddr string) {
	if true {
		client := &http.Client{}
		address := "http://" + projAddr + "/settings"
		jbody := make(map[string]interface{})
		jbody["projector.cpuProfile"] = false
		pbody, err := json.Marshal(jbody)
		if err != nil {
			panic(err)
		}
		req, _ := http.NewRequest("POST", address, bytes.NewBuffer(pbody))
		req.SetBasicAuth(username, password)
		req.Header.Add("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
			logging.Fatalf(address)
			logging.Fatalf("%v", req)
			logging.Fatalf("%v", resp)
			logging.Fatalf("cpuProfile failed for test: %v\n", name)
		}
	}
	// Collect memory profile at the end-of this run
	if true {
		client := &http.Client{}
		address := "http://" + projAddr + "/settings"
		jbody := make(map[string]interface{})
		jbody["projector.memProfFname"] = name + "_mem_end.pprof"
		pbody, err := json.Marshal(jbody)
		if err != nil {
			panic(err)
		}
		req, _ := http.NewRequest("POST", address, bytes.NewBuffer(pbody))
		req.SetBasicAuth(username, password)
		req.Header.Add("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
			logging.Fatalf(address)
			logging.Fatalf("%v", req)
			logging.Fatalf("%v", resp)
			logging.Fatalf("memProfName failed for test: %v\n", name)
		}
	}
	if true {
		client := &http.Client{}
		address := "http://" + projAddr + "/settings"
		jbody := make(map[string]interface{})
		jbody["projector.memProfile"] = false
		pbody, err := json.Marshal(jbody)
		if err != nil {
			panic(err)
		}
		req, _ := http.NewRequest("POST", address, bytes.NewBuffer(pbody))
		req.SetBasicAuth(username, password)
		req.Header.Add("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
			logging.Fatalf(address)
			logging.Fatalf("%v", req)
			logging.Fatalf("%v", resp)
			logging.Fatalf("memProfile failed for test: %v\n", name)
		}
	}
}
