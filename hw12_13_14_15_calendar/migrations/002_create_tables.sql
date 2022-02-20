-- +goose Up
CREATE TABLE events(  
    id text NOT NULL PRIMARY KEY,
    title text not null,
    time_start TIMESTAMP with time zone,
    time_end TIMESTAMP with time zone,
    description text,
    user_id text,
    notify_before TIMESTAMP with time zone,
    notified BIT
);

-- +goose Down
drop table events;

drop DATABASE calendar;
