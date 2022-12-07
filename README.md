# Grid bot

## Installation From source
```
go install
```
## Start
- Generate api keys with the minimal possible privileges. It is recommended to use sub account.
- Create .env file
- Execute `crypto-bot` from the directory of .env

## ADRs

### 1. ClientOrderID format `[account]_[target_price]_[random]`
In case we hit market price (taker) when posting new order (trade_buy->create_sell could hit directly as market order). This will displace the initial buy order. For exa buy@1 sell@1.1 but market has ask@1.2 => actual sell@1.2 =>  new buy@1.1 (wanted buy@1). For this reason I am adding more order metadata to the client order id.
Examples in docs.

- [todo] create sell order 0.0001 above the current bid (order book needed). This will increase profit by not paying taker fee + selling above the desired target

## TODOs
Easy picks:
- get balance in prepare, remove balance flag, add min order check $10
    - consider situation where you don not want the whole stack to be allocated

- make ctx struct
    - add configured logger. remove date & add ms

- concurrent http requests - nonce loses order due to unknown routine execution order
    - test again UUID - last time could be expired/blocked api keys - regen just in case
    - possible solution would be sending all request to channel and execute sequentially from single point. Probably we are doing it exactly like that atm

Backlog:
- fill partial below min order. If I have order for $10 and only $5 got filled PostCounterTrade would fail
    - keep in memory [price]->[amount_filled] and when enough for min order open it
- auto ws reconnect
    - currently the app auto terminates when the connection is unhealthy. Running the app as a service could solve the problem too
    https://medium.com/@benmorel/creating-a-linux-service-with-systemd-611b5c8b91d6
- handle forced reconnection (https://www.bitstamp.net/websocket/v2/)

- progressive profit
    - predefined (poc)
    - integral/progression calculator (prod)
    - the more distant from SMA is the order the bigger profit to the counter trade should be applied - sell with progressive profit
    ```
    SMA 1.00000
    buy@0.99981 -> sell@1.00031 (0.01% profit)
    ...
    buy@0.99800 -> sell@0.99940 (0.10% profit)
    ```
    
    - sell always on the SMA price. Much easier to implement
    ```
    SMA 1.00000 
    buy@0.99990 -> sell@1.00040 (SMA price not breakeven - apply fixed profit 0.01%)
    buy@0.99965 -> sell@1.00015 (SMA price not breakeven - apply fixed profit 0.01%)
    ...
    buy@0.99910 -> sell@1.00000 (0.05% profit, if fixed 0.01% is applied: 0.99960)
    buy@0.99890 -> sell@1.00000 (0.07% profit, if fixed 0.01% is applied: 0.99940)
    ```
- progressive grid
    - Assume 50/50 usdt/usd
    - We have to historically calculate the 0 starting point. Exclude from the calculation sudden spikes (~accruing once a month)
    - Bigger density around 0 and smaller profit
    - Orders which are getting more distant from 0 should have progressive profit
    - Make sheet representation of calculator
   

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
        - trade around 0 point and in case of period when stuck in usdt/usd use thether.
        - range can be easily determined by MVC (moving average channel) and MA
        - bollinger band for upper and lower bound
    - think of on/off ramp for usdt/usd
        - bank acc
        - thether.to - min $100k

## Roadmap
- reduce fees as low as possible
- implement Kraken
- integrate im-mem order book for multiple pairs
- calculate 3way arbs