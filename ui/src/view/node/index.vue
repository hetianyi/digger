<template>
  <div>

    <Row>

      <Row style="margin: 10px 0">
        <Tag type="dot" color="primary">总计：{{total}}</Tag>
        <Tag type="dot" color="success">在线：{{online}}</Tag>
        <Tag type="dot" color="error">下线：{{offline}}</Tag>
      </Row>

      <Table border :columns="columns" :data="nodes" :border="true" :loading="loading" :no-data-text="noDataText">

      </Table>
    </Row>
    <br>

  </div>
</template>

<script>
  import {
    nodeList,
  } from "@/api/node"
  import {
    dateFormat,
  } from "@/libs/util"

  export default {
    name: 'node_list',
    data () {
      return {
        total: 0,
        online: 0,
        offline: 0,
        loading: false,
        noDataText: '空空如也',
        statusSelect: [
          {code: 1, name: '上线'},
          {code: 2, name: '下线'},
        ],
        nodes: [],
        columns: [
          {
            title: 'ID',
            key: 'instance_id',
            width: 80,
            align: 'center',
          },
          {
            title: '分配',
            width: 80,
            key: 'assign',
            align: 'center',
          },
          {
            title: '成功',
            width: 80,
            key: 'success',
            align: 'center',
          },
          {
            title: '错误',
            width: 80,
            key: 'error',
            align: 'center',
          },
          {
            title: '断线',
            width: 80,
            key: 'down',
            align: 'center',
          },
          {
            title: '地址',
            width: 180,
            key: 'address',
            align: 'center',
          },
          {
            title: '当前状态',
            width: 120,
            align: 'center',
            render: (h, params) => {
              return h('div', [
                h('Tag', {
                  props: {
                    type: 'dot',
                    color: params.row.status == '1' ? 'success' : 'error',
                  }
                }, params.row.status == '1' ? '在线' : '下线'),
              ]);
            }
          },
          {
            title: '注册时间',
            width: 150,
            align: 'center',
            render: (h, params) => {
              return h('div', dateFormat(params.row.register_at, 'yyyy-MM-dd hh:mm:ss'));
            }
          },
          {
            title: '标签',
            align: 'center',
            render: (h, params) => {
              if (params.row.labels) {
                let tags = new Array()
                for (let n in params.row.labels) {
                  console.log(n)
                  tags.push(h('Tag', {
                    props: {
                      type: 'dot',
                      color: 'primary',
                    }
                  }, n + ':' + params.row.labels[n]))
                }

                return h('div', tags);
              }
            }
          },
          /*{
            title: '操作',
            slot: 'action',
            width: 200,
            align: 'center'
          }*/
        ],
      }
    },

    computed: {

    },

    methods: {

      async listNodes() {
        this.loading = true
        this.noDataText = '加载中...'
        const { data: data } = await nodeList()
        console.log(data)
        if (data && data.code == 0) {
          let ret = null == data.data ? [] : data.data
          ret.sort((a, b) => {
            return a.instance_id - b.instance_id
          })
          this.nodes = ret
          if (this.nodes.length == 0) {
            this.noDataText = '空空如也'
          }
          this.total = this.nodes.length
          let _online = 0
          let _offline = 0
          this.online = this.nodes.forEach(v => {
            if (v.status == 1) {
              _online++
            } else {
              _offline++
            }
          })
          this.online = _online
          this.offline = _offline
        } else {
          this.$Message.error('加载失败：' + data.msg)
          this.noDataText = '加载失败'
        }
        this.loading = false
      },
    },

    mounted() {
      this.listNodes()
    }
  }
</script>
<style>

</style>
