-- name: CreatePayout :one
INSERT INTO payouts (
  id, batch_id, reference_id, 
  recipient_name, recipient_phone, recipient_email, recipient_tag,
  country_code, bank_code, account_number, bank_name,
  amount, currency, status
) VALUES (
  ?, ?, ?, 
  ?, ?, ?, ?,
  ?, ?, ?, ?,
  ?, ?, ?
)
RETURNING *;

-- name: GetPayout :one
SELECT * FROM payouts 
WHERE id = ? LIMIT 1;

-- name: ListPayouts :many
SELECT * FROM payouts 
ORDER BY created_at DESC;

-- name: UpdatePayoutStatus :exec
UPDATE payouts 
SET status = ?, error_message = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: ListPayoutsByBatchID :many
SELECT * FROM payouts 
WHERE batch_id = ?
ORDER BY created_at DESC;