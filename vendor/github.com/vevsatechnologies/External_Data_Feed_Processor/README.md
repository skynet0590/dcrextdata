# External_Data_Feed_Processor
An external data feed processor takes takes DCR pricing data from various exchanges periodically and stores them .

## Get the required Golang packages.
```
go get -u github.com/vattle/sqlboiler/boil
```
```
go get github.com/spf13/viper
```
```
go get github.com/vevsatechnologies/External_Data_Feed_Processor/models
```
```
go get github.com/lib/pq
```
```
go get gopkg.in/inconshreveable/log15.v2
```

### Set the configration settings for database in config.json file
### Create database *data_feed_processor*


```psql -U postgres data_feed_processor < data_feed_processor.sql
```
```
go generate
```

```
sqlboiler postgres
```


### Run the project

```
go run main.go bittrex.go poloniex.go POS.go POW.go 
```



