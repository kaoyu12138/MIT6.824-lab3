package raft

import(
	"6.824/labgob"
	"bytes"
)


func (rf *Raft) persist() {
	writer := new(bytes.Buffer)
	encoder := labgob.NewEncoder(writer)
	if encoder.Encode(rf.votedFor) == nil &&
		encoder.Encode(rf.currentTerm) == nil &&
		encoder.Encode(rf.log) == nil &&
		encoder.Encode(rf.lastsnapshotIndex) == nil &&
		encoder.Encode(rf.lastsnapshotTerm) == nil{
		rf.persister.SaveStateAndSnapshot(writer.Bytes(), rf.snapshot)
	}
}


func (rf *Raft) readPersist(data []byte) {
	if data == nil || len(data) < 1 {
		return
	}
	decoder := labgob.NewDecoder(bytes.NewBuffer(data))
	var votedFor, currentTerm, lastsnapshotIndex, lastsnapshotTerm int
	var logs Log

	if decoder.Decode(&votedFor) == nil &&
		decoder.Decode(&currentTerm) == nil &&
		decoder.Decode(&logs) == nil &&
		decoder.Decode(&lastsnapshotIndex) == nil &&
		decoder.Decode(&lastsnapshotTerm) == nil{
		rf.votedFor = votedFor
		rf.currentTerm = currentTerm
		rf.log = logs
	
		SnapshotData := rf.persister.ReadSnapshot()
		if len(SnapshotData) > 0 {
			rf.snapshot = SnapshotData
			rf.lastsnapshotIndex = lastsnapshotIndex
			rf.lastsnapshotTerm = lastsnapshotTerm
			if rf.commitIndex < lastsnapshotIndex{
				rf.commitIndex = lastsnapshotIndex
				rf.applyCond.Broadcast()
			}
		}
	}
}