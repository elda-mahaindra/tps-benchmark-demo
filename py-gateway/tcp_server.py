#!/usr/bin/env python
# -*- coding: utf-8 -*-

"""
TCP Server for py-gateway

Handles TCP connections and forwards requests to py-switching.
This provides the TCP interface that the load-tester expects.
"""

import socket
import json
import struct
import threading
import logging
import time
from datetime import datetime

logger = logging.getLogger(__name__)


class PyGatewayTCPServer(object):
    """TCP server for py-gateway"""
    
    def __init__(self, host, port, request_handler):
        """
        Initialize TCP server
        
        Args:
            host (str): Server host
            port (int): Server port
            request_handler (RequestHandler): Request handler instance
        """
        self.host = host
        self.port = port
        self.request_handler = request_handler
        self.server_socket = None
        self.is_running = False
        self.active_connections = 0
        self.connections_lock = threading.Lock()
    
    def start(self):
        """Start the TCP server"""
        try:
            self.server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
            self.server_socket.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
            self.server_socket.bind((self.host, self.port))
            self.server_socket.listen(5)
            self.is_running = True
            
            logger.info("=" * 60)
            logger.info("py-gateway TCP server starting...")
            logger.info("Host: {}".format(self.host))
            logger.info("Port: {}".format(self.port))
            logger.info("Protocol: TCP (persistent connections)")
            logger.info("Python version: 2.7 (deprecated)")
            logger.info("=" * 60)
            
            logger.info("TCP server listening on {}:{}".format(self.host, self.port))
            
            while self.is_running:
                try:
                    client_socket, client_address = self.server_socket.accept()
                    
                    with self.connections_lock:
                        self.active_connections += 1
                    
                    logger.info("New connection from {} - Active connections: {}".format(
                        client_address, self.active_connections))
                    
                    # Handle client in a separate thread
                    client_thread = threading.Thread(
                        target=self.handle_client,
                        args=(client_socket, client_address)
                    )
                    client_thread.daemon = True
                    client_thread.start()
                    
                except socket.error as e:
                    if self.is_running:
                        logger.error("Error accepting connection: {}".format(str(e)))
                    break
                    
        except Exception as e:
            logger.error("Failed to start TCP server: {}".format(str(e)))
            raise e
    
    def stop(self):
        """Stop the TCP server"""
        logger.info("Stopping TCP server...")
        self.is_running = False
        
        if self.server_socket:
            try:
                self.server_socket.close()
            except:
                pass
        
        logger.info("TCP server stopped")
    
    def handle_client(self, client_socket, client_address):
        """Handle a client connection"""
        start_time = time.time()
        
        try:
            while self.is_running:
                # Read request length (2 bytes)
                length_bytes = self._receive_exact(client_socket, 2)
                if not length_bytes:
                    logger.info("Client {} disconnected".format(client_address))
                    break
                
                request_length = struct.unpack('!H', length_bytes)[0]
                
                # Read request data
                request_data = self._receive_exact(client_socket, request_length)
                if not request_data:
                    logger.error("Failed to read request data from {}".format(client_address))
                    break
                
                # Parse JSON request
                try:
                    request = json.loads(request_data.decode('utf-8'))
                except ValueError as e:
                    logger.error("Invalid JSON from {}: {}".format(client_address, str(e)))
                    continue
                
                logger.info("Received request from {}: {}".format(client_address, json.dumps(request)))
                
                # Process request
                response = self.process_request(request)
                
                # Send response
                self._send_response(client_socket, response, request.get('id_message', ''))
                
        except Exception as e:
            logger.error("Error handling client {}: {}".format(client_address, str(e)))
        
        finally:
            # Clean up
            try:
                client_socket.close()
            except:
                pass
            
            with self.connections_lock:
                self.active_connections -= 1
            
            duration = time.time() - start_time
            logger.info("Connection closed: {} - Duration: {:.2f}s - Active: {}".format(
                client_address, duration, self.active_connections))
    
    def process_request(self, request):
        """Process incoming request"""
        try:
            operation = request.get('operation')
            if not operation:
                return self.create_error_response(request, '999', 'missing operation field')
            
            message_id = request.get('id_message') or request.get('message_id')
            logger.info("Processing operation: {} for message: {}".format(operation, message_id))
            
            if operation == 'get_account_by_account_number':
                return self.handle_get_account_by_account_number(request)
            elif operation == 'ping':
                return self.handle_ping(request)
            else:
                return self.create_error_response(request, '999', 'unknown operation: {}'.format(operation))
                
        except Exception as e:
            logger.error("Error processing request: {}".format(str(e)))
            return self.create_error_response(request, '999', 'internal error: {}'.format(str(e)))
    
    def handle_get_account_by_account_number(self, request):
        """Handle get account by account number request"""
        try:
            params = request.get('params', {})
            account_number = params.get('account_number')
            
            if not account_number:
                return self.create_error_response(request, '999', 'missing account_number in params')
            
            logger.info("Getting account: {}".format(account_number))
            
            # Use request handler to get account
            result = self.request_handler.handle_get_account_by_account_number(account_number)
            
            if result.get('error'):
                return self.create_error_response(request, result.get('status', '999'), result.get('message', 'Unknown error'))
            
            return self.create_success_response(request, result.get('data', {}))
            
        except Exception as e:
            logger.error("Error in handle_get_account_by_account_number: {}".format(str(e)))
            return self.create_error_response(request, '999', 'internal error: {}'.format(str(e)))
    
    def handle_ping(self, request):
        """Handle ping request"""
        try:
            result = self.request_handler.handle_ping()
            
            if result.get('error'):
                return self.create_error_response(request, result.get('status', '999'), result.get('message', 'Unknown error'))
            
            return self.create_success_response(request, result.get('data', {}))
            
        except Exception as e:
            logger.error("Error in handle_ping: {}".format(str(e)))
            return self.create_error_response(request, '999', 'internal error: {}'.format(str(e)))
    
    def create_success_response(self, request, data):
        """Create success response"""
        message_id = request.get('id_message') or request.get('message_id')
        response = {
            'id_message': message_id,
            'status': '000',
            'err_info': '',
            'data': data
        }
        return response
    
    def create_error_response(self, request, status, error_message):
        """Create error response"""
        message_id = request.get('id_message') or request.get('message_id', '')
        response = {
            'id_message': message_id,
            'status': status,
            'err_info': error_message
        }
        return response
    
    def _receive_exact(self, sock, num_bytes):
        """Receive exact number of bytes"""
        data = b''
        while len(data) < num_bytes:
            try:
                chunk = sock.recv(num_bytes - len(data))
                if not chunk:
                    return None
                data += chunk
            except socket.timeout:
                continue
            except Exception as e:
                logger.error("Error receiving data: {}".format(str(e)))
                return None
        return data
    
    def _send_response(self, client_socket, response, message_id):
        """Send response to client"""
        try:
            # Convert to JSON
            json_response = json.dumps(response)
            response_bytes = json_response.encode('utf-8')
            
            # Get length
            response_length = len(response_bytes)
            if response_length > 65535:
                raise Exception("Response too large: {} bytes".format(response_length))
            
            # Pack length (2 bytes, big-endian)
            length_bytes = struct.pack('!H', response_length)
            
            # Send length + data
            client_socket.sendall(length_bytes + response_bytes)
            
            logger.info("Sent response to {} for message {}".format(client_socket.getpeername(), message_id))
            
        except Exception as e:
            logger.error("Failed to send response: {}".format(str(e)))
            raise e
