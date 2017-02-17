// Copyright 2010 The Go Authors. All rights reserved.
     2  // Use of this source code is governed by a BSD-style
     3  // license that can be found in the LICENSE file.
     4
     5  // Package pprof serves via its HTTP server runtime profiling data
     6  // in the format expected by the pprof visualization tool.
     7  //
     8  // The package is typically only imported for the side effect of
     9  // registering its HTTP handlers.
    10  // The handled paths all begin with /debug/pprof/.
    11  //
    12  // To use pprof, link this package into your program:
    13  //  import _ "net/http/pprof"
    14  //
    15  // If your application is not already running an http server, you
    16  // need to start one. Add "net/http" and "log" to your imports and
    17  // the following code to your main function:
    18  //
    19  //  go func() {
    20  //    log.Println(http.ListenAndServe("localhost:6060", nil))
    21  //  }()
    22  //
    23  // Then use the pprof tool to look at the heap profile:
    24  //
    25  //  go tool pprof http://localhost:6060/debug/pprof/heap
    26  //
    27  // Or to look at a 30-second CPU profile:
    28  //
    29  //  go tool pprof http://localhost:6060/debug/pprof/profile
    30  //
    31  // Or to look at the goroutine blocking profile, after calling
    32  // runtime.SetBlockProfileRate in your program:
    33  //
    34  //  go tool pprof http://localhost:6060/debug/pprof/block
    35  //
    36  // Or to collect a 5-second execution trace:
    37  //
    38  //  wget http://localhost:6060/debug/pprof/trace?seconds=5
    39  //
    40  // To view all available profiles, open http://localhost:6060/debug/pprof/
    41  // in your browser.
    42  //
    43  // For a study of the facility in action, visit
    44  //
    45  //  https://blog.golang.org/2011/06/profiling-go-programs.html
    46  //
    47  package pprof
    48
    49  import (
    50    "bufio"
    51    "bytes"
    52    "fmt"
    53    "html/template"
    54    "io"
    55    "log"
    56    "net/http"
    57    "os"
    58    "runtime"
    59    "runtime/pprof"
    60    "runtime/trace"
    61    "strconv"
    62    "strings"
    63    "time"
    64  )
    65
    66  func init() {
    67    http.Handle("/debug/pprof/", http.HandlerFunc(Index))
    68    http.Handle("/debug/pprof/cmdline", http.HandlerFunc(Cmdline))
    69    http.Handle("/debug/pprof/profile", http.HandlerFunc(Profile))
    70    http.Handle("/debug/pprof/symbol", http.HandlerFunc(Symbol))
    71    http.Handle("/debug/pprof/trace", http.HandlerFunc(Trace))
    72  }
    73
    74  // Cmdline responds with the running program's
    75  // command line, with arguments separated by NUL bytes.
    76  // The package initialization registers it as /debug/pprof/cmdline.
    77  func Cmdline(w http.ResponseWriter, r *http.Request) {
    78    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    79    fmt.Fprintf(w, strings.Join(os.Args, "\x00"))
    80  }
    81
    82  func sleep(w http.ResponseWriter, d time.Duration) {
    83    var clientGone <-chan bool
    84    if cn, ok := w.(http.CloseNotifier); ok {
    85      clientGone = cn.CloseNotify()
    86    }
    87    select {
    88    case <-time.After(d):
    89    case <-clientGone:
    90    }
    91  }
    92
    93  // Profile responds with the pprof-formatted cpu profile.
    94  // The package initialization registers it as /debug/pprof/profile.
    95  func Profile(w http.ResponseWriter, r *http.Request) {
    96    sec, _ := strconv.ParseInt(r.FormValue("seconds"), 10, 64)
    97    if sec == 0 {
    98      sec = 30
    99    }
   100
   101    // Set Content Type assuming StartCPUProfile will work,
   102    // because if it does it starts writing.
   103    w.Header().Set("Content-Type", "application/octet-stream")
   104    if err := pprof.StartCPUProfile(w); err != nil {
   105      // StartCPUProfile failed, so no writes yet.
   106      // Can change header back to text content
   107      // and send error code.
   108      w.Header().Set("Content-Type", "text/plain; charset=utf-8")
   109      w.WriteHeader(http.StatusInternalServerError)
   110      fmt.Fprintf(w, "Could not enable CPU profiling: %s\n", err)
   111      return
   112    }
   113    sleep(w, time.Duration(sec)*time.Second)
   114    pprof.StopCPUProfile()
   115  }
   116
   117  // Trace responds with the execution trace in binary form.
   118  // Tracing lasts for duration specified in seconds GET parameter, or for 1 second if not specified.
   119  // The package initialization registers it as /debug/pprof/trace.
   120  func Trace(w http.ResponseWriter, r *http.Request) {
   121    sec, err := strconv.ParseFloat(r.FormValue("seconds"), 64)
   122    if sec <= 0 || err != nil {
   123      sec = 1
   124    }
   125
   126    // Set Content Type assuming trace.Start will work,
   127    // because if it does it starts writing.
   128    w.Header().Set("Content-Type", "application/octet-stream")
   129    if err := trace.Start(w); err != nil {
   130      // trace.Start failed, so no writes yet.
   131      // Can change header back to text content and send error code.
   132      w.Header().Set("Content-Type", "text/plain; charset=utf-8")
   133      w.WriteHeader(http.StatusInternalServerError)
   134      fmt.Fprintf(w, "Could not enable tracing: %s\n", err)
   135      return
   136    }
   137    sleep(w, time.Duration(sec*float64(time.Second)))
   138    trace.Stop()
   139  }
   140
   141  // Symbol looks up the program counters listed in the request,
   142  // responding with a table mapping program counters to function names.
   143  // The package initialization registers it as /debug/pprof/symbol.
   144  func Symbol(w http.ResponseWriter, r *http.Request) {
   145    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
   146
   147    // We have to read the whole POST body before
   148    // writing any output. Buffer the output here.
   149    var buf bytes.Buffer
   150
   151    // We don't know how many symbols we have, but we
   152    // do have symbol information. Pprof only cares whether
   153    // this number is 0 (no symbols available) or > 0.
   154    fmt.Fprintf(&buf, "num_symbols: 1\n")
   155
   156    var b *bufio.Reader
   157    if r.Method == "POST" {
   158      b = bufio.NewReader(r.Body)
   159    } else {
   160      b = bufio.NewReader(strings.NewReader(r.URL.RawQuery))
   161    }
   162
   163    for {
   164      word, err := b.ReadSlice('+')
   165      if err == nil {
   166        word = word[0 : len(word)-1] // trim +
   167      }
   168      pc, _ := strconv.ParseUint(string(word), 0, 64)
   169      if pc != 0 {
   170        f := runtime.FuncForPC(uintptr(pc))
   171        if f != nil {
   172          fmt.Fprintf(&buf, "%#x %s\n", pc, f.Name())
   173        }
   174      }
   175
   176      // Wait until here to check for err; the last
   177      // symbol will have an err because it doesn't end in +.
   178      if err != nil {
   179        if err != io.EOF {
   180          fmt.Fprintf(&buf, "reading request: %v\n", err)
   181        }
   182        break
   183      }
   184    }
   185
   186    w.Write(buf.Bytes())
   187  }
   188
   189  // Handler returns an HTTP handler that serves the named profile.
   190  func Handler(name string) http.Handler {
   191    return handler(name)
   192  }
   193
   194  type handler string
   195
   196  func (name handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
   197    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
   198    debug, _ := strconv.Atoi(r.FormValue("debug"))
   199    p := pprof.Lookup(string(name))
   200    if p == nil {
   201      w.WriteHeader(404)
   202      fmt.Fprintf(w, "Unknown profile: %s\n", name)
   203      return
   204    }
   205    gc, _ := strconv.Atoi(r.FormValue("gc"))
   206    if name == "heap" && gc > 0 {
   207      runtime.GC()
   208    }
   209    p.WriteTo(w, debug)
   210    return
   211  }
   212
   213  // Index responds with the pprof-formatted profile named by the request.
   214  // For example, "/debug/pprof/heap" serves the "heap" profile.
   215  // Index responds to a request for "/debug/pprof/" with an HTML page
   216  // listing the available profiles.
   217  func Index(w http.ResponseWriter, r *http.Request) {
   218    if strings.HasPrefix(r.URL.Path, "/debug/pprof/") {
   219      name := strings.TrimPrefix(r.URL.Path, "/debug/pprof/")
   220      if name != "" {
   221        handler(name).ServeHTTP(w, r)
   222        return
   223      }
   224    }
   225
   226    profiles := pprof.Profiles()
   227    if err := indexTmpl.Execute(w, profiles); err != nil {
   228      log.Print(err)
   229    }
   230  }
   231
   232  var indexTmpl = template.Must(template.New("index").Parse(`<html>
   233  <head>
   234  <title>/debug/pprof/</title>
   235  </head>
   236  <body>
   237  /debug/pprof/<br>
   238  <br>
   239  profiles:<br>
   240  <table>
   241  {{range .}}
   242  <tr><td align=right>{{.Count}}<td><a href="{{.Name}}?debug=1">{{.Name}}</a>
   243  {{end}}
   244  </table>
   245  <br>
   246  <a href="goroutine?debug=2">full goroutine stack dump</a><br>
   247  </body>
   248  </html>
   249  `))