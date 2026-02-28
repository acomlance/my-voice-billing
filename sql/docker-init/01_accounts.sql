create table public.accounts
(
    id           bigint                  not null
        constraint accounts_pk
            primary key,
    balance      bigint      default 0   not null,
    reserve      bigint      default 0   not null,
    state        smallint    default 0   not null,
    date_create  timestamp  default now() not null,
    date_update  timestamp  default now() not null,
    constraint balance_check check (balance >= reserve)
);

alter table public.accounts
    owner to postgres;
