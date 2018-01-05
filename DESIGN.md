Present a minimal HTTP interface into a webserver to show bittrex balances etc
    HTTP server + client, can be polymer2?
    Initially just displays balance object
    Refresh to update balances (websocket)

World state should know:
    Current balances (web client can refresh this)
    Current / outstanding conditional orders

Advanced:
    Password protect access to the server (set during server startup?)
