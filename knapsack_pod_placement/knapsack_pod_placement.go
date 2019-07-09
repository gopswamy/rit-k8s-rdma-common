package knapsack_pod_placement

import (
	"log"

	"github.com/cal8384/k8s-rdma-common/rdma_hardware_info"
)

type RdmaInterfaceRequest struct {
        MinTxRate uint `json:"min_tx_rate"`
        MaxTxRate uint `json:"max_tx_rate"`
}

func PlacePod(requested_interfaces []RdmaInterfaceRequest, pfs_available []rdma_hardware_info.PF, debug_logging bool) ([]int, bool) {
	//if no interfaces are required
	if(len(requested_interfaces) <= 0) {
		//request is trivially satisfiable
		return []int{}, true
	}

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
		if(debug_logging) {
			log.Println("Outer loop iteration:")
			log.Println("\tcurrent_requested=", current_requested, " (value=", requested_interfaces[current_requested].MinTxRate, ")")
		}
		//move to next placement for current item
		for placements[current_requested]++; placements[current_requested] < len(pfs_available); placements[current_requested]++ {
			var cur_pf *rdma_hardware_info.PF = &(pfs_available[placements[current_requested]])
			//if the current pf can fit the current requested interface
			if(((*cur_pf).CapacityVFs - (*cur_pf).UsedVFs) > 0) {
				if((int((*cur_pf).CapacityTxRate) - int((*cur_pf).UsedTxRate)) >= int(requested_interfaces[current_requested].MinTxRate)) {
					//add the current interface's bandwidth to the pf's used bw
					(*cur_pf).UsedTxRate += requested_interfaces[current_requested].MinTxRate
					(*cur_pf).UsedVFs += 1
					//break
					break
				}
			}
		}
		if(debug_logging) {
			log.Println("\tplacement=", placements[current_requested])
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
				//reset placement of current item
				placements[current_requested] = -1
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
