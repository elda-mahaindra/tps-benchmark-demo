#!/usr/bin/env python
# -*- coding: utf-8 -*-

"""
py-core - TCP Server Entry Point

Main entry point for py-core TCP server.
Implements persistent connection TCP protocol with generic request/response format.
"""

import json
import os
import sys
import signal
import logging

from tcp_server import TCPServer
from request_handler import RequestHandler
from database import DatabaseService

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)


# Global server instance for signal handling
server_instance = None


def load_config():
    """Load configuration from config.json"""
    config_path = os.path.join(os.path.dirname(__file__), 'config.json')
    with open(config_path, 'r') as f:
        return json.load(f)


def signal_handler(signum, frame):
    """Handle shutdown signals gracefully"""
    logger.info("Received signal {}, shutting down...".format(signum))
    if server_instance:
        server_instance.stop()
    sys.exit(0)


def main():
    """Main entry point"""
    global server_instance
    
    try:
        # Load configuration
        config = load_config()
        
        # Get server configuration
        host = config.get('tcp_host', config.get('host', '0.0.0.0'))
        port = config.get('tcp_port', config.get('port', 5001))
        
        # Initialize database service
        db_service = DatabaseService(config['database'])
        
        # Initialize request handler
        request_handler = RequestHandler(db_service)
        
        # Initialize TCP server
        server_instance = TCPServer(host, port, config, request_handler)
        
        # Setup signal handlers for graceful shutdown
        signal.signal(signal.SIGINT, signal_handler)
        signal.signal(signal.SIGTERM, signal_handler)
        
        # Start server (blocking call)
        server_instance.start()
        
    except KeyboardInterrupt:
        logger.info("Received keyboard interrupt, shutting down...")
        if server_instance:
            server_instance.stop()
        sys.exit(0)
        
    except Exception as e:
        logger.error("Fatal error: {}".format(str(e)))
        sys.exit(1)


if __name__ == '__main__':
    main()

