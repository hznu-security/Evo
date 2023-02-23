# 大屏模块

建立连接：server为client生成一个id，返回id以及基础信息

ping/pong机制：server不定时向客户端发送ping消息，client回复pong，重置计时器超时就要将client移除，重新建立连接。

对于每一个新的连接，起一个协程去和它通信
再起一个协程用来管理所有连接

管理协程和连接协程之间通过通道通信

### 接口

初始化，建立连接时发送

```json
{
  "Type": "init",
  "Data": {
    "Title": "哈哈哈",
    "Time": -7200,
    "Round": 1,
    "Teams": [
      {
        "Id": 1,
        "Name": "aaa队",
        "Rank": 1,
        "Img": "/upload/logo/gIyLSZyTBLbmqJfncuOoTQkKBQSM.png",
        "Score": 0
      },
      {
        "Id": 2,
        "Name": "bbbbs队",
        "Rank": 2,
        "Img": "",
        "Score": 0
      },
      {
        "Id": 3,
        "Name": "ccc队",
        "Rank": 3,
        "Img": "",
        "Score": 0
      }
    ]
  }
}
```

清除某个队伍的状态

```json
{
  "Type": "clear",
  "Data": {
    "Id": 1
  }
}
```

清除所有队伍的状态

```json
{
  "Type": "clearAll",
  "Data": null
}
```

本轮剩余时间，单位 秒

```json
{
  "Type": "time",
  "Data": {
    "Time": 23
  }
}
```

排名

```json
{
  "Type": "rank",
  "Data": {
    "Teams": [
      {
        "Id": 1,
        "Name": "aaa队",
        "Rank": 1,
        "Img": "/upload/logo/gIyLSZyTBLbmqJfncuOoTQkKBQSM.png",
        "Score": 0
      },
      {
        "Id": 2,
        "Name": "bbbbs队",
        "Rank": 2,
        "Img": "",
        "Score": 0
      },
      {
        "Id": 3,
        "Name": "ccc队",
        "Rank": 3,
        "Img": "",
        "Score": 0
      }
    ]
  }
}
```

状态，只有两种，down和attacked

```json
{
  "Type": "status",
  "Data": {
    "Id": 1,
    "Status": "down"
  }
}
```

```json
{
  "Type": "status",
  "Data": {
    "Id": 1,
    "Status": "attacked"
  }
}
```

攻击
```json
{
  "Type": "attack",
  "Data": {
    "From": 1,
    "To": 2
  }
}
```