package main

const (
	requestClientPool = "request-client-pool"
)

func (app *application) API(name string, payload interface{}) (interface{}, error) {
	switch name {
	case requestClientPool:
		return app.clusterList.Pool, nil
	default:
		return nil, newUnknownAPICallError(name, payload)
	}
}
