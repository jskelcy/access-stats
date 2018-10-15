# Access-Stats

Access-Stats watches a file for updates using the [Common Log Format](https://en.wikipedia.org/wiki/Common_Log_Format). Stats are emmitted at a 10 second interval.
If the average number of logs per second is above an alert threshold for more than 2 minutes an alert will be raised.
When the average number of logs per second has been below the alert threshold for more than 2 mintues the alert will recover.

The default file to be watched is `/var/log/access.log`. This is overwritable using the `-src` flag.  
The default alert threshold is 10 logs per second. This is overwritable using the `-alertThreshold` flag.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

### Prerequisites

Go 1.11.1

### Building and running

Installing the dependecies for this project can be done via the `deps` make target.
```
make deps
```

Alternitivly you can use the `build` make target to build the binary in `/out` directory. When using `build`, dependecies are also installed.

```
make build
```

Finally you can install dependencies, build a binary, and run the binary with the `run` make target.

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
Using the default alert threshold of 10 logs per second this script will generate an alert after 2 minutes. After the sample-data script exists the alert will eventually recover, recovery progress can be tracked by watching the ongoing alert notification.

```
make sample-data src=<log-file>
```

As a seperate process run access-stats using the `run` make target:

```
make run src=<log-file> alertThreshold=10
```


## Deployment

A docker image can be built from the docker file using the docker build command.

An example might be:
```
docker build -t jskelcy/acces-stats:latest .
```

## Potential Improvments

There are three areas where this application could be improved:
* Access-stats does very little in the way of log validation. A basic improvment would be to handle malformed logs and emmit some additional metircs about the number of malformed logs.
* There are a number of properties of this app which are not configurable which could be. Those being the aggregation window, which is currently hardcoded at 10 seconds and the alert window, which is currently hardcoded at 2 minutes. Given the logging pattern of some applications these windows may not provide proper signal, so allowing users to customize these windows could allow the application to apply to a broader set of use cases.
* When a log is ingested it is bucketed by a few predefined vectors, for example user, status code, section. These histograms are isoliated so you can get stats based only on the metrics in a specific histogram. For example you could say something like "Give me the 90th percentile of sections which go 4XX responses". However adding a more advance stat tagging system access-stats could support advanced queries like "Give me the 90th percentile of sections which go 4XX responses where the callers name was Lisa". 
