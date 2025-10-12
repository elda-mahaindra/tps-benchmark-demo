\c demo_db;

-- Table definitions

-- Schema: demo

-- Customers table
CREATE TABLE demo.customers (
    customer_id BIGSERIAL PRIMARY KEY,
    customer_number VARCHAR(20) UNIQUE NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    id_number VARCHAR(50) NOT NULL,
    phone_number VARCHAR(20),
    email VARCHAR(100),
    address TEXT,
    date_of_birth DATE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Accounts table for sharia banking
CREATE TABLE demo.accounts (
    account_id BIGSERIAL PRIMARY KEY,
    account_number VARCHAR(20) UNIQUE NOT NULL,
    customer_id BIGINT NOT NULL REFERENCES demo.customers(customer_id),
    account_type VARCHAR(50) NOT NULL, -- e.g., 'WADIAH', 'MUDHARABAH', 'QARD'
    account_status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE', -- ACTIVE, BLOCKED, CLOSED
    balance DECIMAL(18, 2) NOT NULL DEFAULT 0.00,
    currency VARCHAR(3) NOT NULL DEFAULT 'IDR',
    opened_date DATE NOT NULL DEFAULT CURRENT_DATE,
    closed_date DATE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT check_balance_non_negative CHECK (balance >= 0)
);

-- Transaction log table (for audit purposes)
CREATE TABLE demo.transaction_log (
    transaction_id BIGSERIAL PRIMARY KEY,
    account_id BIGINT NOT NULL REFERENCES demo.accounts(account_id),
    transaction_type VARCHAR(50) NOT NULL, -- BALANCE_INQUIRY, DEBIT, CREDIT
    amount DECIMAL(18, 2),
    balance_before DECIMAL(18, 2),
    balance_after DECIMAL(18, 2),
    transaction_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    reference_number VARCHAR(50),
    description TEXT,
    response_time_ms INTEGER
);

-- Create indexes for better query performance
CREATE INDEX idx_customers_number ON demo.customers(customer_number);
CREATE INDEX idx_accounts_number ON demo.accounts(account_number);
CREATE INDEX idx_accounts_customer_id ON demo.accounts(customer_id);
CREATE INDEX idx_accounts_status ON demo.accounts(account_status);
CREATE INDEX idx_transaction_log_account_id ON demo.transaction_log(account_id);
CREATE INDEX idx_transaction_log_time ON demo.transaction_log(transaction_time);
CREATE INDEX idx_transaction_log_type ON demo.transaction_log(transaction_type);