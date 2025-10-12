\c demo_db;

-- Function definitions

-- Create a function to log balance inquiries (useful for load testing)
CREATE OR REPLACE FUNCTION demo.log_balance_inquiry(
    p_account_id BIGINT,
    p_reference_number VARCHAR(50),
    p_response_time_ms INTEGER DEFAULT NULL
)
RETURNS void AS $$
DECLARE
    v_current_balance DECIMAL(18, 2);
BEGIN
    -- Get current balance
    SELECT balance INTO v_current_balance
    FROM demo.accounts
    WHERE account_id = p_account_id;

    -- Log the inquiry
    INSERT INTO demo.transaction_log (
        account_id,
        transaction_type,
        amount,
        balance_before,
        balance_after,
        reference_number,
        description,
        response_time_ms
    ) VALUES (
        p_account_id,
        'BALANCE_INQUIRY',
        NULL,
        v_current_balance,
        v_current_balance,
        p_reference_number,
        'Balance inquiry via API',
        p_response_time_ms
    );
END;
$$ LANGUAGE plpgsql;
