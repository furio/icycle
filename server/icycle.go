package server

import (
    "fmt"
    "net/http"
    "encoding/json"
    "runtime"
    "flag"
    "log"
    "github.com/furio/icycle/idworker"
)

var (
    workerId = flag.Int64("w", 5, "Worker id")
    datacenterId = flag.Int64("d", 1, "Datacenter id")
    port = flag.String("p", "9000", "Port to listen on")
    lastStamp = flag.Int64("t", -1, "Last timestamp in milliseconds")
)
var idGenerator *idworker.IdWorker

type Sequence struct {
    Sequence int64
    Error error
}

func handlerTest(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello World!")
}

func handlerId(w http.ResponseWriter, r *http.Request) {
    seq,err := idGenerator.NextId();

    profile := Sequence{seq, err}

    js, err := json.Marshal(profile)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    log.Printf("%s", js)
    w.Header().Set("Content-Type", "application/json")
    w.Write(js)
}

func handlerWorker(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "%s", idGenerator.String())
}

func Main() {
    if (!initWorker()) {
        log.Fatal("Enable to init worker")
        return
    }

    runtime.GOMAXPROCS(runtime.NumCPU())
    http.HandleFunc("/", handlerTest)

    http.HandleFunc("/id", handlerId)
    http.HandleFunc("/worker", handlerWorker)

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