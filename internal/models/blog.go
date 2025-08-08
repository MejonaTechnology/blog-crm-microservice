package models

import (
	"time"
	"gorm.io/gorm"
)

// Blog represents a blog post with CRM integration
type Blog struct {
	ID                    uint      `json:"id" gorm:"primaryKey"`
	Title                 string    `json:"title" gorm:"not null;size:500"`
	Slug                  string    `json:"slug" gorm:"uniqueIndex;not null;size:255"`
	Content               string    `json:"content" gorm:"type:longtext"`
	Excerpt               string    `json:"excerpt" gorm:"size:1000"`
	Status                string    `json:"status" gorm:"default:'draft';size:20"` // draft, published, archived
	AuthorID              uint      `json:"author_id" gorm:"not null"`
	CategoryID            *uint     `json:"category_id"`
	FeaturedImage         string    `json:"featured_image" gorm:"size:500"`
	MetaTitle             string    `json:"meta_title" gorm:"size:255"`
	MetaDescription       string    `json:"meta_description" gorm:"size:500"`
	FocusKeyword          string    `json:"focus_keyword" gorm:"size:255"`
	ViewsCount            int       `json:"views_count" gorm:"default:0"`
	LikesCount            int       `json:"likes_count" gorm:"default:0"`
	CommentsCount         int       `json:"comments_count" gorm:"default:0"`
	SharesCount           int       `json:"shares_count" gorm:"default:0"`
	PublishedAt           *time.Time `json:"published_at"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
	DeletedAt             gorm.DeletedAt `json:"-" gorm:"index"`

	// CRM Fields for lead generation and analytics
	LeadGenerationCount   int       `json:"lead_generation_count" gorm:"default:0"`
	ConversionRate        float64   `json:"conversion_rate" gorm:"default:0.00;type:decimal(5,2)"`
	EngagementScore       float64   `json:"engagement_score" gorm:"default:0.00;type:decimal(5,2)"`
	SEOScore              float64   `json:"seo_score" gorm:"default:0.00;type:decimal(5,2)"`
	PerformanceStatus     string    `json:"performance_status" gorm:"default:'average';size:20"` // poor, average, good, excellent
	RevenueAttribution    float64   `json:"revenue_attribution" gorm:"default:0.00;type:decimal(10,2)"`
	LeadSource            string    `json:"lead_source" gorm:"default:'organic';size:100"`
	UTMSource             string    `json:"utm_source" gorm:"size:100"`
	UTMMedium             string    `json:"utm_medium" gorm:"size:100"`
	UTMCampaign           string    `json:"utm_campaign" gorm:"size:100"`
}

// AdminUser represents the user who authors/manages blogs
type AdminUser struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Name     string `json:"name" gorm:"not null;size:255"`
	Email    string `json:"email" gorm:"uniqueIndex;not null;size:255"`
	Role     string `json:"role" gorm:"default:'editor';size:50"`
	Active   bool   `json:"active" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}