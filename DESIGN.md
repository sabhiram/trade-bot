Present a minimal HTTP interface into a webserver to show bittrex balances etc
    HTTP server + client, can be polymer2?
    Initially just displays balance object
    Refresh to update balances (websocket)

Create websocket between client and server to send messages and events
    Enable gorilla/ws to send messages to client page

World state should know:
    Current balances (web client can refresh this)
    Current / outstanding conditional orders

JSON file for "db" will be stored on disk - maybe we swap this to bolt et al?

Advanced:
    Password protect access to the server (set during server startup?)
