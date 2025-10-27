#!/usr/bin/env python
# -*- coding: utf-8 -*-

"""
Test REST Client - For testing py-gateway REST API

Simple client to test the py-gateway REST API implementation.
"""

import requests
import json
import sys


def test_health(host, port):
    """Test health endpoint"""
    print "\n" + "=" * 60
    print "TEST: Health Check"
    print "=" * 60
    
    try:
        url = "http://{}:{}/health".format(host, port)
        print "GET {}".format(url)
        
        response = requests.get(url, timeout=10)
        
        print "Status Code: {}".format(response.status_code)
        print "Response: {}".format(json.dumps(response.json(), indent=2))
        
        if response.status_code == 200:
            print "✓ Health check successful!"
        else:
            print "✗ Health check failed!"
            
    except Exception as e:
        print "✗ Health check failed with error: {}".format(str(e))


def test_ping(host, port):
    """Test ping endpoint"""
    print "\n" + "=" * 60
    print "TEST: Ping"
    print "=" * 60
    
    try:
        url = "http://{}:{}/ping".format(host, port)
        print "GET {}".format(url)
        
        response = requests.get(url, timeout=10)
        
        print "Status Code: {}".format(response.status_code)
        print "Response: {}".format(json.dumps(response.json(), indent=2))
        
        if response.status_code == 200:
            print "✓ Ping successful!"
        else:
            print "✗ Ping failed!"
            
    except Exception as e:
        print "✗ Ping failed with error: {}".format(str(e))


def test_get_account(host, port, account_number):
    """Test get account endpoint"""
    print "\n" + "=" * 60
    print "TEST: Get Account by Account Number"
    print "=" * 60
    
    try:
        url = "http://{}:{}/api/v1/accounts".format(host, port)
        params = {'account_number': account_number}
        print "GET {}?account_number={}".format(url, account_number)
        
        response = requests.get(url, params=params, timeout=10)
        
        print "Status Code: {}".format(response.status_code)
        print "Response: {}".format(json.dumps(response.json(), indent=2))
        
        if response.status_code == 200:
            print "✓ Account retrieval successful!"
            data = response.json().get('data', {})
            if 'account' in data:
                print "  Account Number: {}".format(data['account'].get('account_number'))
                print "  Account Type: {}".format(data['account'].get('account_type'))
                print "  Balance: {}".format(data['account'].get('balance'))
                print "  Currency: {}".format(data['account'].get('currency'))
            if 'customer' in data:
                print "  Customer: {}".format(data['customer'].get('full_name'))
                print "  Email: {}".format(data['customer'].get('email'))
        else:
            print "✗ Account retrieval failed!"
            
    except Exception as e:
        print "✗ Account retrieval failed with error: {}".format(str(e))


def main():
    """Main test function"""
    # Configuration
    host = 'localhost'
    port = 8083
    
    if len(sys.argv) > 1:
        host = sys.argv[1]
    if len(sys.argv) > 2:
        port = int(sys.argv[2])
    
    print "=" * 60
    print "py-gateway REST API Client Test"
    print "=" * 60
    print "Target: {}:{}".format(host, port)
    print "=" * 60
    
    try:
        # Test 1: Health Check
        test_health(host, port)
        
        # Test 2: Ping
        test_ping(host, port)
        
        # Test 3: Get Account
        test_get_account(host, port, '1001000000001')
        
        # Test 4: Get another account
        test_get_account(host, port, '2001000000001')
        
        # Test 5: Non-existent account
        test_get_account(host, port, '9999999999999')
        
        print "\n" + "=" * 60
        print "All tests completed!"
        print "=" * 60
        
    except Exception as e:
        print "\n✗ Test failed with error: {}".format(str(e))
        sys.exit(1)


if __name__ == '__main__':
    main()
