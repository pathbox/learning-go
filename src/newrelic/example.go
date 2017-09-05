// go get github.com/newrelic/go-agent
// config := newrelic.NewConfig("Your Application Name", "__YOUR_NEW_RELIC_LICENSE_KEY__")
// app, err := newrelic.NewApplication(config)

// Using NR to instrument calls to your API

import "github.com/newrelic/go-agent"

func createNRApp(licenseKey string)(newrelic.Application, error){
  cfg := newrelic.NewConfig(”YourNRAppName”, licenseKey)
  app, err := newrelic.NewApplication(cfg)

  if err != nil {

        return nil, err

    }

    return app, nil
}

type NewRelicContextKey string

var NRKey = NewRelicContextKey(“NewRelicTxn”)

func newRelicMiddleware(app newrelic.Application) negroni.Handler {

    return negroni.HandlerFunc(func(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {

        txn := app.StartTransaction(”txnName”, w, req) // 这里是Newrelic的真正处理方法
        defer txn.End()

        ctx := context.WithValue(req.Context(), NRKey, txn)
        req = req.WithContext(ctx)

        next(txn, req)

    }
}

/*
Now every request context will have NR txn embedded within it. Strictly speaking we did not have to embed it into the context, if all we wanted to do was instrument calls to our APIs. However, by placing the txn into the context and propagating the request context down your call chain, you can use NR to instrument other aspects of your code, as you will see in the next section.

Using NR to instrument calls to your DB
*/

func NRDS(txn newrelic.Transaction, tableName, operation string) newrelic.DatastoreSegment {
    s := newrelic.DatastoreSegment{
        Product:    newrelic.DatastoreMySQL,
        Collection: tableName,
        Operation:  operation,
    }
    s.StartTime = newrelic.StartSegmentNow(txn)
    return s
}

// With the above helper function at hand, adding instrumentation in your models becomes trivial

func myModelFunc(ctx context.Context, args...) {
  txn, _ := ctx.Value(NRKey).(newrelic.Transaction)

  s := NRDS(txn, "MyTableName", "SELECT")
  defer s.End()

  // your model code
}


// Using NR to instrument external calls

func NRES(txn newrelic.Transaction, url string) newrelic.ExternalSegment {
    return newrelic.ExternalSegment{
        StartTime: newrelic.StartSegmentNow(txn),
        URL:       url,
    }
}

func myExternalFunc(ctx context.Context, url string) {

    txn, _ := ctx.Value(NRKey).(newrelic.Transaction)

    seg := NRES(txn, url)
    defer seg.End()

    // your external call code here

}


// Using NR to instrument background workers

func NRMakeCtxAndTxn(app *newrelic.Application, txnName string) (newrelic.Transaction, context.Context) {
    if app == nil {
        return nil, context.Background()
    }
    txn := (*app).StartTransaction(txnName, nil, nil)
    ctx := context.WithValue(context.Background(), NRKey, txn)
    return txn, ctx
}

func NREndTxn(txn newrelic.Transaction) {
    if txn != nil {
    txn.End()
    }
}

func myWorkerTask(nrApp *newrelic.Application) {

    txn, _ := NRMakeCtxAndTxn(nrApp, “taskName”)
    defer NREndTxn(txn)

    // task code goes here…

}



//ref: http://jimmyislive.tumblr.com/post/164511785370/using-newrelic-with-go?utm_source=golangweekly&utm_medium=email