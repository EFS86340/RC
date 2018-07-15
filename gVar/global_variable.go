package gVar

import "time"

//MagicNumber magic
const MagicNumber byte = 66

//ShardSize is the number of miners in one shard
const ShardSize uint32 = 2

//ShardCnt is the number of shards
const ShardCnt uint32 = 1

//used in rep calculation, scaling factor
const RepTP = 1
const RepTN = 1
const RepFP = 1
const RepFN = 1

//channel

const SlidingWindows = 10

//NumTxListPerEpoch is the number of txblocks in one epoch
const NumTxListPerEpoch = 4 //60

//NumTxBlockForRep is the number of blocks for one rep block
const NumTxBlockForRep = 3 //10

//const GensisAcc = []byte{0}

const GensisAccValue = 2147483647

const TxSendInterval = 10

const NumOfTxForTest = 200

const GeneralSleepTime = 50

var T1 time.Time = time.Now()
