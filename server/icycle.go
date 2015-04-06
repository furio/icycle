package server

import (
    "fmt"
    "net/http"
    "encoding/json"
    "runtime"
    "flag"
    "log"
    "strconv"
    "github.com/furio/icycle/idworker"
    "time"
)

var (
    workerId = flag.Int64("w", 0, "Worker id")
    datacenterId = flag.Int64("d", 0, "Datacenter id")
    port = flag.String("p", "9000", "Port to listen on")
    lastStamp = flag.Int64("t", -1, "Last timestamp in milliseconds")
)

var (
    idGenerator *idworker.IdWorker
    statsAccumulator *Stats
)


func handlerHome(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi from icycle: %s", idGenerator.String())
}

func handlerStats(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "%s", statsAccumulator.TotalStats())
}

func handlerId(w http.ResponseWriter, r *http.Request) {
    defer statsAccumulator.RecordRequest(time.Now())

    seq,err := idGenerator.NextId();
    profile := map[string]interface{}{"sequence": seq, "error": err}

    js, jerr := json.Marshal(profile)
    if jerr != nil {
        http.Error(w, jerr.Error(), http.StatusInternalServerError)
        return
    }

    // log.Printf("%s", js)
    w.Header().Set("Content-Type", "application/json")
    w.Write(js)
}

func handlerIdStr(w http.ResponseWriter, r *http.Request) {
    defer statsAccumulator.RecordRequest(time.Now())

    seq, err := idGenerator.NextId();
    profile := map[string]interface{}{"sequence": strconv.FormatInt(seq, 10), "error": err}

    js, jerr := json.Marshal(profile)
    if jerr != nil {
        http.Error(w, jerr.Error(), http.StatusInternalServerError)
        return
    }

    // log.Printf("%s", js)
    w.Header().Set("Content-Type", "application/json")
    w.Write(js)
}

func handlerWorker(w http.ResponseWriter, r *http.Request) {
    ts := map[string]int64{"workerId": idGenerator.WorkerId(), "datacenterId": idGenerator.DatacenterId()}

    js, jerr := json.Marshal(ts)
    if jerr != nil {
        http.Error(w, jerr.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write(js)
}

func handlerWorkerTimestamp(w http.ResponseWriter, r *http.Request) {
    ts := map[string]string{"timestamp": strconv.FormatInt(idGenerator.Timestamp(), 10)}

    js, jerr := json.Marshal(ts)
    if jerr != nil {
        http.Error(w, jerr.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write(js)
}

func Main() {
    if (!initWorker()) {
        log.Fatal("Unable to init worker")
        return
    }

    statsAccumulator = NewStats()

    runtime.GOMAXPROCS(runtime.NumCPU())
    http.HandleFunc("/", handlerHome)

    http.HandleFunc("/id", handlerId)
    http.HandleFunc("/id/str", handlerIdStr)
    http.HandleFunc("/worker", handlerWorker)
    http.HandleFunc("/worker/timestamp", handlerWorkerTimestamp)
    http.HandleFunc("/stats", handlerStats)

    log.Printf("Serving on port :%s", *port)
    http.ListenAndServe(":" + *port, nil)
}

func initWorker() bool {
    flag.Parse()

    var err error = nil
    idGenerator,err = idworker.NewIdWorker(*workerId,*datacenterId,*lastStamp)

    if (idGenerator == nil) {
        log.Fatal(err)
        return false
    }

    return true
}