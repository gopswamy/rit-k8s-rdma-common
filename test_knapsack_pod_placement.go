package main

import (
	"log"

	"github.com/rit-k8s-rdma/rit-k8s-rdma-common/knapsack_pod_placement"
	"github.com/rit-k8s-rdma/rit-k8s-rdma-common/rdma_hardware_info"
)

/*
	Allocate the requested interface in bandwidths: 3, 3, 8
	Onto PFs with the following remaining free bandwidth: 9, 6
*/
func test_case_1() {
	var pfs []rdma_hardware_info.PF = make([]rdma_hardware_info.PF, 2)
	pfs[0].UsedTxRate = 0
	pfs[0].CapacityTxRate = 9
	pfs[0].UsedVFs = 1
	pfs[0].CapacityVFs = 120
	pfs[1].UsedTxRate = 0
	pfs[1].CapacityTxRate = 6
	pfs[1].UsedVFs = 1
	pfs[1].CapacityVFs = 120

	var req []knapsack_pod_placement.RdmaInterfaceRequest = make([]knapsack_pod_placement.RdmaInterfaceRequest, 3)
	req[0].MinTxRate = 3
	req[1].MinTxRate = 3
	req[2].MinTxRate = 8

	allocation, possible := knapsack_pod_placement.PlacePod(req, pfs, true)

	log.Println("Possible: ", possible)
	log.Println("Allocation: ", allocation)
}

/*
	Allocate the requested interface in bandwidths: 3, 3, 8
	Onto PFs with the following remaining free bandwidth: 9, 6
	Where the second PF only has one VF available.
*/
func test_case_2() {
	var pfs []rdma_hardware_info.PF = make([]rdma_hardware_info.PF, 2)
	pfs[0].UsedTxRate = 0
	pfs[0].CapacityTxRate = 9
	pfs[0].UsedVFs = 1
	pfs[0].CapacityVFs = 120
	pfs[1].UsedTxRate = 0
	pfs[1].CapacityTxRate = 6
	pfs[1].UsedVFs = 119
	pfs[1].CapacityVFs = 120

	var req []knapsack_pod_placement.RdmaInterfaceRequest = make([]knapsack_pod_placement.RdmaInterfaceRequest, 3)
	req[0].MinTxRate = 3
	req[1].MinTxRate = 3
	req[2].MinTxRate = 8

	allocation, possible := knapsack_pod_placement.PlacePod(req, pfs, true)

	log.Println("Possible: ", possible)
	log.Println("Allocation: ", allocation)
}

/*
	Allocate the requested interface in bandwidths: 8, 3, 3, 3, 3
	Onto PFs with the following remaining free bandwidth: 12, 8
*/
func test_case_3() {
	var pfs []rdma_hardware_info.PF = make([]rdma_hardware_info.PF, 2)
	pfs[0].UsedTxRate = 0
	pfs[0].CapacityTxRate = 12
	pfs[0].UsedVFs = 1
	pfs[0].CapacityVFs = 120
	pfs[1].UsedTxRate = 0
	pfs[1].CapacityTxRate = 8
	pfs[1].UsedVFs = 1
	pfs[1].CapacityVFs = 120

	var req []knapsack_pod_placement.RdmaInterfaceRequest = make([]knapsack_pod_placement.RdmaInterfaceRequest, 5)
	req[0].MinTxRate = 8
	req[1].MinTxRate = 3
	req[2].MinTxRate = 3
	req[3].MinTxRate = 3
	req[4].MinTxRate = 3

	allocation, possible := knapsack_pod_placement.PlacePod(req, pfs, true)

	log.Println("Possible: ", possible)
	log.Println("Allocation: ", allocation)
}

/*
	Allocate the requested interface in bandwidths: 8, 3, 3, 3, 3
	Onto PFs with the following remaining free bandwidth: 10, 11
*/
func test_case_4() {
	var pfs []rdma_hardware_info.PF = make([]rdma_hardware_info.PF, 2)
	pfs[0].UsedTxRate = 0
	pfs[0].CapacityTxRate = 10
	pfs[0].UsedVFs = 1
	pfs[0].CapacityVFs = 120
	pfs[1].UsedTxRate = 0
	pfs[1].CapacityTxRate = 11
	pfs[1].UsedVFs = 1
	pfs[1].CapacityVFs = 120

	var req []knapsack_pod_placement.RdmaInterfaceRequest = make([]knapsack_pod_placement.RdmaInterfaceRequest, 5)
	req[0].MinTxRate = 8
	req[1].MinTxRate = 3
	req[2].MinTxRate = 3
	req[3].MinTxRate = 3
	req[4].MinTxRate = 3

	allocation, possible := knapsack_pod_placement.PlacePod(req, pfs, true)

	log.Println("Possible: ", possible)
	log.Println("Allocation: ", allocation)
}

/*
	Test program for RDMA bacndwidth allocation algorithm.
*/
func main() {
	log.Println("TEST CASE 1:")
	test_case_1()
	log.Println("\n\n")

	log.Println("TEST CASE 2:")
	test_case_2()
	log.Println("\n\n")

	log.Println("TEST CASE 3:")
	test_case_3()
	log.Println("\n\n")

	log.Println("TEST CASE 4:")
	test_case_4()
	log.Println("\n\n")
}
