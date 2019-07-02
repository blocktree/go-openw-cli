# openw-cli
[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)

## 概述

[TOC]

## Build development environment

The requirements to build OpenWallet are:

- Golang version 1.11 or later
- Properly configured Go language environment
- Golang supported operating system

## 源码工具

### 安装xgo（支持跨平台编译C代码）

[官方github（目前还不支持go module）](https://github.com/karalabe/xgo)
[支持go module的xgo fork](https://github.com/gythialy/xgo)

xgo的使用依赖docker。并且把要跨平台编译的项目文件加入到File sharing。

```shell

# 官方的目前还不支持go module编译，所以我们用了别人改造后能够给支持的fork版
$ go get -u github.com/gythialy/xgo
...
$ xgo -h
...

# 本地系统编译
$ make clean build

# 跨平台编译wmd，更多可选平台可修改Makefile的$TARGETS变量
$ make clean openw-cli


```
---

## openw-cli介绍

openw-cli是一款命令行工具，通过调用openw-server钱包服务API实现主机客户端下的钱包管理。

---

## 功能详细设计 `openw-cli`

### 配置文件

使用openw-cli需要依赖配置文件，样例如下：

```ini

# Remote Server
remoteserver = "api.openwallet.cn"

# API Version
version = "1.0.0"

# App ID
appid = "1234qwer"

# App key
appkey = "qwer1234"

# Log file path
logdir = "/usr/logs/"

# Data directory, store keys, databases, backups
datadir = "/usr/data/"

# Wallet Summary Period
summaryperiod = "1h"

# The custom name of local node
localname = "blocktree"

# Be trusted client server
trustedserver = "client.blocktree.top"

# Enable client server request local transfer
enablerequesttransfer = false

# Enable client server execute summary task
enableexecutesummarytask = false

# Enable client server edit wallet summary settings
enableeditsummarysettings = false

# Enable key agreement on local node communicate with client server
enablekeyagreement = true

# Enable https or wss
enablessl = false

# Network request timeout, unit: second
requesttimeout = 60

# Terminal print log of debug 
logdebug = false

# Enable trusted server connect with https or wss
enabletrustserverssl = false

```

我们提供命令行工具openw-cli，以下功能点作为管理资产的【子命令】，附加以下参数变量。

### 全局可选参数

| 参数变量    | 描述                                                  |
|-------------|-----------------------------------------------------|
| -s, -symbol | 币种标识符，其后带值[symbol]，如btc，ltc，eth，ada，btm，sc  |
| -c, -conf   | 工具配置文件路径。                                     |
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

> 命令输入结构: openw-cli [配置文件] [子命令] [可选参数...]
> 如：openw-cli -c ./node.ini newwallet -s btc

### 命令示例

```shell

# 临时生成随机的OWTP通信证书，不保存到本地缓存
$ ./openw-cli genkeychain

# 通过-c或-conf设置工具的配置文件路径
$ ./openw-cli -c=./node.ini

#### 节点相关 ####

# 登记到openw-server，成为应用的授权节点。
$ ./openw-cli -c=./node.ini noderegister

# 查看节点的信息
$ ./openw-cli -c=./node.ini nodeinfo

# 更新区块链资料
$ ./openw-cli -c=./node.ini updateinfo

#### 钱包相关 ####

# 创建钱包
$ ./openw-cli -c=./node.ini newwallet

# Enter wallet's name: MyWallet         //输入钱包名字
# Enter wallet password:                //输入钱包密码
# //创建成功后，显示钱包种子文件
# Wallet create successfully, key path: openw/data/key/NASSUM-W6zkTDtnWZWFd2SQPms9F62BBPfuqU2ETg.key

# 查看节点本地已创建的钱包
$ ./openw-cli -c=./node.ini listwallet

# 创建钱包资产账户，先选择钱包
$ ./openw-cli -c=./node.ini newaccount


# [Please select a wallet]
# Enter wallet No.: 0               //输入No.序号，选择本地已有钱包
# Enter wallet password:            //输入密码，解锁钱包
# Enter account's name: NASSUM      //输入需要创建资产账户名
# Enter account's symbol: NAS       //输入需要创建的币种symbol
# //创建成功，默认显示资产账户ID和默认接收地址
# create [NAS] account successfully
# new accountID: 9HqxxcNSMxdt225Dis3mdnzT18egbV7Cg3R85y6AUPx8
# new address: n1EZVYXBx5tQ41L6QRyEhpqV4TpH6NwPrPE

# 查看钱包资产账户，先选择钱包
$ ./openw-cli -c=./node.ini listaccount

# 创建新地址，先选择钱包，再选择资产账户
$ ./openw-cli -c=./node.ini newaddress


# [Please select a wallet]
# Enter wallet No.: 0                               //输入No.序号，选择本地已有钱包
# [Please select a account]
# Enter account No.: 0                              //输入No.序号，选择钱包已有的资产账户
# Enter the number of addresses you want: 100       //输入需要创建地址的数量
# create [100] addresses successfully
# //创建地址成功，把新地址导出以下文件路径。
# addresses has been exported into: openw/data/export/address/[9HqxxcNSMxdt225Dis3mdnzT18egbV7Cg3R85y6AUPx8]-20190313163227.txt

# 创建新地址，先选择钱包，再选择资产账户，输入offset和limit查询地址列表
$ ./openw-cli -c=./node.ini searchaddress


# 选择资产账户，发起转账交易
$ ./openw-cli -c=./node.ini transfer

# [Please select a wallet]                                          //输入No.序号，选择本地已有钱包
# Enter wallet No.: 0
# [Please select a account]                                         //输入No.序号，选择钱包已有的资产账户
# Enter account No.: 0
# Enter contract address:                                           //如果是代币合约转账，输入合约地址，默认空
# Enter received address: AR8LWKndC2ztfLoCobZhHEwkQCUZk1yKEsF       //输入转账的目标地址
# Enter amount to send: 2.8                                         //输入转账数量
# Enter fee rate:                                                   //输入手续费率，默认空（推荐）
# Enter memo: hello                                                 //输入备注，适用于可添加备注的交易
# Enter wallet password:                                            //输入钱包解锁密码
# -----------------------------------------------                   //以下为转账日志信息
# [VSYS  Transfer]
# From Account: 33yYCwSeBump7V6AFX8r2KFsXqkJ7zxkg36UxyvYgy1o
# To Address: AR8LWKndC2ztfLoCobZhHEwkQCUZk1yKEsF
# Send Amount: 2.8
# Fees: 0.1
# FeeRate: 10000000
# Memo: hlls
# -----------------------------------------------
# send transaction successfully.
# transaction id: GbB1oQkXQTSDudTEhKwdhyvnUvunnHKngZqGE9Xfa3tn

# 选择资产账户，转账账户下所有地址的资产到目标地址
$ ./openw-cli -c=./node.ini transferall

# [Please select a account]
# Enter account No.: 15
# Enter contract address:
# Enter received address: TRJJ9Mq4aMjdmKWpTDJAgbYNoY2P9Facg5        //目标地址
# Enter fee rate:
# Enter memo:
# Enter wallet password:
# Summary account[9Z4ivqhr5mniEBzpxkRQn7w5r1YgeoB6qTLLas3VV6Z] Symbol: TRX, token:
# Summary account[9Z4ivqhr5mniEBzpxkRQn7w5r1YgeoB6qTLLas3VV6Z] Current Balance = 9.9
# Summary account[9Z4ivqhr5mniEBzpxkRQn7w5r1YgeoB6qTLLas3VV6Z] Summary Address = TRJJ9Mq4aMjdmKWpTDJAgbYNoY2P9Facg5
# Summary account[9Z4ivqhr5mniEBzpxkRQn7w5r1YgeoB6qTLLas3VV6Z] Start Create Summary Transaction
# Create Summary Transaction in address range [0...200]
# [Success] txid: 4106a2a1dff1647d4e12b14d181ed45d4c847d710e6685588125674a481c42af
# Save summary task log successfully

# 设置汇总，先选择钱包，再选择资产账户
$ ./openw-cli -c=./node.ini setsum


# [Please select a wallet]
# Enter wallet No.: 0                               //输入No.序号，选择本地已有钱包
# [Please select a account]
# Enter account No.: 0                                                          //输入No.序号，选择钱包已有的资产账户
# Enter account's summary address: n1EZVYXBx5tQ41L6QRyEhpqV4TpH6NwPrPE          //输入钱包汇总转账到的地址
# Enter account's summary threshold: 2                                          //输入汇总阈值，账户总余额超过此值，执行汇总交易
# Enter address's minimum transfer amount: 0.01                                 //输入地址最低转账额，地址余额超过此值才发起转账
# Enter address's retained balance: 0                                           //输入地址保留余额，地址转账时需要剩余部分余额
# Enter how many confirms can transfer: 1                                       //输入地址未花得到多少确认后才可用于转账交易
# setup summary info successfully

# 查看已有账户的汇总设置信息
$ ./openw-cli -c=./node.ini listsuminfo

# 启动汇总定时器
$ ./openw-cli -c=./node.ini startsum

# Enter summary task json file path:            //输入汇总任务json文件，如果为空，则提供选择钱包和资产账户启动汇总


# [Please select a wallet]
# Enter wallet No.: 0                               //输入No.序号，选择本地已有钱包
# [Please select a account]
# Enter account No.: 0                              //输入No.序号，选择钱包已有的资产账户

# //汇总启动成功，定时执行任务
# The timer for summary task start now. Execute by every 10 seconds.
# [Summary Task Start]------2019-03-13 16:43:33
# Summary account[9HqxxcNSMxdt225Dis3mdnzT18egbV7Cg3R85y6AUPx8] Symbol: NAS start
# Summary account[9HqxxcNSMxdt225Dis3mdnzT18egbV7Cg3R85y6AUPx8] Current Balance: 0, below threshold: 2
# Summary account[9HqxxcNSMxdt225Dis3mdnzT18egbV7Cg3R85y6AUPx8] Symbol: NAS end
# [Summary Task End]------2019-03-13 16:43:34

# 启动汇总定时器，通过文件加载需要汇总的钱包和资产账户
$ ./openw-cli -c=./node.ini startsum -f=/usr/to/sum.json

```

```json

`汇总样例JSON`

> 汇总任务配置，可重新制定每个账户的汇总信息(汇总地址只能通过setsum更改，安全考虑)。

{
    "wallets": [
        {
            "walletID": "1234qwer", //钱包ID
            "password": "12345678", //钱包解锁密码，不填时，执行命令时要求输入
            "accounts": [ //需要汇总的账户列表
                {
                    "accountID": "123",               //资产账户ID
                    "threshold": "1000",              //账户总阈值
                    "minTransfer": "1000",            //地址最低转账额
                    "retainedBalance": "0",           //地址保留余额
                    "confirms": 1,                    //未花大于该确认次数汇总中
                    "feeRate": "0.0001",              //交易费率，填空为推荐费率
                    "onlyContracts": false,           //只汇总代币, 为true时，contracts数组必须有值
                    "switchSymbol": "ETH",            //强制切换symbol（默认不设置，用于解决代币打错地址，需要切换网络汇总）
                    "memo": "hello",                  //备注，适用于可添加备注的交易单
                    "contracts": {                            //汇总代币合约
                        "all": {                              //全部合约
                            "threshold": "1000",              //账户总阈值
                            "minTransfer": "1000",            //地址最低转账额
                            "retainedBalance": "0"            //地址保留余额
                        },         
                        "0x1234dcba": {                        //指定的合约地址或编号
                            "threshold": "1000",      
                            "minTransfer": "1000",
                            "retainedBalance": "0"
                        },
                    },
                    "feesSupportAccount": {           //主币余额不足时，可选择一个同钱包下的账户提供手续费
                        "accountID": "12323",         //同钱包下的账户ID
                        "lowBalanceWarning": "0.5",   //手续费账户余额过低警告阈值，打印警告
                        "lowBalanceStop": "0.1",      //手续费账户余额过低停止工作阈值
                        "fixSupportAmount": "2",      //手续费不足时，提供固定的数量支持
                        "feesScale": "2"              //手续费不足时，提供所缺手续费 * 倍数
                        "isTokenContract": true,      //是否用代币合约做手续费
                        "contractAddress": "0200000000000000000000000000000000000000"   //代币合约地址，ONT使用ONG作为手续费
                    }
                }
            ],
        }
    ]
}

```

```shell

# 查询主链列表
$ ./openw-cli -c=./node.ini listsymbol

# 查询主链下的合约列表
$ ./openw-cli -c=./node.ini listtokencontract

# 选择钱包及账户，查看账户下拥有的代币余额
$ ./openw-cli -c=./node.ini listtokenbalance

# 启动后台托管钱包服务，执行时会要求是否解锁钱包
$ ./openw-cli -c=./node.ini trustserver

```