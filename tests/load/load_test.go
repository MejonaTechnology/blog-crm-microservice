package load

import (
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	// Test configuration
	testServerURL = "http://localhost:8082"
	testDuration  = 30 * time.Second
	rampUpTime    = 5 * time.Second
)

// LoadTestResult represents the results of a load test
type LoadTestResult struct {
	TotalRequests   int
	SuccessfulReqs  int
	FailedRequests  int
	AverageResponse time.Duration
	MinResponse     time.Duration
	MaxResponse     time.Duration
	RequestsPerSec  float64
	ErrorRate       float64
	StatusCodes     map[int]int
}

// RequestResult represents the result of a single HTTP request
type RequestResult struct {
	StatusCode   int
	ResponseTime time.Duration
	Error        error
}

// LoadTester manages load testing operations
type LoadTester struct {
	URL         string
	Duration    time.Duration
	Concurrency int
	RampUp      time.Duration
}

// NewLoadTester creates a new load tester
func NewLoadTester(url string, duration time.Duration, concurrency int) *LoadTester {
	return &LoadTester{
		URL:         url,
		Duration:    duration,
		Concurrency: concurrency,
		RampUp:      rampUpTime,
	}
}

// ExecuteLoadTest runs a load test with the specified parameters
func (lt *LoadTester) ExecuteLoadTest(endpoint string) (*LoadTestResult, error) {
	fullURL := lt.URL + endpoint
	results := make(chan RequestResult, lt.Concurrency*100)

	// WaitGroup to coordinate workers
	var wg sync.WaitGroup

	// Start time
	startTime := time.Now()
	endTime := startTime.Add(lt.Duration)

	// Launch workers with ramp-up
	workersStarted := 0
	for workersStarted < lt.Concurrency {
		wg.Add(1)
		go lt.worker(fullURL, results, &wg, startTime, endTime)
		workersStarted++

		// Ramp-up delay
		if lt.RampUp > 0 && workersStarted < lt.Concurrency {
			time.Sleep(lt.RampUp / time.Duration(lt.Concurrency))
		}
	}

	// Close results channel when all workers finish
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	return lt.collectResults(results, startTime)
}

// worker performs HTTP requests for the duration of the test
func (lt *LoadTester) worker(url string, results chan<- RequestResult, wg *sync.WaitGroup, startTime, endTime time.Time) {
	defer wg.Done()

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	for time.Now().Before(endTime) {
		reqStart := time.Now()
		resp, err := client.Get(url)
		reqDuration := time.Since(reqStart)

		result := RequestResult{
			ResponseTime: reqDuration,
			Error:        err,
		}

		if err == nil {
			result.StatusCode = resp.StatusCode
			resp.Body.Close()
		}

		select {
		case results <- result:
		default:
			// Channel full, skip this result
		}

		// Small delay to prevent overwhelming the server
		time.Sleep(10 * time.Millisecond)
	}
}

// collectResults processes the results from all workers
func (lt *LoadTester) collectResults(results <-chan RequestResult, startTime time.Time) (*LoadTestResult, error) {
	result := &LoadTestResult{
		StatusCodes: make(map[int]int),
		MinResponse: time.Hour, // Initialize to high value
		MaxResponse: 0,
	}

	var totalResponseTime time.Duration

	for res := range results {
		result.TotalRequests++

		if res.Error != nil {
			result.FailedRequests++
			continue
		}

		result.SuccessfulReqs++
		result.StatusCodes[res.StatusCode]++

		// Track response times
		totalResponseTime += res.ResponseTime
		if res.ResponseTime < result.MinResponse {
			result.MinResponse = res.ResponseTime
		}
		if res.ResponseTime > result.MaxResponse {
			result.MaxResponse = res.ResponseTime
		}
	}

	// Calculate metrics
	testDuration := time.Since(startTime)

	if result.SuccessfulReqs > 0 {
		result.AverageResponse = totalResponseTime / time.Duration(result.SuccessfulReqs)
	}

	if result.TotalRequests > 0 {
		result.RequestsPerSec = float64(result.TotalRequests) / testDuration.Seconds()
		result.ErrorRate = float64(result.FailedRequests) / float64(result.TotalRequests) * 100
	}

	return result, nil
}

// PrintResults displays load test results
func (result *LoadTestResult) PrintResults(t *testing.T) {
	t.Logf("=== LOAD TEST RESULTS ===")
	t.Logf("Total Requests: %d", result.TotalRequests)
	t.Logf("Successful Requests: %d", result.SuccessfulReqs)
	t.Logf("Failed Requests: %d", result.FailedRequests)
	t.Logf("Error Rate: %.2f%%", result.ErrorRate)
	t.Logf("Requests/sec: %.2f", result.RequestsPerSec)
	t.Logf("Average Response Time: %v", result.AverageResponse)
	t.Logf("Min Response Time: %v", result.MinResponse)
	t.Logf("Max Response Time: %v", result.MaxResponse)

	t.Logf("Status Code Distribution:")
	for code, count := range result.StatusCodes {
		percentage := float64(count) / float64(result.SuccessfulReqs) * 100
		t.Logf("  %d: %d (%.1f%%)", code, count, percentage)
	}
}

// TestHealthEndpointLoad tests the health endpoint under load
func TestHealthEndpointLoad(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping load tests in short mode")
	}

	tester := NewLoadTester(testServerURL, testDuration, 10)

	t.Log("Starting load test on /health endpoint...")
	result, err := tester.ExecuteLoadTest("/health")
	assert.NoError(t, err)

	result.PrintResults(t)

	// Assertions
	assert.Greater(t, result.TotalRequests, 0, "Should have made requests")
	assert.Less(t, result.ErrorRate, 5.0, "Error rate should be less than 5%")
	assert.Greater(t, result.RequestsPerSec, 10.0, "Should handle at least 10 requests/sec")
	assert.Less(t, result.AverageResponse, 1000*time.Millisecond, "Average response time should be under 1s")

	// Check that most responses are 200 OK
	okCount := result.StatusCodes[200]
	assert.Greater(t, okCount, result.TotalRequests*8/10, "At least 80% responses should be 200 OK")
}

// TestAPIEndpointLoad tests the API endpoint under load
func TestAPIEndpointLoad(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping load tests in short mode")
	}

	tester := NewLoadTester(testServerURL, testDuration, 8)

	t.Log("Starting load test on /api/v1/test endpoint...")
	result, err := tester.ExecuteLoadTest("/api/v1/test")
	assert.NoError(t, err)

	result.PrintResults(t)

	// Assertions for API endpoint
	assert.Greater(t, result.TotalRequests, 0, "Should have made requests")
	assert.Less(t, result.ErrorRate, 10.0, "Error rate should be less than 10%")
	assert.Greater(t, result.RequestsPerSec, 5.0, "Should handle at least 5 requests/sec")
	assert.Less(t, result.AverageResponse, 2000*time.Millisecond, "Average response time should be under 2s")
}

// TestMixedEndpointsLoad tests multiple endpoints simultaneously
func TestMixedEndpointsLoad(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping load tests in short mode")
	}

	endpoints := []string{"/health", "/api/v1/test", "/status", "/ready"}
	concurrency := 3 // per endpoint
	duration := 20 * time.Second

	// Channel to collect all results
	allResults := make(chan *LoadTestResult, len(endpoints))

	// Start load tests for each endpoint concurrently
	var wg sync.WaitGroup
	for _, endpoint := range endpoints {
		wg.Add(1)
		go func(ep string) {
			defer wg.Done()

			tester := NewLoadTester(testServerURL, duration, concurrency)
			result, err := tester.ExecuteLoadTest(ep)
			assert.NoError(t, err)

			t.Logf("Results for %s:", ep)
			result.PrintResults(t)

			allResults <- result
		}(endpoint)
	}

	// Wait for all tests to complete
	go func() {
		wg.Wait()
		close(allResults)
	}()

	// Collect and analyze combined results
	totalRequests := 0
	totalSuccessful := 0
	totalFailed := 0

	for result := range allResults {
		totalRequests += result.TotalRequests
		totalSuccessful += result.SuccessfulReqs
		totalFailed += result.FailedRequests
	}

	t.Logf("=== COMBINED RESULTS ===")
	t.Logf("Total Requests Across All Endpoints: %d", totalRequests)
	t.Logf("Total Successful: %d", totalSuccessful)
	t.Logf("Total Failed: %d", totalFailed)

	if totalRequests > 0 {
		overallErrorRate := float64(totalFailed) / float64(totalRequests) * 100
		t.Logf("Overall Error Rate: %.2f%%", overallErrorRate)

		// Assertions for mixed load
		assert.Less(t, overallErrorRate, 15.0, "Overall error rate should be less than 15%")
		assert.Greater(t, totalRequests, 100, "Should have made significant number of requests")
	}
}

// TestStressTest performs a stress test with high concurrency
func TestStressTest(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping stress tests in short mode")
	}

	// Stress test parameters
	stressConcurrency := 50
	stressDuration := 15 * time.Second

	tester := NewLoadTester(testServerURL, stressDuration, stressConcurrency)

	t.Log("Starting stress test on /health endpoint...")
	result, err := tester.ExecuteLoadTest("/health")
	assert.NoError(t, err)

	result.PrintResults(t)

	// More lenient assertions for stress test
	assert.Greater(t, result.TotalRequests, 0, "Should have made requests")
	assert.Less(t, result.ErrorRate, 25.0, "Error rate should be less than 25% under stress")

	// Log if performance is concerning
	if result.AverageResponse > 5*time.Second {
		t.Logf("WARNING: High average response time under stress: %v", result.AverageResponse)
	}

	if result.RequestsPerSec < 5 {
		t.Logf("WARNING: Low throughput under stress: %.2f req/sec", result.RequestsPerSec)
	}
}

// TestGradualLoadIncrease tests service behavior as load gradually increases
func TestGradualLoadIncrease(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping gradual load tests in short mode")
	}

	concurrencyLevels := []int{1, 5, 10, 20, 30}
	testDur := 10 * time.Second

	for _, concurrency := range concurrencyLevels {
		t.Run(fmt.Sprintf("Concurrency_%d", concurrency), func(t *testing.T) {
			tester := NewLoadTester(testServerURL, testDur, concurrency)

			result, err := tester.ExecuteLoadTest("/health")
			assert.NoError(t, err)

			t.Logf("Concurrency %d - Req/sec: %.2f, Avg Response: %v, Error Rate: %.2f%%",
				concurrency, result.RequestsPerSec, result.AverageResponse, result.ErrorRate)

			// Basic assertions
			assert.Greater(t, result.TotalRequests, 0)
			assert.Less(t, result.ErrorRate, 20.0) // Allow higher error rate at high concurrency
		})
	}
}

// BenchmarkHealthEndpoint benchmarks the health endpoint
func BenchmarkHealthEndpoint(b *testing.B) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	url := testServerURL + "/health"

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			resp, err := client.Get(url)
			if err != nil {
				b.Errorf("Request failed: %v", err)
				continue
			}
			resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				b.Errorf("Expected status 200, got %d", resp.StatusCode)
			}
		}
	})
}

// BenchmarkAPIEndpoint benchmarks the API test endpoint
func BenchmarkAPIEndpoint(b *testing.B) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	url := testServerURL + "/api/v1/test"

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			resp, err := client.Get(url)
			if err != nil {
				b.Errorf("Request failed: %v", err)
				continue
			}
			resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				b.Errorf("Expected status 200, got %d", resp.StatusCode)
			}
		}
	})
}
