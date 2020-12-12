<template>
  <div>
    <Row style="margin: 15px 0">

      <Col span="24" style="text-align: right">
        <Button type="primary" shape="circle" icon="md-add" @click="prepareInsert">新增推送源</Button>
      </Col>

      <h3>推送源设置</h3>

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
        <template slot-scope="{ row }" slot="enable_retry">
          <strong>{{ row.enable_retry ? '是' : '否' }}</strong>
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
                  @click="doDeletePushes([row.id])"
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


    <!-- 推送源编辑对话框 -->
    <Modal v-model="showPushEditModal"
           :mask-closable="false"
           title="推送源编辑">
      <Form :model="tempInsertPushSource"
            label-position="right"
            :label-width="120">
        <FormItem label="推送地址">
          <Input v-model="tempInsertPushSource.url" placeholder="推送源接口地址"></Input>
        </FormItem>
        <FormItem label="推送方法">
          <Select v-model="tempInsertPushSource.method" placeholder="POST, PUT">
            <Option v-for="item in ['POST', 'PUT', 'PATCH']" :value="item" :key="item">{{ item }}</Option>
          </Select>
        </FormItem>
        <FormItem label="批量推送大小">
          <Input v-model="tempInsertPushSource.push_size" placeholder="50"></Input>
        </FormItem>
        <FormItem label="推送间隔时间(ms)">
          <Input v-model="tempInsertPushSource.push_interval" placeholder="0"></Input>
        </FormItem>
        <FormItem label="是否重试">
          <i-switch
            :true-value="true"
            :false-value="false"
            v-model="tempInsertPushSource.enable_retry">
            <span slot="open">开</span>
            <span slot="close">关</span>
          </i-switch>
        </FormItem>
      </Form>

      <div slot="footer">
        <Button type="success" :loading="savePushModalLoading" @click="doSavePush(false)">保存</Button>
      </div>
    </Modal>

    <!-- 新增推送源 -->
    <Modal v-model="showPushInsertModal"
           :mask-closable="false"
           title="新增推送源">
      <Form :model="tempInsertPushSource"
            label-position="right"
            :label-width="120">
        <FormItem label="推送地址">
          <Input v-model="tempInsertPushSource.url" placeholder="推送源接口地址"></Input>
        </FormItem>
        <FormItem label="推送方法">
          <!--<Input v-model="tempInsertPushSource.method" placeholder="POST, PUT"></Input>-->
          <Select v-model="tempInsertPushSource.method" placeholder="POST, PUT">
            <Option v-for="item in ['POST', 'PUT', 'PATCH']" :value="item" :key="item">{{ item }}</Option>
          </Select>

        </FormItem>
        <FormItem label="批量推送大小">
          <Input v-model="tempInsertPushSource.push_size" placeholder="50"></Input>
        </FormItem>
        <FormItem label="推送间隔时间(ms)">
          <Input v-model="tempInsertPushSource.push_interval" placeholder="0"></Input>
        </FormItem>
        <FormItem label="是否重试">
          <i-switch
            :true-value="true"
            :false-value="false"
            v-model="tempInsertPushSource.enable_retry">
            <span slot="open">开</span>
            <span slot="close">关</span>
          </i-switch>
        </FormItem>
      </Form>

      <div slot="footer">
        <Button type="success" :loading="savePushModalLoading" @click="doSavePush(true)">保存</Button>
      </div>
    </Modal>

  </div>
</template>

<script>
  import {
    pushList,
    deletePushes,
    savePush,
  } from "@/api/push"

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
    name: 'push',
    components: {},
    data() {
      return {
        page: 1,
        pageSize: 10,
        total: 0,
        loading: false,
        noDataText: '空空如也',
        showPushEditModal: false,
        showPushInsertModal: false,
        savePushModalLoading: false,
        tempEditPushSource: {},
        tempInsertPushSource: {
          id: 0,
          url: 'http://',
          method: 'POST',
          push_size: 50,
          enable_retry: false,
          push_interval: 0,
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
            title: '推送地址',
            key: 'url',
            width: 250,
            align: 'center',
          },
          {
            title: '推送方法',
            key: 'method',
            width: 250,
            align: 'center',
          },
          {
            title: '批量大小',
            key: 'push_size',
            align: 'left',
          },
          {
            title: '重试',
            slot: 'enable_retry',
            align: 'left',
          },
          {
            title: '重试间隔(ms)',
            key: 'push_interval',
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

      async listPushSources() {
        let that = this
        const {data: data} = await pushList('', this.page, this.pageSize)
        if (data && data.code == 0) {
          this.total = data.data.total
          this.proxyList = data.data.data == null ? [] : data.data.data
          that.loading = false
        } else {
          this.$Message.error('加载失败：' + data.msg)
        }
      },


      prepareEdit(row) {
        this.tempInsertPushSource = {
          id: row.id,
          url: row.url,
          method: row.method,
          push_size: row.push_size,
          enable_retry: row.enable_retry,
          push_interval: row.push_interval,
        }
        this.showPushEditModal = true
      },


      prepareInsert() {
        this.tempInsertPushSource = {
          id: 0,
          url: 'http://',
          method: 'POST',
          push_size: 50,
          enable_retry: true,
          push_interval: 0,
        }
        this.showPushInsertModal = true
      },

      async doSavePush(isNew) {
        if (!isNew) {
          this.savePushModalLoading = true
          this.tempInsertPushSource.push_size = Number(this.tempInsertPushSource.push_size)
          this.tempInsertPushSource.push_interval = Number(this.tempInsertPushSource.push_interval)
          const {data: data} = await savePush(this.tempInsertPushSource)
          this.savePushModalLoading = false
          if (data && data.code == 0) {
            this.$Message.success('保存成功')
            this.showPushEditModal = false
            await this.listPushSources()
          } else {
            this.$Message.error('保存失败：' + data.msg)
          }
        } else {
          this.savePushModalLoading = true
          this.tempInsertPushSource.push_size = Number(this.tempInsertPushSource.push_size)
          this.tempInsertPushSource.push_interval = Number(this.tempInsertPushSource.push_interval)
          const {data: data} = await savePush(this.tempInsertPushSource)
          this.savePushModalLoading = false
          if (data && data.code == 0) {
            this.$Message.success('保存成功')
            this.showPushInsertModal = false
            await this.listPushSources()
          } else {
            this.$Message.error('保存失败：' + data.msg)
          }
        }
      },

      async doDeletePushes(pushIds) {
        let that = this
        this.$Modal.confirm({
          title: '确定删除吗？',
          loading: true,
          onOk: () => {
            deletePushes(pushIds).then(res=>{
              console.log(res)
              if (res.status == 200) {
                this.$Message.success('删除成功');
                that.listPushSources()
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

    },

    mounted() {
      this.listPushSources()
    }
  }
</script>
<style>

</style>
