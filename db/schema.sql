create table account (
		id int unsigned primary key auto_increment,
        owner varchar(255) not null,
        balance bigint not null default 0,
        currency varchar(2) not null,
        created_at timestamp not null default current_timestamp,

        index idx_owner (owner)
);

create table entry (
	id int unsigned primary key auto_increment,
    account_id int unsigned not null,
    amount bigint not null,
    created_at timestamp not null default current_timestamp,

    constraint fk_account foreign key (account_id) references account(id),

    index idx_account (account_id)
);

create table transfer (
	id int unsigned primary key auto_increment,
    from_account_id int unsigned not null,
    to_account_id int unsigned not null,
    amount bigint unsigned not null,
    created_at timestamp not null default current_timestamp,

    constraint fk_from_account foreign key (from_account_id) references account(id),
    constraint fk_to_account foreign key (to_account_id) references account(id),

    index idx_from (from_account_id),
    index idx_to (to_account_id),
    index idx_transfer (from_account_id, to_account_id)
);
