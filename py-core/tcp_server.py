#!/usr/bin/env python
# -*- coding: utf-8 -*-

"""
py-core TCP Server - Persistent Connection Handler

Implements TCP server with persistent connections following the nbbl2 pattern.
Handles multiple concurrent requests on a single connection using message ID correlation.
"""

import socket
import threading
import json
import struct
import logging
from datetime import datetime

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)


class TCPServer(object):
    """TCP Server with persistent connections"""
    
    def __init__(self, host, port, config, request_handler):
        """
        Initialize TCP server
        
        Args:
            host (str): Host to bind to
            port (int): Port to listen on
            config (dict): Configuration dictionary
            request_handler (RequestHandler): Handler for processing requests
        """
        self.host = host
        self.port = port
        self.config = config
        self.request_handler = request_handler
        self.active_connections = 0
        self.server_socket = None
        self.is_running = False
    
    def start(self):
        """Start the TCP server"""
        try:
            self.server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
            self.server_socket.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
            self.server_socket.bind((self.host, self.port))
            self.server_socket.listen(5)
            self.is_running = True
            
            logger.info("=" * 60)
            logger.info("py-core TCP server starting...")
            logger.info("Host: {}".format(self.host))
            logger.info("Port: {}".format(self.port))
            logger.info("Protocol: TCP (persistent connections)")
            logger.info("Python version: 2.7 (deprecated)")
            logger.info("Database: {}@{}:{}".format(
                self.config['database']['database'],
                self.config['database']['host'],
                self.config['database']['port']
            ))
            logger.info("=" * 60)
            
            # Accept connections loop
            while self.is_running:
                try:
                    client_socket, client_address = self.server_socket.accept()
                    self.active_connections += 1
                    
                    logger.info("New connection from {} - Active connections: {}".format(
                        client_address, self.active_connections
                    ))
                    
                    # Handle connection in a separate thread
                    connection_thread = threading.Thread(
                        target=self.handle_connection,
                        args=(client_socket, client_address)
                    )
                    connection_thread.daemon = True
                    connection_thread.start()
                    
                except socket.error as e:
                    if self.is_running:
                        logger.error("Error accepting connection: {}".format(str(e)))
                        
        except Exception as e:
            logger.error("Failed to start TCP server: {}".format(str(e)))
            raise
    
    def stop(self):
        """Stop the TCP server"""
        logger.info("Stopping TCP server...")
        self.is_running = False
        if self.server_socket:
            self.server_socket.close()
    
    def handle_connection(self, client_socket, client_address):
        """
        Handle a single client connection with persistent connection pattern
        
        Args:
            client_socket: Client socket connection
            client_address: Client address tuple (host, port)
        """
        conn_start_time = datetime.now()
        
        try:
            # Main read loop - keep connection alive
            while True:
                try:
                    # Read request length (2 bytes, big-endian)
                    length_bytes = self.receive_exact(client_socket, 2)
                    if not length_bytes:
                        logger.info("Client {} disconnected".format(client_address))
                        break
                    
                    # Unpack length
                    request_length = struct.unpack('!H', length_bytes)[0]
                    
                    # Sanity check
                    if request_length <= 0 or request_length > 65535:
                        logger.error("Invalid request length: {} from {}".format(
                            request_length, client_address
                        ))
                        break
                    
                    # Read JSON data
                    json_data = self.receive_exact(client_socket, request_length)
                    if not json_data:
                        logger.error("Failed to read request data from {}".format(
                            client_address
                        ))
                        break
                    
                    logger.info("Received request from {}: {}".format(
                        client_address, json_data
                    ))
                    
                    # Process request in a separate thread to allow concurrent processing
                    request_thread = threading.Thread(
                        target=self.process_request,
                        args=(client_socket, client_address, json_data)
                    )
                    request_thread.daemon = True
                    request_thread.start()
                    
                except socket.error as e:
                    logger.error("Socket error reading from {}: {}".format(
                        client_address, str(e)
                    ))
                    break
                except Exception as e:
                    logger.error("Error handling request from {}: {}".format(
                        client_address, str(e)
                    ))
                    # Don't break - try to continue processing other requests
                    
        finally:
            # Connection cleanup
            client_socket.close()
            self.active_connections -= 1
            duration = (datetime.now() - conn_start_time).total_seconds()
            
            logger.info("Connection closed: {} - Duration: {:.2f}s - Active: {}".format(
                client_address, duration, self.active_connections
            ))
    
    def process_request(self, client_socket, client_address, json_data):
        """
        Process a single request and send response
        
        Args:
            client_socket: Client socket
            client_address: Client address
            json_data: Raw JSON request data
        """
        try:
            # Parse JSON request
            request = json.loads(json_data)
            
            # Extract message ID
            message_id = request.get('id_message') or request.get('message_id')
            if not message_id:
                logger.error("Request missing id_message from {}".format(client_address))
                response = {
                    'status': '999',
                    'err_info': 'missing id_message or message_id'
                }
                self.send_response(client_socket, response)
                return
            
            # Process request via handler
            response = self.request_handler.handle_request(request)
            
            # Ensure message ID is in response
            if 'id_message' not in response and 'message_id' not in response:
                response['id_message'] = message_id
            
            # Send response
            self.send_response(client_socket, response)
            
            logger.info("Sent response to {} for message {}".format(
                client_address, message_id
            ))
            
        except ValueError as e:
            logger.error("Invalid JSON from {}: {}".format(client_address, str(e)))
            response = {
                'status': '999',
                'err_info': 'invalid JSON format'
            }
            self.send_response(client_socket, response)
        except Exception as e:
            logger.error("Error processing request from {}: {}".format(
                client_address, str(e)
            ))
            response = {
                'status': '999',
                'err_info': 'internal server error: {}'.format(str(e))
            }
            self.send_response(client_socket, response)
    
    def send_response(self, client_socket, response):
        """
        Send response to client
        
        Args:
            client_socket: Client socket
            response: Response dictionary
        """
        try:
            # Convert response to JSON
            json_response = json.dumps(response)
            response_bytes = json_response.encode('utf-8')
            
            # Check length
            response_length = len(response_bytes)
            if response_length > 65535:
                logger.error("Response too large: {} bytes".format(response_length))
                return
            
            # Pack length (2 bytes, big-endian)
            length_bytes = struct.pack('!H', response_length)
            
            # Send length + data
            client_socket.sendall(length_bytes + response_bytes)
            
            logger.debug("Sent response: {}".format(json_response))
            
        except socket.error as e:
            logger.error("Failed to send response: {}".format(str(e)))
        except Exception as e:
            logger.error("Error sending response: {}".format(str(e)))
    
    def receive_exact(self, client_socket, num_bytes):
        """
        Receive exact number of bytes from socket
        
        Args:
            client_socket: Client socket
            num_bytes: Number of bytes to receive
            
        Returns:
            bytes: Received data or None if connection closed
        """
        data = b''
        while len(data) < num_bytes:
            chunk = client_socket.recv(num_bytes - len(data))
            if not chunk:
                return None
            data += chunk
        return data

