#!/usr/bin/env python
# -*- coding: utf-8 -*-

"""
Test TCP Client - For testing py-core TCP server

Simple client to test the TCP server implementation.
"""

import socket
import json
import struct
import sys
import uuid


def send_tcp_request(host, port, request_data):
    """
    Send TCP request and receive response
    
    Args:
        host (str): Server host
        port (int): Server port
        request_data (dict): Request dictionary
        
    Returns:
        dict: Response dictionary
    """
    # Create socket
    client_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    
    try:
        # Connect to server
        print "Connecting to {}:{}...".format(host, port)
        client_socket.connect((host, port))
        print "Connected!"
        
        # Convert request to JSON
        json_request = json.dumps(request_data)
        request_bytes = json_request.encode('utf-8')
        
        # Get length
        request_length = len(request_bytes)
        if request_length > 65535:
            raise Exception("Request too large: {} bytes".format(request_length))
        
        # Pack length (2 bytes, big-endian)
        length_bytes = struct.pack('!H', request_length)
        
        # Send length + data
        print "Sending request ({} bytes): {}".format(request_length, json_request)
        client_socket.sendall(length_bytes + request_bytes)
        
        # Read response length (2 bytes)
        response_length_bytes = recv_exact(client_socket, 2)
        if not response_length_bytes:
            raise Exception("Failed to read response length")
        
        response_length = struct.unpack('!H', response_length_bytes)[0]
        print "Response length: {} bytes".format(response_length)
        
        # Read response data
        response_data = recv_exact(client_socket, response_length)
        if not response_data:
            raise Exception("Failed to read response data")
        
        # Parse JSON response
        response = json.loads(response_data.decode('utf-8'))
        print "Received response: {}".format(json.dumps(response, indent=2))
        
        return response
        
    finally:
        client_socket.close()
        print "Connection closed"


def recv_exact(sock, num_bytes):
    """
    Receive exact number of bytes from socket
    
    Args:
        sock: Socket to read from
        num_bytes: Number of bytes to receive
        
    Returns:
        bytes: Received data or None if connection closed
    """
    data = b''
    while len(data) < num_bytes:
        chunk = sock.recv(num_bytes - len(data))
        if not chunk:
            return None
        data += chunk
    return data


def test_ping(host, port):
    """Test ping operation"""
    print "\n" + "=" * 60
    print "TEST: Ping"
    print "=" * 60
    
    request = {
        'id_message': str(uuid.uuid4()),
        'operation': 'ping',
        'params': {}
    }
    
    response = send_tcp_request(host, port, request)
    
    if response['status'] == '000':
        print "✓ Ping successful!"
        print "  Data: {}".format(response.get('data', {}))
    else:
        print "✗ Ping failed!"
        print "  Error: {}".format(response.get('err_info', 'unknown'))


def test_get_account(host, port, account_number):
    """Test get account operation"""
    print "\n" + "=" * 60
    print "TEST: Get Account by Account Number"
    print "=" * 60
    
    request = {
        'id_message': str(uuid.uuid4()),
        'operation': 'get_account_by_account_number',
        'params': {
            'account_number': account_number
        }
    }
    
    response = send_tcp_request(host, port, request)
    
    if response['status'] == '000':
        print "✓ Account retrieval successful!"
        data = response.get('data', {})
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
        print "  Error: {}".format(response.get('err_info', 'unknown'))


def main():
    """Main test function"""
    # Configuration
    host = 'localhost'
    port = 5001
    
    if len(sys.argv) > 1:
        host = sys.argv[1]
    if len(sys.argv) > 2:
        port = int(sys.argv[2])
    
    print "=" * 60
    print "py-core TCP Client Test"
    print "=" * 60
    print "Target: {}:{}".format(host, port)
    print "=" * 60
    
    try:
        # Test 1: Ping
        test_ping(host, port)
        
        # Test 2: Get Account
        test_get_account(host, port, '1001000000001')
        
        # Test 3: Get another account
        test_get_account(host, port, '2001000000001')
        
        # Test 4: Non-existent account
        test_get_account(host, port, '9999999999999')
        
        print "\n" + "=" * 60
        print "All tests completed!"
        print "=" * 60
        
    except Exception as e:
        print "\n✗ Test failed with error: {}".format(str(e))
        sys.exit(1)


if __name__ == '__main__':
    main()

