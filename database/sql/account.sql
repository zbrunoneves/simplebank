-- name: CreateAccount :execresult
insert into account(
    owner, balance, currency
) values (
    ?, ?, ?
);

-- name: GetAccount :one
select * from account
where id = ?
limit 1
for share;

-- name: GetAccounts :many


-- name: UpdateAccount :execresult
update account
set balance = ?
where id = ?;

-- name: DeleteAccount :execresult
delete from account
where id = ?;