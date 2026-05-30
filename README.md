# SSL Certificate Checker

A simple tool to check certificates for website from command line.

## Purpose

The tool has been developed for learning purpose only. There are tons of better, more reliable tools for this.

The code is written as an exercise in re-writing from scratch in go an old script hacked together in Ruby many years ago. There is no AI involved in the code, part of the challenge is doing this the old way, like the original script. THIS IS **intentionally** a basic rewrite of script with all its quirks and bad design. 

This has been written using: GoLand, terminal with Fish Shell, DuckDuckGo, online documentation.

As this is educational, feel free to open issues, write comments, give guidance or anything, it was my first Go program after reading "Learning Go" and has been done in barely a couple of hours. Actually, writing the README has taken more time than the code.

## Rationale

Suppose you have several customers who keep inquiring: "Are these certificate going to expire soon"?
The point is: there are different agencies producing different certificates with different timespans for different webbsites.
In time, some of the websites / webbapps have been dismissed, some have changed the certificate provider, some have been moved to other machines.

The final goal is to gather as much info as possible, acquire a renewed certificate by the current provider of that certificate and contact the DevOps team with precise information on what needs to be updated where. In formation such as load balancer or "the sites is being dismissed but for now redirect to that other site" are useful to have.

The ideal solution is to have an automated software check out all the domains managed autonomously.
There are tools for this, but as a quick and dirty solution a script used to do this:
- ask for a domain
- check available LOCAL/INTERNAL info 
  - resolve the domain name to get an IP address
  - check the IP and extract a description for the host
  - look for the domain name in a dictionary (manually created by fiddling with data from a worksheet and morphed using some fancy multi cursor black magic) and provide known information ("knowledge" of the internal systems)
- check the certificate for the domain by connecting, getting the certificate and getting live information about it: expiration, SANs, issuers, ecc.
- provide a schematic report on the domain

With this script being a temporary solution, it never got much time spent on it. Planned changes were:
- Tests
- external configuration
- read from CSV file so the Worksheet is exported and there is no fiddling, black magic or similar
- same for domain -> description
- a scan mode to check all the domains in the Worksheet, with alerts for those expiring
- exporters (e.g. generate the text of an email with the information laid out to make contacting the DevOps teams quickly).

Things yet to be ported:
- interactive loop to keep inserting domains until finished
- non blocking error domain (will need to either turn panic into errors or using recover)

The point of the code is:
- compare the differences in writing the code between a language, Ruby, that is built with ease of development (developer before the machine) and abstraction (pure OOP) and another language, Go, that follows a pragmatical approach that is neither functional, nor OOP, nor properly imperative
  - in code statistics (length, time spent, quality of the result, ...)
  - semantics, statistics and so on. 

## Should I use it?

Short answer is no. There are few checks, the code is maily for learning purposes. The code is being hacked together in the night, on al old laptop, lacks tests (the test should have been written before writing the code, but the original script was supposed to be **quick** and _dirty_).

The first version is, by choice, the same design as the original script, so buolt to be some quick and some dirty. The CSV planned feature has been implemented only to get some further insights.

## Future Development

Some updates and features might be planned as experiments. More likely the whole code will be rewritten with proper design, if I still find this toy project to still have value in my learning.

As I am concentrating my attention much more on the web side, this is likely to become abandoned.

## Disclaimer

This code is here as is, may not work, may be broken, may e not working as intended. use it at your own risk. 
