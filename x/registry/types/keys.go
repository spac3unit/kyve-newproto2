package types

import (
	"encoding/binary"
)

const (
	// ModuleName defines the module name
	ModuleName = "registry"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_registry"
)

// registry constants
const (
	MaxFunders        = 50 // maximum amount of funders which are allowed
	MaxStakers        = 50 // maximum amount of stakers which are allowed
	DefaultCommission = "0.9"
	UnbondingTime     = 60 * 60 * 24 * 1
	EmptyBundle       = "KYVE_EMPTY_BUNDLE"
)

// ========== EVENTS ===================
// general event props
const (
	EventName    = "EventName"
	EventPoolId  = "PoolId"
	EventCreator = "Creator"
	EventAmount  = "Amount"
)

// voting
const (
	VoteEventKey      = "Voted"
	VoteEventBundleId = "BundleId"
	VoteEventSupport  = "Support"
)

// slashing
const (
	SlashEventKey = "ReceivedSlash"
	SlashAccount  = "Account"
)

// Activity
const (
	ProposalEventKey          = "ProposalEnded"
	ProposalEventBundleId     = "BundleId"
	ProposalEventByteSize     = "ByteSize"
	ProposalEventUploader     = "Uploader"
	ProposalEventNextUploader = "NextUploader"
	ProposalEventReward       = "BundleReward"
	ProposalEventValid        = "Valid"
	ProposalEventInvalid      = "Invalid"
	ProposalEventFromHeight   = "FromHeight"
	ProposalEventToHeight     = "ToHeight"
	ProposalEventStatus       = "Status"
)

// ============ KV-STORE ===============

func KeyPrefix(p string) []byte {
	return []byte(p)
}

const (
	PoolKey           = "Pool-value-"
	PoolCountKey      = "Pool-count-"
	UnbondingStateKey = "UnbondingState-value-"
)

const (
	// DelegationEntriesKeyPrefix is the prefix to retrieve all DelegationEntries
	DelegationEntriesKeyPrefix = "DelegationEntries/value/"
	// DelegationPoolDataKeyPrefix is the prefix to retrieve all DelegationPoolData
	DelegationPoolDataKeyPrefix = "DelegationPoolData/value/"
	// DelegatorKeyPrefix is the prefix to retrieve all Delegator
	DelegatorKeyPrefix = "Delegator/value/"
	// FunderKeyPrefix is the prefix to retrieve all Funder
	FunderKeyPrefix = "Funder/value/"
	// ProposalKeyPrefix is the prefix to retrieve all Proposal
	ProposalKeyPrefix = "Proposal/value/"
	// StakerKeyPrefix is the prefix to retrieve all Staker
	StakerKeyPrefix = "Staker/value/"
	// UnbondingEntriesKeyPrefix is the prefix to retrieve all UnbondingEntries
	UnbondingEntriesKeyPrefix = "UnbondingEntries/value/"
	// UnbondingEntriesKeyPrefix is the prefix to retrieve all UnbondingEntries
	UnbondingEntriesKeyPrefixByDelegator = "UnbondingEntriesByDelegator/value/"
)

// DelegationEntriesKey returns the store key to retrieve a DelegationEntries from the index fields
func DelegationEntriesKey(poolId uint64, stakerAddress string, kIndex uint64) []byte {
	return keyPrefix{}.aInt(poolId).aString(stakerAddress).aInt(kIndex).key
}

// DelegationPoolDataKey returns the store key to retrieve a DelegationPoolData from the index fields
func DelegationPoolDataKey(poolId uint64, stakerAddress string) []byte {
	return keyPrefix{}.aInt(poolId).aString(stakerAddress).key
}

// DelegatorKey returns the store key to retrieve a Delegator from the index fields
func DelegatorKey(poolId uint64, stakerAddress string, delegatorAddress string) []byte {
	return keyPrefix{}.aInt(poolId).aString(stakerAddress).aString(delegatorAddress).key
}

// FunderKey returns the store key to retrieve a Funder from the index fields
func FunderKey(funder string, poolId uint64) []byte {
	return keyPrefix{}.aString(funder).aInt(poolId).key
}

// ProposalKey returns the store key to retrieve a Proposal from the index fields
func ProposalKey(bundleId string) []byte {
	return keyPrefix{}.aString(bundleId).key
}

// StakerKey returns the store key to retrieve a Staker from the index fields
func StakerKey(staker string, poolId uint64) []byte {
	return keyPrefix{}.aString(staker).aInt(poolId).key
}

// UnbondingEntriesKey returns the store key to retrieve a UnbondingEntries from the index fields
func UnbondingEntriesKey(index uint64) []byte {
	return keyPrefix{}.aInt(index).key
}

// UnbondingEntriesByDelegatorKey returns the store key to retrieve a UnbondingEntries from the index fields
// Index is still needed to make key unique
func UnbondingEntriesByDelegatorKey(delegator string, index uint64) []byte {
	return keyPrefix{}.aString(delegator).aInt(index).key
}

type keyPrefix struct {
	key []byte
}

func (k keyPrefix) aInt(n uint64) keyPrefix {
	indexBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(indexBytes, n)
	k.key = append(k.key, indexBytes...)
	k.key = append(k.key, []byte("/")...)
	return k
}

func (k keyPrefix) aString(s string) keyPrefix {
	k.key = append(k.key, []byte(s)...)
	k.key = append(k.key, []byte("/")...)
	return k
}
