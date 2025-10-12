#!/usr/bin/env python
# -*- coding: utf-8 -*-

"""
py-core - Balance Inquiry Service (Python 2.7)

A simple REST API service for account balance inquiry.
Intentionally uses Python 2.7 without connection pooling to demonstrate
performance differences with modern Go implementation.
"""

import json
import os
import psycopg2
from flask import Flask, jsonify, request
from datetime import datetime

app = Flask(__name__)

# Load configuration
def load_config():
    config_path = os.path.join(os.path.dirname(__file__), 'config.json')
    with open(config_path, 'r') as f:
        return json.load(f)

config = load_config()

# Database connection helper (no pooling - new connection per request)
def get_db_connection():
    """
    Creates a new database connection for each request.
    This is intentionally inefficient to demonstrate the performance
    difference with connection pooling.
    """
    conn = psycopg2.connect(
        host=config['database']['host'],
        port=config['database']['port'],
        database=config['database']['database'],
        user=config['database']['user'],
        password=config['database']['password']
    )
    return conn

# Helper to format date/datetime fields
def format_date(date_obj):
    if date_obj is None:
        return ""
    if isinstance(date_obj, datetime):
        return date_obj.strftime('%Y-%m-%dT%H:%M:%SZ')
    return date_obj.strftime('%Y-%m-%d')

# Health check endpoint
@app.route('/health', methods=['GET'])
def health_check():
    return jsonify({"status": "ok", "service": "py-core"}), 200

# Get account by account number
@app.route('/api/v1/accounts/<account_number>', methods=['GET'])
def get_account_by_account_number(account_number):
    """
    Retrieve account information with customer details by account number.

    Response:
    {
        "account": {
            "account_id": 1,
            "account_number": "1001000000001",
            "customer_id": 1,
            "account_type": "WADIAH",
            "account_status": "ACTIVE",
            "balance": "15750000.00",
            "currency": "IDR",
            "opened_date": "2025-01-01",
            "closed_date": "",
            "created_at": "2025-01-01T00:00:00Z",
            "updated_at": "2025-01-01T00:00:00Z"
        },
        "customer": {
            "customer_number": "CUST0000001",
            "full_name": "Ahmad Hidayat",
            "id_number": "3201012345670001",
            "phone_number": "081234567801",
            "email": "ahmad.hidayat@email.com",
            "address": "Jl. Merdeka No. 123, Jakarta",
            "date_of_birth": "1985-03-15"
        }
    }
    """
    conn = None
    cursor = None

    try:
        # Create new connection (no pooling!)
        conn = get_db_connection()
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
            return jsonify({
                "error": {
                    "code": "NOT_FOUND",
                    "message": "Account not found"
                }
            }), 404

        # Map row to response
        response = {
            "account": {
                "account_id": row[0],
                "account_number": row[1],
                "customer_id": row[2],
                "account_type": row[3],
                "account_status": row[4],
                "balance": str(row[5]),
                "currency": row[6],
                "opened_date": format_date(row[7]),
                "closed_date": format_date(row[8]),
                "created_at": format_date(row[9]),
                "updated_at": format_date(row[10])
            },
            "customer": {
                "customer_number": row[11],
                "full_name": row[12],
                "id_number": row[13],
                "phone_number": row[14] or "",
                "email": row[15] or "",
                "address": row[16] or "",
                "date_of_birth": format_date(row[17])
            }
        }

        return jsonify(response), 200

    except psycopg2.Error as e:
        return jsonify({
            "error": {
                "code": "DATABASE_ERROR",
                "message": str(e)
            }
        }), 500

    except Exception as e:
        return jsonify({
            "error": {
                "code": "INTERNAL_ERROR",
                "message": str(e)
            }
        }), 500

    finally:
        # Always close connection (no pooling!)
        if cursor:
            cursor.close()
        if conn:
            conn.close()

# Error handlers
@app.errorhandler(404)
def not_found(error):
    return jsonify({
        "error": {
            "code": "NOT_FOUND",
            "message": "Endpoint not found"
        }
    }), 404

@app.errorhandler(500)
def internal_error(error):
    return jsonify({
        "error": {
            "code": "INTERNAL_ERROR",
            "message": "Internal server error"
        }
    }), 500

if __name__ == '__main__':
    port = config.get('port', 5000)
    host = config.get('host', '0.0.0.0')

    print "=" * 50
    print "py-core service starting..."
    print "Host: {}".format(host)
    print "Port: {}".format(port)
    print "Python version: 2.7 (deprecated)"
    print "Database: {}@{}:{}".format(
        config['database']['database'],
        config['database']['host'],
        config['database']['port']
    )
    print "=" * 50

    app.run(host=host, port=port, debug=False)
