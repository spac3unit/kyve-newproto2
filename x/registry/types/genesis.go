package types

import (
	"fmt"
)

// DefaultIndex is the default capability global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		PoolList:               []Pool{},
		FunderList:             []Funder{},
		StakerList:             []Staker{},
		DelegatorList:          []Delegator{},
		DelegationPoolDataList: []DelegationPoolData{},
		DelegationEntriesList:  []DelegationEntries{},
		// this line is used by starport scaffolding # genesis/types/default
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// Check for duplicated ID in pool
	poolIdMap := make(map[uint64]bool)
	poolCount := gs.GetPoolCount()
	for _, elem := range gs.PoolList {
		if _, ok := poolIdMap[elem.Id]; ok {
			return fmt.Errorf("duplicated id for pool")
		}
		if elem.Id >= poolCount {
			return fmt.Errorf("pool id should be lower or equal than the last id")
		}
		poolIdMap[elem.Id] = true
	}
	// Check for duplicated index in funder
	funderIndexMap := make(map[string]struct{})

	for _, elem := range gs.FunderList {
		index := string(FunderKey(elem.Account, elem.PoolId))
		if _, ok := funderIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for funder")
		}
		funderIndexMap[index] = struct{}{}
	}
	// Check for duplicated index in staker
	stakerIndexMap := make(map[string]struct{})

	for _, elem := range gs.StakerList {
		index := string(StakerKey(elem.Account, elem.PoolId))
		if _, ok := stakerIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for staker")
		}
		stakerIndexMap[index] = struct{}{}
	}
	// Check for duplicated index in delegator
	delegatorIndexMap := make(map[string]struct{})

	for _, elem := range gs.DelegatorList {
		index := string(DelegatorKey(elem.Id, elem.Staker, elem.Delegator))
		if _, ok := delegatorIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for delegator")
		}
		delegatorIndexMap[index] = struct{}{}
	}
	// Check for duplicated index in delegationPoolData
	delegationPoolDataIndexMap := make(map[string]struct{})

	for _, elem := range gs.DelegationPoolDataList {
		index := string(DelegationPoolDataKey(elem.Id, elem.Staker))
		if _, ok := delegationPoolDataIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for delegationPoolData")
		}
		delegationPoolDataIndexMap[index] = struct{}{}
	}
	// Check for duplicated index in delegationEntries
	delegationEntriesIndexMap := make(map[string]struct{})

	for _, elem := range gs.DelegationEntriesList {
		index := string(DelegationEntriesKey(elem.Id, elem.Staker, elem.KIndex))
		if _, ok := delegationEntriesIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for delegationEntries")
		}
		delegationEntriesIndexMap[index] = struct{}{}
	}
	// Check for duplicated index in proposal
	proposalIndexMap := make(map[string]struct{})

	for _, elem := range gs.ProposalList {
		index := string(ProposalKey(elem.StorageId))
		if _, ok := proposalIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for proposal")
		}
		proposalIndexMap[index] = struct{}{}
	}
	// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}
