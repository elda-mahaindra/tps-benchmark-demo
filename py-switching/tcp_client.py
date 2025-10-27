#!/usr/bin/env python
# -*- coding: utf-8 -*-

"""
TCP Client for py-core Communication

Handles TCP communication with py-core service using the same protocol
as py-core TCP server (length-prefixed JSON messages).
"""

import socket
import json
import struct
import logging
import threading
import time
from datetime import datetime

logger = logging.getLogger(__name__)


class PyCoreTCPClient(object):
    """TCP client for communicating with py-core service"""
    
    def __init__(self, host, port, config):
        """
        Initialize TCP client
        
        Args:
            host (str): py-core host
            port (int): py-core port
            config (dict): Configuration dictionary
        """
        self.host = host
        self.port = port
        self.config = config
        self.conn = None
        self.conn_mutex = threading.Lock()
        self.is_connected = False
        self.pending_requests = {}
        self.requests_mutex = threading.Lock()
        self.response_thread = None
        self.shutdown_event = threading.Event()
    
    def connect(self):
        """Establish connection to py-core"""
        try:
            self.conn = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
            self.conn.connect((self.host, self.port))
            self.is_connected = True
            
            # Start response reading thread
            self.response_thread = threading.Thread(target=self.read_responses)
            self.response_thread.daemon = True
            self.response_thread.start()
            
            logger.info("Connected to py-core at {}:{}".format(self.host, self.port))
            return True
            
        except Exception as e:
            logger.error("Failed to connect to py-core: {}".format(str(e)))
            self.is_connected = False
            return False
    
    def disconnect(self):
        """Disconnect from py-core"""
        self.shutdown_event.set()
        self.is_connected = False
        
        if self.conn:
            try:
                self.conn.close()
            except:
                pass
            self.conn = None
        
        if self.response_thread and self.response_thread.is_alive():
            self.response_thread.join(timeout=1)
        
        logger.info("Disconnected from py-core")
    
    def send_request(self, request_data, timeout=30):
        """
        Send request to py-core and wait for response
        
        Args:
            request_data (dict): Request dictionary
            timeout (int): Timeout in seconds
            
        Returns:
            dict: Response dictionary
        """
        if not self.is_connected:
            if not self.connect():
                raise Exception("Not connected to py-core")
        
        # Generate unique message ID if not present
        message_id = request_data.get('id_message')
        if not message_id:
            message_id = str(int(time.time() * 1000000))
            request_data['id_message'] = message_id
        
        # Create response channel
        response_channel = threading.Event()
        response_data = [None]  # Use list to allow modification in nested function
        
        # Store pending request
        with self.requests_mutex:
            self.pending_requests[message_id] = (response_channel, response_data)
        
        try:
            # Send request
            self._send_message(request_data)
            
            # Wait for response
            if response_channel.wait(timeout):
                return response_data[0]
            else:
                # Timeout - clean up
                with self.requests_mutex:
                    if message_id in self.pending_requests:
                        del self.pending_requests[message_id]
                raise Exception("Request timeout after {} seconds".format(timeout))
                
        except Exception as e:
            # Clean up on error
            with self.requests_mutex:
                if message_id in self.pending_requests:
                    del self.pending_requests[message_id]
            raise e
    
    def _send_message(self, data):
        """Send a message to py-core"""
        try:
            # Convert to JSON
            json_data = json.dumps(data)
            message_bytes = json_data.encode('utf-8')
            
            # Get length
            message_length = len(message_bytes)
            if message_length > 65535:
                raise Exception("Message too large: {} bytes".format(message_length))
            
            # Pack length (2 bytes, big-endian)
            length_bytes = struct.pack('!H', message_length)
            
            # Send length + data
            with self.conn_mutex:
                if self.conn:
                    self.conn.sendall(length_bytes + message_bytes)
            
            logger.debug("Sent request to py-core: {}".format(json_data))
            
        except Exception as e:
            logger.error("Failed to send message to py-core: {}".format(str(e)))
            raise e
    
    def read_responses(self):
        """Continuously read responses from py-core"""
        while not self.shutdown_event.is_set() and self.is_connected:
            try:
                if not self.conn:
                    time.sleep(0.1)
                    continue
                
                # Read response length (2 bytes)
                length_bytes = self._receive_exact(2)
                if not length_bytes:
                    logger.info("Connection closed by py-core")
                    break
                
                response_length = struct.unpack('!H', length_bytes)[0]
                
                # Read response data
                response_data = self._receive_exact(response_length)
                if not response_data:
                    logger.error("Failed to read response data")
                    break
                
                # Parse JSON response
                response = json.loads(response_data.decode('utf-8'))
                
                # Process response
                self._process_response(response)
                
            except Exception as e:
                logger.error("Error reading response from py-core: {}".format(str(e)))
                break
        
        # Clean up on disconnect
        self.is_connected = False
    
    def _receive_exact(self, num_bytes):
        """Receive exact number of bytes"""
        data = b''
        while len(data) < num_bytes:
            try:
                chunk = self.conn.recv(num_bytes - len(data))
                if not chunk:
                    return None
                data += chunk
            except socket.timeout:
                continue
            except Exception as e:
                logger.error("Error receiving data: {}".format(str(e)))
                return None
        return data
    
    def _process_response(self, response):
        """Process response from py-core"""
        try:
            # Extract message ID
            message_id = response.get('id_message')
            if not message_id:
                logger.error("Response missing id_message")
                return
            
            # Find waiting request
            with self.requests_mutex:
                if message_id in self.pending_requests:
                    response_channel, response_data = self.pending_requests[message_id]
                    del self.pending_requests[message_id]
                else:
                    logger.warn("No waiting request for message ID: {}".format(message_id))
                    return
            
            # Set response data and notify
            response_data[0] = response
            response_channel.set()
            
            logger.debug("Processed response for message ID: {}".format(message_id))
            
        except Exception as e:
            logger.error("Error processing response: {}".format(str(e)))
