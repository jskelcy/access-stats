# Access-Stats

Access-Stats watches a file for updates using the [Common Log Format](https://en.wikipedia.org/wiki/Common_Log_Format). Stats are emmitted at a 10 second interval.
If the average number of logs per second is above an `Alert Threshold` for more than 2 minutes an alert will be raised.
When the average number of logs per second has been below the `Alert Threshold` for more than 2 mintues the alert will recover.

The default file to be watched is `/var/log/access.log`. This is overwritable using the `-src` flag.  
The default `Alert Threshold` is 10 logs per second. This is overwritable using the `-alertThreshold` flag.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

### Prerequisites

Go 1.11.1

### Building and running

Installing the dependecies for this project can be done via the `deps` make target.
```
make deps
```

Alternitivly you can use the `build` make target to build the binary in `/out` directory. When using the `build` make dependecies are also installed.

```
make build
```

Finally you can run install dependencies, build a binary and run the binary with the `run` make target.

```
make run
```

To overwrite the default log source file or alert threshold use the appropriate flags. And example could be:

```
make run src=test.log alertThreshold=5
```

## Running the tests

Tests can be run using the `test` make target.

```
make test
```

### Break down into end to end tests

In order to populate a watched file with sample data a script has been provided. The `src` flag is used to pass in the file where the sample data should be written.
Using the default alert threshold of 10 logs per second this script should be able to generate an alert after 2 minutes. 

```
make sample-data src=test.log
```


## Deployment

Add additional notes about how to deploy this on a live system
