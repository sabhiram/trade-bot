# trade-bot

A simple bittrex trade bot to set slightly complex buy / sell orders

## Install

```
go get github.com/sabhiram/trade-bot
```

TODO: Releases etc...

## Usage

The bittrex documentation explains how to acquire a `API_KEY` and `API_SECRET` from your settings panel (once logged in).  Once you have these you can use the `trade-bot` like so:

```
BITTREX_API_KEY=<apikey> BITTREX_SECRET=<secret> trade-bot [options]
```

Or the appropriate env variables can be set beforehand.

