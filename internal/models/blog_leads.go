package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// Blog Lead Models

// BlogLead represents a lead captured from blog content
type BlogLead struct {
	ID                    uint                   `json:"id" gorm:"primaryKey"`
	Email                 string                 `json:"email" gorm:"size:255;not null;index"`
	Name                  string                 `json:"name" gorm:"size:255"`
	FirstName             string                 `json:"first_name" gorm:"size:100"`
	LastName              string                 `json:"last_name" gorm:"size:100"`
	Company               string                 `json:"company" gorm:"size:255"`
	JobTitle              string                 `json:"job_title" gorm:"size:255"`
	Phone                 string                 `json:"phone" gorm:"size:50"`
	Website               string                 `json:"website" gorm:"size:500"`
	LinkedInProfile       string                 `json:"linkedin_profile" gorm:"size:500"`
	
	// Blog-related fields
	BlogID                uint                   `json:"blog_id" gorm:"not null;index"`
	BlogTitle             string                 `json:"blog_title" gorm:"size:500"`
	BlogURL               string                 `json:"blog_url" gorm:"size:500"`
	BlogCategory          string                 `json:"blog_category" gorm:"size:100"`
	
	// Lead source and tracking
	SourceType            string                 `json:"source_type" gorm:"size:50;not null;index"` // cta, contact_form, newsletter, download, social_share, etc.
	SourceDetails         JSONMap                `json:"source_details" gorm:"type:json"`
	CaptureMethod         string                 `json:"capture_method" gorm:"size:50"` // inline_form, popup, exit_intent, scroll_trigger
	
	// Lead scoring and qualification
	LeadScore             int                    `json:"lead_score" gorm:"default:0;index"`
	Status                string                 `json:"status" gorm:"size:50;default:new;index"` // new, contacted, qualified, unqualified, converted, lost, nurturing
	AutoQualification     string                 `json:"auto_qualification" gorm:"size:50"` // hot, warm, cold, unqualified
	ManualQualification   string                 `json:"manual_qualification" gorm:"size:50"`
	QualificationNotes    string                 `json:"qualification_notes" gorm:"type:text"`
	
	// Timestamps
	CapturedAt            time.Time              `json:"captured_at" gorm:"not null;index"`
	QualifiedAt           *time.Time             `json:"qualified_at" gorm:"index"`
	LastContactedAt       *time.Time             `json:"last_contacted_at"`
	ConvertedAt           *time.Time             `json:"converted_at"`
	LastEngagementAt      *time.Time             `json:"last_engagement_at"`
	CreatedAt             time.Time              `json:"created_at"`
	UpdatedAt             time.Time              `json:"updated_at"`
	
	// User assignment and management
	AssignedTo            *uint                  `json:"assigned_to" gorm:"index"`
	QualifiedBy           *uint                  `json:"qualified_by"`
	
	// UTM and tracking parameters
	UTMSource             string                 `json:"utm_source" gorm:"size:100"`
	UTMMedium             string                 `json:"utm_medium" gorm:"size:100"`
	UTMCampaign           string                 `json:"utm_campaign" gorm:"size:100"`
	UTMTerm               string                 `json:"utm_term" gorm:"size:100"`
	UTMContent            string                 `json:"utm_content" gorm:"size:100"`
	
	// Traffic and referrer information
	TrafficSource         string                 `json:"traffic_source" gorm:"size:100"` // organic, direct, social, email, paid
	ReferrerURL           string                 `json:"referrer_url" gorm:"size:1000"`
	ReferrerDomain        string                 `json:"referrer_domain" gorm:"size:255"`
	LandingPage           string                 `json:"landing_page" gorm:"size:1000"`
	
	// Device and location information
	DeviceType            string                 `json:"device_type" gorm:"size:50"` // desktop, mobile, tablet
	Browser               string                 `json:"browser" gorm:"size:100"`
	OperatingSystem       string                 `json:"operating_system" gorm:"size:100"`
	IPAddress             string                 `json:"ip_address" gorm:"size:45"`
	Country               string                 `json:"country" gorm:"size:100"`
	Region                string                 `json:"region" gorm:"size:100"`
	City                  string                 `json:"city" gorm:"size:100"`
	Timezone              string                 `json:"timezone" gorm:"size:100"`
	
	// Engagement and behavior data
	TotalEngagements      int                    `json:"total_engagements" gorm:"default:0"`
	PageViewsBeforeCapture int                   `json:"page_views_before_capture" gorm:"default:0"`
	TimeOnSiteBeforeCapture int                  `json:"time_on_site_before_capture" gorm:"default:0"` // seconds
	ScrollDepthAtCapture  float64                `json:"scroll_depth_at_capture" gorm:"default:0"`
	PreviousVisits        int                    `json:"previous_visits" gorm:"default:0"`
	
	// Conversion and revenue tracking
	ConversionPath        JSONArray              `json:"conversion_path" gorm:"type:json"` // array of touchpoints
	AttributedRevenue     float64                `json:"attributed_revenue" gorm:"type:decimal(10,2);default:0"`
	ConversionValue       float64                `json:"conversion_value" gorm:"type:decimal(10,2);default:0"`
	CustomerLTV           float64                `json:"customer_ltv" gorm:"type:decimal(10,2);default:0"`
	
	// Additional metadata and custom fields
	Tags                  JSONArray              `json:"tags" gorm:"type:json"`
	CustomFields          JSONMap                `json:"custom_fields" gorm:"type:json"`
	Notes                 string                 `json:"notes" gorm:"type:text"`
	
	// Privacy and consent
	ConsentGiven          bool                   `json:"consent_given" gorm:"default:false"`
	ConsentType           string                 `json:"consent_type" gorm:"size:50"` // gdpr, ccpa, general
	ConsentTimestamp      *time.Time             `json:"consent_timestamp"`
	OptedOut              bool                   `json:"opted_out" gorm:"default:false"`
	OptOutTimestamp       *time.Time             `json:"opt_out_timestamp"`
	
	// Relationships
	Blog                  *Blog                  `json:"blog,omitempty" gorm:"foreignKey:BlogID"`
	AssignedUser          *AdminUser             `json:"assigned_user,omitempty" gorm:"foreignKey:AssignedTo"`
	QualifiedByUser       *AdminUser             `json:"qualified_by_user,omitempty" gorm:"foreignKey:QualifiedBy"`
	Activities            []LeadActivity         `json:"activities,omitempty" gorm:"foreignKey:LeadID"`
	Touchpoints           []LeadTouchpoint       `json:"touchpoints,omitempty" gorm:"foreignKey:LeadID"`
}

// TableName specifies the table name for BlogLead
func (BlogLead) TableName() string {
	return "blog_leads"
}

// LeadActivity represents activities performed on or by a lead
type LeadActivity struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	LeadID      uint      `json:"lead_id" gorm:"not null;index"`
	ActivityType string   `json:"activity_type" gorm:"size:50;not null;index"` // email_sent, email_opened, page_view, form_submission, call_made, etc.
	Title       string    `json:"title" gorm:"size:255"`
	Description string    `json:"description" gorm:"type:text"`
	UserID      *uint     `json:"user_id" gorm:"index"` // user who performed the activity
	Metadata    JSONMap   `json:"metadata" gorm:"type:json"`
	CreatedAt   time.Time `json:"created_at" gorm:"index"`
	
	// Relationships
	Lead        *BlogLead  `json:"lead,omitempty" gorm:"foreignKey:LeadID"`
	User        *AdminUser `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// TableName specifies the table name for LeadActivity
func (LeadActivity) TableName() string {
	return "lead_activities"
}

// LeadTouchpoint represents touchpoints in the customer journey
type LeadTouchpoint struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	LeadID          uint      `json:"lead_id" gorm:"not null;index"`
	TouchpointType  string    `json:"touchpoint_type" gorm:"size:50;not null"` // blog_view, email_click, social_share, etc.
	BlogID          *uint     `json:"blog_id" gorm:"index"`
	URL             string    `json:"url" gorm:"size:1000"`
	Title           string    `json:"title" gorm:"size:500"`
	Source          string    `json:"source" gorm:"size:100"`
	Medium          string    `json:"medium" gorm:"size:100"`
	Campaign        string    `json:"campaign" gorm:"size:100"`
	TimeSpent       int       `json:"time_spent"` // seconds
	ScrollDepth     float64   `json:"scroll_depth"`
	Interactions    int       `json:"interactions"`
	ConversionValue float64   `json:"conversion_value" gorm:"type:decimal(10,2)"`
	AttributionWeight float64 `json:"attribution_weight" gorm:"type:decimal(3,2);default:0"`
	CreatedAt       time.Time `json:"created_at" gorm:"index"`
	
	// Relationships
	Lead            *BlogLead `json:"lead,omitempty" gorm:"foreignKey:LeadID"`
	Blog            *Blog     `json:"blog,omitempty" gorm:"foreignKey:BlogID"`
}

// TableName specifies the table name for LeadTouchpoint
func (LeadTouchpoint) TableName() string {
	return "lead_touchpoints"
}

// Request/Response Models

// BlogLeadCaptureRequest represents a request to capture a blog lead
type BlogLeadCaptureRequest struct {
	// Required fields
	Email      string `json:"email" binding:"required,email"`
	BlogID     uint   `json:"blog_id" binding:"required"`
	SourceType string `json:"source_type" binding:"required"`
	
	// Optional personal information
	Name         string `json:"name"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Company      string `json:"company"`
	JobTitle     string `json:"job_title"`
	Phone        string `json:"phone"`
	Website      string `json:"website"`
	LinkedInProfile string `json:"linkedin_profile"`
	
	// Source and tracking information
	SourceDetails   map[string]interface{} `json:"source_details"`
	CaptureMethod   string                 `json:"capture_method"`
	UTMSource       string                 `json:"utm_source"`
	UTMMedium       string                 `json:"utm_medium"`
	UTMCampaign     string                 `json:"utm_campaign"`
	UTMTerm         string                 `json:"utm_term"`
	UTMContent      string                 `json:"utm_content"`
	TrafficSource   string                 `json:"traffic_source"`
	ReferrerURL     string                 `json:"referrer_url"`
	LandingPage     string                 `json:"landing_page"`
	
	// Device and behavior information
	DeviceInfo      map[string]interface{} `json:"device_info"`
	LocationInfo    map[string]interface{} `json:"location_info"`
	EngagementData  map[string]interface{} `json:"engagement_data"`
	
	// Custom fields and metadata
	CustomFields    map[string]interface{} `json:"custom_fields"`
	Tags            []string               `json:"tags"`
	
	// Privacy and consent
	ConsentGiven    bool   `json:"consent_given"`
	ConsentType     string `json:"consent_type"`
}

// BlogLeadResponse represents a blog lead response
type BlogLeadResponse struct {
	ID                  uint                   `json:"id"`
	Email               string                 `json:"email"`
	Name                string                 `json:"name"`
	Company             string                 `json:"company"`
	Phone               string                 `json:"phone"`
	BlogID              uint                   `json:"blog_id"`
	BlogTitle           string                 `json:"blog_title"`
	SourceType          string                 `json:"source_type"`
	LeadScore           int                    `json:"lead_score"`
	Status              string                 `json:"status"`
	AutoQualification   string                 `json:"auto_qualification"`
	ManualQualification string                 `json:"manual_qualification"`
	QualificationNotes  string                 `json:"qualification_notes"`
	CapturedAt          time.Time              `json:"captured_at"`
	QualifiedAt         *time.Time             `json:"qualified_at"`
	UpdatedAt           *time.Time             `json:"updated_at"`
	UTMSource           string                 `json:"utm_source"`
	UTMMedium           string                 `json:"utm_medium"`
	UTMCampaign         string                 `json:"utm_campaign"`
	TrafficSource       string                 `json:"traffic_source"`
}

// DetailedBlogLeadResponse represents a detailed blog lead response
type DetailedBlogLeadResponse struct {
	BlogLeadResponse
	Activities         []LeadActivity         `json:"activities"`
	Touchpoints        []LeadTouchpoint       `json:"touchpoints"`
	ConversionPath     []string               `json:"conversion_path"`
	AttributedRevenue  float64                `json:"attributed_revenue"`
	LastEngagementAt   *time.Time             `json:"last_engagement_at"`
	TotalEngagements   int                    `json:"total_engagements"`
	DeviceInfo         map[string]interface{} `json:"device_info"`
	LocationInfo       map[string]interface{} `json:"location_info"`
	ReferrerInfo       map[string]interface{} `json:"referrer_info"`
}

// LeadQualificationRequest represents a request to qualify a lead
type LeadQualificationRequest struct {
	Status             string                 `json:"status" binding:"required"`
	QualificationNotes string                 `json:"qualification_notes"`
	LeadScore          *int                   `json:"lead_score"`
	Tags               []string               `json:"tags"`
	CustomFields       map[string]interface{} `json:"custom_fields"`
	FollowUpDate       *time.Time             `json:"follow_up_date"`
	AssignedTo         *uint                  `json:"assigned_to"`
	QualifiedBy        uint                   `json:"-"` // Set by handler from context
	QualifiedAt        time.Time              `json:"-"` // Set by handler
}

// LeadStatusUpdateRequest represents a request to update lead status
type LeadStatusUpdateRequest struct {
	Status      string                 `json:"status" binding:"required"`
	Notes       string                 `json:"notes"`
	CustomFields map[string]interface{} `json:"custom_fields"`
	UpdatedBy   uint                   `json:"-"` // Set by handler from context
	UpdatedAt   time.Time              `json:"-"` // Set by handler
}

// BlogLeadFilters represents filters for lead queries
type BlogLeadFilters struct {
	BlogID     uint       `json:"blog_id"`
	Status     string     `json:"status"`
	SourceType string     `json:"source_type"`
	AssignedTo *uint      `json:"assigned_to"`
	StartDate  *time.Time `json:"start_date"`
	EndDate    *time.Time `json:"end_date"`
	Tags       []string   `json:"tags"`
	MinScore   *int       `json:"min_score"`
	MaxScore   *int       `json:"max_score"`
	Sort       string     `json:"sort"`
	Order      string     `json:"order"`
	Page       int        `json:"page"`
	Limit      int        `json:"limit"`
}

// PaginatedBlogLeadsResponse represents paginated blog leads response
type PaginatedBlogLeadsResponse struct {
	Leads      []BlogLeadResponse `json:"leads"`
	Total      int64              `json:"total"`
	Page       int                `json:"page"`
	Limit      int                `json:"limit"`
	TotalPages int                `json:"total_pages"`
}

// Lead Analytics Models

// BlogLeadAnalyticsRequest represents a request for lead analytics
type BlogLeadAnalyticsRequest struct {
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
	BlogIDs      []uint    `json:"blog_ids"`
	LeadSources  []string  `json:"lead_sources"`
	Status       []string  `json:"status"`
	Granularity  string    `json:"granularity"` // day, week, month
}

// BlogLeadAnalyticsResponse represents blog lead analytics response
type BlogLeadAnalyticsResponse struct {
	Period           string                    `json:"period"`
	Summary          LeadAnalyticsSummary      `json:"summary"`
	ConversionFunnel LeadConversionFunnel      `json:"conversion_funnel"`
	SourceBreakdown  []LeadSourceAnalytics     `json:"source_breakdown"`
	BlogBreakdown    []BlogLeadAnalytics       `json:"blog_breakdown"`
	QualityMetrics   LeadQualityMetrics        `json:"quality_metrics"`
	TrendData        []LeadTrendData           `json:"trend_data"`
	Predictions      LeadPredictions           `json:"predictions"`
}

// LeadAnalyticsSummary represents summary lead analytics
type LeadAnalyticsSummary struct {
	TotalLeads         int     `json:"total_leads"`
	NewLeads           int     `json:"new_leads"`
	QualifiedLeads     int     `json:"qualified_leads"`
	ConvertedLeads     int     `json:"converted_leads"`
	ConversionRate     float64 `json:"conversion_rate"`
	AvgLeadScore       float64 `json:"avg_lead_score"`
	AvgTimeToConvert   float64 `json:"avg_time_to_convert"` // hours
	TotalRevenue       float64 `json:"total_revenue"`
	AvgLeadValue       float64 `json:"avg_lead_value"`
	LeadVelocity       float64 `json:"lead_velocity"` // leads per day
	GrowthRate         float64 `json:"growth_rate"`
}

// LeadConversionFunnel represents the lead conversion funnel
type LeadConversionFunnel struct {
	BlogVisitors       int                      `json:"blog_visitors"`
	LeadsCaptured      int                      `json:"leads_captured"`
	LeadsQualified     int                      `json:"leads_qualified"`
	LeadsNurtured      int                      `json:"leads_nurtured"`
	LeadsConverted     int                      `json:"leads_converted"`
	ConversionRates    FunnelConversionRates    `json:"conversion_rates"`
	DropOffAnalysis    FunnelDropOffAnalysis    `json:"drop_off_analysis"`
}

// FunnelConversionRates represents conversion rates at each stage
type FunnelConversionRates struct {
	VisitorToLead      float64 `json:"visitor_to_lead"`
	LeadToQualified    float64 `json:"lead_to_qualified"`
	QualifiedToNurture float64 `json:"qualified_to_nurture"`
	NurtureToConversion float64 `json:"nurture_to_conversion"`
	OverallConversion  float64 `json:"overall_conversion"`
}

// FunnelDropOffAnalysis represents drop-off analysis
type FunnelDropOffAnalysis struct {
	HighestDropOff     string                 `json:"highest_drop_off_stage"`
	DropOffReasons     []DropOffReason        `json:"drop_off_reasons"`
	ImprovementAreas   []string               `json:"improvement_areas"`
}

// DropOffReason represents reasons for lead drop-off
type DropOffReason struct {
	Stage       string  `json:"stage"`
	Reason      string  `json:"reason"`
	Percentage  float64 `json:"percentage"`
	Impact      string  `json:"impact"` // high, medium, low
}

// LeadSourceAnalytics represents analytics by lead source
type LeadSourceAnalytics struct {
	Source           string  `json:"source"`
	LeadCount        int     `json:"lead_count"`
	QualifiedCount   int     `json:"qualified_count"`
	ConvertedCount   int     `json:"converted_count"`
	Revenue          float64 `json:"revenue"`
	ConversionRate   float64 `json:"conversion_rate"`
	AvgLeadScore     float64 `json:"avg_lead_score"`
	AvgLeadValue     float64 `json:"avg_lead_value"`
	CostPerLead      float64 `json:"cost_per_lead"`
	ROI              float64 `json:"roi"`
	QualityScore     int     `json:"quality_score"`
}

// BlogLeadAnalytics represents analytics by blog
type BlogLeadAnalytics struct {
	BlogID           uint    `json:"blog_id"`
	BlogTitle        string  `json:"blog_title"`
	BlogURL          string  `json:"blog_url"`
	LeadCount        int     `json:"lead_count"`
	QualifiedCount   int     `json:"qualified_count"`
	ConvertedCount   int     `json:"converted_count"`
	Revenue          float64 `json:"revenue"`
	ConversionRate   float64 `json:"conversion_rate"`
	AvgLeadScore     float64 `json:"avg_lead_score"`
	LeadGenEfficiency float64 `json:"lead_gen_efficiency"` // leads per view
	PerformanceRank  int     `json:"performance_rank"`
}

// LeadQualityMetrics represents lead quality analysis
type LeadQualityMetrics struct {
	OverallQualityScore   int                        `json:"overall_quality_score"`
	QualityDistribution   LeadQualityDistribution    `json:"quality_distribution"`
	ScoreRangeAnalysis    []ScoreRangeAnalysis       `json:"score_range_analysis"`
	QualityTrends         QualityTrendAnalysis       `json:"quality_trends"`
	QualityFactors        []QualityFactor            `json:"quality_factors"`
}

// LeadQualityDistribution represents quality distribution
type LeadQualityDistribution struct {
	HighQuality    int     `json:"high_quality"`    // score 80-100
	MediumQuality  int     `json:"medium_quality"`  // score 50-79
	LowQuality     int     `json:"low_quality"`     // score 0-49
	Percentages    QualityPercentages `json:"percentages"`
}

// QualityPercentages represents quality percentages
type QualityPercentages struct {
	High   float64 `json:"high"`
	Medium float64 `json:"medium"`
	Low    float64 `json:"low"`
}

// ScoreRangeAnalysis represents analysis by score ranges
type ScoreRangeAnalysis struct {
	ScoreRange     string  `json:"score_range"`
	LeadCount      int     `json:"lead_count"`
	ConversionRate float64 `json:"conversion_rate"`
	AvgRevenue     float64 `json:"avg_revenue"`
	TimeToConvert  float64 `json:"avg_time_to_convert"`
}

// QualityTrendAnalysis represents quality trends over time
type QualityTrendAnalysis struct {
	TrendDirection   string                  `json:"trend_direction"` // improving, declining, stable
	QualityOverTime  []QualityDataPoint      `json:"quality_over_time"`
	SeasonalPatterns map[string]float64      `json:"seasonal_patterns"`
}

// QualityDataPoint represents a quality data point over time
type QualityDataPoint struct {
	Date         time.Time `json:"date"`
	AvgScore     float64   `json:"avg_score"`
	LeadCount    int       `json:"lead_count"`
	ConversionRate float64 `json:"conversion_rate"`
}

// QualityFactor represents factors affecting lead quality
type QualityFactor struct {
	Factor      string  `json:"factor"`
	Impact      string  `json:"impact"` // positive, negative
	Correlation float64 `json:"correlation"` // -1 to 1
	Description string  `json:"description"`
}

// LeadTrendData represents trend data over time
type LeadTrendData struct {
	Date            time.Time `json:"date"`
	LeadsCaptured   int       `json:"leads_captured"`
	LeadsQualified  int       `json:"leads_qualified"`
	LeadsConverted  int       `json:"leads_converted"`
	Revenue         float64   `json:"revenue"`
	AvgLeadScore    float64   `json:"avg_lead_score"`
	ConversionRate  float64   `json:"conversion_rate"`
}

// LeadPredictions represents predictive analytics
type LeadPredictions struct {
	NextMonthLeads     int                    `json:"next_month_leads"`
	NextMonthRevenue   float64                `json:"next_month_revenue"`
	PredictionConfidence float64              `json:"prediction_confidence"`
	TrendPredictions   []LeadTrendPrediction  `json:"trend_predictions"`
	RecommendedActions []string               `json:"recommended_actions"`
}

// LeadTrendPrediction represents trend predictions
type LeadTrendPrediction struct {
	Date           time.Time `json:"date"`
	PredictedLeads int       `json:"predicted_leads"`
	PredictedRevenue float64 `json:"predicted_revenue"`
	ConfidenceRange PredictionRange `json:"confidence_range"`
}

// PredictionRange represents prediction confidence ranges
type PredictionRange struct {
	Lower  float64 `json:"lower"`
	Upper  float64 `json:"upper"`
}

// Custom JSON types for database storage

// JSONMap represents a JSON map for database storage
type JSONMap map[string]interface{}

// Value implements the driver.Valuer interface for database storage
func (j JSONMap) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan implements the sql.Scanner interface for database retrieval
func (j *JSONMap) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("cannot scan %T into JSONMap", value)
	}
	
	return json.Unmarshal(bytes, j)
}

// JSONArray represents a JSON array for database storage
type JSONArray []interface{}

// Value implements the driver.Valuer interface for database storage
func (j JSONArray) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan implements the sql.Scanner interface for database retrieval
func (j *JSONArray) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("cannot scan %T into JSONArray", value)
	}
	
	return json.Unmarshal(bytes, j)
}