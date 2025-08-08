package models

import (
	"time"
)

// Blog Analytics Models

// BlogAnalyticsRequest represents a request for blog analytics
type BlogAnalyticsRequest struct {
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Granularity string    `json:"granularity"` // day, week, month
	Categories  []string  `json:"categories"`
	AuthorIDs   []uint    `json:"author_ids"`
}

// BlogPerformanceResponse represents overall blog performance metrics
type BlogPerformanceResponse struct {
	Period             string                  `json:"period"`
	Summary            BlogPerformanceSummary  `json:"summary"`
	TopPerformers      []BlogPerformanceMetric `json:"top_performers"`
	CategoryMetrics    []CategoryPerformance   `json:"category_metrics"`
	AuthorMetrics      []AuthorPerformance     `json:"author_metrics"`
	TrendData          []BlogTrendData         `json:"trend_data"`
	EngagementInsights EngagementInsights      `json:"engagement_insights"`
}

// BlogPerformanceSummary represents summary statistics
type BlogPerformanceSummary struct {
	TotalPosts       int     `json:"total_posts"`
	TotalViews       int     `json:"total_views"`
	TotalEngagements int     `json:"total_engagements"`
	TotalLeads       int     `json:"total_leads"`
	TotalRevenue     float64 `json:"total_revenue"`
	AvgTimeOnPage    float64 `json:"avg_time_on_page"`
	AvgBounceRate    float64 `json:"avg_bounce_rate"`
	AvgSocialShares  float64 `json:"avg_social_shares"`
	ConversionRate   float64 `json:"conversion_rate"`
	LeadValuePerPost float64 `json:"lead_value_per_post"`
	ROI              float64 `json:"roi"`
	GrowthRate       float64 `json:"growth_rate"`
}

// BlogPerformanceMetric represents individual blog performance
type BlogPerformanceMetric struct {
	BlogID          uint      `json:"blog_id"`
	Title           string    `json:"title"`
	URL             string    `json:"url"`
	PublishedAt     time.Time `json:"published_at"`
	Category        string    `json:"category"`
	AuthorID        uint      `json:"author_id"`
	AuthorName      string    `json:"author_name"`
	Views           int       `json:"views"`
	UniqueVisitors  int       `json:"unique_visitors"`
	TimeOnPage      float64   `json:"time_on_page"`
	BounceRate      float64   `json:"bounce_rate"`
	SocialShares    int       `json:"social_shares"`
	Comments        int       `json:"comments"`
	Leads           int       `json:"leads"`
	Revenue         float64   `json:"revenue"`
	EngagementScore int       `json:"engagement_score"`
	SEOScore        int       `json:"seo_score"`
	PerformanceRank int       `json:"performance_rank"`
}

// CategoryPerformance represents category-wise performance
type CategoryPerformance struct {
	Category         string  `json:"category"`
	PostCount        int     `json:"post_count"`
	TotalViews       int     `json:"total_views"`
	AvgViews         float64 `json:"avg_views"`
	TotalEngagements int     `json:"total_engagements"`
	AvgEngagements   float64 `json:"avg_engagements"`
	TotalLeads       int     `json:"total_leads"`
	ConversionRate   float64 `json:"conversion_rate"`
	Revenue          float64 `json:"revenue"`
	ROI              float64 `json:"roi"`
	PopularityScore  int     `json:"popularity_score"`
}

// AuthorPerformance represents author-wise performance
type AuthorPerformance struct {
	AuthorID         uint    `json:"author_id"`
	AuthorName       string  `json:"author_name"`
	PostCount        int     `json:"post_count"`
	TotalViews       int     `json:"total_views"`
	AvgViews         float64 `json:"avg_views"`
	TotalEngagements int     `json:"total_engagements"`
	EngagementRate   float64 `json:"engagement_rate"`
	TotalLeads       int     `json:"total_leads"`
	LeadConversion   float64 `json:"lead_conversion"`
	Revenue          float64 `json:"revenue"`
	AuthorScore      int     `json:"author_score"`
}

// BlogTrendData represents trend data over time
type BlogTrendData struct {
	Date        time.Time `json:"date"`
	Views       int       `json:"views"`
	Engagements int       `json:"engagements"`
	Leads       int       `json:"leads"`
	Revenue     float64   `json:"revenue"`
	PostCount   int       `json:"post_count"`
}

// EngagementInsights represents engagement analysis
type EngagementInsights struct {
	TopEngagementDrivers []EngagementDriver `json:"top_engagement_drivers"`
	EngagementPatterns   EngagementPatterns `json:"engagement_patterns"`
	AudienceSegments     []AudienceSegment  `json:"audience_segments"`
	ContentPreferences   ContentPreferences `json:"content_preferences"`
}

// EngagementDriver represents factors driving engagement
type EngagementDriver struct {
	Driver      string  `json:"driver"`
	Impact      string  `json:"impact"` // high, medium, low
	Description string  `json:"description"`
	Score       float64 `json:"score"`
}

// EngagementPatterns represents user engagement patterns
type EngagementPatterns struct {
	BestPublishingTime   map[string]float64 `json:"best_publishing_time"` // day of week
	BestPublishingHour   map[int]float64    `json:"best_publishing_hour"` // hour of day
	SeasonalPatterns     map[string]float64 `json:"seasonal_patterns"`    // month
	OptimalContentLength int                `json:"optimal_content_length"`
	OptimalHeadingCount  int                `json:"optimal_heading_count"`
}

// AudienceSegment represents audience segmentation data
type AudienceSegment struct {
	Segment        string   `json:"segment"`
	Size           int      `json:"size"`
	Percentage     float64  `json:"percentage"`
	EngagementRate float64  `json:"engagement_rate"`
	ConversionRate float64  `json:"conversion_rate"`
	RevenuePerUser float64  `json:"revenue_per_user"`
	TopTopics      []string `json:"top_topics"`
}

// ContentPreferences represents content preference analysis
type ContentPreferences struct {
	PreferredLength    []ContentLengthPreference `json:"preferred_length"`
	PreferredFormat    []ContentFormatPreference `json:"preferred_format"`
	PreferredTopics    []TopicPreference         `json:"preferred_topics"`
	PreferredTone      string                    `json:"preferred_tone"`
	OptimalImageCount  int                       `json:"optimal_image_count"`
	OptimalVideoLength int                       `json:"optimal_video_length"`
}

// BlogMetricsRequest represents a request for individual blog metrics
type BlogMetricsRequest struct {
	BlogID             uint      `json:"blog_id"`
	StartDate          time.Time `json:"start_date"`
	EndDate            time.Time `json:"end_date"`
	Period             string    `json:"period"`
	IncludeComparisons bool      `json:"include_comparisons"`
}

// BlogMetricsResponse represents individual blog metrics response
type BlogMetricsResponse struct {
	BlogID           uint                   `json:"blog_id"`
	Title            string                 `json:"title"`
	URL              string                 `json:"url"`
	PublishedAt      time.Time              `json:"published_at"`
	Period           string                 `json:"period"`
	Metrics          BlogDetailedMetrics    `json:"metrics"`
	Comparisons      *BlogMetricsComparison `json:"comparisons,omitempty"`
	TrendAnalysis    BlogTrendAnalysis      `json:"trend_analysis"`
	Recommendations  []string               `json:"recommendations"`
	OptimizationTips []OptimizationTip      `json:"optimization_tips"`
}

// BlogDetailedMetrics represents detailed metrics for a single blog
type BlogDetailedMetrics struct {
	Views             int                   `json:"views"`
	UniqueVisitors    int                   `json:"unique_visitors"`
	PageViews         int                   `json:"page_views"`
	Sessions          int                   `json:"sessions"`
	TimeOnPage        float64               `json:"time_on_page"`
	BounceRate        float64               `json:"bounce_rate"`
	ExitRate          float64               `json:"exit_rate"`
	ScrollDepth       ScrollDepthMetrics    `json:"scroll_depth"`
	SocialShares      SocialShareMetrics    `json:"social_shares"`
	Comments          CommentMetrics        `json:"comments"`
	EmailShares       int                   `json:"email_shares"`
	PrintActions      int                   `json:"print_actions"`
	Downloads         int                   `json:"downloads"`
	ClickThroughRate  float64               `json:"click_through_rate"`
	ConversionMetrics BlogConversionMetrics `json:"conversion_metrics"`
	SEOMetrics        BlogSEOMetrics        `json:"seo_metrics"`
	EngagementScore   int                   `json:"engagement_score"`
	PerformanceScore  int                   `json:"performance_score"`
}

// ScrollDepthMetrics represents scroll depth analysis
type ScrollDepthMetrics struct {
	AvgScrollDepth float64              `json:"avg_scroll_depth"`
	ScrollDepth25  int                  `json:"scroll_depth_25"`  // users reaching 25%
	ScrollDepth50  int                  `json:"scroll_depth_50"`  // users reaching 50%
	ScrollDepth75  int                  `json:"scroll_depth_75"`  // users reaching 75%
	ScrollDepth100 int                  `json:"scroll_depth_100"` // users reaching 100%
	ScrollHeatmap  []ScrollHeatmapPoint `json:"scroll_heatmap"`
}

// ScrollHeatmapPoint represents a point in the scroll heatmap
type ScrollHeatmapPoint struct {
	Position   int     `json:"position"`   // percentage of page
	Engagement float64 `json:"engagement"` // engagement at this position
	TimeSpent  float64 `json:"time_spent"` // average time spent at this position
}

// SocialShareMetrics represents social sharing analysis
type SocialShareMetrics struct {
	TotalShares      int               `json:"total_shares"`
	SharesByPlatform map[string]int    `json:"shares_by_platform"`
	ShareVelocity    float64           `json:"share_velocity"` // shares per hour
	ViralCoefficient float64           `json:"viral_coefficient"`
	InfluencerShares []InfluencerShare `json:"influencer_shares"`
}

// InfluencerShare represents shares by influencers
type InfluencerShare struct {
	Platform      string    `json:"platform"`
	Username      string    `json:"username"`
	FollowerCount int       `json:"follower_count"`
	ShareDate     time.Time `json:"share_date"`
	Engagement    int       `json:"engagement"`
}

// CommentMetrics represents comment analysis
type CommentMetrics struct {
	TotalComments     int              `json:"total_comments"`
	ApprovedComments  int              `json:"approved_comments"`
	RejectedComments  int              `json:"rejected_comments"`
	AvgResponseTime   float64          `json:"avg_response_time"` // hours
	SentimentAnalysis CommentSentiment `json:"sentiment_analysis"`
	TopCommenters     []TopCommenter   `json:"top_commenters"`
}

// CommentSentiment represents comment sentiment analysis
type CommentSentiment struct {
	PositiveComments int     `json:"positive_comments"`
	NegativeComments int     `json:"negative_comments"`
	NeutralComments  int     `json:"neutral_comments"`
	OverallSentiment string  `json:"overall_sentiment"` // positive, negative, neutral
	SentimentScore   float64 `json:"sentiment_score"`   // -1 to 1
}

// TopCommenter represents active commenters
type TopCommenter struct {
	Name         string  `json:"name"`
	Email        string  `json:"email"`
	CommentCount int     `json:"comment_count"`
	AvgSentiment float64 `json:"avg_sentiment"`
	IsInfluencer bool    `json:"is_influencer"`
}

// BlogConversionMetrics represents conversion tracking for blogs
type BlogConversionMetrics struct {
	Leads             int                  `json:"leads"`
	QualifiedLeads    int                  `json:"qualified_leads"`
	Customers         int                  `json:"customers"`
	Revenue           float64              `json:"revenue"`
	ConversionRate    float64              `json:"conversion_rate"`
	LeadQualityScore  int                  `json:"lead_quality_score"`
	ConversionFunnel  BlogConversionFunnel `json:"conversion_funnel"`
	ConversionSources []ConversionSource   `json:"conversion_sources"`
	AttributedRevenue float64              `json:"attributed_revenue"`
	CustomerLTV       float64              `json:"customer_ltv"`
}

// BlogConversionFunnel represents conversion funnel for blog traffic
type BlogConversionFunnel struct {
	Visitors        int             `json:"visitors"`
	Engaged         int             `json:"engaged"` // users who spent >30 seconds
	CTAClicks       int             `json:"cta_clicks"`
	FormViews       int             `json:"form_views"`
	FormSubmissions int             `json:"form_submissions"`
	LeadNurturing   int             `json:"lead_nurturing"`
	SalesQualified  int             `json:"sales_qualified"`
	Customers       int             `json:"customers"`
	ConversionRates ConversionRates `json:"conversion_rates"`
}

// ConversionRates represents conversion rates at each funnel stage
type ConversionRates struct {
	VisitorToEngaged    float64 `json:"visitor_to_engaged"`
	EngagedToCTA        float64 `json:"engaged_to_cta"`
	CTAToForm           float64 `json:"cta_to_form"`
	FormToSubmission    float64 `json:"form_to_submission"`
	SubmissionToNurture float64 `json:"submission_to_nurture"`
	NurtureToSQL        float64 `json:"nurture_to_sql"`
	SQLToCustomer       float64 `json:"sql_to_customer"`
	OverallConversion   float64 `json:"overall_conversion"`
}

// ConversionSource represents conversion source breakdown
type ConversionSource struct {
	Source         string  `json:"source"`
	Leads          int     `json:"leads"`
	Customers      int     `json:"customers"`
	Revenue        float64 `json:"revenue"`
	ConversionRate float64 `json:"conversion_rate"`
	Quality        string  `json:"quality"` // high, medium, low
}

// BlogSEOMetrics represents SEO performance metrics
type BlogSEOMetrics struct {
	OrganicTraffic     int              `json:"organic_traffic"`
	KeywordRankings    []KeywordRanking `json:"keyword_rankings"`
	FeaturedSnippets   int              `json:"featured_snippets"`
	Backlinks          BacklinkMetrics  `json:"backlinks"`
	SearchVisibility   float64          `json:"search_visibility"`
	ClickThroughRate   float64          `json:"click_through_rate"`
	AveragePosition    float64          `json:"average_position"`
	Impressions        int              `json:"impressions"`
	SEOScore           int              `json:"seo_score"`
	TechnicalSEOIssues []SEOIssue       `json:"technical_seo_issues"`
}

// KeywordRanking represents keyword ranking data
type KeywordRanking struct {
	Keyword          string    `json:"keyword"`
	Position         int       `json:"position"`
	PreviousPosition *int      `json:"previous_position"`
	SearchVolume     int       `json:"search_volume"`
	Difficulty       int       `json:"difficulty"`
	Traffic          int       `json:"estimated_traffic"`
	URL              string    `json:"url"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// BacklinkMetrics represents backlink analysis
type BacklinkMetrics struct {
	TotalBacklinks   int              `json:"total_backlinks"`
	ReferringDomains int              `json:"referring_domains"`
	NewBacklinks     int              `json:"new_backlinks"`
	LostBacklinks    int              `json:"lost_backlinks"`
	DomainAuthority  int              `json:"domain_authority"`
	TopReferrers     []BacklinkSource `json:"top_referrers"`
	BacklinkQuality  BacklinkQuality  `json:"backlink_quality"`
}

// BacklinkSource represents a backlink source
type BacklinkSource struct {
	Domain          string    `json:"domain"`
	URL             string    `json:"url"`
	AnchorText      string    `json:"anchor_text"`
	DomainAuthority int       `json:"domain_authority"`
	FirstSeen       time.Time `json:"first_seen"`
	LastSeen        time.Time `json:"last_seen"`
	IsDoFollow      bool      `json:"is_dofollow"`
}

// BacklinkQuality represents backlink quality distribution
type BacklinkQuality struct {
	HighQuality   int `json:"high_quality"`
	MediumQuality int `json:"medium_quality"`
	LowQuality    int `json:"low_quality"`
	ToxicLinks    int `json:"toxic_links"`
	QualityScore  int `json:"quality_score"`
}

// SEOIssue represents technical SEO issues
type SEOIssue struct {
	Type           string `json:"type"`
	Severity       string `json:"severity"` // critical, warning, info
	Description    string `json:"description"`
	Recommendation string `json:"recommendation"`
	Impact         string `json:"impact"`
}

// BlogMetricsComparison represents period-over-period comparison
type BlogMetricsComparison struct {
	PreviousPeriod    BlogBasicMetrics `json:"previous_period"`
	CurrentPeriod     BlogBasicMetrics `json:"current_period"`
	Changes           MetricChanges    `json:"changes"`
	PercentageChanges MetricChanges    `json:"percentage_changes"`
}

// BlogBasicMetrics represents basic metrics for comparison
type BlogBasicMetrics struct {
	Views        int     `json:"views"`
	Engagements  int     `json:"engagements"`
	Leads        int     `json:"leads"`
	Revenue      float64 `json:"revenue"`
	TimeOnPage   float64 `json:"time_on_page"`
	BounceRate   float64 `json:"bounce_rate"`
	SocialShares int     `json:"social_shares"`
}

// MetricChanges represents changes in metrics
type MetricChanges struct {
	Views        float64 `json:"views"`
	Engagements  float64 `json:"engagements"`
	Leads        float64 `json:"leads"`
	Revenue      float64 `json:"revenue"`
	TimeOnPage   float64 `json:"time_on_page"`
	BounceRate   float64 `json:"bounce_rate"`
	SocialShares float64 `json:"social_shares"`
}

// BlogTrendAnalysis represents trend analysis for a blog
type BlogTrendAnalysis struct {
	TrendDirection       string               `json:"trend_direction"` // up, down, stable
	TrendStrength        string               `json:"trend_strength"`  // strong, moderate, weak
	GrowthRate           float64              `json:"growth_rate"`
	Seasonality          SeasonalityInfo      `json:"seasonality"`
	PredictedPerformance PredictedPerformance `json:"predicted_performance"`
	PerformanceFactors   []PerformanceFactor  `json:"performance_factors"`
}

// SeasonalityInfo represents seasonality analysis
type SeasonalityInfo struct {
	HasSeasonality   bool     `json:"has_seasonality"`
	SeasonalPattern  string   `json:"seasonal_pattern"`
	PeakMonths       []string `json:"peak_months"`
	LowMonths        []string `json:"low_months"`
	SeasonalityScore int      `json:"seasonality_score"`
}

// PredictedPerformance represents performance predictions
type PredictedPerformance struct {
	NextMonth       BlogBasicMetrics `json:"next_month"`
	NextQuarter     BlogBasicMetrics `json:"next_quarter"`
	Confidence      float64          `json:"confidence"` // 0-100
	PredictionBasis string           `json:"prediction_basis"`
}

// PerformanceFactor represents factors affecting performance
type PerformanceFactor struct {
	Factor      string  `json:"factor"`
	Impact      string  `json:"impact"`   // positive, negative, neutral
	Strength    string  `json:"strength"` // high, medium, low
	Description string  `json:"description"`
	Score       float64 `json:"score"`
}

// OptimizationTip represents actionable optimization suggestions
type OptimizationTip struct {
	Category    string `json:"category"` // seo, engagement, conversion, social
	Priority    string `json:"priority"` // high, medium, low
	Title       string `json:"title"`
	Description string `json:"description"`
	Impact      string `json:"expected_impact"`
	Effort      string `json:"required_effort"`
	Action      string `json:"action_items"`
}

// Supporting data structures for content preferences

type ContentLengthPreference struct {
	LengthRange    string  `json:"length_range"` // e.g., "500-1000", "1000-2000"
	Percentage     float64 `json:"percentage"`   // percentage of audience preferring this
	EngagementRate float64 `json:"engagement_rate"`
	ConversionRate float64 `json:"conversion_rate"`
}

type ContentFormatPreference struct {
	Format         string  `json:"format"` // e.g., "listicle", "how-to", "case-study"
	Percentage     float64 `json:"percentage"`
	EngagementRate float64 `json:"engagement_rate"`
	ShareRate      float64 `json:"share_rate"`
}

type TopicPreference struct {
	Topic          string  `json:"topic"`
	Interest       float64 `json:"interest_score"` // 0-100
	Engagement     float64 `json:"engagement_rate"`
	ConversionRate float64 `json:"conversion_rate"`
	Trending       bool    `json:"is_trending"`
}
