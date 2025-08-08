package analytics

import (
	"math"
	"time"
)

// PerformanceCalculator handles various performance calculations for blog analytics
type PerformanceCalculator struct{}

// NewPerformanceCalculator creates a new performance calculator
func NewPerformanceCalculator() *PerformanceCalculator {
	return &PerformanceCalculator{}
}

// CalculateEngagementScore calculates an engagement score based on various metrics
func (pc *PerformanceCalculator) CalculateEngagementScore(metrics EngagementMetrics) float64 {
	if metrics.PageViews == 0 {
		return 0
	}

	// Weighted scoring system (0-100 scale)
	var score float64

	// Time on page score (30% weight) - normalized to 0-100
	timeScore := math.Min(float64(metrics.AvgTimeOnPage)/300, 1.0) * 30 // 5 minutes = max score

	// Bounce rate score (25% weight) - inverted since lower is better
	bounceScore := (1.0 - math.Min(metrics.BounceRate/100.0, 1.0)) * 25

	// Scroll depth score (20% weight)
	scrollScore := math.Min(metrics.AvgScrollDepth/100.0, 1.0) * 20

	// Social shares score (15% weight) - logarithmic scaling
	var shareScore float64
	if metrics.SocialShares > 0 {
		shareScore = math.Min(math.Log10(float64(metrics.SocialShares)+1)*10, 15)
	}

	// Comments score (10% weight) - logarithmic scaling
	var commentScore float64
	if metrics.Comments > 0 {
		commentScore = math.Min(math.Log10(float64(metrics.Comments)+1)*5, 10)
	}

	score = timeScore + bounceScore + scrollScore + shareScore + commentScore

	return math.Min(score, 100.0)
}

// CalculateContentQualityScore calculates content quality based on SEO and readability metrics
func (pc *PerformanceCalculator) CalculateContentQualityScore(metrics ContentQualityMetrics) float64 {
	var score float64

	// SEO score (40% weight)
	seoScore := float64(metrics.SEOScore) * 0.4

	// Readability score (25% weight)
	readabilityScore := float64(metrics.ReadabilityScore) * 0.25

	// Word count optimization (15% weight) - optimal range 1000-3000 words
	var wordCountScore float64
	if metrics.WordCount >= 1000 && metrics.WordCount <= 3000 {
		wordCountScore = 15.0
	} else if metrics.WordCount >= 500 && metrics.WordCount < 1000 {
		wordCountScore = 10.0
	} else if metrics.WordCount > 3000 && metrics.WordCount <= 5000 {
		wordCountScore = 12.0
	} else {
		wordCountScore = 5.0
	}

	// Internal links score (10% weight)
	var internalLinksScore float64
	if metrics.InternalLinks >= 3 && metrics.InternalLinks <= 10 {
		internalLinksScore = 10.0
	} else if metrics.InternalLinks >= 1 && metrics.InternalLinks < 3 {
		internalLinksScore = 7.0
	} else if metrics.InternalLinks > 10 {
		internalLinksScore = 6.0
	}

	// External links score (10% weight)
	var externalLinksScore float64
	if metrics.ExternalLinks >= 2 && metrics.ExternalLinks <= 5 {
		externalLinksScore = 10.0
	} else if metrics.ExternalLinks >= 1 && metrics.ExternalLinks < 2 {
		externalLinksScore = 7.0
	} else if metrics.ExternalLinks > 5 {
		externalLinksScore = 8.0
	}

	score = seoScore + readabilityScore + wordCountScore + internalLinksScore + externalLinksScore

	return math.Min(score, 100.0)
}

// CalculateViralityScore calculates how viral a piece of content is
func (pc *PerformanceCalculator) CalculateViralityScore(metrics ViralityMetrics) float64 {
	if metrics.PageViews == 0 {
		return 0
	}

	var score float64

	// Share rate (shares per view) - 40% weight
	shareRate := float64(metrics.SocialShares) / float64(metrics.PageViews)
	shareScore := math.Min(shareRate*1000, 1.0) * 40 // Normalize to 40 max

	// Growth velocity (30% weight) - rate of view increase
	var velocityScore float64
	if metrics.GrowthVelocity > 0 {
		velocityScore = math.Min(math.Log10(metrics.GrowthVelocity+1)*10, 30)
	}

	// Engagement velocity (20% weight) - rate of engagement increase
	var engagementVelocityScore float64
	if metrics.EngagementVelocity > 0 {
		engagementVelocityScore = math.Min(math.Log10(metrics.EngagementVelocity+1)*7, 20)
	}

	// Cross-platform sharing (10% weight)
	platformScore := math.Min(float64(metrics.PlatformReach)*2, 10)

	score = shareScore + velocityScore + engagementVelocityScore + platformScore

	return math.Min(score, 100.0)
}

// CalculateROI calculates return on investment for content
func (pc *PerformanceCalculator) CalculateROI(revenue, cost float64) float64 {
	if cost == 0 {
		return 0
	}
	return ((revenue - cost) / cost) * 100
}

// CalculateConversionRate calculates conversion rate for different actions
func (pc *PerformanceCalculator) CalculateConversionRate(conversions, totalVisitors int) float64 {
	if totalVisitors == 0 {
		return 0
	}
	return (float64(conversions) / float64(totalVisitors)) * 100
}

// CalculateGrowthRate calculates period-over-period growth rate
func (pc *PerformanceCalculator) CalculateGrowthRate(current, previous float64) float64 {
	if previous == 0 {
		return 100.0 // Infinite growth from zero
	}
	return ((current - previous) / previous) * 100
}

// CalculateTrendScore calculates trend score based on recent performance
func (pc *PerformanceCalculator) CalculateTrendScore(dataPoints []TrendDataPoint) float64 {
	if len(dataPoints) < 2 {
		return 50.0 // Neutral score for insufficient data
	}

	// Calculate linear regression slope for trend direction
	n := float64(len(dataPoints))
	var sumX, sumY, sumXY, sumX2 float64

	for i, point := range dataPoints {
		x := float64(i)
		y := point.Value
		sumX += x
		sumY += y
		sumXY += x * y
		sumX2 += x * x
	}

	// Linear regression slope
	slope := (n*sumXY - sumX*sumY) / (n*sumX2 - sumX*sumX)

	// Normalize slope to 0-100 score
	// Positive slope = upward trend (>50), negative = downward trend (<50)
	trendScore := 50.0 + (slope * 10) // Adjust multiplier as needed

	return math.Max(0, math.Min(100, trendScore))
}

// CalculateSeasonalityIndex calculates seasonality index for content performance
func (pc *PerformanceCalculator) CalculateSeasonalityIndex(monthlyData []float64) map[int]float64 {
	seasonalityIndex := make(map[int]float64)

	if len(monthlyData) == 0 {
		return seasonalityIndex
	}

	// Calculate average
	var total float64
	for _, value := range monthlyData {
		total += value
	}
	average := total / float64(len(monthlyData))

	// Calculate seasonality index for each month
	for month, value := range monthlyData {
		if average > 0 {
			seasonalityIndex[month+1] = (value / average) * 100
		} else {
			seasonalityIndex[month+1] = 100.0
		}
	}

	return seasonalityIndex
}

// CalculateCompetitiveScore calculates competitive positioning score
func (pc *PerformanceCalculator) CalculateCompetitiveScore(ownMetrics, competitorMetrics CompetitiveMetrics) float64 {
	var score float64

	// Market share score (30% weight)
	totalMarketShare := ownMetrics.MarketShare + competitorMetrics.MarketShare
	if totalMarketShare > 0 {
		marketShareScore := (ownMetrics.MarketShare / totalMarketShare) * 30
		score += marketShareScore
	}

	// Content volume score (25% weight)
	totalContentVolume := ownMetrics.ContentVolume + competitorMetrics.ContentVolume
	if totalContentVolume > 0 {
		contentVolumeScore := (float64(ownMetrics.ContentVolume) / float64(totalContentVolume)) * 25
		score += contentVolumeScore
	}

	// Engagement quality score (25% weight)
	if competitorMetrics.AvgEngagement > 0 {
		engagementRatio := ownMetrics.AvgEngagement / competitorMetrics.AvgEngagement
		engagementScore := math.Min(engagementRatio, 2.0) * 12.5 // Cap at 25 points
		score += engagementScore
	} else {
		score += 25.0 // Full points if no competitor engagement
	}

	// Innovation score (20% weight) - based on content uniqueness and freshness
	innovationScore := float64(ownMetrics.InnovationScore) * 0.2
	score += innovationScore

	return math.Min(score, 100.0)
}

// Data structures for calculations

type EngagementMetrics struct {
	PageViews      int
	AvgTimeOnPage  int     // seconds
	BounceRate     float64 // percentage
	AvgScrollDepth float64 // percentage
	SocialShares   int
	Comments       int
}

type ContentQualityMetrics struct {
	SEOScore         int
	ReadabilityScore int
	WordCount        int
	InternalLinks    int
	ExternalLinks    int
	ImageCount       int
	VideoCount       int
}

type ViralityMetrics struct {
	PageViews          int
	SocialShares       int
	GrowthVelocity     float64 // views per hour/day
	EngagementVelocity float64 // engagements per hour/day
	PlatformReach      int     // number of platforms shared on
}

type TrendDataPoint struct {
	Date  time.Time
	Value float64
}

type CompetitiveMetrics struct {
	MarketShare     float64 // percentage
	ContentVolume   int     // number of posts
	AvgEngagement   float64 // average engagement per post
	InnovationScore int     // 0-100 score for content innovation
}
