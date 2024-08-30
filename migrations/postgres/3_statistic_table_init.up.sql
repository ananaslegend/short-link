create table if not exists statistic(
    id serial primary key,
    redirect_time_stamp integer,
    link text,
    redirect integer
);