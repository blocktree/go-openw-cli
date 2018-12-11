# openw-cli

## 概述

[TOC]

## 修订信息

| 版本  | 时间       | 修订人 | 修订内容         |
|-------|------------|--------|------------------|
| 1.0.0 | 2018-12-11 | 麦志泉 | 创建文档         |

---

## 1. openw-cli介绍

openw-cli是一款命令行工具，通过调用openw-server钱包服务API实现主机客户端下的钱包管理。

---

## 2. 需求要点

> openw-cli实现了openw-server的API的方法调用，并提供一些openwallet常用方法。

| 需求             | 描述                                                                                      |
|----------------|-----------------------------------------------------------------------------------------|
| **节点管理相关** | 管理通信方面的功能                                                                        |
| 登记节点         | 如果还没生成通信密钥，则先生成密钥对，再向钱包平台登记节点。                                 |
| 查看节点         | 查看通信节点信息，如：节点ID，通信密钥对。                                                    |
| **钱包管理相关** | 管理钱包方面的功能                                                                        |
| 创建钱包流程     | 输入钱包别名，钱包密码，进行创建。                                                           |
| 查看钱包流程     | 查看已创建的钱包列表，如：钱包ID，钱包别名，钱包账户数。                                       |
| 创建资产账户流程 | 选择钱包，输入币种类型，账户别名，进行创建。目前只支持（单签）                                  |
| 查看资产账户流程 | 选择钱包，查看钱包下的资产账户列表，如：账户ID，资产种类，资产余额，地址总数，汇总地址，汇总阈值。 |
| 创建地址流程     | 选择钱包，选择资产账户，输入创建地址数量，批量生成地址。                                      |
| 查询地址         | 输入：地址，显示地址信息：所在钱包ID，账户ID，地址公钥，地址私钥，余额。                          |
| 资产转账         | 选择钱包，选择资产账户，输入：转账数量，目标地址，自定义费率。                                  |
| 设置汇总         | 选择钱包，显示资产账户列表，选择资产账户，设置：汇总地址，汇总阈值。                            |
| 定时汇总         | 启动汇总后台。                                                                             |
| **服务管理相关** | 搭建本地钱包服务（待研究）                                                                  |
| 创建服务         | 本地部署一个openw-server。                                                                 |
| 安装区块链全节点 | 通过简单命令执行，通过docker快速节点。                                                      |
| 启动服务         | 启动openw-server。                                                                         |

---

## 3. 功能详细设计 `openw-cli`

### 配置文件

使用openw-cli需要依赖配置文件，样例如下：

```ini

# Remote Server
remoteserver = "www.openwallet.site"

# API Version
version = "1.0.0"

# App Key
appkey = "1234qwer"

# App Secret
appkey = "qwer1234"

# Log file path
logdir = "/usr/logs/"

# Data directory, store keys, databases, backups
datadir = "/usr/data/"

# Wallet Summary Period
summaryperiod = "1h"

```

我们提供命令行工具openw-cli，以下功能点作为管理资产的【子命令】，附加以下参数变量。

### 全局可选参数

| 参数变量    | 描述                                                  |
|-------------|-----------------------------------------------------|
| -s, -symbol | 币种标识符，其后带值[symbol]，如btc，ltc，eth，ada，btm，sc  |
| -i, -init   | 是否初始化，应用于配置功能时候，是否需要执行初始化流程。 |
| -p, -path   | 指定文件目录。                                         |
| -f, -file   | 指定加载的文件。                                       |
| -debug      | 是否打印debug日志信息。                                |
| -logdir     | 指定日志输出目录。                                     |

### 文件目录结构

使用openw-cli管理区块链资产时会创建以下文件目录，存储一些钱包相关的数据。目录用途说明如下：

| 参数变量                  | 描述                                                                         |
|---------------------------|----------------------------------------------------------------------------|
| {datadir}/key/               | 钱包keystore文件目录，文件命名 [alias]-[WalletID].key                         |
| {datadir}/db/                | 钱包数据库缓存目录，文件命名 [alias]-[WalletID].db                              |
| {datadir}/backup/            | 钱包备份文件导出目录，以文件夹归档备份，文件夹命名 [alias]-[WalletID]-yyyyMMddHHmmss |

> 命令输入结构: openw-cli [命令模块] [子命令] [可选参数...]
> 如：openw-cli wallet newwallet -s btc

### 3.1 节点管理 `node`

#### 3.3.1 登记节点 `register`

1. 初始化通信密钥对，如果已存在，则提示是否需要覆盖。
1. 如果重新创建成功，客户端会重新登记节点。

#### 3.3.2 查看节点 `info`

1. 查看节点ID，通信密钥对。

### 3.2 钱包管理 `wallet`

#### 3.2.1 创建钱包 `newwallet`

1. 执行钱包创建流程，要求用户输入钱包别名，钱包密码等信息。
1. 创建成功后，钱包种子保存成keystore文件到 {datadir}/key/目录下。

#### 3.2.2 查看钱包列表 `listwallet`

1. 只查看本地节点创建的拥有钱包种子的所有钱包。
1. 显示钱包列表，显示钱包别名，钱包ID，钱包账户总数。

#### 3.2.3 创建资产账户 `newaccount`

1. 选择一个已创建的钱包。
1. 输入：账户别名，币种类型，下一步创建资产账户。

#### 3.2.4 查看资产账户列表 `listaccount`

1. 选择一个已创建的钱包。
1. 显示资产账户列表，显示账户别名，账户ID，币种类型，账户余额，地址总数，汇总地址，汇总阈值。

#### 3.2.5 创建地址 `newaddress`

1. 选择一个已创建的钱包，及已创建的账户。
1. 输入：地址数量，下一步批量创建新地址。

#### 3.2.6 创建地址 `searchaddress`

1. 输入一个地址，后台去查询地址。显示地址所在的钱包ID，账户ID，地址余额，地址公钥，地址私钥。

#### 3.2.7 资产转账 `transfer`

1. 选择一个已创建的钱包，及已创建的账户。
1. 输入：转账数目，转账地址，钱包密码。
1. 构建交易单，计算手续费，完成交易签名，广播交易。

#### 3.2.8 设置汇总 `setsum`

1. 选择一个已创建的钱包，及已创建的账户。
1. 输入：汇总地址，汇总阈值，保存设置。

#### 3.2.9 定时汇总 `startsum`

1. 选择需要汇总的钱包，输入要汇总的资产账户数组，输入钱包密码，启动汇总定时程序。
1. 通过-file 加载json文件，指定多个汇总钱包，接着输入各个钱包的密码，启动汇总定时程序。

`汇总样例JSON`

```json

{
    "wallets": [
        {
            "alias": "bit",
            "walletID": "1234qwer",
            "accounts": [
                "123",
                "4567",
            ]
        },
        {
            "alias": "bitw",
            "walletID": "1234qwer",
            "accounts": [
                "123",
                "4567",
            ]
        }
    ]
}

```

### 3.3 服务管理 `service`

自行搭建openw-server，待设计

---

## 4. openw-cli应用说明

### 4.1 编译openw-cli工具

```shell

# 进入目录
$ $GOPATH/src/github.com/blocktree/go-openw-cli

# 全部平台版本编译
$ xgo .

# 或自编译某个系统的版本
$ xgo --targets=linux/amd64 .

```

### 4.2 命令示例

```shell

#### 节点相关 ####

# 登记到openw-server，成为应用的授权节点。
$ ./openw-cli node register

# 查看节点的信息
$ ./openw-cli node info

# 执行来自配置文件关闭[symbol]节点的命令
$ ./openw-cli node stop -s [symbol]

# 移除节点docker容器
$ ./openw-cli node remove -s [symbol]

# 查看与[symbol]节点相关的信息
$ ./openw-cli node status -s [symbol]

# 查看./conf/[symbol].ini文件中与节点相关的配置信息
$ ./openw-cli node config -s [symbol]


#### 钱包相关 ####

# 创建钱包
$ ./openw-cli wallet newwallet

# 查看节点本地已创建的钱包
$ ./openw-cli wallet listwallet

# 创建钱包资产账户，先选择钱包
$ ./openw-cli wallet newaccount

# 查看钱包资产账户，先选择钱包
$ ./openw-cli wallet listaccount

# 创建新地址，先选择钱包，再选择资产账户
$ ./openw-cli wallet newaddress

# 查询地址信息
$ ./openw-cli wallet searchaddress

# 设置汇总，先选择钱包，再选择资产账户
$ ./openw-cli wallet setsum

# 启动汇总定时器
$ ./openw-cli wallet startsum

# 启动汇总定时器，通过文件加载需要汇总的钱包和资产账户
$ ./openw-cli wallet startsum -f /usr/to/sum.json

```