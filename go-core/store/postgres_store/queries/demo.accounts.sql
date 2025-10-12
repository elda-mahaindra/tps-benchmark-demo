-- name: GetAccountByAccountNumber :one
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
WHERE a.account_number = $1;