create table reminders(
      id              serial       primary key,
      user_id         varchar(255) not null,
      text            text         not null,
      created_at      timestamptz  default now(),
      remind_at       timestamptz  not null,
      remind_original varchar(255),
      snoozed         boolean      default false
);

create index uid_time on reminders(user_id, remind_at, snoozed);
