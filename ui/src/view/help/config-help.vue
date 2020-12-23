<template>
    <div style="height: 100%;">
      <VueMarkdown :source="source"/>
      <br>
    </div>
</template>

<script>

  import VueMarkdown from 'vue-markdown'

  export default {
    name: 'config_help',
    components: {
      VueMarkdown
    },
    data () {
      return {
        source:`
### 一、流程示意图
略
### 二、概念
- stage
顾名思义代表一个爬虫的阶段，每个阶段都会对应一个输入（Input）

在程序中这个输入即为Queue，

- field
要从当前stage的请求结果中提取的字段

### 三、配置解读

\`\`\`yaml
start_url:
- https://music.163.com/discover/playlist # 爬虫起始URL
start_stage: list # 起始stage
stages:
- name: list # stage的name
  is_list: true # 是否是列表类型
  is_unique: false # 页面URL是否是唯一（暂时可以忽略）
  list_xpath: # 列表的xpath选择器表达式
  list_css: ul#m-pl-container>li # 列表的css选择器表达式
  page_xpath: # 分页按钮的xpath选择器表达式（如果有分页）
  page_css: a.znxt # 分页按钮的css选择器表达式（如果有分页）
  page_attr: href # 分页按钮的url标签属性（通常是href）
  plugin: "" # 插件，请参考插件一节
  fields:
  - name: cover # 字段name
    is_array: false # 指示该字段是否是一个数组，如标签，组图等一个字段需要匹配多个值的场景
    is_html: false # 指示该字段是否是匹配标签下的原始html内容，对于字段需要提取原始html内容的场景非常有用
    xpath: "" # xpath选择器
    css: div.u-cover>a # 字段的css选择器表达式
    attr: href # 字段的标签属性
    plugin: "" # 插件，请参考插件一节
    remark: 歌单地址 # 字段备注
    next_stage: detail # 下一阶段，将该字段的结果作为下一阶段的输入，例如：列表页提取的详情页URL，下阶段可以是详情页的stage
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
  EXPORT_PAGE_SIZE: "1000" # 导出时每次从数据库查询的分页大小，影响导出速度和内存占用
  FOLLOW_ROBOTS_TXT: "false" # 是否遵循robots指令
node_affinity: # 节点亲和标签列表，标签能够匹配相应的worker
- "key=value"
\`\`\`

### 四、开发插件
#### 关于
---
插件开发能解决很多爬虫过程中需要定制化的功能和业务逻辑处理问题。

比如URL替换，发送POST请求，正则表达式，处理复杂分页等等。

插件的引入能够使用很轻量简便的方式，少量的javascript代码就能够影响爬虫程序的运行。

---

插件在设计中引入了5个插槽点(slot)：
- **\`\`\`s1\`\`\`**
请求之前拦截URL，返回处理之后的URL
- **\`\`\`sr\`\`\`**
请求中，能够自定义请求，默认是GET方法，可以在此自定义POST请求
- **\`\`\`s2\`\`\`**
请求之后，引擎处理之前，输入为http请求结果，输出为处理后的结果
- **\`\`\`s3\`\`\`**
TODO 处理中，处理使用的引擎，默认为goquery，否则为自定义的插件处理引擎
- **\`\`\`s4\`\`\`**
引擎为默认的goquery时，解析得到字段值(包括分页)之后，可以用来修正数据，如去空格，剪切等

stage可以引入多个插件，field只能引入一个插件，引入格式为：
\`\`\`插件1@s1,插件2@s2...\`\`\`


#### 插件的开发
---

插件开发有一个预定义的样板代码：
\`\`\`javascript
(function(){
  // start here...
})()
\`\`\`

开发插件的过程中可以直接使用一些内置的方法：

内置函数列表：
- \`\`\`LEN(str)\`\`\`
返回字符串长度，返回值类型int
- \`\`\`STARTS_WITH(source, target)\`\`\`
判断字符串\`\`\`source\`\`\`是否有前缀\`\`\`target\`\`\`，返回值类型\`\`\`boolean\`\`\`
- \`\`\`ENDS_WITH(source, target)\`\`\`
判断字符串\`\`\`source\`\`\`是否有后缀\`\`\`target\`\`\`，返回值类型\`\`\`boolean\`\`\`
- \`\`\`SUBSTR(source, start, end)\`\`\`
获取字符串\`\`\`source\`\`\`的子串，位于\`\`\`start\`\`\`, \`\`\`end\`\`\`之间，返回值类型\`\`\`string\`\`\`
- \`\`\`CONTAINS(source, target)\`\`\`
判断字符串\`\`\`source\`\`\`是否包含字符串\`\`\`target\`\`\`，返回值类型\`\`\`boolean\`\`\`
- \`\`\`REPLACE(source, old, new)\`\`\`
将字符串\`\`\`source\`\`\`中的字符串\`\`\`old\`\`\`替换为\`\`\`new\`\`\`并返回替换后的字符串
- \`\`\`REGEXP_GROUP_FIND(regexp, source, target)\`\`\`
正则表达式匹配组替换，例如\`\`\`REGEXP_GROUP_FIND(".*([0-9]+).*", "abc123mn", "$1")\`\`\`将得到返回结果\`\`\`123\`\`\`
- \`\`\`MD5(source)\`\`\`
计算字符串\`\`\`source\`\`\`的md5值
- \`\`\`TRIM(source)\`\`\`
去除字符串\`\`\`source\`\`\`首尾空格
- \`\`\`ENV(key)\`\`\`
获取环境值，目前可用的key有：\`\`\`currentFieldName\`\`\`，\`\`\`currentFieldValue\`\`\`

- \`\`\`MIDDLE_DATA()\`\`\`
  获取中间值，可以获取父级stage里的field值和本级stage里其他field的值。例如：\`\`\`MIDDLE_DATA().field_name1\`\`\`

- \`\`\`FROM_JSON(string)\`\`\`
  将字符串解析为js对象

- \`\`\`TO_JSON()\`\`\`
  将js对象转成json字符串
- \`\`\`RESPONSE_DATA()\`\`\`
获取http请求响应结果
- \`\`\`SET_RESPONSE_DATA(data)\`\`\`
如果是自定义AJAX请求，可以通过该函数将响应结果设置到上下文中供go程序使用
- \`\`\`QUEUE()\`\`\`
获取当前任务实体类信息，Queue的 go struct 定义如下：
\`\`\`golang
type Queue struct {
	Id         int64  \`json:"id" gorm:"column:id;primary_key"\`
	TaskId     int    \`json:"task_id" gorm:"column:task_id"\`
	StageName  string \`json:"stage_name" gorm:"column:stage_name"\`
	Url        string \`json:"url" gorm:"column:url"\`
	MiddleData string \`json:"middle_data" gorm:"column:middle_data"\`
	Expire     int64  \`json:"expire" gorm:"column:expire"\`
}
\`\`\`
例如，可以通过\`\`\`QUEUE().Url\`\`\`获取当前任务的Url
- \`\`\`ABS(url)\`\`\`
将相对URL转化为绝对URL
- \`\`\`ADD_QUEUE(obj)\`\`\`
添加任务，适用于需要从当前任务派生出子任务的场景，如根据尾页码计算所有分页的URL，并手动添加至队列。对象\`\`\`obj\`\`\`格式：\`\`\`{stage: "", url: "", middle_data: {}}\`\`\`

- \`\`\`AJAX(method, url, headers, querys, body)\`\`\`
发送AJAX请求，例如：
\`\`\`shell
var result = AJAX("POST",
             "https://demo.com/some/page",
             {
               "X-TOKEN": "xxx"
             },
             {
               "page": "1",
             },
             "name=zhangsan&sex=1")
\`\`\`
相当于
\`\`\`shell
curl -X POST -H 'X-TOKEN:xxx' \\
     -d '{\"field1\":\"value1\"}' \\
     'https://demo.com/some/page?page=1'
\`\`\`

返回值result：
\`\`\`json
{
 status: 200, # 请求http响应码
  data: "", # 请求响应
}
\`\`\`


- \`\`\`LOG("text1", "text2", ...)\`\`\`
打印日志，能够在调试阶段展示在界面上，帮助排错
例如：
\`\`\`LOG("ABC", "123");\`\`\`
输出：ABC123


- \`\`\`LOGF("%s:%s", "text1", "text2", ...)\`\`\`
格式化打印日志，能够在调试阶段展示在界面上，帮助排错
例如：
\`\`\`LOGF("%s:%s", "localhost", "8080");\`\`\`
输出：localhost:8080




`
      }
    },
    methods: {
    },

    watch: {
    },

    mounted() {
    }
  }
</script>
<style>
</style>
