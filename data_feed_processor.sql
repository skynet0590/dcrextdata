create table data (
id serial primary key not null
);


create table historic_data (
id serial primary key not null ,
exchangeID integer,
globaltradeid varchar (20),
tradeid varchar (20),
timestamp timestamp,
quantity  varchar (30),
price varchar (20), 
total varchar (20),
fill_type varchar (20),
order_type varchar (20)
);

create table chart_data (
    id serial primary key not null ,
    exchangeID integer ,
    date       timestamp,
    high    varchar (20),
    low     varchar (20),
    open    varchar (20),
    close   varchar (20),
    volume  varchar (20),
    quotevolume varchar (20),
    weightedaverage varchar (20)
);