# Blog CRM Microservice - API Documentation

![Version](https://img.shields.io/badge/version-1.0.0-green) ![Go](https://img.shields.io/badge/go-1.23+-blue) ![Port](https://img.shields.io/badge/port-8082-orange)

## Overview

The Blog CRM Management Microservice provides comprehensive blog management capabilities with integrated lead generation and analytics features. This API enables content creation, lead tracking, and performance monitoring for the Mejona Technology blog platform.

## Base URLs

- **Production**: `http://65.1.94.25:8082`
- **Nginx Proxy**: `http://65.1.94.25/api/v1/blog`
- **Development**: `http://localhost:8082`

## Authentication

### JWT Bearer Authentication

Protected endpoints require JWT Bearer token authentication:

```bash
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     http://65.1.94.25:8082/api/v1/blogs
```

### Token Structure

```json
{
  "user_id": "uuid",
  "email": "user@mejona.com",
  "role": "admin",
  "permissions": ["blog:read", "blog:write", "analytics:read"],
  "exp": 1704721200,
  "iat": 1704634800
}
```

## Rate Limiting

| Endpoint Type | Rate Limit | Burst |
|---------------|------------|-------|
| General API | 100 req/min | 20 |
| Health Checks | 30 req/min | 10 |
| Documentation | 10 req/min | 5 |
| Authentication | 5 req/min | 3 |

## Response Format

### Success Response

```json
{
  "success": true,
  "data": {
    // Response data
  },
  "meta": {
    "timestamp": "2025-01-08T12:00:00Z",
    "request_id": "req_123456789"
  }
}
```

### Error Response

```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid request parameters",
    "details": {
      "field": "email",
      "reason": "invalid_format"
    }
  },
  "meta": {
    "timestamp": "2025-01-08T12:00:00Z",
    "request_id": "req_123456789",
    "path": "/api/v1/blogs"
  }
}
```

## API Endpoints

### 🏥 Health & Monitoring

#### GET /health
Basic health check endpoint.

```bash
curl http://65.1.94.25:8082/health
```

**Response:**
```json
{
  "status": "healthy",
  "service": "blog-service",
  "version": "1.0.0",
  "timestamp": "2025-01-08T12:00:00Z",
  "uptime": 3600
}
```

#### GET /health/deep
Comprehensive health check with dependency status.

```bash
curl http://65.1.94.25:8082/health/deep
```

**Response:**
```json
{
  "status": "healthy",
  "service": "blog-service",
  "version": "1.0.0",
  "timestamp": "2025-01-08T12:00:00Z",
  "uptime": 3600,
  "checks": {
    "database": {
      "status": "connected",
      "response_time": "5ms"
    },
    "redis": {
      "status": "connected", 
      "response_time": "2ms"
    }
  },
  "system": {
    "memory": {
      "used": "45MB",
      "available": "512MB"
    },
    "goroutines": 25
  }
}
```

#### GET /status
Quick service status check.

```bash
curl http://65.1.94.25:8082/status
```

#### GET /ready
Kubernetes-style readiness probe.

```bash
curl http://65.1.94.25:8082/ready
```

#### GET /alive
Kubernetes-style liveness probe.

```bash
curl http://65.1.94.25:8082/alive
```

#### GET /metrics
Prometheus-style metrics endpoint.

```bash
curl http://65.1.94.25:8082/metrics
```

### 🧪 Testing

#### GET /api/v1/test
API functionality test endpoint.

```bash
curl http://65.1.94.25:8082/api/v1/test
```

**Response:**
```json
{
  "success": true,
  "message": "Blog service test endpoint working",
  "data": {
    "service": "Blog CRM Management Microservice",
    "version": "1.0.0",
    "status": "operational",
    "port": "8082",
    "timestamp": "2025-01-08T12:00:00Z"
  }
}
```

### 📝 Blog Management

#### GET /api/v1/blogs
List all blog posts with pagination and filtering.

**Parameters:**
- `page` (integer): Page number (default: 1)
- `limit` (integer): Items per page (default: 20, max: 100)
- `status` (string): Filter by status (`draft`, `published`, `archived`)
- `category` (string): Filter by category
- `search` (string): Search in title and content

```bash
curl -H "Authorization: Bearer $TOKEN" \
     "http://65.1.94.25:8082/api/v1/blogs?page=1&limit=10&status=published"
```

**Response:**
```json
{
  "blogs": [
    {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "title": "How to Build Modern Web Applications",
      "slug": "how-to-build-modern-web-applications",
      "content": "In this comprehensive guide...",
      "excerpt": "Learn the essentials...",
      "status": "published",
      "category": "Technology",
      "tags": ["web-development", "javascript", "react"],
      "author": {
        "id": "author-uuid",
        "name": "Sarah Johnson",
        "email": "sarah@mejona.com"
      },
      "seo": {
        "meta_title": "How to Build Modern Web Applications",
        "meta_description": "Learn the essentials...",
        "keywords": ["web development", "React"]
      },
      "metadata": {
        "read_time": 8,
        "word_count": 1500,
        "featured_image": "https://mejona.in/images/blog/featured.jpg"
      },
      "analytics": {
        "views": 1250,
        "shares": 45,
        "leads_generated": 12
      },
      "created_at": "2025-01-01T10:00:00Z",
      "updated_at": "2025-01-01T10:00:00Z",
      "published_at": "2025-01-01T12:00:00Z"
    }
  ],
  "pagination": {
    "current_page": 1,
    "per_page": 10,
    "total_pages": 5,
    "has_next": true,
    "has_previous": false
  },
  "total": 47
}
```

#### POST /api/v1/blogs
Create a new blog post.

```bash
curl -X POST \
     -H "Authorization: Bearer $TOKEN" \
     -H "Content-Type: application/json" \
     -d '{
       "title": "Advanced React Patterns",
       "slug": "advanced-react-patterns",
       "content": "React has evolved significantly...",
       "excerpt": "Explore advanced React patterns...",
       "status": "draft",
       "category": "Technology",
       "tags": ["react", "javascript", "patterns"],
       "seo": {
         "meta_title": "Advanced React Patterns | Mejona",
         "meta_description": "Explore advanced React patterns...",
         "keywords": ["React", "patterns", "advanced"]
       },
       "metadata": {
         "featured_image": "https://mejona.in/images/blog/react-patterns.jpg"
       }
     }' \
     http://65.1.94.25:8082/api/v1/blogs
```

#### GET /api/v1/blogs/{id}
Get a specific blog post by ID.

```bash
curl http://65.1.94.25:8082/api/v1/blogs/123e4567-e89b-12d3-a456-426614174000
```

#### PUT /api/v1/blogs/{id}
Update an existing blog post.

```bash
curl -X PUT \
     -H "Authorization: Bearer $TOKEN" \
     -H "Content-Type: application/json" \
     -d '{
       "title": "Updated Title",
       "status": "published"
     }' \
     http://65.1.94.25:8082/api/v1/blogs/123e4567-e89b-12d3-a456-426614174000
```

#### DELETE /api/v1/blogs/{id}
Delete a blog post.

```bash
curl -X DELETE \
     -H "Authorization: Bearer $TOKEN" \
     http://65.1.94.25:8082/api/v1/blogs/123e4567-e89b-12d3-a456-426614174000
```

### 🎯 Lead Generation

#### POST /api/v1/blogs/{id}/leads
Capture a lead from a blog post interaction.

```bash
curl -X POST \
     -H "Content-Type: application/json" \
     -d '{
       "email": "john.doe@example.com",
       "name": "John Doe",
       "phone": "+91 9876543210",
       "company": "Tech Solutions Inc.",
       "source_type": "newsletter",
       "utm_source": "google",
       "utm_medium": "organic",
       "utm_campaign": "blog-promotion",
       "additional_data": {
         "interest": "web-development",
         "budget": "10k-50k"
       }
     }' \
     http://65.1.94.25:8082/api/v1/blogs/123e4567-e89b-12d3-a456-426614174000/leads
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "lead-uuid",
    "email": "john.doe@example.com",
    "name": "John Doe",
    "phone": "+91 9876543210",
    "company": "Tech Solutions Inc.",
    "source_blog_id": "123e4567-e89b-12d3-a456-426614174000",
    "source_type": "newsletter",
    "utm_data": {
      "source": "google",
      "medium": "organic",
      "campaign": "blog-promotion"
    },
    "lead_score": 75,
    "status": "new",
    "created_at": "2025-01-08T12:00:00Z"
  }
}
```

#### GET /api/v1/leads
List captured leads with filtering and pagination.

```bash
curl -H "Authorization: Bearer $TOKEN" \
     "http://65.1.94.25:8082/api/v1/leads?page=1&limit=20&source=blog-uuid"
```

### 📊 Analytics

#### GET /api/v1/analytics/blogs/{id}
Get detailed analytics for a specific blog post.

**Parameters:**
- `period` (string): Time period (`24h`, `7d`, `30d`, `90d`)

```bash
curl -H "Authorization: Bearer $TOKEN" \
     "http://65.1.94.25:8082/api/v1/analytics/blogs/123e4567-e89b-12d3-a456-426614174000?period=7d"
```

**Response:**
```json
{
  "blog_id": "123e4567-e89b-12d3-a456-426614174000",
  "period": "7d",
  "metrics": {
    "views": {
      "total": 1250,
      "unique": 980,
      "daily_breakdown": [
        {
          "date": "2025-01-01",
          "views": 180
        }
      ]
    },
    "engagement": {
      "average_time_on_page": "4m 32s",
      "bounce_rate": 35.5,
      "scroll_depth": 78.2
    },
    "leads": {
      "total_generated": 12,
      "conversion_rate": 0.96
    },
    "social": {
      "shares": 45,
      "comments": 8
    },
    "seo": {
      "search_ranking": 3,
      "organic_traffic": 750,
      "keyword_positions": [
        {
          "keyword": "modern web development",
          "position": 3
        }
      ]
    }
  }
}
```

#### GET /api/v1/analytics/dashboard
Get overall analytics dashboard data.

```bash
curl -H "Authorization: Bearer $TOKEN" \
     "http://65.1.94.25:8082/api/v1/analytics/dashboard?period=30d"
```

**Response:**
```json
{
  "period": "30d",
  "overview": {
    "total_blogs": 45,
    "published_blogs": 32,
    "draft_blogs": 13,
    "total_views": 15750,
    "total_leads": 234,
    "conversion_rate": 1.48
  },
  "top_performing_blogs": [
    {
      "id": "blog-uuid",
      "title": "How to Build Modern Web Applications",
      "views": 2500,
      "leads": 45,
      "conversion_rate": 1.8
    }
  ],
  "traffic_sources": {
    "organic": 8750,
    "direct": 3200,
    "social": 2100,
    "referral": 1700
  },
  "monthly_trends": [
    {
      "month": "2025-01",
      "views": 15750,
      "leads": 234,
      "blogs_published": 8
    }
  ]
}
```

## Error Codes

| Code | HTTP Status | Description |
|------|-------------|-------------|
| `VALIDATION_ERROR` | 400 | Request validation failed |
| `UNAUTHORIZED` | 401 | Authentication required |
| `FORBIDDEN` | 403 | Insufficient permissions |
| `NOT_FOUND` | 404 | Resource not found |
| `CONFLICT` | 409 | Resource already exists |
| `RATE_LIMITED` | 429 | Too many requests |
| `INTERNAL_ERROR` | 500 | Internal server error |
| `DATABASE_ERROR` | 503 | Database connection failed |

## Status Codes

| Code | Description |
|------|-------------|
| `200` | Success |
| `201` | Resource created |
| `400` | Bad request |
| `401` | Unauthorized |
| `403` | Forbidden |
| `404` | Not found |
| `409` | Conflict |
| `429` | Too many requests |
| `500` | Internal server error |
| `503` | Service unavailable |

## SDK Examples

### JavaScript/Node.js

```javascript
const axios = require('axios');

class BlogCRMClient {
  constructor(baseURL, token) {
    this.client = axios.create({
      baseURL,
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    });
  }

  async getBlogs(params = {}) {
    const response = await this.client.get('/api/v1/blogs', { params });
    return response.data;
  }

  async createBlog(blogData) {
    const response = await this.client.post('/api/v1/blogs', blogData);
    return response.data;
  }

  async captureLead(blogId, leadData) {
    const response = await this.client.post(`/api/v1/blogs/${blogId}/leads`, leadData);
    return response.data;
  }

  async getAnalytics(blogId, period = '7d') {
    const response = await this.client.get(`/api/v1/analytics/blogs/${blogId}`, {
      params: { period }
    });
    return response.data;
  }
}

// Usage
const client = new BlogCRMClient('http://65.1.94.25:8082', 'your-jwt-token');
```

### Python

```python
import requests

class BlogCRMClient:
    def __init__(self, base_url, token):
        self.base_url = base_url
        self.session = requests.Session()
        self.session.headers.update({
            'Authorization': f'Bearer {token}',
            'Content-Type': 'application/json'
        })

    def get_blogs(self, **params):
        response = self.session.get(f'{self.base_url}/api/v1/blogs', params=params)
        return response.json()

    def create_blog(self, blog_data):
        response = self.session.post(f'{self.base_url}/api/v1/blogs', json=blog_data)
        return response.json()

    def capture_lead(self, blog_id, lead_data):
        response = self.session.post(
            f'{self.base_url}/api/v1/blogs/{blog_id}/leads', 
            json=lead_data
        )
        return response.json()

    def get_analytics(self, blog_id, period='7d'):
        response = self.session.get(
            f'{self.base_url}/api/v1/analytics/blogs/{blog_id}',
            params={'period': period}
        )
        return response.json()

# Usage
client = BlogCRMClient('http://65.1.94.25:8082', 'your-jwt-token')
```

### curl Examples

```bash
# Create a blog post
curl -X POST \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d @blog-post.json \
  http://65.1.94.25:8082/api/v1/blogs

# Capture a lead
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "name": "John Doe",
    "source_type": "newsletter"
  }' \
  http://65.1.94.25:8082/api/v1/blogs/$BLOG_ID/leads

# Get analytics
curl -H "Authorization: Bearer $TOKEN" \
  "http://65.1.94.25:8082/api/v1/analytics/blogs/$BLOG_ID?period=30d"
```

## Webhooks (Future)

The service will support webhooks for real-time notifications:

- `blog.published` - Blog post published
- `lead.captured` - New lead captured
- `analytics.milestone` - Analytics milestone reached

## OpenAPI/Swagger

- **Swagger UI**: http://65.1.94.25:8082/swagger/index.html
- **OpenAPI Spec**: http://65.1.94.25:8082/swagger/doc.json
- **Yaml Spec**: Available in `/docs/swagger.yaml`

## Testing

Test the API using the provided endpoints:

```bash
# Health check
curl http://65.1.94.25:8082/health

# API test
curl http://65.1.94.25:8082/api/v1/test

# Through nginx proxy
curl http://65.1.94.25/blog/health
```

## Support

- **Email**: support@mejona.com
- **Documentation**: http://65.1.94.25:8082/swagger
- **GitHub Issues**: Submit issues with API logs
- **Status Page**: http://65.1.94.25:8082/status

---

**Last Updated**: 2025-01-08  
**API Version**: 1.0.0  
**Service Port**: 8082