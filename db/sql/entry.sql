-- name: CreateEntry :execresult
insert into entry(
    account_id, amount
) values (
    ?, ?
);

-- name: GetEntry :one
select * from entry
where id = ?
limit 1;

-- name: GetEntries :many
select * from entry
order by id
limit ? offset ?;

-- name: DeleteEntry :execresult
delete from entry
where id = ?;