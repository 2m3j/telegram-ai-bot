create table ai_conversation
(
    id         varchar(36)     not null primary key,
    user_id    bigint unsigned not null,
    started_at datetime        not null,
    ended_at   datetime        null,
    constraint ai_conversation_cst_id_fk
        foreign key (user_id) references user (id)
);

create index ai_conversation_ended_at_indx
    on ai_conversation (ended_at);

create index ai_conversation_started_at_indx
    on ai_conversation (started_at);

create table ai_message
(
    id                 varchar(36)  not null primary key,
    conversation_id    varchar(36)  not null,
    status             varchar(10)  not null,
    user_message       mediumtext   not null,
    assistant_platform varchar(100) not null,
    assistant_model    varchar(100) not null,
    assistant_message  mediumtext   null,
    updated_at         datetime     not null,
    created_at         datetime     not null,
    constraint ai_message_conv_id_fk
        foreign key (conversation_id) references ai_conversation (id)
);

create index ai_message_conversation_id_indx
    on ai_message (conversation_id);

create index ai_message_created_at_indx
    on ai_message (created_at);

create index ai_message_status_indx
    on ai_message (status);
