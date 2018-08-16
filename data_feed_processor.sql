

create table historic_data (
id serial primary key not null ,
exchange_name varchar (50),
globaltradeid numeric,
tradeid numeric,
created_on varchar (100),
quantity  NUMERIC,
price numeric, 
total numeric,
fill_type varchar (20),
order_type varchar (20)

);

create table chart_data (
    id serial primary key not null ,
    exchange_name VARCHAR (50),
    created_on    varchar (100),
    high    numeric,
    low     numeric,
    opening    numeric,
    closing   numeric,
    volume  numeric,
    quotevolume numeric,
    basevolume numeric,
    weightedaverage NUMERIC
);

create table pos_data(
    id serial primary key not null ,
    posid VARCHAR (20),
    apienabled varchar (20),
    network VARCHAR (20),
    network_url VARCHAR (100),
    launched VARCHAR (20),
    last_updated VARCHAR (20),
    immature VARCHAR (20),
    live VARCHAR (20),
    voted NUMERIC,
    missed NUMERIC,
    poolfees NUMERIC,
    proportionlive NUMERIC,
    proportionmissed numeric ,
    usercount NUMERIC,
    usercountactive NUMERIC,
    created_on TIME

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
    created_time numeric ,
    networkdifficulty numeric ,
    coinprice numeric,
    btcprice numeric ,
    est numeric ,
    created_on numeric ,
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

