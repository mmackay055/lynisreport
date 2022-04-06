# Lynis Report Tool

This is a tool written in [Go](https://go.dev/) that is used to process the
[Lynis](https://cisofy.com/lynis/) auditing tool for Unix system. For now it
will just be used to create summary logs with all the warnings, and suggestions
contained in a JSON string so it can be easily processed by platforms such as
[Elasticsearch](https://www.elastic.co/), ect.

## Build

From **src** directory issue command:
`go build`
which will build **lynisreport** binary.

## Install

From **src** directory issue command:
`go install`
which will install **lynisreport** in go bin.

## Run

Run binary like so
`lynisreport`
which will search for Lynis report in **/var/log/lynis-report.dat** and will 
print results to console.

Use **-h** option to review other options.

## Elastic helper script

Script **scripts/elasticsearch** is an example of how the Lynis tool and the 
**lynisreport** tool can be used in a cron script to periodcally perform scans and 
generate logs that can be ingested by Elasticsearch.

