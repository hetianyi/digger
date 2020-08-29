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
start_url: # 爬虫起始页面
start_stage: # 爬虫起始stage
stages: # 定义stage列表
- name: # stage名称
  is_list: # 结果是否是列表
  is_unique: #
  list_css: # # 列表的css样式(如果is_list=true)
  page_css: # 分页css样式
  page_attr: # 分页css属性
  plugin: # stage的插件列表，格式plugin_name@s1,plugin_name@s1,...
  fields: # 定义字段列表
  - name: # 字段名称
    is_array: # 字段结果是否是数组
    is_html: # 字段结果是否是html
    css: # 字段值css样式
    attr: # 字段css属性
    plugin: # 插件(字段只能有一个插件)
    remark: # 字段备注
    next_stage: # 下一stage名称
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
- **\`\`\`LEN(str)\`\`\`**
- **\`\`\`STARTS_WITH(source, target)\`\`\`**
- **\`\`\`END_WITH(source, target)\`\`\`**
- **\`\`\`SUBSTR(source, start, end)\`\`\`**
- **\`\`\`CONTAINS(source, target)\`\`\`**
- **\`\`\`REPLACE(source, old, new)\`\`\`**
- **\`\`\`REGEXP_GROUP_FIND(regexp, source, target)\`\`\`**
- **\`\`\`MD5(source)\`\`\`**
- **\`\`\`TRIM(source)\`\`\`**
- **\`\`\`ENV(key)\`\`\`**
- **\`\`\`RESPONSE_DATA()\`\`\`**
- **\`\`\`SET_RESPONSE_DATA(data)\`\`\`**
- **\`\`\`QUEUE()\`\`\`**
- **\`\`\`ABS(url)\`\`\`**
- **\`\`\`ADD_QUEUE(url)\`\`\`**
- **\`\`\`AJAX(method, url, headers, querys, body)\`\`\`**



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
