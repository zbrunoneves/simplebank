-- name: CreateAccount :execresult
insert into account(
    owner, balance, currency
) values (
    ?, ?, ?
);

-- name: GetAccount :one
select * from account
where id = ?
limit 1;

-- name: GetAccounts :many
select * from account
order by id
limit ? offset ?;

-- name: UpdateAccount :execresult
update account
set balance = ?
where id = ?;

-- name: DeleteAccount :execresult
delete from account
where id = ?;