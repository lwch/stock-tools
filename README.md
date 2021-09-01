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
    # 初选完毕，能够筛选到90只个股（截止到2021/9/1的数据）
    训练算法模型。。。还没实现

## download

数据下载工具，参数如下

1. `code`: 股票代号，cn_开头
2. `begin`: 开始时间，负数为向前追溯N天，否则为yyyymmdd格式

最终文件下载到当前工作目录下的cn_<股票代号>\_\<yyyymmdd>\_\<yyyymmdd>.csv

## info

获取单只个股的统计信息，传入csv文件所在路径，统计结果如下

    开盘价:
      均值=3.74, 中位数=3.40, 最后=3.40
      最小值=3.40, 最大值=4.50, 标准差=0.15
      P10=3.40, P70=3.40, P90=3.40
    收盘价:
      均值=3.75, 中位数=3.41, 最后=3.42
      最小值=3.41, 最大值=4.46, 标准差=0.16
      P10=3.41, P70=3.41, P90=3.41
      总天数: 244, 上涨天数: 108(44.26%), 下跌天数: 94(38.52%)

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