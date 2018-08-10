create table data (
id serial primary key not null
);

create table New_table (
        id serial primary key not null ,
        Exchangeid numeric,
		Globaltradeid numeric,
		Tradeid numeric,
		Timestamping Timestamp,
	    Quantity numeric,
		Price numeric,
		Total numeric,
		FillType VARCHAR (20),
		OrderType varchar (20)
);


create table historic_data (
id serial primary key not null ,
exchangeID numeric  null,
globaltradeid numeric null,
tradeid numeric  null,
timestamp timestamp null,
quantity  varchar (30) null,
price numeric null, 
total numeric null,
fill_type varchar (20) null,
order_type varchar (20) null
);

create table chart_data (
    id serial primary key not null ,
    exchangeID integer ,
    date       timestamp,
    high    varchar (20),
    low     varchar (20),
    open1   varchar (20),
    close1   varchar (20),
    volume  varchar (20),
    quotevolume varchar (20),
    weightedaverage varchar (20)
);

create table POSData (
    id serial primary key not null ,
    POSid varchar (20),
    Apienabled varchar (10) ,
    APIVersionsSupported varchar (20) ,
    Network varchar (20),
    URL varchar (50) ,
    Launched numeric,
    LastUpdated numeric,
    Immature numeric,
    Live numeric,
    Voted numeric,
    Missed numeric,
    PoolFees numeric,
    ProportionLive numeric,
    ProportionMissed numeric,
    UserCount numeric,
    UserCountActive numeric

)

