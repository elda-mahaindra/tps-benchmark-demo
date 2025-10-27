#!/usr/bin/env python
# -*- coding: utf-8 -*-

"""
Request Handler - Business Logic for py-switching TCP Server

Handles incoming requests and delegates to py-core service.
Implements the switching layer functionality similar to go-switching.
"""

import logging
from datetime import datetime
from tcp_client import PyCoreTCPClient

logger = logging.getLogger(__name__)


class RequestHandler(object):
    """Handles request processing and routing to py-core"""
    
    def __init__(self, py_core_client):
        """
        Initialize request handler
        
        Args:
            py_core_client (PyCoreTCPClient): TCP client for py-core communication
        """
        self.py_core_client = py_core_client
    
    def handle_request(self, request):
        """
        Handle incoming request and return response
        
        Args:
            request (dict): Request dictionary with operation and params
            
        Returns:
            dict: Response dictionary with status, err_info, and data
        """
        try:
            # Extract operation
            operation = request.get('operation')
            if not operation:
                return self.create_error_response(
                    request,
                    '999',
                    'missing operation field'
                )
            
            # Extract message ID for response
            message_id = request.get('id_message') or request.get('message_id')
            
            logger.info("Processing operation: {} for message: {}".format(
                operation, message_id
            ))
            
            # Route to appropriate handler
            if operation == 'get_account_by_account_number':
                return self.handle_get_account_by_account_number(request)
            elif operation == 'ping':
                return self.handle_ping(request)
            else:
                return self.create_error_response(
                    request,
                    '999',
                    'unknown operation: {}'.format(operation)
                )
                
        except Exception as e:
            logger.error("Error handling request: {}".format(str(e)))
            return self.create_error_response(
                request,
                '999',
                'internal error: {}'.format(str(e))
            )
    
    def handle_get_account_by_account_number(self, request):
        """
        Handle get account by account number operation
        
        Args:
            request (dict): Request with params containing account_number
            
        Returns:
            dict: Response with account and customer data
        """
        try:
            # Extract parameters
            params = request.get('params', {})
            account_number = params.get('account_number')
            
            if not account_number:
                return self.create_error_response(
                    request,
                    '999',
                    'missing account_number in params'
                )
            
            logger.info("Getting account: {}".format(account_number))
            
            # Forward request to py-core
            py_core_request = {
                'id_message': request.get('id_message'),
                'operation': 'get_account_by_account_number',
                'params': {
                    'account_number': account_number
                }
            }
            
            # Send request to py-core
            py_core_response = self.py_core_client.send_request(py_core_request)
            
            # Check if py-core returned an error
            if py_core_response.get('status') != '000':
                return self.create_error_response(
                    request,
                    py_core_response.get('status', '999'),
                    py_core_response.get('err_info', 'py-core error')
                )
            
            # Extract data from py-core response
            data = py_core_response.get('data', {})
            
            # Create success response
            response = self.create_success_response(request, data)
            
            logger.info("Successfully retrieved account: {}".format(account_number))
            
            return response
            
        except Exception as e:
            logger.error("Error in get_account_by_account_number: {}".format(str(e)))
            return self.create_error_response(
                request,
                '999',
                'py-core communication error: {}'.format(str(e))
            )
    
    def handle_ping(self, request):
        """
        Handle ping operation
        
        Args:
            request (dict): Request dictionary
            
        Returns:
            dict: Pong response
        """
        response_data = {
            'message': 'pong',
            'service': 'py-switching',
            'timestamp': datetime.now().strftime('%Y-%m-%dT%H:%M:%SZ')
        }
        
        return self.create_success_response(request, response_data)
    
    def create_success_response(self, request, data):
        """
        Create success response
        
        Args:
            request (dict): Original request
            data (dict): Response data
            
        Returns:
            dict: Success response
        """
        message_id = request.get('id_message') or request.get('message_id')
        
        response = {
            'id_message': message_id,
            'status': '000',
            'err_info': '',
            'data': data
        }
        
        return response
    
    def create_error_response(self, request, status, error_message):
        """
        Create error response
        
        Args:
            request (dict): Original request
            status (str): Error status code (999=error, 998=timeout, 997=dropped)
            error_message (str): Error message
            
        Returns:
            dict: Error response
        """
        message_id = request.get('id_message') or request.get('message_id', '')
        
        response = {
            'id_message': message_id,
            'status': status,
            'err_info': error_message
        }
        
        return response
