

create table historic_data (
id serial primary key not null ,
exchangeID integer,
globaltradeid varchar (20),
tradeid varchar (20),
timest timestamp,
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
    opening    varchar (20),
    closing   varchar (20),
    volume  varchar (20),
    quotevolume varchar (20),
    weightedaverage varchar (20)
);

create table pos_data(
    id serial primary key not null ,
    Posid VARCHAR (20),
    Apienabled varchar (20),
    Apiversionssupported numeric,
    Network VARCHAR (20),
    URL VARCHAR (100),
    Launched VARCHAR (20),
    Lastupdated VARCHAR (20),
    Immature VARCHAR (20),
    Live VARCHAR (20),
    Voted NUMERIC,
    Missed NUMERIC,
    Poolfees NUMERIC,
    Proportionlive NUMERIC,
    Proportionmissed numeric ,
    Usercount NUMERIC,
    Usercountactive NUMERIC,
    Timestamp TIME

);

create table pow_data(

    id serial primary key not null,
    powid numeric ,
    hashrate numeric ,
    efficiency numeric,
    progress numeric ,
    workers numeric,
    currentnetworkblock numeric,
    nextnetworkblock numeric ,
    lastblock numeric ,
    networkdiff numeric,
    esttime numeric ,
    estshare numeric ,
    timesincelast numeric ,
    nethashrate numeric,
    blocksfound numeric,
    totalminers numeric,
    time numeric ,
    networkdifficulty numeric ,
    coinprice numeric,
    btcprice numeric ,
    est numeric ,
    date numeric ,
    blocksper numeric ,
    luck numeric ,
    ppshare numeric ,
    totalkickback numeric,
    success VARCHAR (20),
    lastupdate numeric ,
    name VARCHAR (20),
    port numeric ,
    fees numeric ,
    estimatecurrent numeric,
    estimatelast24h numeric,
    actual124h numeric,
    mbtcmhfactor numeric,
    hashratelast24h numeric,
    rentalcurrent numeric ,
    height numeric ,
    blocks24h numeric ,
    btc24h numeric,
    currentheight numeric ,
    total numeric ,
    pos numeric ,
    pow numeric ,
    dev numeric 


);

