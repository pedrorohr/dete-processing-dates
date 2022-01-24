[![DETE Processing Dates Pipeline](https://github.com/pedrorohr/dete-processing-dates/actions/workflows/dete-processing-dates.yml/badge.svg?branch=main)](https://github.com/pedrorohr/dete-processing-dates/actions/workflows/dete-processing-dates.yml)

# DETE Processing Dates

This repository implements a system that notifies updates in current processing dates for employment permits processed by the Irish Department of Enterprise, Trade and Employment (DETE).

You can check the updates in [this Telegram channel](https://t.me/deteprocessingdates).

My intention was to learn AWS Lambda while building something useful for me, since I've been waiting for my Employment Permit to be processed by the Department.

# Irish Employment Permits

In order to work in Ireland a non-EEA National, unless they are exempted, must hold a valid [Employment Permit](https://enterprise.gov.ie/en/What-We-Do/Workplace-and-Skills/Employment-Permits/).
The Department of Enterprise, Trade and Employment (DETE) is responsible for processing Employment Permit applications.

They process the applications in a FIFO approach with a queue for each Employer Type.
Applications coming from [Trusted Partner](https://enterprise.gov.ie/en/What-We-Do/Workplace-and-Skills/Employment-Permits/Trusted-Partner-Initiative/) Employers line up in a higher priority queue, while the other applications are defined as Standard and line up in another queue.

The backlog of applications to process is getting bigger and bigger.
It's taking more than 3 months for the application to be processed.
DETE publishes the [Current Processing Dates](https://enterprise.gov.ie/en/What-We-Do/Workplace-and-Skills/Employment-Permits/Current-Application-Processing-Dates/) in its website so you can have an idea of how far your application is from being processed.

# Architecture

The system is composed by:
* two AWS Lambdas
* one AWS CloudWatch Event
* one AWS DynamoDB Table
* one AWS DynamoDB Stream
* one Telegram Bot
* one Telegram Channel

The project uses Terraform to define and provide the infrastructure.
A GitHub Action is responsible by the project's pipeline.

## Scraper

It's an AWS Lambda written in Golang.
The Lambda is invoked by a CloudWatch Event scheduled to happen every five minutes.
It uses [Colly](http://go-colly.org/) to scrap the Current Processing Dates page to get the Processing Dates.
It checks if there's a new Processing Date using a DynamoDB Table.
If there's a new one, it saves on the database.

## Notifier

It's an AWS Lambda written in Golang.
The Lambda is invoked by a DynamoDB Stream Event that happens every time a new Processing Date is saved.
It notifies the update to [this Telegram channel](https://t.me/deteprocessingdates).