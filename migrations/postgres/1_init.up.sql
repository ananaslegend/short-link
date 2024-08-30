create table if not exists link(
   id serial primary key,
   alias text not null unique,
   link text not null
);