#!/usr/bin/env python
# -*- coding: utf-8 -*-

"""
Request Handler - Business Logic for TCP Server

Handles incoming requests and delegates to appropriate service operations.
"""

import logging
from datetime import datetime
from database import DatabaseService

logger = logging.getLogger(__name__)


class RequestHandler(object):
    """Handles request processing and routing"""
    
    def __init__(self, db_service):
        """
        Initialize request handler
        
        Args:
            db_service (DatabaseService): Database service instance
        """
        self.db_service = db_service
    
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
            
            # Query database
            result = self.db_service.get_account_by_account_number(account_number)
            
            if not result:
                return self.create_error_response(
                    request,
                    '999',
                    'account not found'
                )
            
            # Create success response
            response = self.create_success_response(request, result)
            
            logger.info("Successfully retrieved account: {}".format(account_number))
            
            return response
            
        except Exception as e:
            logger.error("Error in get_account_by_account_number: {}".format(str(e)))
            return self.create_error_response(
                request,
                '999',
                'database error: {}'.format(str(e))
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
            'service': 'py-core',
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

