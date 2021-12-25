package main

type worker struct {
	ss secretStore // file or redis or postgres or embeded db or ...
	ds dataSource  // 増やせるだけ...
}

func newWorker(ss secretStore, ds dataSource) *worker {
	return &worker{
		ss: ss,
		ds: ds,
	}
}

type secretStore interface {
	Store(dataSourceType, settings string) error
	Load(dataSourceType string) (interface{}, error)
}

type dataSource interface {
	Parse(string) error
	Connect(interface{}) error
	Execute(string) (string, error)
	Close() error
}

func (w *worker) runSettings(typ, settings string) error {
	// data sourceへの接続確認
	if err := w.ds.Connect(settings); err != nil {
		return err
	}
	defer w.ds.Close()

	return w.ss.Store(typ, settings)
}

func (w *worker) runQuery(typ, query string) (string, error) {
	secret, err := w.ss.Load(typ)
	if err != nil {
		return "", err
	}

	if err := w.ds.Parse(query); err != nil {
		return "", err
	}

	if err := w.ds.Connect(secret); err != nil {
		return "", err
	}
	defer w.ds.Close()

	result, err := w.ds.Execute(query)
	if err != nil {
		return "", err
	}

	return result, nil
}
