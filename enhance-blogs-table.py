#!/usr/bin/env python3
"""
Blog Table CRM Enhancement Script
Adds 57 CRM fields to existing blogs table
"""

import pymysql
import sys

def enhance_blogs_table():
    try:
        # Connect to MySQL database
        connection = pymysql.connect(
            host='65.1.94.25',
            user='phpmyadmin',
            password='mFVarH2LCrQK',
            database='mejona_unified',
            charset='utf8mb4'
        )
        
        cursor = connection.cursor()
        
        print("Connected to production database successfully!")
        print("Enhancing blogs table with 57 CRM fields...")
        print("=" * 60)
        
        # CRM Enhancement Fields
        crm_fields = [
            # Lead Generation & Attribution
            "ADD COLUMN lead_source VARCHAR(100) DEFAULT 'organic'",
            "ADD COLUMN utm_source VARCHAR(100)",
            "ADD COLUMN utm_medium VARCHAR(100)", 
            "ADD COLUMN utm_campaign VARCHAR(100)",
            "ADD COLUMN utm_term VARCHAR(255)",
            "ADD COLUMN utm_content VARCHAR(255)",
            "ADD COLUMN referrer_url TEXT",
            "ADD COLUMN lead_generation_count INT DEFAULT 0",
            "ADD COLUMN conversion_rate DECIMAL(5,2) DEFAULT 0.00",
            "ADD COLUMN qualified_leads_count INT DEFAULT 0",
            
            # Performance & Analytics
            "ADD COLUMN engagement_score DECIMAL(4,2) DEFAULT 0.00",
            "ADD COLUMN bounce_rate DECIMAL(5,2) DEFAULT 0.00",
            "ADD COLUMN time_on_page_avg INT DEFAULT 0",
            "ADD COLUMN scroll_depth_avg DECIMAL(5,2) DEFAULT 0.00",
            "ADD COLUMN social_shares_count INT DEFAULT 0",
            "ADD COLUMN email_shares_count INT DEFAULT 0",
            "ADD COLUMN click_through_rate DECIMAL(5,2) DEFAULT 0.00",
            "ADD COLUMN page_views_unique INT DEFAULT 0",
            "ADD COLUMN returning_visitors INT DEFAULT 0",
            
            # SEO & Content Optimization
            "ADD COLUMN seo_score DECIMAL(4,2) DEFAULT 0.00",
            "ADD COLUMN content_quality_score DECIMAL(4,2) DEFAULT 0.00",
            "ADD COLUMN keyword_density_score DECIMAL(4,2) DEFAULT 0.00",
            "ADD COLUMN readability_score DECIMAL(4,2) DEFAULT 0.00",
            "ADD COLUMN meta_title VARCHAR(60)",
            "ADD COLUMN focus_keyword VARCHAR(255)",
            "ADD COLUMN secondary_keywords JSON",
            "ADD COLUMN internal_links_count INT DEFAULT 0",
            "ADD COLUMN external_links_count INT DEFAULT 0",
            "ADD COLUMN image_alt_tags_count INT DEFAULT 0",
            
            # Revenue & ROI
            "ADD COLUMN revenue_attribution DECIMAL(10,2) DEFAULT 0.00",
            "ADD COLUMN cost_per_acquisition DECIMAL(8,2) DEFAULT 0.00",
            "ADD COLUMN lifetime_value_generated DECIMAL(12,2) DEFAULT 0.00",
            "ADD COLUMN roi_percentage DECIMAL(6,2) DEFAULT 0.00",
            "ADD COLUMN conversion_value_avg DECIMAL(8,2) DEFAULT 0.00",
            
            # Audience & Demographics  
            "ADD COLUMN audience_segments JSON",
            "ADD COLUMN geographic_performance JSON",
            "ADD COLUMN device_performance JSON",
            "ADD COLUMN age_group_performance JSON",
            "ADD COLUMN gender_performance JSON",
            
            # Content Performance
            "ADD COLUMN performance_status ENUM('excellent','good','average','below_average','poor') DEFAULT 'average'",
            "ADD COLUMN optimization_score DECIMAL(4,2) DEFAULT 0.00",
            "ADD COLUMN content_freshness_score DECIMAL(4,2) DEFAULT 0.00",
            "ADD COLUMN user_engagement_events JSON",
            "ADD COLUMN conversion_events JSON",
            
            # Campaign Integration
            "ADD COLUMN campaign_id INT",
            "ADD COLUMN campaign_type VARCHAR(50)",
            "ADD COLUMN campaign_budget_allocated DECIMAL(10,2) DEFAULT 0.00",
            "ADD COLUMN campaign_spend DECIMAL(10,2) DEFAULT 0.00",
            "ADD COLUMN a_b_test_variant VARCHAR(10)",
            "ADD COLUMN a_b_test_results JSON",
            
            # Social Media Integration
            "ADD COLUMN social_media_performance JSON",
            "ADD COLUMN influencer_mentions_count INT DEFAULT 0",
            "ADD COLUMN viral_coefficient DECIMAL(6,4) DEFAULT 0.0000",
            
            # Advanced Analytics
            "ADD COLUMN heatmap_data JSON",
            "ADD COLUMN user_journey_data JSON",
            "ADD COLUMN conversion_funnel_data JSON",
            "ADD COLUMN predictive_performance_score DECIMAL(4,2) DEFAULT 0.00",
            "ADD COLUMN content_recommendation_score DECIMAL(4,2) DEFAULT 0.00"
        ]
        
        # Apply each field enhancement
        success_count = 0
        for i, field_sql in enumerate(crm_fields, 1):
            try:
                sql = f"ALTER TABLE blogs {field_sql}"
                cursor.execute(sql)
                print(f"  [{i:2d}/57] Added: {field_sql.split('ADD COLUMN')[1].split()[0]}")
                success_count += 1
            except pymysql.Error as e:
                if "Duplicate column name" in str(e):
                    print(f"  [{i:2d}/57] Exists: {field_sql.split('ADD COLUMN')[1].split()[0]}")
                    success_count += 1
                else:
                    print(f"  [{i:2d}/57] Failed: {field_sql.split('ADD COLUMN')[1].split()[0]} - {str(e)}")
        
        # Create indexes for performance
        print("\nCreating performance indexes...")
        indexes = [
            "CREATE INDEX idx_blogs_lead_source ON blogs(lead_source)",
            "CREATE INDEX idx_blogs_engagement_score ON blogs(engagement_score)",
            "CREATE INDEX idx_blogs_seo_score ON blogs(seo_score)", 
            "CREATE INDEX idx_blogs_performance_status ON blogs(performance_status)",
            "CREATE INDEX idx_blogs_campaign_id ON blogs(campaign_id)",
            "CREATE INDEX idx_blogs_conversion_rate ON blogs(conversion_rate)",
        ]
        
        for idx_sql in indexes:
            try:
                cursor.execute(idx_sql)
                index_name = idx_sql.split('CREATE INDEX')[1].split(' ON')[0].strip()
                print(f"  Created index: {index_name}")
            except pymysql.Error as e:
                if "Duplicate key name" in str(e):
                    index_name = idx_sql.split('CREATE INDEX')[1].split(' ON')[0].strip()
                    print(f"  Index exists: {index_name}")
                else:
                    print(f"  Index failed: {str(e)}")
        
        # Commit changes
        connection.commit()
        
        # Verify enhancements
        print("\nVerifying table structure...")
        cursor.execute("DESCRIBE blogs")
        columns = cursor.fetchall()
        print(f"Total columns now: {len(columns)} (was 20)")
        
        # Show sample of new CRM fields
        print("\nSample of new CRM fields:")
        new_fields = [col[0] for col in columns if col[0] not in [
            'id', 'title', 'slug', 'content', 'excerpt', 'author_id', 
            'banner_image_url', 'banner_image_prompt', 'content_images', 
            'tags', 'status', 'is_featured', 'meta_description', 
            'meta_keywords', 'reading_time_minutes', 'views_count', 
            'likes_count', 'created_at', 'updated_at', 'published_at'
        ]]
        
        for field in new_fields[:10]:  # Show first 10 new fields
            print(f"  - {field}")
        
        if len(new_fields) > 10:
            print(f"  ... and {len(new_fields) - 10} more CRM fields")
        
        cursor.close()
        connection.close()
        
        print(f"\nBlog table enhancement completed!")
        print(f"Successfully added {success_count}/57 CRM fields")
        print("Database is now ready for blog CRM service!")
        
        return True
        
    except Exception as e:
        print(f"Database enhancement failed: {str(e)}")
        return False

if __name__ == "__main__":
    if enhance_blogs_table():
        print("\nBlog CRM database enhancement successful!")
    else:
        print("\nDatabase enhancement failed!")
        sys.exit(1)