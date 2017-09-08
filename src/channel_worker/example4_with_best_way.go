// Go channel的常用编程模式。使用一个两级channel系统，一个用来存放任务队列，另一个用来控制处理任务队列的并发量

type PayloadCollection struct {
    WindowsVersion  string    `json:"version"`
    Token           string    `json:"token"`
    Payloads        []Payload `json:"data"`
}

type Payload struct {
    // [redacted]
}

func (p *Payload) UploadToS3() error {
    // the storageFolder method ensures that there are no name collision in
    // case we get same timestamp in the key name
    storage_path := fmt.Sprintf("%v/%v", p.storageFolder, time.Now().UnixNano())

    bucket := S3Bucket

    b := new(bytes.Buffer)
    encodeErr := json.NewEncoder(b).Encode(payload)
    if encodeErr != nil {
        return encodeErr
    }

    // Everything we post to the S3 bucket should be marked 'private'
    var acl = s3.Private
    var contentType = "application/octet-stream"

    return bucket.PutReader(storage_path, b, int64(b.Len()), contentType, acl, s3.Options{})
}

var (
  MaxWorker = os.Getenv("MAX_WORKERS")
  MaxQueue = os.Getenv("MAX_QUEUE")
)

// Job represents the job to be run
type Job struct {
  Payload Payload
}

// A buffered channel that we can send work requests on
var JobQueue chan Job

// Worker represents the worker that executes the job

type Worker struct {
  WorkerPool chan chan Job // WorkerPool 是一个chan，处理chan Job
  JobChannel chan Job
  quit chan bool
}

func NewWorker(workerPool chan chan Job) Worker {
  return Worker{
    WorkerPool: workerPool,
    JobChannel: make(chan Job),
    quit: make(chan bool),
  }
}

 // Start method starts the run loop for the worker, listening for a quit channel in
// case we need to stop it
// 上传操作
func (w Worker) Start(){
  go func() {
    for {
       // register the current worker into the worker queue.
      w.WorlerPool <- w.JobChannel

      select{
      case job := <-w.JobChannel:
        if err := job.Payload.UploadToS3(); err != nil {
          log.Errorf("Error uploading to S3: %s", err.Error())
        }
      case <-w.quit:
        // we have received a signal to stop
        return
      }
    }
  }()
}

// Stop signals the worker to stop listening for work requests.
func (w Worker) Stop() {
  go func() {
    w.quit <- true
  }()
}

// 我们修改了HTTP请求处理函数来创建一个含有载荷（payload）的Job 结构，然后将它送到一个叫JobQueue 的channel。worker会对它们进行处理。
// 处理请求过来的http request
func payloadHandler(w http.ResponseWrite, r *http.Request) {
  if r.Method != "POST" {
    w.WriteHeader(http.StatusMethodNotAllowed)
    return
  }

  // Read the body into a string for json decoding
  var content = &PayloadCollection{}
  err := json.NewDecoder(io.LimitReader(r.Body, MaxLength)).Decode(&content)

  if err != nil {
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusBadRequest)
    return
  }

   // Go through each payload and queue items individually to be posted to S3
  for _, payload := range content.Payloads {

    work := Job{Payload: payload}

    JobQueue <- work
  }

  w.WriteHeader(http.StatusOK)
}


// 在初始化服务的时候，我们创建了一个Dispatcher 并且调用了Run() 函数来创建worker池。这些worker会监听JobQueue 上是否有新的任务并进行处理。

dispatcher := NewDispatcher(MaxWorker)
dispatcher.Run()

type Dispatcher struct {
  // A pool of workers channels that are registered with the dispatcher(A simple pool is a queue )
  WorkerPool chan chan Job
}

func NewDispatcher(maxWorkers int) *Dispatcher{
  pool := make(chan chan Job, maxWorkers)
  return &Dispatcher{WorkerPool: pool}
}

func (d *Dispatcher) Run() {
  // starting n number of workers
  for i := 0; i < d.maxWorkers; i++ {
    worker := NewWorker(d.pool)
    worker.Start()
  }

  go d.dispatch()
}

func (d *Dispatcher) dispatch() {
  for {
    select{
    case job := <-JobQueue:
      go func(job Job) {
        // try to obtain a worker job channel that is available.
        // this will block until a worker is idle
        jobChannel := <-d.WorkerPool

        jobChannel <- job
      }(job)
    }
  }
}

