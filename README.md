# Digger

![](ui/src/assets/images/logo.png)

# 简介

[Digger](https://github/hetianyi/digger)是用纯[Golang](https://golang.org)开发的配置式分布式跨平台爬虫系统，支持使用Javascript编写插件来实现各种你想要达到的目标。Digger及相关组件能够以极低的资源开销运行在各种廉价服务器和开发板上，如树莓派。

Digger没有复杂的依赖，部署十分简单，支持Linux和Windows平台，目前支持的CPU架构有：```amd64```, ```arm```, ```arm64```

## 功能简介

- 支持使用Css选择器和Xpath选择器
- 支持多种结果类型：plain text，html，array等
- web端爬虫配置编辑器
- 在线调试爬虫配置，精准定位问题
- 支持插件功能
- 实时浏览爬虫日志
- 结果在线浏览、导出，一键生成数据库schema（postgres和mysql）
- 定时任务
- 支持暂停任务
- 分布式worker实例，有效避免爬虫被block
- 支持任务和worker标签匹配调度功能
- 支持配置导入导出
- 邮件通知功能
- 钉钉通知功能（TODO）
- DiggerHub支持爬虫配置分享（TODO）


>  更多功能敬请探索！

# 安装
## 使用docker安装（推荐）
使用docker部署digger是最简单的方式，你可以在[Docker Hub](https://hub.docker.com/repository/docker/hehety/tiller)下载digger的镜像。

这里以单机作为演示部署完整的digger服务。

**部署环境**
- 操作系统：CentOS 7
- Docker Version: 19.03.11
- 关闭防火墙Firewalld
- 确保9012端口没有被占用


克隆git仓库：
```shell
git clone https://github.com/hetianyi/digger
```
使用docker swarm部署集群：
```shell
cd digger
docker stack deploy -c docker-compose.yml --prune digger
```
如果集群启动成功，就可以在9012端口使用浏览器访问digger控制台。

## 使用二进制安装
从二进制可执行程序安装比docker方式稍微繁琐一些，需要自行手动部署postgres数据库，redis以及manager，ui，worker等
postgres和redis安装过程此处不再赘述，这里假设这两个服务已存在。
首先，需要下载可执行程序```digger```和前端ui压缩包，下载地址在[此处](https://github.com/hetianyi/digger/releases)，
选择版本下载完成后，你应该得到以下两个文件：

```shell
│ ..
├ digger_0.0.1_linux_amd64.tar.gz
└ ui_0.0.1.tar.gz
```
将这两个文件解压，得到一个可执行程序<kbd>digger</kbd>和一个目录<kbd>ui</kbd>，将可执行文件复制到```/usr/local/bin```目录下，然后将ui.tar.gz解压，放到```/usr/share/digger```下：
```shell
tar xzf digger_0.0.1_linux_amd64.tar.gz
cp digger /usr/local/bin
tar xzf ui_0.0.1.tar.gz
mkdir -p /usr/share/digger
cp ui /usr/share/digger/
```

**启动manager**
```shell
digger manager --ui-dir /usr/share/digger/ui -d "postgres://user:password@127.0.0.1:5432/digger?sslmode=disable" -r "redispassword@127.0.0.1:6379#0"
```
> 这里```-d```指示的为postgres数据库连接字符串，```-r```为redis连接字符串，需要将相关参数替换为你的实际值

**启动worker**
```shell
digger worker --id=1 -m http://127.0.0.1:9012
```

manager和worker都启动成功，那么就可以使用浏览器通过9012端口访问web控制台了。

## 使用Kubernetes安装（TODO）

### 普通安装

### Helm安装




# 第一个爬虫项目
> 该节介绍如何开始你的第一个爬虫项目，请确保有一个搭建完整的digger服务

浏览器访问web控制台，进入首页，点击新建项目

![20200828111557.png](docs/screenshot/20200828111557.png)

在弹窗里填写项目基本信息

![20200828111730.png](docs/screenshot/20200828111730.png)

填写完成，点击确认按钮提交，项目创建完成！接下来，我们需要编写项目的爬虫配置文件。点击按钮编辑配置
![20200828113339.png](docs/screenshot/20200828113339.png)

在弹出页面中即可开始进行配置编辑，digger的爬虫配置使用```yaml```格式，如下图
![20200828113511.png](docs/screenshot/20200828113511.png)

这里我们以爬取网易云音乐的[歌单](https://music.163.com/#/discover/playlist)为例，在编辑页面底部，点击导入按钮，导入我们的示例配置，可以在这里下载我们提供的[示例配置文件](docs/demo/hello_world.cfg.json)。

自此，示例项目已配置完成，回到项目页面，点击“+”按钮开始一个新任务。
![20200828130659.png](docs/screenshot/20200828130659.png)
任务创建成功之后页面自动跳转到任务详情页
![20200828131350.png](docs/screenshot/20200828131350.png)
在任务页面点击"眼睛图标"的按钮，查看爬虫结果。
![20200828131723.png](docs/screenshot/20200828131723.png)
> 恭喜，你的第一个爬虫项目已成功运行，此时，你可以在结果页面将结果导出为各种格式。更多功能敬请探索！


# Digger启动参数
## 配置manager
查看manager启动帮助文档
```shell
~ # digger manager --help

Usage: digger manager [command options] [arguments...]

start as manager

Options:

   --secret value, -s value     custom global secret (default: "123456")
   --port value, -p value       server http port (default: 9012)
   --log-level value, -l value  set log level, available options:
                                (trace|debug|info|warn|error|fatal)
   --log-dir value              log directory (default: "/var/log/digger")
   --database value, -d value   postgres connection string, format:
                                postgres://<user>:<password>@<host>:<port>/<db>?sslmode=disable
   --redis value, -r value      redis connection string, format:
                                <password>@<host>:<port>#<db> (default: "@127.0.0.1#1")
   --labels value               node labels, format:
                                key1=value1<,key2=value2>
   --ui-dir value               ui dir (default: "/var/www/html")
   
```

| 参数 | 缩写 | 类型 | 说明 |
| -- | -- | -- | -- |
| ```--secret``` | ```-s``` | string | 密码用于加密token |
| ```--port``` | ```-p``` | int | manager监听端口 |
| ```--log-level``` | ```-l``` | string | 日志级别，可选值：```debug```, ```info```, ```warn```, ```error```, ```fatal``` |
| ```--log-dir``` | - | string | 日志文件目录 |
| ```--database``` | ```-d``` | string | postgres数据库连接，例如：```postgres://user:pass@127.0.0.1:5432/digger?sslmode=disable``` |
| ```--redis``` | ```-r``` | string | redis连接，例如：```pass@demo.com:6379#0``` |
| ```--ui-dir``` | - | string | 前端ui目录 |



## 配置worker
查看worker启动帮助文档
```shell
~ # digger worker --help

Usage: digger worker [command options] [arguments...]

start as worker

Options:

   --id value                 unique instance id (default: 1)
   --log-level value          set log level, available options:
                              (trace|debug|info|warn|error|fatal)
   --log-dir value            log directory (default: "/var/log/digger")
   --manager value, -m value  manager address, format:
                              <host>:<port> (default: "localhost:9012")
   --labels value             node labels, format:
                              key1=value1<,key2=value2>
   
```

| 参数 | 缩写 | 类型 | 说明 |
| -- | -- | -- | -- |
| ```--id``` | - | int | 标识该worker的唯一id |
| ```--log-level``` | ```-l``` | string | 日志级别，可选值：```debug```, ```info```, ```warn```, ```error```, ```fatal``` |
| ```--log-dir``` | - | string | 日志文件目录 |
| ```--manager``` | ```-m``` | string | manager地址，例如：```http://localhost:9012``` |
| ```--labels``` | - | string | worker的标签，用于任务匹配，例如：```app=demo,zone=east-cn``` |

# 爬虫项目配置文件
## 爬虫配置文件介绍
digger爬虫采用yaml文件配置，yaml文件结构清晰，易于编写。更多yaml入门知识，请参阅[菜鸟入门教程](https://www.runoob.com/w3cnote/yaml-intro.html)。
digger爬虫配置文件主要由6部分组成：
- start: 起始部分，标识爬虫从哪里开始，可配置多个start url
```yaml
start_urls:
- https://music.163.com/discover/playlist
start_stage: list
```
- stage: 阶段，表示爬虫进行到了哪个环节，或者页面层次
```yaml
stages:
- name: 
  is_list: 
  is_unique: 
  list_css: 
  page_css: 
  page_attr: 
  plugin: 
  fields:
```
- fields: 字段，字段表示最终要提取的目标数据对象
```yaml
fields:
- name: cover
  is_array: false
  is_html: false
  css: 
  attr: 
  plugin: ""
  remark: 
  next_stage: 
```
- headers: 全局http头部，用于爬虫发送http请求时附带
```yaml
headers:
  User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML,like
    Gecko) Chrome/78.0.3904.108 Safari/537.36
```
- settings: 全局配置，用于控制爬虫行为
```yaml
settings:
  CONCURRENT_REQUESTS: "5"
```
- node_affinity: 节点调度亲和标签，只有和worker的的标签匹配才能将任务调度到该worker
```yaml
node_affinity:
  "zone": "east-cn"
```


## 爬虫配置参数释义
```yaml
start_url: 
- https://music.163.com/discover/playlist # 爬虫起始URL
start_stage: list # 起始stage
stages: 
- name: list # stage的name
  is_list: true # 是否是列表类型
  is_unique: false # 页面URL是否是唯一（暂时可以忽略）
  list_css: ul#m-pl-container>li # css选择器表达式
  page_css: a.znxt # 分页按钮的css选择器表达式（如果有分页）
  page_attr: href # 分页按钮的url标签属性（通常是href）
  plugin: "" # 插件，请参考插件一节
  fields:
  - name: cover # 字段name
    is_array: false # 指示该字段是否是一个数组，如标签，组图等一个字段需要匹配多个值的场景
    is_html: false # 指示该字段是否是匹配标签下的原始html内容，对于字段需要提取原始html内容的场景非常有用
    css: div.u-cover>a # 字段的css选择器表达式
    attr: href # 字段的标签属性
    plugin: "" # 插件，请参考插件一节
    remark: 歌单地址 # 字段备注
    next_stage: detail # 下一阶段，将该字段的结果作为下一阶段的输入，例如：列表页提取的详情页URL，下阶段可以是详情页的stage
- name: detail
  is_list: false
  is_unique: false
  list_css: ""
  page_css: ""
  page_attr: ""
  plugin: ""
  fields:
  - name: cover
    is_array: false
    is_html: false
    css: div.cover>img
    attr: src
    plugin: ""
    remark: 歌单封面
    next_stage: "" 
headers: # 爬虫请求时会携带的http头部
  User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML,like
    Gecko) Chrome/78.0.3904.108 Safari/537.36
settings: 
  CONCURRENT_REQUESTS: "5" # 全局最大并发请求数
  FOLLOW_REDIRECT: "false" # 是否跟随重定向
  REQUEST_TIMEOUT: "60" # 请求超时时间(s)
  RETRY_COUNT: "3" # 重试次数（单节点，非全局）
  RETRY_WAIT: "0" # 重试间隔时间(s)
  SKIP_TLS_VERIFY: "false" # 是否跳过tls验证，解决自谦证书问题
node_affinity: # 节点亲和标签列表
  "": ""

```
## 调试爬虫
爬虫调试是一个非常实用的功能，能够帮助你在编写爬虫配置文件的过程中快速判断某一个阶段的配置是否编写正确，而不必等到实际运行时才能判断是否有问题。
如果没有调试功能，那么对于一些大型或长期任务可能是个灾难。以“第一个爬虫项目”里的示例配置为例，打开爬虫配置编辑页面，在右下角点击调试按钮

![20200828203733.png](docs/screenshot/20200828203733.png)

在弹出的调试对话框里输入调试参数，通常是一个stage的页面URL，然后在下拉框中选择调试的stage，这里选择```list```。
> 注意：输入的调试参数URL必须属于选择的stage，否则无法得到预期的调试结果。

![20200828205258.png](docs/screenshot/20200828205258.png)

点击“调试按钮”，在文本框里会出现调试得到的结果，该结果是未经处理的json字符串，其中```next```表示得到的下一阶段或下一页的结果，```result```表示爬虫结果，
通过调试结果中各个字段的对应结果就能判断css选择器表达式或其他配置是否正确。

## 爬虫配置文件导入导出

在爬虫配置编辑页面的右下角有两个按钮，能够帮助你将外部配置文件导入到当前项目，也可以将当前配置导出为json文件进行备份。
> 注意：导入配置会覆盖当前项目的所有配置和插件，请谨慎操作！

# 插件
## 插件介绍
插件是digger爬虫中引入的一个非常新颖和实用的一个功能，能够帮助爬虫工程师在不修改平台源代码的情况下，通过简单的javascript代码来实现各种你想要达到的效果。
插件的应用场景涵盖非常多的方面，例如数据清洗，纠正，剪切，转换，甚至自定义爬虫请求。
插件能够应用于stage和field之上，能够在一个请求的各个阶段影响爬虫结果。每个阶段我们称之为一个插槽点(slot)，插件目前具有以下几个插槽点：
- s1
http请求之前拦截URL，返回处理之后的URL
- sr
声明的sr阶段插件能够替换为自定义AJAX请求，默认是GET方法，可以在此POST请求等
- s2
http请求之后，引擎处理之前，输入为http请求结果，输出为处理后的结果，适用于一些结果需要预处理的场景，例如，一些restful接口返回的结果json对象中包含html，需要将这个html提取出来再由结果引擎处理。
- s3
TODO 处理中，处理使用的引擎，默认为goquery，否则为自定义的插件处理引擎（待开发）
- s4
在解析得到field值(包括stage的分页值)之后，可以用来修正数据，如去空格，剪切等

插件声明在```stage```和```field```的```plugin```参数中，格式为：```plugin1@s1,plugin2@sr,plugin3@s4...```

> 注意：一个```stage```可以引入多个插件，而```field```只能引入一个插件

下面是插件作用示意图

![digger插件.png](docs/screenshot/digger插件.png)

## 编写自定义插件
本节介绍如何编写自定义插件。

以“第一个爬虫项目”里的示例配置为例，爬取的歌单信息结果里，字段```create_date```理想情况下为```2020-03-27```，尾部却多了“创建”，我们可以通过编写一个非常简单的插件来去除它。

![20200828220056.png](docs/screenshot/20200828220056.png)

打开爬虫配置编辑页面，切换到插件页面，点击“+”按钮，在弹框中输入插件名称，在这里我们将插件命名为```remove_suffix```，然后点击“添加”按钮确认添加

![20200828221148.png](docs/screenshot/20200828221148.png)

然后在插件编辑器里输入插件脚本，插件脚本使用javascript来实现，并且具有固定的初始化模板：

![20200828221633.png](docs/screenshot/20200828221633.png)

接下来，我们只需要在```// Start here```后输入脚本即可，此处为了将后缀文字```"创建"```删除，我们可以使用以下脚本来实现：

```javascript
(function(){
  // Plugin Name: remove_suffix
  // Start here...
  return TRIM(REPLACE(ENV("currentFieldValue"), "创建", ""));
})()
```
这里进行了三个操作： 
- ENV("currentFieldValue")&nbsp;&nbsp;-->&nbsp;&nbsp;获取当前字段的值
 - REPLACE()&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-->&nbsp;&nbsp;当前字段的值中的“创建”替换为空字符串
 - TRIM()&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-->&nbsp;&nbsp;将替换后的结果去除头尾部空格

然后将结果返回作为字段create_date最终结果。
ENV, REPLACE, TRIM为digger内置函数，可以直接在插件中调用，更多内置函数参见“内置函数”一节。

> 编写完插件要记得保存！

插件编写完成之后需要将其应用于field上，将create_date的field属性```plugin```设为```remove_suffix@s4```，表示插件将在s4插槽作用在该字段上。如下：

![20200828223501.png](docs/screenshot/20200828223501.png)

保存配置，回到项目列表页面，点击“+”按钮再次一个开启新任务，这时可以看到结果里的日期已经正确处理。

![20200828223407.png](docs/screenshot/20200828223407.png)

## 内置函数
digger的内置函数包含常见的字符处理和其他常用函数，如果需要新的函数，请提交PR

内置函数列表：
- ```LEN(str)```
返回字符串长度，返回值类型int
- ```STARTS_WITH(source, target)```
判断字符串```source```是否有前缀```target```，返回值类型```boolean```
- ```END_WITH(source, target)```
判断字符串```source```是否有后缀```target```，返回值类型```boolean```
- ```SUBSTR(source, start, end)```
获取字符串```source```的子串，位于```start```, ```end```之间，返回值类型```string```
- ```CONTAINS(source, target)```
判断字符串```source```是否包含字符串```target```，返回值类型```boolean```
- ```REPLACE(source, old, new)```
将字符串```source```中的字符串```old```替换为```new```并返回替换后的字符串
- ```REGEXP_GROUP_FIND(regexp, source, target)```
正则表达式匹配组替换，例如```REGEXP_GROUP_FIND(".\*([0-9]+).\*", "abc123mn", "$1")```将得到返回结果```123```
- ```MD5(source)```
计算字符串```source```的md5值
- ```TRIM(source)```
去除字符串```source```首尾空格
- ```ENV(key)```
获取环境值，目前可用的key有：```currentFieldName```，```currentFieldValue```
- ```RESPONSE_DATA()```
获取http请求响应结果
- ```SET_RESPONSE_DATA(data)```
如果是自定义AJAX请求，可以通过该函数将响应结果设置到上下文中供go程序使用
- ```QUEUE()```
获取当前任务实体类信息，Queue的 go struct 定义如下：
```golang
type Queue struct {
	Id         int64  `json:"id" gorm:"column:id;primary_key"`
	TaskId     int    `json:"task_id" gorm:"column:task_id"`
	StageName  string `json:"stage_name" gorm:"column:stage_name"`
	Url        string `json:"url" gorm:"column:url"`
	MiddleData string `json:"middle_data" gorm:"column:middle_data"`
	Expire     int64  `json:"expire" gorm:"column:expire"`
}
```
例如，可以通过```QUEUE().Url```获取当前任务的Url
- ```ABS(url)```
将相对URL转化为绝对URL
- ```ADD_QUEUE(url)```
添加任务，适用于需要从当前任务派生出子任务的场景，如根据尾页码计算所有分页的URL，并手动添加至队列。
注意，派生出的任务stage和当前任务的stage相同。
- ```AJAX(method, url, headers, querys, body)```
发送AJAX请求，例如：
```shell
AJAX("POST", "https://demo.com/some/page", {
	"X-TOKEN": "xxx"
}, {
	"page": "1",
}, {
	"field1": "value1"
})

```

相当于
```shell
curl -X POST -H "X-TOKEN:xxx" -d "{\"field1\":\"value1\"}" "https://demo.com/some/page?page=1"
```

# 常见问题
1. ```x509: certificate signed by unknown authority```
可以尝试在爬虫配置的settings下添加配置项：
```yaml
settings:
	SKIP_TLS_VERIFY: "true"
```


# 加入讨论组
如果您在使用过程中遇到任何问题，欢迎加作者微信，由作者邀请加入微信讨论群。
> 加好友请备注“digger”

![20200828131723.png](docs/author_qr.jpg)

# 支持
如果您觉得该软件解决了您的问题，请考虑打赏以支持作者

![20200828131723.png](docs//tip.jpg)

