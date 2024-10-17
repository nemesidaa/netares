package domhttp

import (
	"errors"
	"io"
	"log"
	"sync"
	"time"
)

// * Domrouter is a struct provided for routing XPath queries
type DomRouter struct {
	queries map[string]string // ! Format: name>query
}

func NewRouter(queries map[string]string) *DomRouter {
	return &DomRouter{
		queries: queries,
	}
}

var (
	ErrTimedOut = errors.New("timed out")
)

// * Solve parses XPath and returns the result
func (dr *DomRouter) Solve(body io.ReadCloser, timeout time.Duration) (map[string]string, error) {
	results := make(map[string]string)
	timeoutTrigger := time.NewTimer(timeout)
	successTrigger := make(chan struct{})
	timedOut := false

	mx := sync.Mutex{}
	wg := sync.WaitGroup{}

	rawHtmlBody, err := io.ReadAll(body)
	if err != nil {
		log.Printf("error: %v", err)
		return nil, err
	}

	defer body.Close()

	for name, query := range dr.queries {
		wg.Add(1)
		go func(name, query string) {
			defer wg.Done()

			retDir, err := find(rawHtmlBody, query, name)
			if err != nil {
				log.Printf("error: %v", err)
			}

			mx.Lock()
			defer mx.Unlock()

			if timedOut {
				return
			}
			results[retDir.name] = retDir.value

		}(name, query)
	}

	go func() {
		wg.Wait()
		close(successTrigger)
	}()

	select {
	case <-timeoutTrigger.C:
		mx.Lock()
		timedOut = true
		mx.Unlock()
		_ = timeoutTrigger.Stop()
		return results, ErrTimedOut
	case <-successTrigger:
		// * Happy path.
	}

	return results, nil
}

// this incapsulated method does logic of finding.
// It uses bodycutting method to find the query and to
func find(body []byte, route, target string) (returnDir, error) {
	qd, err := ParseQuery(route, body)
	if err != nil {
		return returnDir{}, err
	}
	return returnDir{name: target, value: qd.Result()}, nil
}

type returnDir struct {
	name  string
	value string
}
