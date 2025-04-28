-- name: CreateOperation :execresult
insert into operation(
    account_id, amount
) values (
    ?, ?
);

-- name: GetOperation :one
select * from operation
where id = ?
limit 1;

-- name: GetOperations :many
select * from operation
order by id
limit ? offset ?;

-- name: DeleteOperation :execresult
delete from operation
where id = ?;