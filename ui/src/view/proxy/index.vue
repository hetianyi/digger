<template>
  <div>
    <Row style="margin: 15px 0">

      <Col span="24" style="text-align: right">
        <Button type="primary" shape="circle" icon="md-add" @click="prepareInsert">新增代理</Button>
      </Col>

      <h3>代理设置</h3>

    </Row>

    <Row>
      <Table border
             :columns="columns"
             :data="proxyList"
             :border="false"
             :loading="loading"
             :no-data-text="noDataText">
        <template slot-scope="{ row }" slot="name">
          <strong>{{ row.name }}</strong>
        </template>
        <template slot-scope="{ row, index }" slot="action">
          <Button type="primary"
                  shape="circle"
                  ghost
                  icon="md-create"
                  style="margin-right: 5px"
                  @click="prepareEdit(row)"
                  title="编辑"></Button>

          <Button type="error"
                  shape="circle"
                  ghost
                  icon="md-trash"
                  style="margin-right: 5px"
                  @click="doDeleteProxies([row.id])"
                  title="删除"></Button>
        </template>
      </Table>
    </Row>
    <br>
    <Row style="margin: 10px 0; text-align: right">
      <Page :total="total"
            :page-size="pageSize"
            :current.sync="page"
            show-total/>
    </Row>


    <!-- 代理编辑对话框 -->
    <Modal v-model="showProxyEditModal"
           :mask-closable="false"
           title="代理编辑服务器">
      <Form :model="tempEditProxy"
            label-position="right"
            :label-width="80">
        <FormItem label="代理地址">
          <Input v-model="tempEditProxy.address" placeholder="host:port"></Input>
        </FormItem>
        <FormItem label="备注">
          <Input v-model="tempEditProxy.remark"></Input>
        </FormItem>
        <FormItem label="启用脚本">
          <i-switch v-model="tempEditProxy.enable_script">
            <span slot="open">开</span>
            <span slot="close">关</span>
          </i-switch>
          <span style="margin-left: 20px">(面向代理商，自定义生成header)</span>
        </FormItem>
      </Form>

      <textarea ref="proxy_script_textarea"/>

      <div slot="footer">
        <Button type="success" :loading="saveProxyModalLoading" @click="doSaveProxy(false)">保存</Button>
      </div>
    </Modal>

    <!-- 新增代理对话框 -->
    <Modal v-model="showProxyInsertModal"
           :mask-closable="false"
           title="新增代理服务器">
      <Form :model="tempInsertProxy"
            label-position="right"
            :label-width="80">
        <FormItem label="代理地址">
          <Input v-model="tempInsertProxy.address" placeholder="host:port"></Input>
        </FormItem>
        <FormItem label="备注">
          <Input v-model="tempInsertProxy.remark"></Input>
        </FormItem>
        <FormItem label="启用脚本">
          <i-switch v-model="tempInsertProxy.enable_script">
            <span slot="open">开</span>
            <span slot="close">关</span>
          </i-switch>
          <span style="margin-left: 20px">(面向代理商，自定义生成header)</span>
        </FormItem>
      </Form>

      <textarea ref="insert_proxy_script_textarea"/>

      <div slot="footer">
        <Button type="success" :loading="saveProxyModalLoading" @click="doSaveProxy(true)">保存</Button>
      </div>
    </Modal>

  </div>
</template>

<script>
  import {
    proxyList,
    deleteProxies,
    saveProxy,
  } from "@/api/proxy"

  import CodeMirror from 'codemirror'
  import 'codemirror/addon/lint/lint.css'
  import 'codemirror/lib/codemirror.css'

  import 'codemirror/theme/darcula.css'
  import 'codemirror/theme/dracula.css'
  import 'codemirror/theme/material.css'
  import 'codemirror/theme/material-darker.css'
  import 'codemirror/theme/material-ocean.css'
  import 'codemirror/theme/material-palenight.css'
  import 'codemirror/theme/mbo.css'
  import 'codemirror/theme/midnight.css'

  import 'codemirror/mode/yaml/yaml'
  import 'codemirror/mode/javascript/javascript'
  import 'codemirror/addon/lint/lint'
  import 'codemirror/addon/lint/yaml-lint'
  import 'codemirror/keymap/sublime.js'
  import 'codemirror/addon/hint/javascript-hint'

  export default {
    name: 'proxy',
    components: {},
    data() {
      return {
        page: 1,
        pageSize: 10,
        total: 0,
        loading: false,
        noDataText: '空空如也',
        showProxyEditModal: false,
        showProxyInsertModal: false,
        saveProxyModalLoading: false,
        tempEditProxy: {},
        tempInsertProxy: {
          address: '',
          remark: '',
          enable_script: false,
        },

        proxyList: [],

        scriptFileEditor: false,
        scriptInsertFileEditor: false,

        columns: [
          {
            title: 'ID',
            key: 'id',
            width: 80,
            align: 'center',
          },
          {
            title: '地址',
            key: 'address',
            width: 250,
            align: 'center',
          },
          {
            title: '备注',
            key: 'remark',
            align: 'left',
          },
          {
            title: '操作',
            slot: 'action',
            width: 280,
            align: 'center'
          }
        ],
      }
    },

    methods: {

      async listProxies() {
        let that = this
        const {data: data} = await proxyList('', this.page, this.pageSize)
        if (data && data.code == 0) {
          this.total = data.data.total
          this.proxyList = data.data.data == null ? [] : data.data.data
          that.loading = false
        } else {
          this.$Message.error('加载失败：' + data.msg)
        }
      },


      prepareEdit(row) {
        if (row.proxy_gen_script == "") {
          this.scriptFileEditor.setValue(`(function(){
  // generate header here...
  // do something...
  var genToken = "";
  // return header
  return {
    "Proxy-Authorization": genToken
  };
})()`)
        } else {
          this.scriptFileEditor.setValue(row.proxy_gen_script)
        }
        this.tempEditProxy = {
          id: row.id,
          address: row.address,
          remark: row.remark,
          enable_script: row.enable_script,
        }
        this.scriptFileEditor.refresh()
        this.showProxyEditModal = true
      },


      prepareInsert() {
        this.scriptInsertFileEditor.setValue(`(function(){
  // generate header here...
  // do something...
  var genToken = "";
  // return header
  return {
    "Proxy-Authorization": genToken
  };
})()`)
        this.tempInsertProxy = {
          address: '',
          remark: '',
          enable_script: false,
        }
        this.scriptInsertFileEditor.refresh()
        this.showProxyInsertModal = true
      },

      async doSaveProxy(isNew) {
        if (!isNew) {
          this.saveProxyModalLoading = true
          this.tempEditProxy.proxy_gen_script = this.scriptFileEditor.getValue()
          const {data: data} = await saveProxy(this.tempEditProxy)
          this.saveProxyModalLoading = false
          if (data && data.code == 0) {
            this.$Message.success('保存成功')
            this.showProxyEditModal = false
            await this.listProxies()
          } else {
            this.$Message.error('保存失败：' + data.msg)
          }
        } else {
          this.saveProxyModalLoading = true
          this.tempInsertProxy.proxy_gen_script = this.scriptInsertFileEditor.getValue()
          const {data: data} = await saveProxy(this.tempInsertProxy)
          this.saveProxyModalLoading = false
          if (data && data.code == 0) {
            this.$Message.success('保存成功')
            this.showProxyInsertModal = false
            await this.listProxies()
          } else {
            this.$Message.error('保存失败：' + data.msg)
          }
        }
      },

      async doDeleteProxies(proxyIds) {
        let that = this
        this.$Modal.confirm({
          title: '确定删除吗？',
          loading: true,
          onOk: () => {
            deleteProxies(proxyIds).then(res=>{
              console.log(res)
              if (res.status == 200) {
                this.$Message.success('删除成功');
                that.listProxies()
              } else {
                this.$Message.error('删除失败');
              }
            }).catch(err=> {
              this.$Message.error('删除失败');
            })
            this.$Modal.remove();
          }
        });
      },

      initScriptEditor() {
        // 初始化配置文件编辑器
        this.scriptFileEditor = CodeMirror.fromTextArea(this.$refs.proxy_script_textarea, {
          lineNumbers: true, // 显示行号
          mode: 'text/javascript', // 语法model
          gutters: ['CodeMirror-lint-markers'],  // 语法检查器
          theme: 'darcula', // 编辑器主题
          tabSize: 2,
          autofocus: true,
          lint: false // 开启语法检查
        })
        this.scriptFileEditor.setValue('')
        this.scriptFileEditor.getWrapperElement().style.fontSize = '15px'
        this.scriptFileEditor.getWrapperElement().style.fontFamily = 'Consolas'
        this.scriptFileEditor.getWrapperElement().style.width = '100%'
        this.scriptFileEditor.getWrapperElement().style.height = '200px'
        this.scriptFileEditor.refresh()


        this.scriptInsertFileEditor = CodeMirror.fromTextArea(this.$refs.insert_proxy_script_textarea, {
          lineNumbers: true, // 显示行号
          mode: 'text/javascript', // 语法model
          gutters: ['CodeMirror-lint-markers'],  // 语法检查器
          theme: 'darcula', // 编辑器主题
          tabSize: 2,
          autofocus: true,
          lint: false // 开启语法检查
        })
        this.scriptInsertFileEditor.setValue('')
        this.scriptInsertFileEditor.getWrapperElement().style.fontSize = '15px'
        this.scriptInsertFileEditor.getWrapperElement().style.fontFamily = 'Consolas'
        this.scriptInsertFileEditor.getWrapperElement().style.width = '100%'
        this.scriptInsertFileEditor.getWrapperElement().style.height = '200px'
        this.scriptInsertFileEditor.refresh()
      },


    },

    mounted() {
      this.initScriptEditor()
      this.listProxies()
    }
  }
</script>
<style>

</style>
