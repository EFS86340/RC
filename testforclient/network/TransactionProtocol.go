package network

import (
	"math/rand"
	"time"

	"github.com/uchihatmtkinu/RC/gVar"

	"github.com/uchihatmtkinu/RC/shard"

	"github.com/uchihatmtkinu/RC/basic"
)

// sendRepPowMessage send reputation block
func sendTxMessage(addr string, command string, message []byte) {
	request := append(commandToBytes(command), message...)
	sendData(addr, request)
}

//TxGeneralLoop is the normall loop of transaction cache
func TxGeneralLoop() {
	tmp := 0
	flag := false
	for {
		tmp++
		time.Sleep(time.Second * 10)

		CacheDbRef.Mu.Lock()
		CacheDbRef.SignTDS(0)
		CacheDbRef.Mu.Unlock()
		CacheDbRef.Mu.RLock()

		data1 := new([]byte)
		CacheDbRef.TLS[CacheDbRef.ShardNum].Encode(data1)
		go SendTxList(data1)
		if flag {
			data2 := make([][]byte, gVar.ShardCnt)
			for i := uint32(0); i < gVar.ShardCnt; i++ {
				CacheDbRef.TDSCache[0][i].Encode(&data2[i])
			}
			go SendTxDecSet(data2)
		}
		if tmp == 3 {
			flag = true
			tmp = 0
			CacheDbRef.GenerateTxBlock()
			data3 := new([]byte)
			CacheDbRef.TxB.Encode(data3, 0)
			go SendTxBlock(data3)
		}

		CacheDbRef.Mu.RUnlock()

		CacheDbRef.Mu.Lock()

		CacheDbRef.BuildTDS()
		CacheDbRef.NewTxList()
		CacheDbRef.Release()
		CacheDbRef.Mu.RUnlock()
	}
}

//SendTxList is sending txlist
func SendTxList(data *[]byte) {
	for i := uint32(0); i < gVar.ShardSize; i++ {
		xx := shard.ShardToGlobal[CacheDbRef.ShardNum][i]
		if xx != int(CacheDbRef.ID) {
			sendTxMessage(shard.GlobalGroupMems[xx].Address, "TxList", *data)
		}
	}
}

//SendTxDecSet is sending txDecSet
func SendTxDecSet(data [][]byte) {
	for i := uint32(0); i < gVar.ShardSize; i++ {
		xx := shard.ShardToGlobal[CacheDbRef.ShardNum][i]
		if xx != int(CacheDbRef.ID) {
			sendTxMessage(shard.GlobalGroupMems[xx].Address, "TxDecSetM", data[CacheDbRef.ShardNum])
		}
	}
	for i := uint32(0); i < gVar.ShardCnt; i++ {
		xx := rand.Int()%(int(gVar.ShardSize)-1) + 1
		if i != CacheDbRef.ShardNum {
			sendTxMessage(shard.GlobalGroupMems[shard.ShardToGlobal[i][xx]].Address, "TxDecSet", data[i])
		}
	}
}

//SendTxBlock is sending txBlock
func SendTxBlock(data *[]byte) {
	for i := uint32(0); i < gVar.ShardSize; i++ {
		xx := shard.ShardToGlobal[CacheDbRef.ShardNum][i]
		if xx != int(CacheDbRef.ID) {
			sendTxMessage(shard.GlobalGroupMems[xx].Address, "TxB", *data)
		}
	}
}

/************************Miner***************************/

//HandleTx when receives a tx
func HandleTx(data []byte) error {
	tmp := new(basic.Transaction)
	err := tmp.Decode(&data)
	if err != nil {
		return err
	}
	CacheDbRef.Mu.Lock()
	CacheDbRef.GetTx(tmp)
	CacheDbRef.Mu.Unlock()
	return nil
}

//HandleAndSendTx when receives a tx
func HandleAndSendTx(data []byte) error {
	tmp := new(basic.Transaction)
	err := tmp.Decode(&data)
	if err != nil {
		return err
	}
	CacheDbRef.Mu.Lock()
	CacheDbRef.GetTx(tmp)
	CacheDbRef.Mu.Unlock()
	for i := uint32(0); i < gVar.ShardSize; i++ {
		xx := shard.ShardToGlobal[CacheDbRef.ShardNum][i]
		if xx != int(CacheDbRef.ID) {
			sendTxMessage(shard.GlobalGroupMems[xx].Address, "TxM", data)
		}
	}
	return nil
}

//HandleTxList when receives a txlist
func HandleTxList(data []byte) error {
	tmp := new(basic.TxList)
	err := tmp.Decode(&data)
	if err != nil {
		return err
	}
	CacheDbRef.Mu.RLock()
	CacheDbRef.PreTxList(tmp)
	CacheDbRef.Mu.RUnlock()
	CacheDbRef.Mu.Lock()
	CacheDbRef.ProcessTL(tmp)
	var sent []byte
	CacheDbRef.TLSent.Encode(&sent)
	CacheDbRef.Mu.Unlock()
	sendTxMessage(shard.GlobalGroupMems[tmp.ID].Address, "TxDec", sent)
	return nil
}

//HandleTxDecSet when receives a txdecset
func HandleTxDecSet(data []byte) error {
	tmp := new(basic.TxDecSet)
	err := tmp.Decode(&data)
	if err != nil {
		return err
	}
	CacheDbRef.Mu.RLock()
	CacheDbRef.PreTxDecSet(tmp)
	CacheDbRef.Mu.RUnlock()
	CacheDbRef.Mu.Lock()
	CacheDbRef.GetTDS(tmp)
	CacheDbRef.Mu.Unlock()
	return nil
}

//HandleAndSentTxDecSet when receives a txdecset
func HandleAndSentTxDecSet(data []byte) error {
	tmp := new(basic.TxDecSet)
	err := tmp.Decode(&data)
	if err != nil {
		return err
	}
	CacheDbRef.Mu.RLock()
	CacheDbRef.PreTxDecSet(tmp)
	CacheDbRef.Mu.RUnlock()
	CacheDbRef.Mu.Lock()
	CacheDbRef.GetTDS(tmp)
	CacheDbRef.Mu.Unlock()
	for i := uint32(0); i < gVar.ShardSize; i++ {
		xx := shard.ShardToGlobal[CacheDbRef.ShardNum][i]
		if xx != int(CacheDbRef.ID) {
			sendTxMessage(shard.GlobalGroupMems[xx].Address, "TxDecSetM", data)
		}
	}

	return nil
}

//HandleTxBlock when receives a txblock
func HandleTxBlock(data []byte) error {
	tmp := new(basic.TxBlock)
	err := tmp.Decode(&data, 0)
	if err != nil {
		return err
	}
	CacheDbRef.Mu.RLock()
	CacheDbRef.PreTxBlock(tmp)
	CacheDbRef.Mu.RUnlock()
	CacheDbRef.Mu.Lock()
	CacheDbRef.GetTxBlock(tmp)
	CacheDbRef.Mu.Unlock()
	return nil
}

//HandleFinalTxBlock when receives a txblock
func HandleFinalTxBlock(data []byte) error {
	tmp := new(basic.TxBlock)
	err := tmp.Decode(&data, 1)
	if err != nil {
		return err
	}
	CacheDbRef.Mu.Lock()
	CacheDbRef.GetFinalTxBlock(tmp)
	CacheDbRef.Mu.Unlock()
	return nil
}

//HandleAndSentFinalTxBlock when receives a txblock
func HandleAndSentFinalTxBlock(data []byte) error {
	tmp := new(basic.TxBlock)
	err := tmp.Decode(&data, 1)
	if err != nil {
		return err
	}
	CacheDbRef.Mu.Lock()
	CacheDbRef.GetFinalTxBlock(tmp)
	CacheDbRef.Mu.Unlock()
	xx := shard.MyMenShard.InShardId
	for i := uint32(0); i < gVar.ShardCnt; i++ {
		if i != CacheDbRef.ShardNum {
			sendTxMessage(shard.GlobalGroupMems[shard.ShardToGlobal[i][xx]].Address, "FinalTxBM", data)
		}
	}
	return nil
}

/*************************Leader**************************/

//HandleTxLeader when receives a tx
func HandleTxLeader(data []byte) error {
	tmp := new(basic.Transaction)
	err := tmp.Decode(&data)
	if err != nil {
		return err
	}
	CacheDbRef.Mu.Lock()
	CacheDbRef.MakeTXList(tmp)
	CacheDbRef.Mu.Unlock()
	return nil
}

//HandleTxDecLeader when receives a txdec
func HandleTxDecLeader(data []byte) error {
	tmp := new(basic.TxDecision)
	err := tmp.Decode(&data)
	if err != nil {
		return err
	}
	CacheDbRef.Mu.Lock()
	CacheDbRef.PreTxDecision(tmp, tmp.HashID)
	CacheDbRef.UpdateTXCache(tmp)
	CacheDbRef.Mu.Unlock()
	return nil
}

//HandleTxDecSetLeader when receives a txdecset
func HandleTxDecSetLeader(data []byte) error {
	tmp := new(basic.TxDecSet)
	err := tmp.Decode(&data)
	if err != nil {
		return err
	}
	CacheDbRef.Mu.RLock()
	CacheDbRef.PreTxDecSet(tmp)
	CacheDbRef.Mu.RUnlock()
	CacheDbRef.Mu.Lock()
	CacheDbRef.ProcessTDS(tmp)
	CacheDbRef.Mu.Unlock()
	return nil
}
