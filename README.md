# Eventbridge Gate ⭐

A CLI tool for monitoring and debugging AWS Eventbridge events.


## Introduction
### The Problem
Developing event-driven systems on AWS is exciting, but misconfiguration may happen and debugging interactions between components is not always as easy as one would wish.
Sometimes you want to react to a given event but you are unsure of how exactly it is structured. Other times you are building events from the ground up and want to be sure that
the events you are emitting are correct, while also make sure to create correct event patterns for catching said events with a rule.

### The Solution
**Eventbridge Gate** is a simple tool that allows you to hook up your event bus with a sniffer and monitor all the events passing through, so you can look at their structure and easily debug your system.

This tool will allow you to setup (and tear down) monitors on any bus in your AWS account, that will sniff any event passing through and make it available for you to watch in realtime.

This tool is based on an existing reference architecture described on this [ServerlessLand Page]()


## Usage

### 🚨 READ THIS BEFORE CONTINUING 🚨
This tool will require to create resources in your AWS account. 
This basically means two things:
* The AWS credentials you provide must have enough permissions for creating and destroying the needed resources
* There will be (minimal, possibly 0) **costs** associated to the usage of this tool, which are directly related to the amount of events passing through the monitors


### Install with [Homebrew](brew.sh)

To install the tool you can use Homebrew with both linux-based or macOs operating systems.
```
brew tap aka-somix/eventbridge-gate
brew install egate
```

### Expor

### Set up your first monitor
TBD

### List active monitors
TBD

### Watch logs of existing monitor
TBD

### Remove monitor
TBD



## Contribution

TBD

## Costs associated