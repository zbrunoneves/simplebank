-- name: CreateTransfer :execresult
insert into transfer(
    from_account_id, to_account_id, amount
) values (
    ?, ?, ?
);

-- name: GetTransfer :one
select * from transfer
where id = ?
limit 1;

-- name: GetTransfers :many
select * from transfer
order by id
limit ? offset ?;

-- name: DeleteTransfer :execresult
delete from transfer
where id = ?;