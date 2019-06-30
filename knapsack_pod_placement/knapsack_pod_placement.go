package knapsack_pod_placement

import (
	"github.com/cal8384/k8s-rdma-common/rdma_hardware_info"
)

type RdmaInterfaceRequest struct {
        MinTxRate uint `json:"min_tx_rate"`
        MaxTxRate uint `json:"max_tx_rate"`
}

func PlacePod(requested_interfaces []RdmaInterfaceRequest, pfs_available []rdma_hardware_info.PF) ([]int, bool) {
	//current_requested = 0
	var current_requested int = 0
	//placements = []int (initialize to all -1)
	var placements = make([]int, len(requested_interfaces))
	for index, _ := range placements {
		placements[index] = -1
	}
	//all_interfaces_sucessfully_placed := false
	var all_interfaces_sucessfully_placed bool = false

	//while(not done)
	for {
		//move to next placement for current item
		for placements[current_requested]++; placements[current_requested] < len(pfs_available); placements[current_requested]++ {
			var cur_pf *rdma_hardware_info.PF = &(pfs_available[placements[current_requested]])
			//if the current pf can fit the current requested interface
			if(((*cur_pf).CapacityVFs - (*cur_pf).UsedVFs) > 0) {
				if((int((*cur_pf).CapacityTxRate) - int((*cur_pf).UsedTxRate)) > int(requested_interfaces[current_requested].MinTxRate)) {
					//add the current interface's bandwidth to the pf's used bw
					(*cur_pf).UsedTxRate += requested_interfaces[current_requested].MinTxRate
					(*cur_pf).UsedVFs += 1
					//break
					break
				}
			}
		}

		//if there was no next placement
		if(placements[current_requested] >= len(pfs_available)) {
			//if the current item was item #0
			if(current_requested == 0) {
				//FAIL
				all_interfaces_sucessfully_placed = false
				break
			//else
			} else {
				//decrement current item
				current_requested--
				//subtract the bw and vf of current item from the pf it was allocated to
				pfs_available[placements[current_requested]].UsedTxRate -= requested_interfaces[current_requested].MinTxRate
				pfs_available[placements[current_requested]].UsedVFs -= 1
				//continue
				continue
			}

		//else
		} else {
			//if current item was the last one
			if(current_requested == (len(requested_interfaces) - 1)) {
				//SUCCEED
				all_interfaces_sucessfully_placed = true
				break

			//else
			} else {
				//increment current item
				current_requested++
				//continue
				continue
			}
		}
	}

	//if the request could be satisfied
	if(all_interfaces_sucessfully_placed) {
		//return the allocation that satisfied it
		return placements, true
	}

	//request could not be satisfied, just return empty allocation
	return []int{}, false
}
