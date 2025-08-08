package analytics

import (
	"math"
	"time"
)

// ROICalculator handles return on investment calculations for blog content
type ROICalculator struct{}

// NewROICalculator creates a new ROI calculator
func NewROICalculator() *ROICalculator {
	return &ROICalculator{}
}

// CalculateContentROI calculates comprehensive ROI for blog content
func (rc *ROICalculator) CalculateContentROI(metrics ContentROIMetrics) ContentROIResult {
	result := ContentROIResult{
		ContentID:    metrics.ContentID,
		Title:        metrics.Title,
		PublishedAt:  metrics.PublishedAt,
		Period:       metrics.Period,
	}

	// Calculate total investment
	result.TotalInvestment = rc.calculateTotalInvestment(metrics.Investment)

	// Calculate direct revenue
	result.DirectRevenue = rc.calculateDirectRevenue(metrics.DirectConversions)

	// Calculate indirect revenue (attributed)
	result.IndirectRevenue = rc.calculateIndirectRevenue(metrics.AttributedConversions, metrics.AttributionModel)

	// Calculate total revenue
	result.TotalRevenue = result.DirectRevenue + result.IndirectRevenue

	// Calculate ROI percentage
	if result.TotalInvestment > 0 {
		result.ROIPercentage = ((result.TotalRevenue - result.TotalInvestment) / result.TotalInvestment) * 100
	}

	// Calculate payback period
	result.PaybackPeriod = rc.calculatePaybackPeriod(result.TotalInvestment, result.TotalRevenue, metrics.Period)

	// Calculate customer lifetime value impact
	result.CLVImpact = rc.calculateCLVImpact(metrics.NewCustomers, metrics.AverageCLV)

	// Calculate lead value
	result.LeadValue = rc.calculateLeadValue(metrics.Leads, result.TotalRevenue)

	// Calculate cost metrics
	result.CostPerLead = rc.calculateCostPerLead(result.TotalInvestment, metrics.Leads)
	result.CostPerAcquisition = rc.calculateCostPerAcquisition(result.TotalInvestment, metrics.DirectConversions)

	// Calculate engagement value
	result.EngagementValue = rc.calculateEngagementValue(metrics.Engagement, metrics.EngagementValue)

	// Calculate brand value impact
	result.BrandValue = rc.calculateBrandValue(metrics.BrandMetrics)

	return result
}

// calculateTotalInvestment calculates total investment in content creation and promotion
func (rc *ROICalculator) calculateTotalInvestment(investment ContentInvestment) float64 {
	total := investment.CreationCost + investment.PromotionCost + investment.ToolsCost

	// Add time-based costs
	total += investment.TimeInvested * investment.HourlyRate

	// Add opportunity cost
	total += investment.OpportunityCost

	return total
}

// calculateDirectRevenue calculates revenue directly attributed to the content
func (rc *ROICalculator) calculateDirectRevenue(conversions []DirectConversion) float64 {
	var total float64
	for _, conversion := range conversions {
		total += conversion.Revenue
	}
	return total
}

// calculateIndirectRevenue calculates attributed revenue based on attribution model
func (rc *ROICalculator) calculateIndirectRevenue(conversions []AttributedConversion, model string) float64 {
	var total float64

	for _, conversion := range conversions {
		switch model {
		case "first_touch":
			if conversion.IsFirstTouch {
				total += conversion.Revenue
			}
		case "last_touch":
			if conversion.IsLastTouch {
				total += conversion.Revenue
			}
		case "linear":
			total += conversion.Revenue * conversion.AttributionWeight
		case "time_decay":
			total += conversion.Revenue * rc.calculateTimeDecayWeight(conversion.DaysFromTouch)
		case "position_based":
			total += conversion.Revenue * rc.calculatePositionBasedWeight(conversion.TouchPosition, conversion.TotalTouches)
		default:
			total += conversion.Revenue * conversion.AttributionWeight
		}
	}

	return total
}

// calculateTimeDecayWeight calculates time decay attribution weight
func (rc *ROICalculator) calculateTimeDecayWeight(daysFromTouch int) float64 {
	// Exponential decay: more recent touches get more credit
	decayRate := 0.1 // Adjust as needed
	return math.Exp(-decayRate * float64(daysFromTouch))
}

// calculatePositionBasedWeight calculates position-based attribution weight
func (rc *ROICalculator) calculatePositionBasedWeight(position, total int) float64 {
	if total == 1 {
		return 1.0
	}

	// 40% for first touch, 20% for last touch, 40% distributed evenly among middle touches
	if position == 1 {
		return 0.4
	} else if position == total {
		return 0.2
	} else {
		middleTouches := total - 2
		if middleTouches > 0 {
			return 0.4 / float64(middleTouches)
		}
	}
	return 0.0
}

// calculatePaybackPeriod calculates how long it takes to recover the investment
func (rc *ROICalculator) calculatePaybackPeriod(investment, revenue float64, periodDays int) float64 {
	if revenue <= 0 {
		return -1 // No payback
	}

	dailyRevenue := revenue / float64(periodDays)
	if dailyRevenue <= 0 {
		return -1
	}

	return investment / dailyRevenue
}

// calculateCLVImpact calculates customer lifetime value impact
func (rc *ROICalculator) calculateCLVImpact(newCustomers int, averageCLV float64) float64 {
	return float64(newCustomers) * averageCLV
}

// calculateLeadValue calculates average value per lead
func (rc *ROICalculator) calculateLeadValue(leads int, totalRevenue float64) float64 {
	if leads == 0 {
		return 0
	}
	return totalRevenue / float64(leads)
}

// calculateCostPerLead calculates cost per lead generated
func (rc *ROICalculator) calculateCostPerLead(investment float64, leads int) float64 {
	if leads == 0 {
		return 0
	}
	return investment / float64(leads)
}

// calculateCostPerAcquisition calculates cost per customer acquisition
func (rc *ROICalculator) calculateCostPerAcquisition(investment float64, conversions []DirectConversion) float64 {
	if len(conversions) == 0 {
		return 0
	}
	return investment / float64(len(conversions))
}

// calculateEngagementValue calculates monetary value of engagement
func (rc *ROICalculator) calculateEngagementValue(engagement EngagementMetrics, valueMetrics EngagementValueMetrics) float64 {
	var value float64

	// Page view value
	value += float64(engagement.PageViews) * valueMetrics.PageViewValue

	// Social share value
	value += float64(engagement.SocialShares) * valueMetrics.ShareValue

	// Comment value
	value += float64(engagement.Comments) * valueMetrics.CommentValue

	// Download value
	value += float64(engagement.Downloads) * valueMetrics.DownloadValue

	// Newsletter signup value
	value += float64(engagement.NewsletterSignups) * valueMetrics.SignupValue

	return value
}

// calculateBrandValue calculates brand value impact
func (rc *ROICalculator) calculateBrandValue(metrics BrandMetrics) float64 {
	var value float64

	// Brand awareness value
	value += metrics.BrandMentions * metrics.MentionValue

	// Domain authority improvement value
	value += metrics.BacklinkValue

	// Search visibility value
	value += metrics.SearchVisibilityValue

	// Thought leadership value
	value += metrics.ThoughtLeadershipValue

	return value
}

// CalculateContentPortfolioROI calculates ROI for a portfolio of content
func (rc *ROICalculator) CalculateContentPortfolioROI(portfolioMetrics []ContentROIMetrics) PortfolioROIResult {
	result := PortfolioROIResult{
		ContentCount: len(portfolioMetrics),
	}

	var totalInvestment, totalRevenue float64
	var allLeads, allConversions int

	for _, metrics := range portfolioMetrics {
		contentROI := rc.CalculateContentROI(metrics)
		
		totalInvestment += contentROI.TotalInvestment
		totalRevenue += contentROI.TotalRevenue
		allLeads += metrics.Leads
		allConversions += len(metrics.DirectConversions)

		result.ContentResults = append(result.ContentResults, contentROI)
	}

	result.TotalInvestment = totalInvestment
	result.TotalRevenue = totalRevenue
	result.TotalLeads = allLeads
	result.TotalConversions = allConversions

	// Calculate portfolio-level metrics
	if totalInvestment > 0 {
		result.PortfolioROI = ((totalRevenue - totalInvestment) / totalInvestment) * 100
		result.AverageCostPerLead = totalInvestment / float64(allLeads)
		result.AverageCostPerConversion = totalInvestment / float64(allConversions)
	}

	if len(portfolioMetrics) > 0 {
		result.AverageROI = result.PortfolioROI / float64(len(portfolioMetrics))
	}

	// Find best and worst performing content
	bestROI := -math.Inf(1)
	worstROI := math.Inf(1)

	for _, contentResult := range result.ContentResults {
		if contentResult.ROIPercentage > bestROI {
			bestROI = contentResult.ROIPercentage
			result.BestPerformingContent = contentResult
		}
		if contentResult.ROIPercentage < worstROI {
			worstROI = contentResult.ROIPercentage
			result.WorstPerformingContent = contentResult
		}
	}

	return result
}

// CalculateROITrends calculates ROI trends over time periods
func (rc *ROICalculator) CalculateROITrends(historicalData []PeriodROIData) ROITrendAnalysis {
	analysis := ROITrendAnalysis{
		Periods: len(historicalData),
	}

	if len(historicalData) < 2 {
		return analysis
	}

	// Calculate period-over-period changes
	for i := 1; i < len(historicalData); i++ {
		current := historicalData[i]
		previous := historicalData[i-1]

		change := ROIChange{
			Period:           current.Period,
			CurrentROI:       current.ROI,
			PreviousROI:      previous.ROI,
			AbsoluteChange:   current.ROI - previous.ROI,
		}

		if previous.ROI != 0 {
			change.PercentageChange = ((current.ROI - previous.ROI) / previous.ROI) * 100
		}

		analysis.Changes = append(analysis.Changes, change)
	}

	// Calculate overall trend
	if len(analysis.Changes) > 0 {
		firstROI := historicalData[0].ROI
		lastROI := historicalData[len(historicalData)-1].ROI
		
		if firstROI != 0 {
			analysis.OverallTrend = ((lastROI - firstROI) / firstROI) * 100
		}

		// Calculate average period change
		var totalChange float64
		for _, change := range analysis.Changes {
			totalChange += change.PercentageChange
		}
		analysis.AveragePeriodChange = totalChange / float64(len(analysis.Changes))
	}

	// Determine trend direction
	if analysis.OverallTrend > 5 {
		analysis.TrendDirection = "strongly_positive"
	} else if analysis.OverallTrend > 0 {
		analysis.TrendDirection = "positive"
	} else if analysis.OverallTrend < -5 {
		analysis.TrendDirection = "strongly_negative"
	} else if analysis.OverallTrend < 0 {
		analysis.TrendDirection = "negative"
	} else {
		analysis.TrendDirection = "stable"
	}

	return analysis
}

// Data structures for ROI calculations

type ContentROIMetrics struct {
	ContentID              uint
	Title                  string
	PublishedAt            time.Time
	Period                 int // Days
	Investment             ContentInvestment
	DirectConversions      []DirectConversion
	AttributedConversions  []AttributedConversion
	AttributionModel       string
	Leads                  int
	NewCustomers           int
	AverageCLV             float64
	Engagement             EngagementMetrics
	EngagementValue        EngagementValueMetrics
	BrandMetrics           BrandMetrics
}

type ContentInvestment struct {
	CreationCost     float64
	PromotionCost    float64
	ToolsCost        float64
	TimeInvested     float64 // Hours
	HourlyRate       float64
	OpportunityCost  float64
}

type DirectConversion struct {
	CustomerID   uint
	Revenue      float64
	ConvertedAt  time.Time
	ProductType  string
}

type AttributedConversion struct {
	CustomerID        uint
	Revenue           float64
	AttributionWeight float64
	IsFirstTouch      bool
	IsLastTouch       bool
	TouchPosition     int
	TotalTouches      int
	DaysFromTouch     int
	ConvertedAt       time.Time
}

type EngagementValueMetrics struct {
	PageViewValue   float64
	ShareValue      float64
	CommentValue    float64
	DownloadValue   float64
	SignupValue     float64
}

type BrandMetrics struct {
	BrandMentions            float64
	MentionValue             float64
	BacklinkValue            float64
	SearchVisibilityValue    float64
	ThoughtLeadershipValue   float64
}

type ContentROIResult struct {
	ContentID         uint      `json:"content_id"`
	Title             string    `json:"title"`
	PublishedAt       time.Time `json:"published_at"`
	Period            int       `json:"period"`
	TotalInvestment   float64   `json:"total_investment"`
	DirectRevenue     float64   `json:"direct_revenue"`
	IndirectRevenue   float64   `json:"indirect_revenue"`
	TotalRevenue      float64   `json:"total_revenue"`
	ROIPercentage     float64   `json:"roi_percentage"`
	PaybackPeriod     float64   `json:"payback_period"` // Days
	CLVImpact         float64   `json:"clv_impact"`
	LeadValue         float64   `json:"lead_value"`
	CostPerLead       float64   `json:"cost_per_lead"`
	CostPerAcquisition float64  `json:"cost_per_acquisition"`
	EngagementValue   float64   `json:"engagement_value"`
	BrandValue        float64   `json:"brand_value"`
}

type PortfolioROIResult struct {
	ContentCount             int                `json:"content_count"`
	TotalInvestment          float64            `json:"total_investment"`
	TotalRevenue             float64            `json:"total_revenue"`
	TotalLeads               int                `json:"total_leads"`
	TotalConversions         int                `json:"total_conversions"`
	PortfolioROI             float64            `json:"portfolio_roi"`
	AverageROI               float64            `json:"average_roi"`
	AverageCostPerLead       float64            `json:"average_cost_per_lead"`
	AverageCostPerConversion float64            `json:"average_cost_per_conversion"`
	BestPerformingContent    ContentROIResult   `json:"best_performing_content"`
	WorstPerformingContent   ContentROIResult   `json:"worst_performing_content"`
	ContentResults           []ContentROIResult `json:"content_results"`
}

type PeriodROIData struct {
	Period string  `json:"period"`
	ROI    float64 `json:"roi"`
}

type ROIChange struct {
	Period           string  `json:"period"`
	CurrentROI       float64 `json:"current_roi"`
	PreviousROI      float64 `json:"previous_roi"`
	AbsoluteChange   float64 `json:"absolute_change"`
	PercentageChange float64 `json:"percentage_change"`
}

type ROITrendAnalysis struct {
	Periods             int         `json:"periods"`
	Changes             []ROIChange `json:"changes"`
	OverallTrend        float64     `json:"overall_trend"`
	AveragePeriodChange float64     `json:"average_period_change"`
	TrendDirection      string      `json:"trend_direction"`
}