package httpc

import (
	"fmt"
	"io"
	"net/http"
	"netares/internal/parser"
	outforms "netares/internal/parser/output_forms"
	"strings"
	"sync"
	"time"
)

// * HTTPClient provides to make requests asynchronously, if needed.
// * It also processes the results and handles the key role in the software.
// * It uses the OutputForm to format the output.
// ! There is some excuses, especially in the bad connection statements, where you can't make a request.
// ! It's a temporary solution, but it works. Will be fixed in the future.
type HTTPClient struct {
	masks          []*parser.ParsedMask // ? Slice of masks for parsing.
	outputFormType string               // ? Form to format the output.
	target         string               // ? Target URL.
	retries        int                  // ? Number of retries for fetching dynamic content.
	delay          time.Duration        // ? Delay between retries.
}

// ? NewHTTPClient initializes and returns an HTTPClient instance.
func NewHTTPClient(masks []*parser.ParsedMask, form string, target string, retries int, delay time.Duration) *HTTPClient {
	return &HTTPClient{
		masks:          masks,
		outputFormType: form,
		target:         target,
		retries:        retries,
		delay:          delay,
	}
}

// ? fetchWithRetries attempts to fetch data with retries.
// ! It's a temporary solution, but it works. Will be fixed in the future.
func (httpc *HTTPClient) fetchWithRetries(mask *parser.ParsedMask) (string, error) {
	// * There is an interesting choice for trying to use incapsulated memory for the method
	// * but not to create much more instances moreover using := and cycle.
	var (
		resp      *http.Response
		err       error
		parseBody *parser.ParsedBody
		form      outforms.OutputForm
	)

	for attempt := 0; attempt < httpc.retries; attempt++ {
		// ? Send GET request.
		resp, err = http.Get(mask.CreateTargetLink())
		if err == nil {
			// ? Process the response body if no error.
			// ! Body closes in the next methods.
			parseBody = parser.NewParsedBody(mask.TargetName, mask.Fields)
			if err := parseBody.Parse(resp.Body); err != nil {
				resp.Body.Close()
				fmt.Printf("Error parsing response body from %s: %v\n", mask.CreateTargetLink(), err)
				return "", err
			}
			resp.Body.Close()
			form = outforms.NewOutputForm(httpc.outputFormType)
			return form.Review(parseBody)
		}

		// ? Log the error and retry after the delay.
		fmt.Printf("Attempt %d: Error fetching data from %s: %v\n", attempt+1, mask.CreateTargetLink(), err)
		time.Sleep(httpc.delay)
	}
	// ? Return error after all retries failed.
	return "", fmt.Errorf("failed to fetch data from %s after %d attempts", mask.CreateTargetLink(), httpc.retries)
}

// ? Research performs a GET request to the target URL with the given params, and returns the response body.
// ? It can effectively can help us to understand the ways of finding needable data.
func (httpc *HTTPClient) Research() ([]byte, error) {
	resp, err := http.Get(httpc.target)
	if err != nil {
		fmt.Printf("Error fetching data from %s: %v\n", httpc.target, err)
		return nil, err
	}

	// ? Process the response body.
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		fmt.Printf("Error reading response body from %s: %v\n", httpc.target, err)
		return nil, err
	}
	return body, nil
}

// ? Do performs the requests concurrently and processes the results.
func (httpc *HTTPClient) Do() string {
	forms := make([]string, len(httpc.masks))
	var wg sync.WaitGroup
	var mu sync.Mutex

	for idx, v := range httpc.masks {
		wg.Add(1)
		go func(v *parser.ParsedMask) {
			defer wg.Done()
			v.TargetName = httpc.target

			// ? Fetch data with retries.
			result, err := httpc.fetchWithRetries(v)
			if err != nil {
				fmt.Printf("Error processing mask %s: %v\n", v.CreateTargetLink(), err)
				return
			}

			// ! Safely append to the forms slice.
			// TODO: Why does I need a mu there? I can`t see any other ways, where I need to read over him.
			mu.Lock()
			forms[idx] = result
			mu.Unlock()
		}(v)
	}

	// ? Wait for all goroutines to finish.
	wg.Wait()
	return strings.Join(forms, "\n\n\n")
}
