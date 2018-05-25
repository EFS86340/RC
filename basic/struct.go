package basic

import (
	"crypto/ecdsa"
	"math/big"
)

//RCSign is the signature type in our design
type RCSign struct {
	R *big.Int
	S *big.Int
}

//Miner is the miner
type Miner struct {
	ID        string
	Rep       int
	Prk       ecdsa.PublicKey
	LastGroup int
}

//OutType is the format of the output address data in the transaction
type OutType struct {
	Value   uint32
	Address [32]byte
}

//InType is the format of the input address data in the transaction
type InType struct {
	PrevTx [32]byte
	Index  uint32
	Sig    RCSign
	PukX   *big.Int
	PukY   *big.Int
}

//Transaction is the transaction data which sent by the sender
type Transaction struct {
	Timestamp uint64
	TxinCnt   uint32
	In        []InType
	TxoutCnt  uint32
	Out       []OutType
	Kind      uint32
	Locktime  uint32
	Hash      [32]byte
}

//TxList is the list of tx sent by Leader to miner for their verification
type TxList struct {
	ID       [32]byte
	HashID   [32]byte
	PrevHash [32]byte
	TxCnt    uint32
	TxArray  []Transaction
	Sig      RCSign
}

//TxDecision is the decisions based on given TxList
type TxDecision struct {
	ID       [32]byte
	HashID   [32]byte
	TxCnt    uint32
	Decision []byte
	Sig      RCSign
}

//TxDecSet is the set of all decisions from one shard, signed by leader
type TxDecSet struct {
	ID       [32]byte
	HashID   [32]byte
	PrevHash [32]byte
	MemCnt   uint32
	MemD     []TxDecision
	TxCnt    uint32
	TxArray  [][32]byte
	Sig      RCSign
}

//TxDPure is the pure struct of the TxDecision
type TxDPure struct {
	ID       [32]byte
	Decision []byte
	Sig      RCSign
}

//TDSHeader is The Header part of TxDecSS
type TDSHeader struct {
	ID       [32]byte
	HashID   [32]byte
	PrevHash [32]byte
	TxCnt    uint32
	TxIndex  []uint32
	MemCnt   uint32
	MemD     []TxDPure
	Sig      RCSign
}

//TxDecSS is the set of TxDecSet
type TxDecSS struct {
	ShardNum uint32
	Header   []TDSHeader
	TxCnt    uint32
	Tx       [][32]byte
}

//TxBlock introduce the struct of the transaction block
type TxBlock struct {
	ID         [32]byte
	PrevHash   [32]byte
	HashID     [32]byte
	MerkleRoot [32]byte
	Timestamp  int64
	Height     uint32
	TxCnt      uint32
	TxArray    []Transaction
	Sig        RCSign
}

//TxDB is the database of cache
type TxDB struct {
	ID   [32]byte
	Used []uint32
	/*0 not checked(the first time),
	1: Correct part in the shard,
	-1: fail due to not correct format
	-2: fail due to double spend
	*/
	InCheck    []bool
	Res        int8
	InCheckSum int
}

//UserClient is the struct for miner and client
type UserClient struct {
	IPaddress string
	Prk       ecdsa.PublicKey
	kind      int
}

//AccCache is the cache of account
type AccCache struct {
	ID    [32]byte
	Value uint32
}
