package rdma_hardware_info

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)


func QueryNode(node_address string, port string, timeout_ms int) ([]PF, error) {
	http_client := http.Client {
		Timeout: time.Duration(time.Duration(timeout_ms) * time.Millisecond),
	}

	resp, err := http_client.Get(fmt.Sprintf("http://%s:%s/%s", node_address, port, RdmaInfoUrl))
	if(err != nil) {
		return []PF{}, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if(err != nil) {
		return []PF{}, err
	}

	var pfs []PF
	err = json.Unmarshal(data, &pfs)
	if(err != nil) {
		return []PF{}, err
	}

	return pfs, nil
}
