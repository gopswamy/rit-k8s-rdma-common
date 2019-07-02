package rdma_hardware_info

import (
	"net/http"
)


func QueryNode(node_address string, port string, timeout_ms int) ([]rdma_hardware_info.PF, error) {
	http_client := http.Client {
		Timeout: time.Duration(timeout_ms * time.Millisecond),
	}

	resp, err := http_client.Get(fmt.Sprintf("http://%s:%s/%s", node_address, port, RdmaInfoUrl))
	if(err != nil) {
		return []rdma_hardware_info.PF{}, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if(err != nil) {
		return []rdma_hardware_info.PF{}, err
	}

	var pfs []rdma_hardware_info.PF
	err = json.Unmarshal(data, &pfs)
	if(err != nil) {
		return []rdma_hardware_info.PF{}, err
	}

	return pfs, nil
}
