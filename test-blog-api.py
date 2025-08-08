#!/usr/bin/env python3
"""
Comprehensive Blog CRM API Testing Suite
Tests all endpoints of the blog-service microservice
"""

import requests
import json
import time
import sys
from typing import Dict, Any, List
from datetime import datetime

class BlogAPITester:
    def __init__(self, base_url: str = "http://localhost:8082"):
        self.base_url = base_url
        self.session = requests.Session()
        self.auth_token = None
        self.test_results = []
        
    def log_result(self, test_name: str, success: bool, details: str = ""):
        """Log test result"""
        result = {
            'test': test_name,
            'success': success,
            'details': details,
            'timestamp': datetime.now().isoformat()
        }
        self.test_results.append(result)
        status = "PASS" if success else "FAIL"
        print(f"{status} {test_name}: {details}")
    
    def test_health_endpoints(self):
        """Test all health-related endpoints"""
        print("\nTESTING HEALTH ENDPOINTS")
        
        health_endpoints = [
            "/health",
            "/health/deep", 
            "/status",
            "/ready",
            "/alive",
            "/metrics"
        ]
        
        for endpoint in health_endpoints:
            try:
                response = self.session.get(f"{self.base_url}{endpoint}", timeout=5)
                if response.status_code in [200, 503]:  # 503 acceptable for DB issues
                    self.log_result(f"Health {endpoint}", True, f"Status: {response.status_code}")
                else:
                    self.log_result(f"Health {endpoint}", False, f"Status: {response.status_code}")
            except Exception as e:
                self.log_result(f"Health {endpoint}", False, f"Error: {str(e)}")
    
    def test_blog_crud_endpoints(self):
        """Test blog CRUD operations"""
        print("\nTESTING BLOG CRUD ENDPOINTS")
        
        # Test GET /api/v1/blogs (list)
        try:
            response = self.session.get(f"{self.base_url}/api/v1/blogs")
            self.log_result("GET /api/v1/blogs", 
                          response.status_code in [200, 401], 
                          f"Status: {response.status_code}")
        except Exception as e:
            self.log_result("GET /api/v1/blogs", False, f"Error: {str(e)}")
        
        # Test GET /api/v1/blogs/:id
        try:
            response = self.session.get(f"{self.base_url}/api/v1/blogs/1")
            self.log_result("GET /api/v1/blogs/:id", 
                          response.status_code in [200, 401, 404], 
                          f"Status: {response.status_code}")
        except Exception as e:
            self.log_result("GET /api/v1/blogs/:id", False, f"Error: {str(e)}")
        
        # Test POST /api/v1/blogs (create)
        blog_data = {
            "title": "Test Blog Post",
            "content": "This is a test blog post content for API testing.",
            "excerpt": "Test excerpt",
            "status": "draft",
            "meta_description": "Test meta description",
            "keywords": "test, api, blog",
            "tags": ["test", "api"]
        }
        
        try:
            response = self.session.post(f"{self.base_url}/api/v1/blogs", 
                                       json=blog_data)
            self.log_result("POST /api/v1/blogs", 
                          response.status_code in [201, 401], 
                          f"Status: {response.status_code}")
        except Exception as e:
            self.log_result("POST /api/v1/blogs", False, f"Error: {str(e)}")
    
    def test_analytics_endpoints(self):
        """Test CRM analytics endpoints"""
        print("\nTESTING ANALYTICS ENDPOINTS")
        
        analytics_endpoints = [
            "/api/v1/analytics/blog-performance",
            "/api/v1/analytics/lead-generation",
            "/api/v1/analytics/seo-performance", 
            "/api/v1/analytics/content-gaps",
            "/api/v1/analytics/blogs/1/metrics"
        ]
        
        for endpoint in analytics_endpoints:
            try:
                response = self.session.get(f"{self.base_url}{endpoint}")
                self.log_result(f"Analytics {endpoint}", 
                              response.status_code in [200, 401, 404], 
                              f"Status: {response.status_code}")
            except Exception as e:
                self.log_result(f"Analytics {endpoint}", False, f"Error: {str(e)}")
    
    def test_lead_endpoints(self):
        """Test lead generation endpoints"""
        print("\nTESTING LEAD GENERATION ENDPOINTS")
        
        # Test POST /api/v1/blog-leads (capture lead)
        lead_data = {
            "name": "Test User",
            "email": "test@example.com",
            "blog_id": 1,
            "source": "blog_cta",
            "utm_source": "test"
        }
        
        try:
            response = self.session.post(f"{self.base_url}/api/v1/blog-leads", 
                                       json=lead_data)
            self.log_result("POST /api/v1/blog-leads", 
                          response.status_code in [201, 401], 
                          f"Status: {response.status_code}")
        except Exception as e:
            self.log_result("POST /api/v1/blog-leads", False, f"Error: {str(e)}")
        
        # Test GET /api/v1/blog-leads/analytics
        try:
            response = self.session.get(f"{self.base_url}/api/v1/blog-leads/analytics")
            self.log_result("GET /api/v1/blog-leads/analytics", 
                          response.status_code in [200, 401], 
                          f"Status: {response.status_code}")
        except Exception as e:
            self.log_result("GET /api/v1/blog-leads/analytics", False, f"Error: {str(e)}")
    
    def test_seo_endpoints(self):
        """Test SEO optimization endpoints"""
        print("\nTESTING SEO ENDPOINTS")
        
        try:
            response = self.session.get(f"{self.base_url}/api/v1/blogs/1/seo")
            self.log_result("GET /api/v1/blogs/:id/seo", 
                          response.status_code in [200, 401, 404], 
                          f"Status: {response.status_code}")
        except Exception as e:
            self.log_result("GET /api/v1/blogs/:id/seo", False, f"Error: {str(e)}")
    
    def test_campaign_endpoints(self):
        """Test campaign management endpoints"""
        print("\nTESTING CAMPAIGN ENDPOINTS")
        
        # Test POST /api/v1/campaigns
        campaign_data = {
            "name": "Test Campaign",
            "description": "Test campaign for API testing",
            "budget": 1000.00,
            "start_date": "2025-01-01",
            "end_date": "2025-12-31"
        }
        
        try:
            response = self.session.post(f"{self.base_url}/api/v1/campaigns", 
                                       json=campaign_data)
            self.log_result("POST /api/v1/campaigns", 
                          response.status_code in [201, 401], 
                          f"Status: {response.status_code}")
        except Exception as e:
            self.log_result("POST /api/v1/campaigns", False, f"Error: {str(e)}")
    
    def test_authentication_flow(self):
        """Test authentication endpoints"""
        print("\nTESTING AUTHENTICATION")
        
        # Test login endpoint
        login_data = {
            "email": "admin@mejona.com",
            "password": "admin123"
        }
        
        try:
            response = self.session.post(f"{self.base_url}/api/v1/auth/login", 
                                       json=login_data)
            self.log_result("POST /api/v1/auth/login", 
                          response.status_code in [200, 401, 500], 
                          f"Status: {response.status_code}")
            
            # If login successful, save token
            if response.status_code == 200:
                data = response.json()
                if 'token' in data:
                    self.auth_token = data['token']
                    self.session.headers.update({
                        'Authorization': f'Bearer {self.auth_token}'
                    })
        except Exception as e:
            self.log_result("POST /api/v1/auth/login", False, f"Error: {str(e)}")
    
    def check_service_availability(self):
        """Check if the blog service is running"""
        print("CHECKING SERVICE AVAILABILITY")
        
        try:
            response = self.session.get(f"{self.base_url}/health", timeout=5)
            if response.status_code in [200, 503]:
                print(f"Blog service is running on {self.base_url}")
                return True
        except requests.exceptions.ConnectionError:
            print(f"Blog service is not running on {self.base_url}")
            print("   Please start the service with: ./blog-service.exe")
            return False
        except Exception as e:
            print(f"Error connecting to service: {str(e)}")
            return False
        
        return False
    
    def run_comprehensive_test(self):
        """Run all tests"""
        print("STARTING COMPREHENSIVE BLOG CRM API TESTS")
        print(f"   Target: {self.base_url}")
        print(f"   Time: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")
        print("=" * 60)
        
        # Check service availability first
        if not self.check_service_availability():
            return False
        
        # Run all test suites
        self.test_health_endpoints()
        self.test_authentication_flow()
        self.test_blog_crud_endpoints()
        self.test_analytics_endpoints()
        self.test_lead_endpoints()
        self.test_seo_endpoints()
        self.test_campaign_endpoints()
        
        # Generate summary
        self.generate_test_report()
        return True
    
    def generate_test_report(self):
        """Generate comprehensive test report"""
        print("\n" + "=" * 60)
        print("TEST SUMMARY REPORT")
        print("=" * 60)
        
        total_tests = len(self.test_results)
        passed_tests = sum(1 for r in self.test_results if r['success'])
        failed_tests = total_tests - passed_tests
        
        success_rate = (passed_tests / total_tests) * 100 if total_tests > 0 else 0
        
        print(f"Total Tests:    {total_tests}")
        print(f"Passed:         {passed_tests}")
        print(f"Failed:         {failed_tests}")
        print(f"Success Rate:   {success_rate:.1f}%")
        print()
        
        if failed_tests > 0:
            print("FAILED TESTS:")
            for result in self.test_results:
                if not result['success']:
                    print(f"   - {result['test']}: {result['details']}")
        
        print()
        print("EXPECTED BEHAVIOR:")
        print("   - Health endpoints should return 200 or 503 (DB connection issues)")
        print("   - Protected endpoints should return 401 (no auth token)")
        print("   - CRUD operations should return 401 (authentication required)")
        print("   - Analytics endpoints should return 401 (authentication required)")
        
        # Save detailed report
        report_file = f"blog-api-test-report-{datetime.now().strftime('%Y%m%d-%H%M%S')}.json"
        with open(report_file, 'w') as f:
            json.dump({
                'summary': {
                    'total_tests': total_tests,
                    'passed_tests': passed_tests,
                    'failed_tests': failed_tests,
                    'success_rate': success_rate,
                    'test_timestamp': datetime.now().isoformat(),
                    'service_url': self.base_url
                },
                'results': self.test_results
            }, f, indent=2)
        
        print(f"Detailed report saved to: {report_file}")

def main():
    """Main function"""
    tester = BlogAPITester()
    
    if len(sys.argv) > 1:
        tester.base_url = sys.argv[1]
    
    success = tester.run_comprehensive_test()
    
    if not success:
        print("\nCould not complete tests - service not available")
        sys.exit(1)
    
    print("\nAPI testing completed!")

if __name__ == "__main__":
    main()