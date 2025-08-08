package analytics

import (
	"math"
	"strings"
	"time"
)

// LeadScorer handles lead scoring calculations for blog-generated leads
type LeadScorer struct{}

// NewLeadScorer creates a new lead scorer
func NewLeadScorer() *LeadScorer {
	return &LeadScorer{}
}

// CalculateLeadScore calculates a comprehensive lead score (0-100)
func (ls *LeadScorer) CalculateLeadScore(profile LeadProfile) int {
	var totalScore float64

	// Demographic scoring (25% weight)
	demographicScore := ls.calculateDemographicScore(profile.Demographics)
	totalScore += demographicScore * 0.25

	// Behavioral scoring (35% weight)
	behavioralScore := ls.calculateBehavioralScore(profile.Behavior)
	totalScore += behavioralScore * 0.35

	// Firmographic scoring (25% weight) - for B2B leads
	firmographicScore := ls.calculateFirmographicScore(profile.Company)
	totalScore += firmographicScore * 0.25

	// Intent scoring (15% weight)
	intentScore := ls.calculateIntentScore(profile.Intent)
	totalScore += intentScore * 0.15

	return int(math.Min(totalScore, 100))
}

// calculateDemographicScore scores based on demographic information
func (ls *LeadScorer) calculateDemographicScore(demo Demographics) float64 {
	var score float64

	// Job title scoring (40% of demographic score)
	score += ls.scoreJobTitle(demo.JobTitle) * 0.4

	// Industry scoring (30% of demographic score)
	score += ls.scoreIndustry(demo.Industry) * 0.3

	// Location scoring (20% of demographic score)
	score += ls.scoreLocation(demo.Location) * 0.2

	// Experience level scoring (10% of demographic score)
	score += ls.scoreExperienceLevel(demo.ExperienceLevel) * 0.1

	return score
}

// calculateBehavioralScore scores based on user behavior
func (ls *LeadScorer) calculateBehavioralScore(behavior Behavior) float64 {
	var score float64

	// Engagement level (30% of behavioral score)
	score += ls.scoreEngagementLevel(behavior) * 0.3

	// Content consumption (25% of behavioral score)
	score += ls.scoreContentConsumption(behavior) * 0.25

	// Website activity (25% of behavioral score)
	score += ls.scoreWebsiteActivity(behavior) * 0.25

	// Recency (20% of behavioral score)
	score += ls.scoreRecency(behavior.LastActivity) * 0.2

	return score
}

// calculateFirmographicScore scores based on company information
func (ls *LeadScorer) calculateFirmographicScore(company Company) float64 {
	if company.Name == "" {
		return 50.0 // Neutral score for missing company data
	}

	var score float64

	// Company size scoring (40% of firmographic score)
	score += ls.scoreCompanySize(company.Size) * 0.4

	// Industry fit scoring (30% of firmographic score)
	score += ls.scoreIndustryFit(company.Industry) * 0.3

	// Revenue scoring (20% of firmographic score)
	score += ls.scoreRevenue(company.Revenue) * 0.2

	// Technology stack scoring (10% of firmographic score)
	score += ls.scoreTechnologyStack(company.TechnologyStack) * 0.1

	return score
}

// calculateIntentScore scores based on purchase intent signals
func (ls *LeadScorer) calculateIntentScore(intent Intent) float64 {
	var score float64

	// Source type scoring (30% of intent score)
	score += ls.scoreSourceType(intent.SourceType) * 0.3

	// Content type engagement (25% of intent score)
	score += ls.scoreContentTypeEngagement(intent.ContentTypes) * 0.25

	// CTA interaction (25% of intent score)
	score += ls.scoreCTAInteraction(intent.CTAInteractions) * 0.25

	// Form completions (20% of intent score)
	score += ls.scoreFormCompletions(intent.FormCompletions) * 0.2

	return score
}

// Individual scoring methods

func (ls *LeadScorer) scoreJobTitle(title string) float64 {
	title = strings.ToLower(title)

	// High-value titles (80-100 points)
	highValueTitles := []string{"ceo", "cto", "cfo", "cmo", "vp", "vice president", "director", "head of", "chief"}
	for _, hvt := range highValueTitles {
		if strings.Contains(title, hvt) {
			return 90.0
		}
	}

	// Medium-value titles (60-79 points)
	mediumValueTitles := []string{"manager", "lead", "senior", "principal", "architect", "consultant"}
	for _, mvt := range mediumValueTitles {
		if strings.Contains(title, mvt) {
			return 70.0
		}
	}

	// Entry-level titles (40-59 points)
	entryTitles := []string{"developer", "engineer", "analyst", "specialist", "coordinator", "associate"}
	for _, et := range entryTitles {
		if strings.Contains(title, et) {
			return 50.0
		}
	}

	return 30.0 // Unknown or low-value title
}

func (ls *LeadScorer) scoreIndustry(industry string) float64 {
	industry = strings.ToLower(industry)

	// High-fit industries (technology services company)
	highFitIndustries := []string{"technology", "software", "saas", "fintech", "healthtech", "edtech", "startup"}
	for _, hfi := range highFitIndustries {
		if strings.Contains(industry, hfi) {
			return 90.0
		}
	}

	// Medium-fit industries
	mediumFitIndustries := []string{"finance", "healthcare", "education", "retail", "ecommerce", "manufacturing"}
	for _, mfi := range mediumFitIndustries {
		if strings.Contains(industry, mfi) {
			return 70.0
		}
	}

	return 50.0 // Other industries
}

func (ls *LeadScorer) scoreLocation(location string) float64 {
	location = strings.ToLower(location)

	// High-value locations (target markets)
	highValueLocations := []string{"india", "usa", "canada", "uk", "australia", "singapore", "germany", "france"}
	for _, hvl := range highValueLocations {
		if strings.Contains(location, hvl) {
			return 85.0
		}
	}

	return 60.0 // Other locations
}

func (ls *LeadScorer) scoreExperienceLevel(experience string) float64 {
	experience = strings.ToLower(experience)

	if strings.Contains(experience, "senior") || strings.Contains(experience, "lead") {
		return 80.0
	} else if strings.Contains(experience, "mid") || strings.Contains(experience, "intermediate") {
		return 70.0
	} else if strings.Contains(experience, "junior") || strings.Contains(experience, "entry") {
		return 50.0
	}

	return 60.0 // Default
}

func (ls *LeadScorer) scoreEngagementLevel(behavior Behavior) float64 {
	var score float64

	// Page views scoring
	if behavior.PageViews >= 10 {
		score += 30.0
	} else if behavior.PageViews >= 5 {
		score += 20.0
	} else if behavior.PageViews >= 2 {
		score += 10.0
	}

	// Time on site scoring
	if behavior.TotalTimeOnSite >= 1800 { // 30+ minutes
		score += 30.0
	} else if behavior.TotalTimeOnSite >= 900 { // 15+ minutes
		score += 20.0
	} else if behavior.TotalTimeOnSite >= 300 { // 5+ minutes
		score += 10.0
	}

	// Repeat visits scoring
	if behavior.VisitCount >= 5 {
		score += 40.0
	} else if behavior.VisitCount >= 3 {
		score += 25.0
	} else if behavior.VisitCount >= 2 {
		score += 15.0
	}

	return math.Min(score, 100.0)
}

func (ls *LeadScorer) scoreContentConsumption(behavior Behavior) float64 {
	var score float64

	// Blog posts read
	if behavior.BlogPostsRead >= 5 {
		score += 30.0
	} else if behavior.BlogPostsRead >= 3 {
		score += 20.0
	} else if behavior.BlogPostsRead >= 1 {
		score += 10.0
	}

	// Downloads
	score += float64(behavior.Downloads) * 15.0

	// Video engagement
	if behavior.VideoWatchTime >= 600 { // 10+ minutes
		score += 25.0
	} else if behavior.VideoWatchTime >= 300 { // 5+ minutes
		score += 15.0
	} else if behavior.VideoWatchTime >= 60 { // 1+ minute
		score += 5.0
	}

	// Social engagement
	score += float64(behavior.SocialEngagements) * 5.0

	return math.Min(score, 100.0)
}

func (ls *LeadScorer) scoreWebsiteActivity(behavior Behavior) float64 {
	var score float64

	// Depth of visit (pages per session)
	avgPagesPerSession := float64(behavior.PageViews) / float64(behavior.VisitCount)
	if avgPagesPerSession >= 5 {
		score += 25.0
	} else if avgPagesPerSession >= 3 {
		score += 15.0
	} else if avgPagesPerSession >= 2 {
		score += 10.0
	}

	// Service/pricing page visits
	if behavior.ServicePagesVisited {
		score += 30.0
	}
	if behavior.PricingPagesVisited {
		score += 35.0
	}

	// Contact page visits
	if behavior.ContactPagesVisited {
		score += 20.0
	}

	// Search behavior
	score += float64(behavior.SearchQueries) * 5.0

	return math.Min(score, 100.0)
}

func (ls *LeadScorer) scoreRecency(lastActivity time.Time) float64 {
	if lastActivity.IsZero() {
		return 0.0
	}

	daysSinceActivity := time.Since(lastActivity).Hours() / 24

	if daysSinceActivity <= 1 {
		return 100.0
	} else if daysSinceActivity <= 7 {
		return 80.0
	} else if daysSinceActivity <= 30 {
		return 60.0
	} else if daysSinceActivity <= 90 {
		return 40.0
	} else {
		return 20.0
	}
}

func (ls *LeadScorer) scoreCompanySize(size string) float64 {
	size = strings.ToLower(size)

	if strings.Contains(size, "enterprise") || strings.Contains(size, "large") {
		return 90.0
	} else if strings.Contains(size, "medium") || strings.Contains(size, "mid") {
		return 80.0
	} else if strings.Contains(size, "small") || strings.Contains(size, "startup") {
		return 70.0
	}

	return 60.0
}

func (ls *LeadScorer) scoreIndustryFit(industry string) float64 {
	// Use same logic as demographic industry scoring
	return ls.scoreIndustry(industry)
}

func (ls *LeadScorer) scoreRevenue(revenue string) float64 {
	revenue = strings.ToLower(revenue)

	if strings.Contains(revenue, "100m+") || strings.Contains(revenue, "billion") {
		return 95.0
	} else if strings.Contains(revenue, "50m") || strings.Contains(revenue, "10m") {
		return 85.0
	} else if strings.Contains(revenue, "1m") || strings.Contains(revenue, "5m") {
		return 75.0
	} else if strings.Contains(revenue, "500k") || strings.Contains(revenue, "1m") {
		return 65.0
	}

	return 50.0
}

func (ls *LeadScorer) scoreTechnologyStack(stack []string) float64 {
	if len(stack) == 0 {
		return 50.0
	}

	relevantTech := []string{"react", "node", "python", "go", "aws", "azure", "gcp", "kubernetes", "docker"}
	matchCount := 0

	for _, tech := range stack {
		techLower := strings.ToLower(tech)
		for _, relevant := range relevantTech {
			if strings.Contains(techLower, relevant) {
				matchCount++
				break
			}
		}
	}

	return math.Min(50.0+float64(matchCount)*10.0, 100.0)
}

func (ls *LeadScorer) scoreSourceType(sourceType string) float64 {
	switch strings.ToLower(sourceType) {
	case "contact_form":
		return 95.0 // Highest intent
	case "download":
		return 85.0 // High intent
	case "newsletter":
		return 70.0 // Medium intent
	case "cta":
		return 80.0 // High intent
	case "social_share":
		return 60.0 // Lower intent
	default:
		return 50.0
	}
}

func (ls *LeadScorer) scoreContentTypeEngagement(contentTypes []string) float64 {
	if len(contentTypes) == 0 {
		return 40.0
	}

	var totalScore float64
	for _, contentType := range contentTypes {
		switch strings.ToLower(contentType) {
		case "case_study":
			totalScore += 90.0
		case "whitepaper":
			totalScore += 85.0
		case "webinar":
			totalScore += 80.0
		case "tutorial":
			totalScore += 70.0
		case "blog":
			totalScore += 60.0
		default:
			totalScore += 50.0
		}
	}

	return math.Min(totalScore/float64(len(contentTypes)), 100.0)
}

func (ls *LeadScorer) scoreCTAInteraction(interactions int) float64 {
	if interactions >= 5 {
		return 100.0
	} else if interactions >= 3 {
		return 80.0
	} else if interactions >= 1 {
		return 60.0
	}
	return 20.0
}

func (ls *LeadScorer) scoreFormCompletions(completions int) float64 {
	if completions >= 3 {
		return 100.0
	} else if completions >= 2 {
		return 85.0
	} else if completions >= 1 {
		return 70.0
	}
	return 30.0
}

// AutoQualifyLead determines if a lead should be automatically qualified
func (ls *LeadScorer) AutoQualifyLead(leadScore int, profile LeadProfile) string {
	if leadScore >= 80 {
		return "hot"
	} else if leadScore >= 60 {
		return "warm"
	} else if leadScore >= 40 {
		return "cold"
	} else {
		return "unqualified"
	}
}

// Data structures for lead scoring

type LeadProfile struct {
	Demographics Demographics
	Behavior     Behavior
	Company      Company
	Intent       Intent
}

type Demographics struct {
	JobTitle        string
	Industry        string
	Location        string
	ExperienceLevel string
}

type Behavior struct {
	PageViews           int
	TotalTimeOnSite     int // seconds
	VisitCount          int
	BlogPostsRead       int
	Downloads           int
	VideoWatchTime      int // seconds
	SocialEngagements   int
	ServicePagesVisited bool
	PricingPagesVisited bool
	ContactPagesVisited bool
	SearchQueries       int
	LastActivity        time.Time
}

type Company struct {
	Name            string
	Size            string
	Industry        string
	Revenue         string
	TechnologyStack []string
}

type Intent struct {
	SourceType      string
	ContentTypes    []string
	CTAInteractions int
	FormCompletions int
}
