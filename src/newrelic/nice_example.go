
func NewNewRelicApp() newrelic.Application {
	towerGoCnf := config.GetTowerGoCnf()
	newrelicConfig := newrelic.NewConfig(towerGoCnf.NewRelicApp, towerGoCnf.NewRelicSecret)
	app, err := newrelic.NewApplication(newrelicConfig)
	if err != nil {
		log.Error("newrelic NewApplication error", err)
	}

	return app
}


// Custom transaction
func TraceCustomTransaction() {
	txn := NewRelicApp.StartTransaction("BaiduTest", nil, nil) // new custom Go transaction
	r, _ := myExternalFunc(txn)                                // for external services
	log.Println(r.StatusCode)
	defer txn.End()
}

// external services
func useNewRoundTripper(txn newrelic.Transaction) (*http.Response, error) {
	client := &http.Client{}
	client.Transport = newrelic.NewRoundTripper(txn, nil)
	resp, err := client.Get("http://www.baidu.com/")
	return resp, err
}

func NewRES(txn newrelic.Transaction, url string) newrelic.ExternalSegment {
	return newrelic.ExternalSegment{
		StartTime: newrelic.StartSegmentNow(txn),
		URL:       url,
	}
}

// external services
func myExternalFunc(txn newrelic.Transaction) (*http.Response, error) {

	url := "http://www.baidu.com/"
	seg := NewRES(txn, url)
	defer seg.End()

	client := &http.Client{}
	client.Transport = newrelic.NewRoundTripper(txn, nil)
	resp, err := client.Get(url)
	return resp, err

}