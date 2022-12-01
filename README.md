## Installation From source
```
go install
```
## Start
Generate api keys with the minimal possible privileges. It is recommended to use sub account.
```
crypto-bot --key=[api_key] --secret=[api_secret]
```

### TODOs
- auto reconnect, increase retry time after each fail, handle health check reconnect message
- concurrent http requests - nonce loses order due to unknown routine execution
- get balance in prepare, remove balance flag, add min order check $10 
- integrate im-mem order book for multiple pairs
- calculate 3way arbs
- console override with the last health check
