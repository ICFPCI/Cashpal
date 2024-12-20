-- USERS

-- name: GetUser :one
SELECT * FROM Users
WHERE id = $1 LIMIT 1;

-- name: GetUserByUsername :one
SELECT * FROM Users
WHERE username = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM Users
ORDER BY id;

-- name: CreateUser :one
INSERT INTO Users (
  username, password
) VALUES (
  $1, $2
)
RETURNING *;

-- name: UpdateUser :one
UPDATE Users
  set password = $2, updated_at = NOW() AT TIME ZONE 'utc'
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM Users
WHERE id = $1;

-- ACCOUNTS

-- name: GetAccount :one
SELECT * FROM Accounts
WHERE id = $1 LIMIT 1;

-- name: GetAccountWithUserCheck :one
SELECT 
    acc.*,
    CASE 
        WHEN mem.user_id IS NOT NULL THEN 1
        ELSE 0
    END AS is_member
FROM Accounts acc
JOIN Members mem ON acc.id = mem.account_id AND mem.user_id = $2
WHERE acc.id = $1;

-- name: ListAccount :many
SELECT * FROM Accounts
ORDER BY id;

-- name: ListAccountByUser :many
SELECT 
  	acc.*
FROM Accounts acc
JOIN Members mem ON acc.id = mem.account_id
WHERE mem.user_id = $1;

-- name: CreateAccount :one
INSERT INTO Accounts (
  account_name, account_type
) VALUES (
  $1, $2
)
RETURNING *;

-- name: UpdateAccount :one
UPDATE Accounts
  set account_name = $2, account_type = $3, updated_at = NOW() AT TIME ZONE 'utc'
WHERE id = $1
RETURNING *;

-- name: DeleteAccount :exec
DELETE FROM Accounts
WHERE id = $1;

-- ACCOUNT_EVENTS

-- name: GetAccountEvent :one
SELECT * FROM Account_Events
WHERE id = $1 LIMIT 1;

-- name: ListAccountEvent :many
SELECT * FROM Account_Events
ORDER BY id;

-- name: ListAccountEventByAccount :many
SELECT * FROM Account_Events
WHERE account_id = $1
ORDER BY id;

-- name: CreateAccountEvent :one
INSERT INTO Account_Events (
  account_id, event_type_id, description
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: UpdateAccountEvent :one
UPDATE Account_Events
  set description = $2, updated_at = NOW() AT TIME ZONE 'utc'
WHERE id = $1
RETURNING *;

-- name: DeleteAccountEvent :exec
DELETE FROM Account_Events
WHERE id = $1;

-- MEMBERS

-- name: GetMember :one
SELECT * FROM Members
WHERE account_id = $1 and user_id = $2 LIMIT 1;

-- name: GetMemberWithUserCheck :one
SELECT m1.*
FROM Members as m1
JOIN Members as m2
ON m1.account_id = m2.account_id AND m2.user_id = $3
WHERE m1.account_id = $1 AND m1.user_id = $2;

-- name: ListMember :many
SELECT * FROM Members
ORDER BY id;

-- name: ListMemberByAccount :many
SELECT * FROM Members
WHERE account_id = $1
ORDER BY id;

-- name: ListMemberByAccountWithUserCheck :many
SELECT m1.*
FROM Members as m1
JOIN Members as m2
ON m1.account_id = m2.account_id AND m2.user_id = $2
WHERE m1.account_id = $1;

-- name: CreateMember :one
INSERT INTO Members (
  account_id, user_id, member_role_id
)
VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: UpdateMember :one
UPDATE Members
  set member_role_id = $3, updated_at = NOW() AT TIME ZONE 'utc'
WHERE account_id = $1 and user_id = $2
RETURNING *;

-- name: DeleteMember :exec
DELETE FROM Members
WHERE account_id = $1 and user_id = $2;

-- TRANSACTIONS

-- name: GetTransaction :one
SELECT * FROM Transactions
WHERE id = $1 LIMIT 1;

-- name: GetTransactionWithCheck :one
SELECT t.*
FROM transactions AS t
WHERE t.account_id = $1 and t.id = $2 AND EXISTS (
	SELECT 1
	FROM members AS m
	WHERE m.account_id = t.account_id AND m.user_id = $3
);

-- name: ListTransaction :many
SELECT * FROM Transactions
ORDER BY id;

-- name: ListTransactionByAccount :many
SELECT t.*
FROM transactions AS t
WHERE t.account_id = $1 AND EXISTS (
	SELECT 1
	FROM members AS m
	WHERE m.account_id = t.account_id AND m.user_id = $2
);

-- SELECT * FROM Transactions
-- WHERE account_id = $1
-- ORDER BY id;

-- name: CreateTransaction :one
INSERT INTO Transactions (
  account_id, user_id, transaction_date, transaction_type_id, amount, description
)
VALUES(
  $1, $2, $3, $4, $5, $6
)
returning *;

-- name: UpdateTransaction :one
UPDATE Transactions
  SET amount = $2, description = $3, updated_at = NOW() AT TIME ZONE 'utc'
  WHERE id = $1
  RETURNING *;