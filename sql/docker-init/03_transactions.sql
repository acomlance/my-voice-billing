-- status: 0 in_progress, 1 approved, 2 cancelled, 3 failed
-- payment_type: 1 deposit, 2 charge, 3 refund, 4 bonus
-- payment_method: 1 card, 2 bank
create table public.transactions
(
    id                  bigserial
        constraint transactions_pk
            primary key,
    account_id          bigint                                     not null
        constraint transactions_accounts_id_fk
            references public.accounts,
    status              smallint     default 0                     not null,
    amount              bigint       default 0                     not null,
    tokens              bigint       default 0                     not null,
    tokens_after        bigint,
    payment_type        smallint                                   not null,
    payment_method      smallint                                   not null,
    payment_description varchar(255) default ''::character varying not null,
    payment_data        text         default ''::text              not null,
    date_create         timestamp    default now()                 not null,
    date_update         timestamp
);

alter table public.transactions
    owner to postgres;

create index transactions_account_id_index
    on public.transactions (account_id);

create index transactions_status_index
    on public.transactions (status);
