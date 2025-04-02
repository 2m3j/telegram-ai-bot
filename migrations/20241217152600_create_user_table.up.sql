create table user
(
    id            bigint unsigned not null primary key,
    ai_platform   varchar(50)  not null,
    ai_model      varchar(50)  not null,
    username      varchar(255) not null,
    first_name    varchar(255) null,
    last_name     varchar(255) null,
    language_code varchar(10)  not null,
    updated_at    datetime     not null,
    created_at    datetime     not null
);
