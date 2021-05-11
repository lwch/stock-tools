# stock-tools

股票分析工具

## download

数据下载工具，参数如下

1. `code`: 股票代号，cn_开头
2. `begin`: 开始时间，负数为向前追溯N天，否则为yyyymmdd格式

最终文件下载到当前工作目录下的cn_<股票代号>\_\<yyyymmdd>\_\<yyyymmdd>.csv

## calc

数据统计，传入csv文件所在路径，统计结果如下

    open:
      avg=6.11, mean=5.90
      min=4.47, max=9.20, stddev=1.37
      P10=4.78, P70=6.54, P90=8.32
    close:
      avg=6.13, mean=5.91
      min=4.50, max=9.19, stddev=1.39
      P10=4.78, P70=6.54, P90=8.49

open、close分别表示开盘和收盘数据

avg: 平均值
mean: 中位数
min: 最小值
max: 最大值
stddev: 标准差
P(N): 分位数