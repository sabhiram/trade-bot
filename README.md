# trade-bot

A simple bittrex trade bot to set slightly complex buy / sell orders.

Code is WIP and does not actually function as intended yet.

## Disclaimer

This is open source software.  Use at your own risk.

Running this application requires you to get an `api key` and `secret` from bittrex.  If this was compromised it could lead to loss of all your coins.  Always use caution, paranoia and copious common sense when running this code.

## Install

```
go get github.com/sabhiram/trade-bot
```

## Usage

The bittrex documentation explains how to acquire a `BITTREX_API_KEY` and `BITTREX_SECRET` from your settings panel (once logged in).  Once you have these you can use the `trade-bot` like so:

```
  $ BITTREX_API_KEY=<key> BITTREX_SECRET=<secret> \
      trade-bot -pair <pair> [-refresh 5s] <command>

  All 'command's will trigger once a target is hit.  Each command will
  query the user for all required parameters and confirm before being
  deployed.  Once deployed the script will query the last trade value
  for the 'pair' every 'refresh' seconds (default 5s).  Once any/all of
  the conditions are met, the appropriate action is taken.

  Where 'command' can include:

    limit-sell          -   regular on-limit sale
    stop-loss           -   simple stop loss
    high-low            -   stop-loss + limit-sell (first one wins)

  For authorizing transactions and conducting market queries, the API
  key and secret need to be provided as environment variables. The two
  required ones are:

    BITTREX_API_KEY     -   api-key from bittrex
    BITTREX_SECRET      -   secret for above key

  The api-key and secret can be accessed once you have logged in to
  bittrex in the "Settings" section under "Manage API Keys".

  Note: You must have 2-factor authentication enabled to make new keys

```

## Issues

If you find this software useful, help out by filing issues or suggestions here: https://github.com/sabhiram/trade-bot/issues.
