create table public.tasks
(
    id              bigserial
        constraint tasks_pk
            primary key,
    token           varchar(128) default ''::character varying not null,
    account_id      bigint      default 0 not null
        constraint tasks_accounts_id_fk
            references public.accounts
            on delete cascade,
    reserved_tokens bigint      default 0 not null
        constraint reserve_check check (reserved_tokens > 0),
    date_create     timestamp   default now() not null
);

alter table public.tasks
    owner to postgres;

create unique index tasks_token_uindex
    on public.tasks (token);

create index tasks_account_id_index
    on public.tasks (account_id);
