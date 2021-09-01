# stock-tools

股票分析选股工具

## 选股示例

    # 编译所有程序
    make
    # 准备操作目录
    mkdir data
    cd data
    # 获取所有个股最近1年数据
    ../bin/list|xargs -L1 -I{} ../bin/download -code cn_{} -begin -365
    # 筛选标准差小于0.2的个股（标准差越小说明股价波动越小，风险越小）
    ../bin/filter -begin -365 -func stdev -lt 0.2 -out close_stdev_lt_0.2 stock
    # 筛选出均价小于10元的个股
    ../bin/filter -func avg -lt 10 -out close_avg_lt_10 close_stdev_lt_0.2
    # 初选完毕，能够筛选到90只个股（截至到2021/9/1的数据）

## download

数据下载工具，参数如下

1. `code`: 股票代号，cn_开头
2. `begin`: 开始时间，负数为向前追溯N天，否则为yyyymmdd格式

最终文件下载到当前工作目录下的cn_<股票代号>\_\<yyyymmdd>\_\<yyyymmdd>.csv

## info

获取单只个股的统计信息，传入csv文件所在路径，统计结果如下

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

## filter

过滤工具，根据给定过滤条件筛选出符合条件的文件，参数如下：

1. `begin`: 开始时间，负数为向前追溯N天，否则为yyyymmdd格式
2. `column`: 筛选字段，支持open,close,low,high
3. `func`: 聚合函数，支持max,min,sum,avg,stdev,p\<n\>，其中p\<n\>表示分位数
4. `gt`: 筛选聚合后大于该值的数据
5. `lt`: 筛选聚合后小于该值的数据
6. `out`: 数据结果目录
7. `其他`: 传入原始数据路径

示例一、筛选最近30天收盘均价小于10元的股票

    ./filter -begin -30 -column close -func avg -lt 10 -out <输出路径> <输入路径>

示例二、筛选最近30天收盘价的标准差小于0.1的股票

    ./filter -begin -30 -column close -func stdev -lt 0.1 -out <输出路径> <输入路径>

示例三、筛选最近30天收盘价的中位数小于10元的股票

    ./filter -begin -30 -column close -func p50 -lt 10 -out <输出路径> <输入路径>