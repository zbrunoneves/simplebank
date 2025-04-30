select id, owner, balance, currency, created_at from account
order by id
limit ? offset ?;