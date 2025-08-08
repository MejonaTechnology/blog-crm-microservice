package analytics

import (
	"fmt"
	"math"
	"sort"
	"time"
)

// TrendAnalyzer analyzes trends and patterns in blog analytics data
type TrendAnalyzer struct{}

// NewTrendAnalyzer creates a new trend analyzer
func NewTrendAnalyzer() *TrendAnalyzer {
	return &TrendAnalyzer{}
}

// AnalyzeTrends performs comprehensive trend analysis on time series data
func (ta *TrendAnalyzer) AnalyzeTrends(data []TrendDataPoint) TrendAnalysis {
	if len(data) < 2 {
		return TrendAnalysis{
			Status: "insufficient_data",
			Points: len(data),
		}
	}

	// Sort data by date
	sort.Slice(data, func(i, j int) bool {
		return data[i].Date.Before(data[j].Date)
	})

	analysis := TrendAnalysis{
		Points:     len(data),
		StartDate:  data[0].Date,
		EndDate:    data[len(data)-1].Date,
		StartValue: data[0].Value,
		EndValue:   data[len(data)-1].Value,
		DataPoints: data,
	}

	// Calculate basic statistics
	analysis.MinValue, analysis.MaxValue, analysis.AverageValue = ta.calculateBasicStats(data)

	// Calculate linear regression for overall trend
	analysis.LinearRegression = ta.calculateLinearRegression(data)

	// Determine trend direction and strength
	analysis.TrendDirection = ta.determineTrendDirection(analysis.LinearRegression.Slope)
	analysis.TrendStrength = ta.calculateTrendStrength(analysis.LinearRegression.RSquared)

	// Calculate growth metrics
	analysis.TotalGrowth = ta.calculateTotalGrowth(analysis.StartValue, analysis.EndValue)
	analysis.AverageGrowthRate = ta.calculateAverageGrowthRate(data)

	// Detect seasonality patterns
	analysis.SeasonalityAnalysis = ta.detectSeasonality(data)

	// Detect anomalies
	analysis.Anomalies = ta.detectAnomalies(data)

	// Calculate volatility
	analysis.Volatility = ta.calculateVolatility(data)

	// Generate forecasting data
	analysis.Forecast = ta.generateForecast(data, 30) // 30 days forecast

	// Provide insights and recommendations
	analysis.Insights = ta.generateInsights(analysis)

	return analysis
}

// calculateBasicStats calculates min, max, and average values
func (ta *TrendAnalyzer) calculateBasicStats(data []TrendDataPoint) (min, max, avg float64) {
	if len(data) == 0 {
		return
	}

	min = data[0].Value
	max = data[0].Value
	sum := 0.0

	for _, point := range data {
		if point.Value < min {
			min = point.Value
		}
		if point.Value > max {
			max = point.Value
		}
		sum += point.Value
	}

	avg = sum / float64(len(data))
	return
}

// calculateLinearRegression calculates linear regression for trend line
func (ta *TrendAnalyzer) calculateLinearRegression(data []TrendDataPoint) LinearRegression {
	n := float64(len(data))
	if n < 2 {
		return LinearRegression{}
	}

	// Convert dates to numeric values (days since start)
	baseTime := data[0].Date
	sumX, sumY, sumXY, sumX2, sumY2 := 0.0, 0.0, 0.0, 0.0, 0.0

	for _, point := range data {
		x := float64(point.Date.Sub(baseTime).Hours() / 24) // Days since start
		y := point.Value

		sumX += x
		sumY += y
		sumXY += x * y
		sumX2 += x * x
		sumY2 += y * y
	}

	// Calculate slope and intercept
	slope := (n*sumXY - sumX*sumY) / (n*sumX2 - sumX*sumX)
	intercept := (sumY - slope*sumX) / n

	// Calculate R-squared (coefficient of determination)
	meanY := sumY / n
	ssRes := 0.0 // Sum of squares of residuals
	ssTot := 0.0 // Total sum of squares

	for _, point := range data {
		x := float64(point.Date.Sub(baseTime).Hours() / 24)
		predicted := slope*x + intercept
		residual := point.Value - predicted

		ssRes += residual * residual
		ssTot += (point.Value - meanY) * (point.Value - meanY)
	}

	rSquared := 1.0 - (ssRes / ssTot)
	if ssTot == 0 {
		rSquared = 0
	}

	return LinearRegression{
		Slope:     slope,
		Intercept: intercept,
		RSquared:  rSquared,
	}
}

// determineTrendDirection determines the overall trend direction
func (ta *TrendAnalyzer) determineTrendDirection(slope float64) string {
	if slope > 0.1 {
		return "strongly_increasing"
	} else if slope > 0.01 {
		return "increasing"
	} else if slope < -0.1 {
		return "strongly_decreasing"
	} else if slope < -0.01 {
		return "decreasing"
	} else {
		return "stable"
	}
}

// calculateTrendStrength calculates trend strength based on R-squared
func (ta *TrendAnalyzer) calculateTrendStrength(rSquared float64) string {
	if rSquared >= 0.9 {
		return "very_strong"
	} else if rSquared >= 0.7 {
		return "strong"
	} else if rSquared >= 0.5 {
		return "moderate"
	} else if rSquared >= 0.3 {
		return "weak"
	} else {
		return "very_weak"
	}
}

// calculateTotalGrowth calculates total growth percentage
func (ta *TrendAnalyzer) calculateTotalGrowth(startValue, endValue float64) float64 {
	if startValue == 0 {
		return 0
	}
	return ((endValue - startValue) / startValue) * 100
}

// calculateAverageGrowthRate calculates average period-over-period growth rate
func (ta *TrendAnalyzer) calculateAverageGrowthRate(data []TrendDataPoint) float64 {
	if len(data) < 2 {
		return 0
	}

	totalGrowthRate := 0.0
	validPeriods := 0

	for i := 1; i < len(data); i++ {
		current := data[i].Value
		previous := data[i-1].Value

		if previous != 0 {
			growthRate := ((current - previous) / previous) * 100
			totalGrowthRate += growthRate
			validPeriods++
		}
	}

	if validPeriods == 0 {
		return 0
	}

	return totalGrowthRate / float64(validPeriods)
}

// detectSeasonality analyzes seasonal patterns in the data
func (ta *TrendAnalyzer) detectSeasonality(data []TrendDataPoint) SeasonalityAnalysis {
	analysis := SeasonalityAnalysis{}

	if len(data) < 7 {
		analysis.HasSeasonality = false
		return analysis
	}

	// Group data by day of week
	dayOfWeekData := make(map[time.Weekday][]float64)
	for _, point := range data {
		weekday := point.Date.Weekday()
		dayOfWeekData[weekday] = append(dayOfWeekData[weekday], point.Value)
	}

	// Calculate averages for each day of week
	dayAverages := make(map[time.Weekday]float64)
	overallAverage := 0.0
	totalPoints := 0

	for day, values := range dayOfWeekData {
		sum := 0.0
		for _, value := range values {
			sum += value
			totalPoints++
		}
		if len(values) > 0 {
			dayAverages[day] = sum / float64(len(values))
			overallAverage += sum
		}
	}

	if totalPoints > 0 {
		overallAverage /= float64(totalPoints)
	}

	// Check for significant variations (threshold: 20% deviation)
	significantVariations := 0
	for _, avg := range dayAverages {
		deviation := math.Abs((avg - overallAverage) / overallAverage)
		if deviation > 0.2 { // 20% threshold
			significantVariations++
		}
	}

	analysis.HasSeasonality = significantVariations >= 2
	analysis.DayOfWeekPattern = dayAverages

	// Simple monthly seasonality check if data spans multiple months
	if len(data) > 30 {
		monthlyData := make(map[int][]float64)
		for _, point := range data {
			month := int(point.Date.Month())
			monthlyData[month] = append(monthlyData[month], point.Value)
		}

		monthAverages := make(map[int]float64)
		for month, values := range monthlyData {
			if len(values) > 0 {
				sum := 0.0
				for _, value := range values {
					sum += value
				}
				monthAverages[month] = sum / float64(len(values))
			}
		}

		analysis.MonthlyPattern = monthAverages
	}

	return analysis
}

// detectAnomalies identifies unusual data points
func (ta *TrendAnalyzer) detectAnomalies(data []TrendDataPoint) []Anomaly {
	if len(data) < 3 {
		return nil
	}

	// Calculate moving average and standard deviation
	windowSize := 7 // 7-day window
	if len(data) < windowSize {
		windowSize = len(data) / 2
	}

	var anomalies []Anomaly

	for i := windowSize; i < len(data)-windowSize; i++ {
		// Calculate moving average and standard deviation for window
		sum := 0.0
		count := 0
		for j := i - windowSize/2; j <= i+windowSize/2; j++ {
			if j != i { // Exclude current point
				sum += data[j].Value
				count++
			}
		}

		if count == 0 {
			continue
		}

		mean := sum / float64(count)

		// Calculate standard deviation
		sumSquaredDiff := 0.0
		for j := i - windowSize/2; j <= i+windowSize/2; j++ {
			if j != i {
				diff := data[j].Value - mean
				sumSquaredDiff += diff * diff
			}
		}

		stdDev := math.Sqrt(sumSquaredDiff / float64(count-1))

		// Check if current point is an anomaly (beyond 2 standard deviations)
		currentValue := data[i].Value
		zScore := math.Abs(currentValue - mean)
		if stdDev > 0 {
			zScore = zScore / stdDev
		}

		if zScore > 2.0 { // 2 sigma threshold
			anomalies = append(anomalies, Anomaly{
				Date:      data[i].Date,
				Value:     currentValue,
				Expected:  mean,
				Deviation: currentValue - mean,
				ZScore:    zScore,
				Type:      ta.classifyAnomalyType(currentValue, mean),
			})
		}
	}

	return anomalies
}

// classifyAnomalyType classifies the type of anomaly
func (ta *TrendAnalyzer) classifyAnomalyType(actual, expected float64) string {
	if actual > expected*1.5 {
		return "spike"
	} else if actual < expected*0.5 {
		return "dip"
	} else if actual > expected {
		return "positive_outlier"
	} else {
		return "negative_outlier"
	}
}

// calculateVolatility calculates the volatility of the data series
func (ta *TrendAnalyzer) calculateVolatility(data []TrendDataPoint) float64 {
	if len(data) < 2 {
		return 0
	}

	// Calculate period-over-period percentage changes
	changes := make([]float64, 0, len(data)-1)
	for i := 1; i < len(data); i++ {
		if data[i-1].Value != 0 {
			change := ((data[i].Value - data[i-1].Value) / data[i-1].Value) * 100
			changes = append(changes, change)
		}
	}

	if len(changes) == 0 {
		return 0
	}

	// Calculate standard deviation of changes
	mean := 0.0
	for _, change := range changes {
		mean += change
	}
	mean /= float64(len(changes))

	variance := 0.0
	for _, change := range changes {
		diff := change - mean
		variance += diff * diff
	}
	variance /= float64(len(changes))

	return math.Sqrt(variance)
}

// generateForecast generates simple forecast using linear regression
func (ta *TrendAnalyzer) generateForecast(data []TrendDataPoint, days int) []ForecastPoint {
	if len(data) < 2 {
		return nil
	}

	regression := ta.calculateLinearRegression(data)
	baseTime := data[0].Date
	lastDate := data[len(data)-1].Date

	forecast := make([]ForecastPoint, days)

	for i := 0; i < days; i++ {
		forecastDate := lastDate.AddDate(0, 0, i+1)
		x := float64(forecastDate.Sub(baseTime).Hours() / 24)

		predictedValue := regression.Slope*x + regression.Intercept

		// Add confidence interval based on historical volatility
		volatility := ta.calculateVolatility(data)
		confidenceInterval := volatility * 1.96 // 95% confidence interval

		forecast[i] = ForecastPoint{
			Date:               forecastDate,
			PredictedValue:     predictedValue,
			ConfidenceInterval: confidenceInterval,
			LowerBound:         predictedValue - confidenceInterval,
			UpperBound:         predictedValue + confidenceInterval,
		}
	}

	return forecast
}

// generateInsights generates actionable insights based on trend analysis
func (ta *TrendAnalyzer) generateInsights(analysis TrendAnalysis) []string {
	var insights []string

	// Trend direction insights
	switch analysis.TrendDirection {
	case "strongly_increasing":
		insights = append(insights, "Strong upward trend detected. Content performance is improving significantly.")
	case "increasing":
		insights = append(insights, "Positive trend observed. Content performance is gradually improving.")
	case "strongly_decreasing":
		insights = append(insights, "Strong downward trend detected. Immediate attention needed to reverse declining performance.")
	case "decreasing":
		insights = append(insights, "Negative trend observed. Consider reviewing content strategy.")
	case "stable":
		insights = append(insights, "Performance is stable. Look for optimization opportunities.")
	}

	// Volatility insights
	if analysis.Volatility > 50 {
		insights = append(insights, "High volatility detected. Performance is inconsistent - consider standardizing content quality.")
	} else if analysis.Volatility < 10 {
		insights = append(insights, "Low volatility suggests consistent performance. Good content strategy execution.")
	}

	// Seasonality insights
	if analysis.SeasonalityAnalysis.HasSeasonality {
		insights = append(insights, "Seasonal patterns detected. Consider timing content releases based on historical performance.")
	}

	// Anomaly insights
	if len(analysis.Anomalies) > 0 {
		insights = append(insights, fmt.Sprintf("Found %d anomalies in the data. Investigate unusual spikes or dips for insights.", len(analysis.Anomalies)))
	}

	// Growth insights
	if analysis.TotalGrowth > 50 {
		insights = append(insights, "Excellent growth performance! Current strategy is highly effective.")
	} else if analysis.TotalGrowth > 20 {
		insights = append(insights, "Good growth performance. Consider scaling successful tactics.")
	} else if analysis.TotalGrowth < -20 {
		insights = append(insights, "Declining performance requires immediate strategic review and optimization.")
	}

	return insights
}

// Data structures for trend analysis

type TrendAnalysis struct {
	Points              int                 `json:"points"`
	Status              string              `json:"status"`
	StartDate           time.Time           `json:"start_date"`
	EndDate             time.Time           `json:"end_date"`
	StartValue          float64             `json:"start_value"`
	EndValue            float64             `json:"end_value"`
	MinValue            float64             `json:"min_value"`
	MaxValue            float64             `json:"max_value"`
	AverageValue        float64             `json:"average_value"`
	LinearRegression    LinearRegression    `json:"linear_regression"`
	TrendDirection      string              `json:"trend_direction"`
	TrendStrength       string              `json:"trend_strength"`
	TotalGrowth         float64             `json:"total_growth"`
	AverageGrowthRate   float64             `json:"average_growth_rate"`
	SeasonalityAnalysis SeasonalityAnalysis `json:"seasonality_analysis"`
	Anomalies           []Anomaly           `json:"anomalies"`
	Volatility          float64             `json:"volatility"`
	Forecast            []ForecastPoint     `json:"forecast"`
	Insights            []string            `json:"insights"`
	DataPoints          []TrendDataPoint    `json:"data_points"`
}

type LinearRegression struct {
	Slope     float64 `json:"slope"`
	Intercept float64 `json:"intercept"`
	RSquared  float64 `json:"r_squared"`
}

type SeasonalityAnalysis struct {
	HasSeasonality   bool                     `json:"has_seasonality"`
	DayOfWeekPattern map[time.Weekday]float64 `json:"day_of_week_pattern"`
	MonthlyPattern   map[int]float64          `json:"monthly_pattern"`
}

type Anomaly struct {
	Date      time.Time `json:"date"`
	Value     float64   `json:"value"`
	Expected  float64   `json:"expected"`
	Deviation float64   `json:"deviation"`
	ZScore    float64   `json:"z_score"`
	Type      string    `json:"type"`
}

type ForecastPoint struct {
	Date               time.Time `json:"date"`
	PredictedValue     float64   `json:"predicted_value"`
	ConfidenceInterval float64   `json:"confidence_interval"`
	LowerBound         float64   `json:"lower_bound"`
	UpperBound         float64   `json:"upper_bound"`
}
