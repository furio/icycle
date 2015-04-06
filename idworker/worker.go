package idworker

import (
    "log"
    "sync"
    "time"
    "fmt"
)

const (
    golangMillis = int64(time.Millisecond)
    twitterEpoch = int64(1288834974657)

    workerIdBits = uint64(5)
    datacenterIdBits = uint64(5)
    maxWorkerId = int64(-1) ^ (int64(-1) << workerIdBits)
    maxDatacenterId = int64(-1) ^ (int64(-1) << datacenterIdBits)
    sequenceBits = uint64(12)
    workerIdShift = sequenceBits
    datacenterIdShift = sequenceBits + workerIdBits
    timestampLeftShift = sequenceBits + workerIdBits + datacenterIdBits
    sequenceMask = int64(-1) ^ (int64(-1) << sequenceBits)
)

type IdWorker struct {
    workerId int64
    datacenterId int64
    sequence int64
    lastTimestamp int64
    mutex sync.Mutex
}

func NewIdWorker(wId int64, dId int64, lts int64) (*IdWorker, error) {
    if (wId > maxWorkerId || wId < 0) {
        return nil, fmt.Errorf("WorkerId %d is out of bounds", wId)
    }

    if (dId > maxDatacenterId || dId < 0) {
        return nil, fmt.Errorf("DatacenterId %d is out of bounds", dId)
    }

    if (lts > millis()) {
        log.Printf("Worker will generate id after a while.")
    }

    w := &IdWorker{}
    w.workerId = wId
    w.datacenterId = dId
    w.lastTimestamp = lts
    w.sequence = int64(0)

    return w,nil
}

func millis() int64 {
    return time.Now().UnixNano() / golangMillis
}

func (w *IdWorker) String() string {
    return fmt.Sprintf("WorkerId: %d, DatacenterId: %d", w.workerId, w.datacenterId)
}

func (w *IdWorker) Timestamp() int64 {
    return millis();
}

func (w *IdWorker) WorkerId() int64 {
    return w.workerId;
}

func (w *IdWorker) DatacenterId() int64 {
    return w.datacenterId;
}

func (w *IdWorker) NextId() (int64, error) {
    w.mutex.Lock()
    defer w.mutex.Unlock()

    timeStamp := millis()
    if (timeStamp < w.lastTimestamp) {
        return 0,fmt.Errorf("Clock moved backwards. Refusing to generate id for %d milliseconds", w.lastTimestamp - timeStamp)
    }

    if (w.lastTimestamp == timeStamp) {
        w.sequence = (w.sequence + 1) & sequenceMask
        if (w.sequence == 0) {
            for (timeStamp <= w.lastTimestamp) {
                timeStamp = millis()
            }
        }
    } else {
        w.sequence = 0
    }

    w.lastTimestamp = timeStamp

    return ((timeStamp - twitterEpoch) << timestampLeftShift) |
            (w.datacenterId << datacenterIdShift) |
            (w.workerId << workerIdShift) |
            w.sequence, nil
}