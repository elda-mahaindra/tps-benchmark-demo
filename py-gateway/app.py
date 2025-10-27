#!/usr/bin/env python
# -*- coding: utf-8 -*-

"""
py-gateway REST API Server

Flask-based REST API server that acts as the gateway layer.
Handles REST requests and forwards them to py-switching via TCP.
"""

import json
import os
import sys
import signal
import logging
import threading
from datetime import datetime
from flask import Flask, request, jsonify
from request_handler import RequestHandler
from tcp_client import PySwitchingTCPClient
from tcp_server import PyGatewayTCPServer

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

# Global instances for signal handling
app = None
request_handler = None
py_switching_client = None
tcp_server = None


def load_config():
    """Load configuration from config.json"""
    config_path = os.path.join(os.path.dirname(__file__), 'config.json')
    with open(config_path, 'r') as f:
        return json.load(f)


def signal_handler(signum, frame):
    """Handle shutdown signals gracefully"""
    logger.info("Received signal {}, shutting down...".format(signum))
    if tcp_server:
        tcp_server.stop()
    if py_switching_client:
        py_switching_client.disconnect()
    sys.exit(0)


def create_app():
    """Create Flask application"""
    global app, request_handler, py_switching_client, tcp_server
    
    # Load configuration
    config = load_config()
    
    # Get server configuration
    rest_host = config.get('rest_host', config.get('host', '0.0.0.0'))
    rest_port = config.get('rest_port', config.get('port', 8083))
    tcp_host = config.get('tcp_host', '0.0.0.0')
    tcp_port = config.get('tcp_port', 8084)
    
    # Get py-switching configuration
    py_switching_config = config['external_service']['py_switching']
    py_switching_host = py_switching_config['host']
    py_switching_port = py_switching_config['port']
    
    # Initialize py-switching TCP client
    py_switching_client = PySwitchingTCPClient(py_switching_host, py_switching_port, config)
    
    # Initialize request handler
    request_handler = RequestHandler(py_switching_client)
    
    # Initialize TCP server
    tcp_server = PyGatewayTCPServer(tcp_host, tcp_port, request_handler)
    
    # Create Flask app
    app = Flask(__name__)
    
    # Setup signal handlers for graceful shutdown
    signal.signal(signal.SIGINT, signal_handler)
    signal.signal(signal.SIGTERM, signal_handler)
    
    # Define routes
    @app.route('/health', methods=['GET'])
    def health():
        """Health check endpoint"""
        return jsonify({
            'status': 'healthy',
            'service': 'py-gateway',
            'timestamp': datetime.now().strftime('%Y-%m-%dT%H:%M:%SZ')
        })
    
    @app.route('/api/v1/accounts', methods=['GET'])
    def get_account():
        """Get account by account number"""
        try:
            # Get account number from query parameters
            account_number = request.args.get('account_number')
            if not account_number:
                return jsonify({
                    'error': True,
                    'message': 'account_number parameter is required'
                }), 400
            
            # Handle request
            result = request_handler.handle_get_account_by_account_number(account_number)
            
            if result.get('error'):
                return jsonify({
                    'error': True,
                    'message': result.get('message', 'Unknown error'),
                    'status': result.get('status', '999')
                }), 500
            
            return jsonify({
                'error': False,
                'data': result.get('data', {})
            })
            
        except Exception as e:
            logger.error("Error in get_account endpoint: {}".format(str(e)))
            return jsonify({
                'error': True,
                'message': 'Internal server error: {}'.format(str(e))
            }), 500
    
    @app.route('/ping', methods=['GET'])
    def ping():
        """Ping endpoint"""
        try:
            # Handle request
            result = request_handler.handle_ping()
            
            if result.get('error'):
                return jsonify({
                    'error': True,
                    'message': result.get('message', 'Unknown error'),
                    'status': result.get('status', '999')
                }), 500
            
            return jsonify({
                'error': False,
                'data': result.get('data', {})
            })
            
        except Exception as e:
            logger.error("Error in ping endpoint: {}".format(str(e)))
            return jsonify({
                'error': True,
                'message': 'Internal server error: {}'.format(str(e))
            }), 500
    
    logger.info("=" * 60)
    logger.info("py-gateway REST API server starting...")
    logger.info("REST Host: {}".format(rest_host))
    logger.info("REST Port: {}".format(rest_port))
    logger.info("TCP Host: {}".format(tcp_host))
    logger.info("TCP Port: {}".format(tcp_port))
    logger.info("Protocol: REST API + TCP")
    logger.info("Python version: 2.7 (deprecated)")
    logger.info("Target py-switching: {}:{}".format(py_switching_host, py_switching_port))
    logger.info("=" * 60)
    
    return app, rest_host, rest_port


def main():
    """Main entry point"""
    try:
        app, rest_host, rest_port = create_app()
        
        # Start TCP server in a separate thread
        tcp_thread = threading.Thread(target=tcp_server.start)
        tcp_thread.daemon = True
        tcp_thread.start()
        
        # Start Flask server
        app.run(host=rest_host, port=rest_port, debug=False, threaded=True)
        
    except KeyboardInterrupt:
        logger.info("Received keyboard interrupt, shutting down...")
        if tcp_server:
            tcp_server.stop()
        if py_switching_client:
            py_switching_client.disconnect()
        sys.exit(0)
        
    except Exception as e:
        logger.error("Fatal error: {}".format(str(e)))
        sys.exit(1)


if __name__ == '__main__':
    main()
