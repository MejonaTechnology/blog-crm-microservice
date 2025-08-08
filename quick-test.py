#!/usr/bin/env python3
"""
Quick Blog Service API Test
"""

import requests
import json
from datetime import datetime

def test_endpoints():
    base_url = "http://localhost:8082"
    
    print(f"TESTING BLOG CRM API ENDPOINTS")
    print(f"Target: {base_url}")
    print(f"Time: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")
    print("=" * 50)
    
    # Test health endpoints
    health_endpoints = ["/health", "/health/deep", "/status", "/ready", "/alive", "/metrics"]
    
    for endpoint in health_endpoints:
        try:
            response = requests.get(f"{base_url}{endpoint}", timeout=5)
            print(f"GET {endpoint:15} -> {response.status_code} {response.reason}")
            
            if endpoint == "/health" and response.status_code == 200:
                try:
                    data = response.json()
                    if 'data' in data and 'service' in data['data']:
                        print(f"    Service: {data['data']['service']}")
                        if 'uptime' in data['data']:
                            print(f"    Uptime: {data['data']['uptime']}")
                        if 'database' in data['data']:
                            print(f"    Database: {data['data']['database']}")
                except:
                    pass
                    
        except requests.exceptions.RequestException as e:
            print(f"GET {endpoint:15} -> ERROR: {str(e)}")
    
    # Test basic API endpoint
    try:
        response = requests.get(f"{base_url}/api/v1/test", timeout=5)
        print(f"GET /api/v1/test    -> {response.status_code} {response.reason}")
        if response.status_code == 200:
            try:
                data = response.json()
                print(f"    Message: {data.get('message', 'No message')}")
            except:
                pass
    except requests.exceptions.RequestException as e:
        print(f"GET /api/v1/test    -> ERROR: {str(e)}")
    
    # Test blogs endpoint (should require auth)
    try:
        response = requests.get(f"{base_url}/api/v1/blogs", timeout=5)
        print(f"GET /api/v1/blogs   -> {response.status_code} {response.reason}")
        if response.status_code == 401:
            print(f"    Expected: Authentication required")
    except requests.exceptions.RequestException as e:
        print(f"GET /api/v1/blogs   -> ERROR: {str(e)}")
    
    print("\nTest completed!")

if __name__ == "__main__":
    test_endpoints()