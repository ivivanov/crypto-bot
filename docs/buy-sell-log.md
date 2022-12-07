Setup:
Type: buy & sell
crypto-bot prepare --price 1.00011 --bank 96.6 --type 2 --grid 0.01 --orders 8
crypto-bot run --pair usdtusd --profit: 0.01

Prepare:
```
2022/12/05 22:06:59 Order->order_created-> buy: 12.08000 @ 1.00001 [testing_1.00001_119682624]
2022/12/05 22:06:59 Order->order_created-> buy: 12.08000 @ 0.99991 [testing_0.99991_837943528]
2022/12/05 22:07:00 Order->order_created-> buy: 12.08000 @ 0.99981 [testing_0.99981_1078055919]
2022/12/05 22:07:00 Order->order_created-> buy: 12.08000 @ 0.99971 [testing_0.99971_1499264013]
2022/12/05 22:07:00 Order->order_created-> sell: 12.08000 @ 1.00021 [testing_1.00021_1677184475]
2022/12/05 22:07:00 Order->order_created-> sell: 12.08000 @ 1.00031 [testing_1.00031_105845707]
2022/12/05 22:07:00 Order->order_created-> sell: 12.08000 @ 1.00041 [testing_1.00041_1692475190]
2022/12/05 22:07:00 Order->order_created-> sell: 12.08000 @ 1.00051 [testing_1.00051_393381479]
```

Log:
```
2022/12/05 22:26:57 Trade-> buy: 12.08000 @ 1.00001 [testing_1.00001_119682624]
2022/12/05 22:26:57 Order->order_deleted-> buy: 0 @ 1.00001 [testing_1.00001_119682624]
2022/12/05 22:27:02 Order->order_created-> sell: 12.08000 @ 1.00051 [testing_1.00051_1159997099]
2022/12/06 05:37:09 Order->order_deleted-> sell: 0 @ 1.00021 [testing_1.00021_1677184475]
2022/12/06 05:37:09 Trade-> sell: 12.08000 @ 1.00021 [testing_1.00021_1677184475]
2022/12/06 05:37:09 Order->order_created-> buy: 12.08000 @ 0.99971 [testing_0.99971_2085421470]
2022/12/06 09:53:24 Trade-> buy: 12.08000 @ 0.99991 [testing_0.99991_837943528]
2022/12/06 09:53:24 Trade-> buy: 12.08000 @ 0.99981 [testing_0.99981_1078055919]
2022/12/06 09:53:24 Trade-> buy: 12.08000 @ 0.99971 [testing_0.99971_1499264013]
2022/12/06 09:53:24 Trade-> buy: 12.08000 @ 0.99971 [testing_0.99971_2085421470]
2022/12/06 09:53:24 Order->order_deleted-> buy: 0 @ 0.99991 [testing_0.99991_837943528]
2022/12/06 09:53:24 Order->order_deleted-> buy: 0 @ 0.99981 [testing_0.99981_1078055919]
2022/12/06 09:53:24 Order->order_deleted-> buy: 0 @ 0.99971 [testing_0.99971_1499264013]
2022/12/06 09:53:24 Order->order_deleted-> buy: 0 @ 0.99971 [testing_0.99971_2085421470]
2022/12/06 09:53:24 Order->order_created-> sell: 12.08000 @ 1.00041 [testing_1.00041_612978178]
2022/12/06 09:53:24 Order->order_created-> sell: 12.08000 @ 1.00031 [testing_1.00031_359472267]
2022/12/06 09:53:24 Order->order_created-> sell: 12.08000 @ 1.00021 [testing_1.00021_1224302906]
2022/12/06 09:53:24 Order->order_created-> sell: 12.08000 @ 1.00021 [testing_1.00021_1509141783]

```

Outcome:
Orders are stacking on each other. What we prevent by having both sell & buy is just catching the initial move direction. After it the setup is identical to having only buy/sell prepare