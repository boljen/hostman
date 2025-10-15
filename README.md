# Hostman

[![Tests](https://github.com/boljen/hostman/actions/workflows/tests.yml/badge.svg)](https://github.com/boljen/hostman/actions/workflows/tests.yml)

A utility to manage hosts files for project-specific development/testing environment setup.

## Features

* Set and sync a hosts file from a configuration file
* Watch the configuration file for changes and automatically update the hosts file
* Use an HTTP endpoint for dynamic updates based on the current application state

## Installation

Head over to the releases page or install from source;

    go install github.com/boljen/hostman@latest

## Usage

IMPORTANT! The binary requires elevated privileges to modify the hosts file.

Step 1: Create a configuration file "hostman.hcl" in the current directory or any of the root directories:

    project = "my-project"

    sources = [
        "main"
    ]

    static "main" {
        host = "example.com"
        ip = "127.0.0.1"
    }

Step 2: Run the binary as administrator (you can enable "sudo" in Windows in the developer settings):

    sudo hostman apply

There should be a new section in your hosts file that maps example.com to 127.0.0.1.

Alternatively you can use watch mode which will refresh every 5 seconds;

    sudo hostman apply -w

## Advanced configuration

### Multiple hosts

You can rewrite the static block as follows to allow multiple hosts:

    static "main" {
        hosts = [
            "example.com",
            "example.org"
        ]
        ip = "127.0.0.1"
    }

### Dynamic hosts

You can specify an "http" block which will use an http endpoint to determine the hosts to map:

    http "main" {
        endpoint = "https://raw.githubusercontent.com/boljen/hostman/refs/heads/master/examples/http/response.json"
    }

This is useful when you are working with a multi-tenant application where a tenant lives in a
dynamically generated subdomain or brings their own domain.