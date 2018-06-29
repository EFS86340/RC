package network

import (
	"github.com/uchihatmtkinu/RC/shard"
	"github.com/uchihatmtkinu/RC/gVar"
	"fmt"
)
var readymask	[]byte
func ShardProcess(){
	var beginShard	shard.Instance
	CurrentEpoch++
	shard.StartFlag = true
	fmt.Println(shard.PreviousSyncBlockHash)
	beginShard.GenerateSeed(&shard.PreviousSyncBlockHash)
	beginShard.Sharding(&shard.GlobalGroupMems, &shard.ShardToGlobal)
	//shard.MyMenShard = &shard.GlobalGroupMems[MyGlobalID]
	myShard := shard.GlobalGroupMems[MyGlobalID].Shard
	LeaderAddr = shard.GlobalGroupMems[shard.ShardToGlobal[myShard][0]].Address
	//intilizeMaskBit(&readymask, (shard.NumMems+7)>>3,false)

	if shard.MyMenShard.Role == 1{
		MinerReadyProcess()
	}	else {
		LeaderReadyProcess(&shard.GlobalGroupMems)
	}

}
func LeaderReadyProcess(ms *[]shard.MemShard){
	//var readaddr string
	readyCount := 0

	for readyCount <= int(gVar.ShardSize/2*3) {
		<-readyCh
		readyCount++
		//setMaskBit((*ms)[GlobalAddrMapToInd[readaddr]].InShardId, cosi.Enabled, &readymask)
	}

}
func MinerReadyProcess(){

	SendShardReadyMessage(LeaderAddr, "shardReady", nil)
}


func SendShardReadyMessage(addr string, command string, message interface{}) {
	payload := gobEncode(message)
	request := append(commandToBytes(command), payload...)
	sendData(addr, request)
}


func HandleShardReady(addr string) {
	readyCh <- addr

}