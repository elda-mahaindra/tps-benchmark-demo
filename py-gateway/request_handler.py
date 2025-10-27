#!/usr/bin/env python
# -*- coding: utf-8 -*-

"""
Request Handler - Business Logic for py-gateway REST API

Handles incoming REST requests and delegates to py-switching service.
Implements the gateway layer functionality similar to go-gateway.
"""

import logging
import time
from datetime import datetime
from tcp_client import PySwitchingTCPClient

logger = logging.getLogger(__name__)


class RequestHandler(object):
    """Handles request processing and routing to py-switching"""
    
    def __init__(self, py_switching_client):
        """
        Initialize request handler
        
        Args:
            py_switching_client (PySwitchingTCPClient): TCP client for py-switching communication
        """
        self.py_switching_client = py_switching_client
    
    def handle_get_account_by_account_number(self, account_number):
        """
        Handle get account by account number request
        
        Args:
            account_number (str): Account number to retrieve
            
        Returns:
            dict: Response dictionary with account and customer data
        """
        try:
            logger.info("Getting account: {}".format(account_number))
            
            # Create request for py-switching
            request = {
                'id_message': str(int(time.time() * 1000000)),
                'operation': 'get_account_by_account_number',
                'params': {
                    'account_number': account_number
                }
            }
            
            # Send request to py-switching
            response = self.py_switching_client.send_request(request)
            
            # Check if py-switching returned an error
            if response.get('status') != '000':
                return {
                    'error': True,
                    'status': response.get('status', '999'),
                    'message': response.get('err_info', 'py-switching error')
                }
            
            # Extract data from py-switching response
            data = response.get('data', {})
            
            # Create success response
            result = {
                'error': False,
                'data': data
            }
            
            logger.info("Successfully retrieved account: {}".format(account_number))
            
            return result
            
        except Exception as e:
            logger.error("Error in get_account_by_account_number: {}".format(str(e)))
            return {
                'error': True,
                'status': '999',
                'message': 'py-switching communication error: {}'.format(str(e))
            }
    
    def handle_ping(self):
        """
        Handle ping request
        
        Returns:
            dict: Pong response
        """
        try:
            # Create request for py-switching
            request = {
                'id_message': str(int(time.time() * 1000000)),
                'operation': 'ping',
                'params': {}
            }
            
            # Send request to py-switching
            response = self.py_switching_client.send_request(request)
            
            # Check if py-switching returned an error
            if response.get('status') != '000':
                return {
                    'error': True,
                    'status': response.get('status', '999'),
                    'message': response.get('err_info', 'py-switching error')
                }
            
            # Extract data from py-switching response
            data = response.get('data', {})
            
            # Create success response
            result = {
                'error': False,
                'data': data
            }
            
            return result
            
        except Exception as e:
            logger.error("Error in ping: {}".format(str(e)))
            return {
                'error': True,
                'status': '999',
                'message': 'py-switching communication error: {}'.format(str(e))
            }
