package seo

import (
	"time"
)

// SEO Analysis Models

// ContentData represents the input data for SEO analysis
type ContentData struct {
	ID                uint          `json:"id"`
	Title             string        `json:"title"`
	URL               string        `json:"url"`
	MetaDescription   string        `json:"meta_description"`
	Content           string        `json:"content"`
	PrimaryKeyword    string        `json:"primary_keyword"`
	SecondaryKeywords []string      `json:"secondary_keywords"`
	Headings          []HeadingData `json:"headings"`
	InternalLinks     []LinkData    `json:"internal_links"`
	ExternalLinks     []LinkData    `json:"external_links"`
	Images            []ImageData   `json:"images"`
	SchemaMarkup      string        `json:"schema_markup"`
	CanonicalURL      string        `json:"canonical_url"`
	LoadTime          float64       `json:"load_time"` // in seconds
	MobileResponsive  bool          `json:"mobile_responsive"`
}

// HeadingData represents heading structure information
type HeadingData struct {
	Level int    `json:"level"` // 1-6 for H1-H6
	Text  string `json:"text"`
}

// LinkData represents link information
type LinkData struct {
	URL        string `json:"url"`
	AnchorText string `json:"anchor_text"`
	IsDoFollow bool   `json:"is_dofollow"`
	IsInternal bool   `json:"is_internal"`
}

// ImageData represents image information
type ImageData struct {
	URL      string `json:"url"`
	FileName string `json:"file_name"`
	AltText  string `json:"alt_text"`
	Title    string `json:"title"`
	Size     int64  `json:"size"` // in bytes
}

// SEOAnalysis represents the complete SEO analysis result
type SEOAnalysis struct {
	ContentID          uint                `json:"content_id"`
	Title              string              `json:"title"`
	URL                string              `json:"url"`
	AnalyzedAt         time.Time           `json:"analyzed_at"`
	OverallScore       int                 `json:"overall_score"`
	TitleAnalysis      TitleAnalysis       `json:"title_analysis"`
	MetaAnalysis       MetaAnalysis        `json:"meta_analysis"`
	StructureAnalysis  StructureAnalysis   `json:"structure_analysis"`
	KeywordAnalysis    KeywordAnalysis     `json:"keyword_analysis"`
	ReadabilityAnalysis ReadabilityAnalysis `json:"readability_analysis"`
	TechnicalAnalysis  TechnicalAnalysis   `json:"technical_analysis"`
	LinkAnalysis       LinkAnalysis        `json:"link_analysis"`
	ImageAnalysis      ImageAnalysis       `json:"image_analysis"`
	Recommendations    []string            `json:"recommendations"`
	Opportunities      []Opportunity       `json:"opportunities"`
}

// TitleAnalysis represents title tag analysis
type TitleAnalysis struct {
	Title                  string   `json:"title"`
	Length                 int      `json:"length"`
	LengthScore            int      `json:"length_score"`
	LengthStatus           string   `json:"length_status"` // optimal, good, too_short, too_long
	ContainsPrimaryKeyword bool     `json:"contains_primary_keyword"`
	KeywordPosition        string   `json:"keyword_position"` // beginning, middle, end
	KeywordScore           int      `json:"keyword_score"`
	PowerWords             []string `json:"power_words"`
	PowerWordScore         float64  `json:"power_word_score"`
}

// MetaAnalysis represents meta description analysis
type MetaAnalysis struct {
	MetaDescription        string `json:"meta_description"`
	Length                 int    `json:"length"`
	LengthScore            int    `json:"length_score"`
	LengthStatus           string `json:"length_status"` // optimal, good, too_short, too_long
	ContainsPrimaryKeyword bool   `json:"contains_primary_keyword"`
	KeywordScore           int    `json:"keyword_score"`
	CallToAction           bool   `json:"call_to_action"`
	CTAScore               int    `json:"cta_score"`
}

// StructureAnalysis represents content structure analysis
type StructureAnalysis struct {
	WordCount             int     `json:"word_count"`
	WordCountScore        int     `json:"word_count_score"`
	WordCountStatus       string  `json:"word_count_status"` // optimal, good, too_short, very_long
	Headings              []HeadingData `json:"headings"`
	H1Count               int     `json:"h1_count"`
	H1Score               int     `json:"h1_score"`
	H1Status              string  `json:"h1_status"` // optimal, missing, multiple
	H2Count               int     `json:"h2_count"`
	H2Score               int     `json:"h2_score"`
	H2Status              string  `json:"h2_status"` // good, minimal, missing, many
	H3Count               int     `json:"h3_count"`
	AvgWordsPerParagraph  float64 `json:"avg_words_per_paragraph"`
	ParagraphScore        int     `json:"paragraph_score"`
}

// KeywordAnalysis represents keyword optimization analysis
type KeywordAnalysis struct {
	PrimaryKeyword           string                  `json:"primary_keyword"`
	PrimaryKeywordCount      int                     `json:"primary_keyword_count"`
	PrimaryKeywordDensity    float64                 `json:"primary_keyword_density"`
	PrimaryKeywordScore      int                     `json:"primary_keyword_score"`
	PrimaryKeywordStatus     string                  `json:"primary_keyword_status"` // optimal, good, too_low, too_high
	PrimaryKeywordInIntro    bool                    `json:"primary_keyword_in_intro"`
	IntroKeywordScore        int                     `json:"intro_keyword_score"`
	SecondaryKeywords        []string                `json:"secondary_keywords"`
	SecondaryKeywordData     []SecondaryKeywordData  `json:"secondary_keyword_data"`
	LSIKeywords              []string                `json:"lsi_keywords"`
	LSIScore                 float64                 `json:"lsi_score"`
}

// SecondaryKeywordData represents secondary keyword analysis
type SecondaryKeywordData struct {
	Keyword string  `json:"keyword"`
	Count   int     `json:"count"`
	Density float64 `json:"density"`
}

// ReadabilityAnalysis represents content readability analysis
type ReadabilityAnalysis struct {
	SentenceCount           int                     `json:"sentence_count"`
	WordCount               int                     `json:"word_count"`
	SyllableCount           int                     `json:"syllable_count"`
	AvgWordsPerSentence     float64                 `json:"avg_words_per_sentence"`
	AvgSyllablesPerWord     float64                 `json:"avg_syllables_per_word"`
	FleschScore             float64                 `json:"flesch_score"`
	FleschKincaidGrade      float64                 `json:"flesch_kincaid_grade"`
	ReadingLevel            string                  `json:"reading_level"`
	ReadabilityScore        int                     `json:"readability_score"`
	SentenceLengthAnalysis  SentenceLengthAnalysis  `json:"sentence_length_analysis"`
	TransitionWords         []string                `json:"transition_words"`
	TransitionWordScore     float64                 `json:"transition_word_score"`
}

// SentenceLengthAnalysis represents sentence length distribution analysis
type SentenceLengthAnalysis struct {
	AverageLength     float64 `json:"average_length"`
	ShortestSentence  int     `json:"shortest_sentence"`
	LongestSentence   int     `json:"longest_sentence"`
	ShortSentences    int     `json:"short_sentences"`    // <= 10 words
	MediumSentences   int     `json:"medium_sentences"`   // 11-20 words
	LongSentences     int     `json:"long_sentences"`     // > 20 words
}

// TechnicalAnalysis represents technical SEO analysis
type TechnicalAnalysis struct {
	URLAnalysis         URLAnalysis `json:"url_analysis"`
	HasSchemaMarkup     bool        `json:"has_schema_markup"`
	SchemaScore         int         `json:"schema_score"`
	HasCanonicalURL     bool        `json:"has_canonical_url"`
	CanonicalScore      int         `json:"canonical_score"`
	LoadTimeScore       int         `json:"load_time_score"`
	LoadTimeStatus      string      `json:"load_time_status"` // excellent, good, fair, poor
	IsMobileResponsive  bool        `json:"is_mobile_responsive"`
	MobileScore         int         `json:"mobile_score"`
}

// URLAnalysis represents URL structure analysis
type URLAnalysis struct {
	URL             string `json:"url"`
	Length          int    `json:"length"`
	LengthScore     int    `json:"length_score"`
	LengthStatus    string `json:"length_status"` // optimal, good, too_long
	ContainsKeyword bool   `json:"contains_keyword"`
	KeywordScore    int    `json:"keyword_score"`
	Structure       string `json:"structure"` // clean, complex
	StructureScore  int    `json:"structure_score"`
}

// LinkAnalysis represents link structure analysis
type LinkAnalysis struct {
	InternalLinks        []LinkData          `json:"internal_links"`
	ExternalLinks        []LinkData          `json:"external_links"`
	InternalLinkCount    int                 `json:"internal_link_count"`
	ExternalLinkCount    int                 `json:"external_link_count"`
	InternalLinkScore    int                 `json:"internal_link_score"`
	InternalLinkStatus   string              `json:"internal_link_status"` // optimal, good, missing, too_many
	ExternalLinkScore    int                 `json:"external_link_score"`
	ExternalLinkStatus   string              `json:"external_link_status"` // optimal, good, none, too_many
	AnchorTextAnalysis   AnchorTextAnalysis  `json:"anchor_text_analysis"`
}

// AnchorTextAnalysis represents anchor text optimization analysis
type AnchorTextAnalysis struct {
	AnchorTextVariety      int      `json:"anchor_text_variety"`
	MostUsedAnchorText     string   `json:"most_used_anchor_text"`
	MaxAnchorFrequency     int      `json:"max_anchor_frequency"`
	OverOptimizedAnchors   []string `json:"over_optimized_anchors"`
	DiversityScore         float64  `json:"diversity_score"`
}

// ImageAnalysis represents image optimization analysis
type ImageAnalysis struct {
	Images               []ImageData `json:"images"`
	ImageCount           int         `json:"image_count"`
	ImagesWithAltText    int         `json:"images_with_alt_text"`
	ImagesWithTitle      int         `json:"images_with_title"`
	OptimizedFileNames   int         `json:"optimized_file_names"`
	AltTextScore         float64     `json:"alt_text_score"`
	TitleScore           float64     `json:"title_score"`
	FileNameScore        float64     `json:"file_name_score"`
	ImageScore           float64     `json:"image_score"`
}

// Opportunity represents an SEO optimization opportunity
type Opportunity struct {
	Type        string `json:"type"`        // title_optimization, keyword_optimization, etc.
	Priority    string `json:"priority"`    // high, medium, low
	Impact      string `json:"impact"`      // high, medium, low
	Effort      string `json:"effort"`      // high, medium, low
	Title       string `json:"title"`
	Description string `json:"description"`
	Action      string `json:"action"`
}

// Competitor Analysis Models

// CompetitorAnalysisRequest represents a request for competitor analysis
type CompetitorAnalysisRequest struct {
	Competitors       []string `json:"competitors"`
	AnalysisDepth     string   `json:"analysis_depth"`     // basic, detailed, comprehensive
	ContentCategories []string `json:"content_categories"`
	Timeframe         string   `json:"timeframe"`          // 30d, 90d, 180d, 1y
}

// CompetitorAnalysisResponse represents competitor analysis results
type CompetitorAnalysisResponse struct {
	RequestedAt     time.Time            `json:"requested_at"`
	Competitors     []CompetitorProfile  `json:"competitors"`
	MarketPosition  MarketPosition       `json:"market_position"`
	ContentGaps     []ContentGap         `json:"content_gaps"`
	KeywordGaps     []KeywordGap         `json:"keyword_gaps"`
	Recommendations []string             `json:"recommendations"`
}

// CompetitorProfile represents a competitor's profile and metrics
type CompetitorProfile struct {
	Domain            string                 `json:"domain"`
	DomainAuthority   int                    `json:"domain_authority"`
	ContentVolume     ContentVolumeMetrics   `json:"content_volume"`
	KeywordMetrics    CompetitorKeywords     `json:"keyword_metrics"`
	ContentStrategy   ContentStrategyAnalysis `json:"content_strategy"`
	SocialPresence    SocialPresenceMetrics  `json:"social_presence"`
	TechnicalMetrics  TechnicalMetrics       `json:"technical_metrics"`
}

// ContentVolumeMetrics represents content volume analysis
type ContentVolumeMetrics struct {
	TotalPosts        int                    `json:"total_posts"`
	PostsPerMonth     float64                `json:"posts_per_month"`
	AvgPostLength     int                    `json:"avg_post_length"`
	ContentCategories map[string]int         `json:"content_categories"`
	PublishingPattern PublishingPatternData  `json:"publishing_pattern"`
}

// CompetitorKeywords represents competitor keyword analysis
type CompetitorKeywords struct {
	TotalKeywords      int                 `json:"total_keywords"`
	RankingKeywords    int                 `json:"ranking_keywords"`
	TopKeywords        []KeywordRanking    `json:"top_keywords"`
	KeywordOverlap     []string            `json:"keyword_overlap"`
	UniqueKeywords     []string            `json:"unique_keywords"`
	AvgKeywordRank     float64             `json:"avg_keyword_rank"`
}

// ContentStrategyAnalysis represents content strategy analysis
type ContentStrategyAnalysis struct {
	PrimaryTopics      []TopicAnalysis     `json:"primary_topics"`
	ContentTypes       map[string]int      `json:"content_types"`
	TargetAudience     AudienceAnalysis    `json:"target_audience"`
	ContentQuality     ContentQualityScore `json:"content_quality"`
	EngagementMetrics  EngagementSummary   `json:"engagement_metrics"`
}

// SocialPresenceMetrics represents social media presence analysis
type SocialPresenceMetrics struct {
	Platforms         map[string]SocialPlatformData `json:"platforms"`
	TotalFollowers    int                           `json:"total_followers"`
	EngagementRate    float64                       `json:"engagement_rate"`
	SocialShares      SocialSharesData              `json:"social_shares"`
	InfluencerScore   int                           `json:"influencer_score"`
}

// TechnicalMetrics represents technical performance metrics
type TechnicalMetrics struct {
	SiteSpeed         SiteSpeedMetrics    `json:"site_speed"`
	MobileOptimization MobileMetrics      `json:"mobile_optimization"`
	SEOHealth         SEOHealthMetrics    `json:"seo_health"`
	SecurityMetrics   SecurityMetrics     `json:"security_metrics"`
}

// Supporting data structures

type PublishingPatternData struct {
	DayOfWeekPattern map[string]int `json:"day_of_week_pattern"`
	HourlyPattern    map[int]int    `json:"hourly_pattern"`
	SeasonalPattern  map[string]int `json:"seasonal_pattern"`
}

type KeywordRanking struct {
	Keyword      string  `json:"keyword"`
	Rank         int     `json:"rank"`
	SearchVolume int     `json:"search_volume"`
	Difficulty   int     `json:"difficulty"`
	Traffic      int     `json:"estimated_traffic"`
}

type TopicAnalysis struct {
	Topic       string  `json:"topic"`
	PostCount   int     `json:"post_count"`
	Engagement  float64 `json:"avg_engagement"`
	KeywordRank float64 `json:"avg_keyword_rank"`
}

type AudienceAnalysis struct {
	Demographics    map[string]interface{} `json:"demographics"`
	Interests       []string               `json:"interests"`
	BehaviorProfile map[string]float64     `json:"behavior_profile"`
}

type ContentQualityScore struct {
	OverallScore    int     `json:"overall_score"`
	ReadabilityScore int    `json:"readability_score"`
	SEOScore        int     `json:"seo_score"`
	EngagementScore int     `json:"engagement_score"`
	OriginalityScore int    `json:"originality_score"`
}

type EngagementSummary struct {
	AvgTimeOnPage   float64 `json:"avg_time_on_page"`
	BounceRate      float64 `json:"bounce_rate"`
	SocialShares    int     `json:"avg_social_shares"`
	Comments        int     `json:"avg_comments"`
	BacklinksEarned int     `json:"avg_backlinks_earned"`
}

type SocialPlatformData struct {
	Platform    string  `json:"platform"`
	Followers   int     `json:"followers"`
	Engagement  float64 `json:"engagement_rate"`
	PostFreq    float64 `json:"posts_per_week"`
}

type SocialSharesData struct {
	TotalShares     int            `json:"total_shares"`
	SharesByPlatform map[string]int `json:"shares_by_platform"`
	ViralContent    []ContentPiece `json:"viral_content"`
}

type SiteSpeedMetrics struct {
	PageLoadTime      float64 `json:"page_load_time"`
	FirstContentPaint float64 `json:"first_content_paint"`
	TimeToInteractive float64 `json:"time_to_interactive"`
	SpeedScore        int     `json:"speed_score"`
}

type MobileMetrics struct {
	MobileFriendly    bool    `json:"mobile_friendly"`
	MobileSpeed       float64 `json:"mobile_speed"`
	ResponsiveDesign  bool    `json:"responsive_design"`
	MobileUsability   int     `json:"mobile_usability_score"`
}

type SEOHealthMetrics struct {
	CrawlErrors       int     `json:"crawl_errors"`
	IndexedPages      int     `json:"indexed_pages"`
	SitemapPresent    bool    `json:"sitemap_present"`
	RobotsTxtPresent  bool    `json:"robots_txt_present"`
	SSLCertificate    bool    `json:"ssl_certificate"`
	OverallSEOScore   int     `json:"overall_seo_score"`
}

type SecurityMetrics struct {
	SSLRating         string `json:"ssl_rating"`
	SecurityHeaders   int    `json:"security_headers_count"`
	VulnerabilityScore int   `json:"vulnerability_score"`
	TrustScore        int    `json:"trust_score"`
}

type ContentPiece struct {
	Title       string    `json:"title"`
	URL         string    `json:"url"`
	PublishedAt time.Time `json:"published_at"`
	Shares      int       `json:"shares"`
	Engagement  float64   `json:"engagement"`
}

// Market Position and Gap Analysis

type MarketPosition struct {
	OverallRank       int                    `json:"overall_rank"`
	MarketShare       float64                `json:"market_share"`
	StrengthsWeakness StrengthsWeaknessAnalysis `json:"strengths_weaknesses"`
	CompetitiveAdvantage string              `json:"competitive_advantage"`
	Threats           []string               `json:"threats"`
	Opportunities     []string               `json:"opportunities"`
}

type StrengthsWeaknessAnalysis struct {
	Strengths  []CompetitiveStrength  `json:"strengths"`
	Weaknesses []CompetitiveWeakness  `json:"weaknesses"`
}

type CompetitiveStrength struct {
	Area        string `json:"area"`
	Description string `json:"description"`
	Impact      string `json:"impact"` // high, medium, low
}

type CompetitiveWeakness struct {
	Area        string `json:"area"`
	Description string `json:"description"`
	Priority    string `json:"priority"` // high, medium, low
}

type ContentGap struct {
	Topic           string   `json:"topic"`
	Opportunity     string   `json:"opportunity"`
	CompetitorCount int      `json:"competitor_count"`
	SearchVolume    int      `json:"estimated_search_volume"`
	Difficulty      string   `json:"difficulty"`
	Priority        string   `json:"priority"`
	SuggestedAngles []string `json:"suggested_angles"`
}

type KeywordGap struct {
	Keyword         string  `json:"keyword"`
	YourRank        int     `json:"your_rank"`        // 0 if not ranking
	CompetitorRank  int     `json:"competitor_rank"`
	SearchVolume    int     `json:"search_volume"`
	Difficulty      int     `json:"difficulty"`
	Opportunity     string  `json:"opportunity"`      // high, medium, low
	TrafficPotential int    `json:"traffic_potential"`
}