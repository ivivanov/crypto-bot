# Grid bot

## Installation From source
```
go install
```
## Start
- Generate api keys with the minimal possible privileges. It is recommended to use sub account.
- Create .env file
- Execute `crypto-bot` from the directory of .env
  
### TODOs
 - In case we hit market price (taker) when posting new order (trade_buy->create_sell could hit directly as market order)
    - this will displace the initial buy order. buy@1 sell@1.1 but market has ask@1.2 => actual sell@1.2 =>  new buy@1.1 (wanted buy@1)
    - store data where should be the initial buy
    - create sell order 0.0001 above the current bid (order book needed). This will increase profit by
    not paying taker fee + selling above the desired target


- add more trade data to clientOrderId - price, target profit, etc
- console output
    - add ClientOrderID
    - log remove date, add ms

- consider that trade_buy->create_sell could hit directly as market order

- auto reconnect, increase retry time after each fail
    - currently the app auto terminates when the connection is unhealthy. Running the app as a service could solve the problem too
    https://medium.com/@benmorel/creating-a-linux-service-with-systemd-611b5c8b91d6
- handle forced reconnection (https://www.bitstamp.net/websocket/v2/)
- get balance in prepare, remove balance flag, add min order check $10 
- cancel all cmd
- dump on run start config - profit, etc ...
- progressive grid
    - Assume 50/50 usdt/usd
    - We have to historically calculate the 0 starting point. Exclude from the calculation sudden spikes (~accruing once a month)
    - Better clientOrderId tagging
        - [subacc]-[target_profit]-[nonce]
    - Bigger density around 0 and smaller profit
    - Orders which are getting more distant from 0 should have progressive profit
    - Make sheet representation of calculator
    - progressive profit
        - predefined (poc)
        - integral/progression calculator (prod)

Progresive grid:
```
------------------- 8
 
 
------------------- 7
 
 
------------------- 6
 
------------------- 5
 
------------------- 4
 
------------------- 3
------------------- 2
------------------- 1
------------------- 0  - starting point
------------------- -1 - acc-0.01%-nonce
------------------- -2 - ...
------------------- -3 - ...
 
------------------- -4 - acc-0.03%-nonce
 
------------------- -5 - ...
 
------------------- -6 - ...
 
 
------------------- -7 - acc-0.1%-nonce
 
 
------------------- -8 - acc-0.15%-nonce
```
- capital efficiency
    - do not cover the whole spectrum of the price fluctuations
        - trade around 0 point and in case of period when stuck in usdt/usd use thether.to 
        - machine learning for the most frequently visited ranges
    - think of on/off ramp for usdt/usd
        - bank acc
        - thether.to - min $100k
- concurrent http requests - nonce loses order due to unknown routine execution order
    - possible solution would be sending all request to channel and execute sequentially from single point. Probably we are doing it exactly like that atm

## Roadmap
- integrate im-mem order book for multiple pairs
- calculate 3way arbs