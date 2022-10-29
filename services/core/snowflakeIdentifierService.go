package core

import (
    "errors"
    "sync"
    "time"
)

const (
    workerIDBits     = uint64(5) // 5bit workerID out of 10bit worker machine ID
    dataCenterIDBits = uint64(5) // 5bit workerID out of 10bit worker dataCenterID
    sequenceBits     = uint64(12)

    maxWorkerID     = int64(-1) ^ (int64(-1) << workerIDBits) // The maximum value of the node ID used to prevent overflow
    maxDataCenterID = int64(-1) ^ (int64(-1) << dataCenterIDBits)
    maxSequence     = int64(-1) ^ (int64(-1) << sequenceBits)

    timeLeft = uint8(22) // timeLeft = workerIDBits + sequenceBits // Timestamp offset left
    dataLeft = uint8(17) // dataLeft = dataCenterIDBits + sequenceBits
    workLeft = uint8(12) // workLeft = sequenceBits // Node IDx offset to the left

    twepoch = int64(1666543203000) // constant timestamp (milliseconds)
)

type Worker struct {
    mu           sync.Mutex
    LastStamp    int64 // Record the timestamp of the last ID
    WorkerID     int64 // the ID of the node
    DataCenterID int64 // The data center ID of the node
    Sequence     int64 // ID sequence numbers that have been generated in the current millisecond (accumulated from 0) A maximum of 4096 IDs are generated within 1 millisecond
}

func NewWorker(workerID, dataCenterID int64) *Worker {
    return &Worker{
        WorkerID:     workerID,
        LastStamp:    0,
        Sequence:     0,
        DataCenterID: dataCenterID,
    }
}

func (w *Worker) NextID() (uint64, error) {

    w.mu.Lock()
    defer w.mu.Unlock()

    return w.nextID()
}

func (w *Worker) nextID() (uint64, error) {

    timeStamp := w.getMilliSeconds()
    if timeStamp < w.LastStamp {
        return 0, errors.New("time is moving backwards, waiting until")
    }

    if w.LastStamp == timeStamp {

        w.Sequence = (w.Sequence + 1) & maxSequence

        if w.Sequence == 0 {
            for timeStamp <= w.LastStamp {
                timeStamp = w.getMilliSeconds()
            }
        }
    } else {
        w.Sequence = 0
    }

    w.LastStamp = timeStamp
    id := ((timeStamp - twepoch) << timeLeft) |
        (w.DataCenterID << dataLeft) |
        (w.WorkerID << workLeft) |
        w.Sequence

    return uint64(id), nil
}

func (w *Worker) getMilliSeconds() int64 {
    return time.Now().UnixMilli()
}
