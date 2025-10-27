#!/usr/bin/env python
# -*- coding: utf-8 -*-

"""
Database Service - Database Operations

Handles all database interactions with PostgreSQL.
Intentionally creates new connection per request (no pooling) for demo purposes.
"""

import psycopg2
import logging
from datetime import datetime

logger = logging.getLogger(__name__)


class DatabaseService(object):
    """Database service for account operations"""
    
    def __init__(self, config):
        """
        Initialize database service
        
        Args:
            config (dict): Database configuration
        """
        self.config = config
    
    def get_db_connection(self):
        """
        Create new database connection (no pooling!)
        
        This is intentionally inefficient to demonstrate the performance
        difference with connection pooling in Go implementation.
        
        Returns:
            connection: PostgreSQL connection
        """
        conn = psycopg2.connect(
            host=self.config['host'],
            port=self.config['port'],
            database=self.config['database'],
            user=self.config['user'],
            password=self.config['password']
        )
        return conn
    
    def format_date(self, date_obj):
        """
        Format date/datetime objects to string
        
        Args:
            date_obj: Date or datetime object
            
        Returns:
            str: Formatted date string
        """
        if date_obj is None:
            return ""
        if isinstance(date_obj, datetime):
            return date_obj.strftime('%Y-%m-%dT%H:%M:%SZ')
        return date_obj.strftime('%Y-%m-%d')
    
    def get_account_by_account_number(self, account_number):
        """
        Retrieve account information with customer details by account number
        
        Args:
            account_number (str): Account number to look up
            
        Returns:
            dict: Dictionary containing account and customer information
            None: If account not found
        """
        conn = None
        cursor = None
        
        try:
            # Create new connection (no pooling!)
            logger.debug("Creating new database connection for account: {}".format(
                account_number
            ))
            conn = self.get_db_connection()
            cursor = conn.cursor()
            
            # Execute query with JOIN
            query = """
                SELECT
                    a.account_id,
                    a.account_number,
                    a.customer_id,
                    a.account_type,
                    a.account_status,
                    a.balance,
                    a.currency,
                    a.opened_date,
                    a.closed_date,
                    a.created_at,
                    a.updated_at,
                    c.customer_number,
                    c.full_name,
                    c.id_number,
                    c.phone_number,
                    c.email,
                    c.address,
                    c.date_of_birth
                FROM demo.accounts a
                INNER JOIN demo.customers c ON a.customer_id = c.customer_id
                WHERE a.account_number = %s
            """
            
            cursor.execute(query, (account_number,))
            row = cursor.fetchone()
            
            if row is None:
                logger.info("Account not found: {}".format(account_number))
                return None
            
            # Map row to response structure
            result = {
                "account": {
                    "account_id": row[0],
                    "account_number": row[1],
                    "customer_id": row[2],
                    "account_type": row[3],
                    "account_status": row[4],
                    "balance": str(row[5]),
                    "currency": row[6],
                    "opened_date": self.format_date(row[7]),
                    "closed_date": self.format_date(row[8]),
                    "created_at": self.format_date(row[9]),
                    "updated_at": self.format_date(row[10])
                },
                "customer": {
                    "customer_number": row[11],
                    "full_name": row[12],
                    "id_number": row[13],
                    "phone_number": row[14] or "",
                    "email": row[15] or "",
                    "address": row[16] or "",
                    "date_of_birth": self.format_date(row[17])
                }
            }
            
            logger.info("Successfully retrieved account: {}".format(account_number))
            
            return result
            
        except psycopg2.Error as e:
            logger.error("Database error: {}".format(str(e)))
            raise Exception("Database error: {}".format(str(e)))
            
        except Exception as e:
            logger.error("Error getting account: {}".format(str(e)))
            raise
            
        finally:
            # Always close connection (no pooling!)
            if cursor:
                cursor.close()
            if conn:
                conn.close()
            logger.debug("Database connection closed for account: {}".format(
                account_number
            ))

