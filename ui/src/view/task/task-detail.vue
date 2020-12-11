<template>
  <div>
    <Row style="margin: 15px 0">

      <Col span="24">

        <RadioGroup v-model="exportFormat">
          <Radio v-for="item in exportFormats" :label="item"></Radio>
        </RadioGroup>

        <Button icon="ios-download-outline" type="primary" @click="exportData">导出</Button>
        <Button icon="ios-download-outline" type="primary" @click="exportDML(1)" style="margin-left: 10px">导出DDL(Postgres)</Button>
        <Button icon="ios-download-outline" type="primary" @click="exportDML(2)" style="margin-left: 10px">导出DDL(MySQL)</Button>


        <Button type="success" style="float: right" icon="md-sync" @click="listResult">刷新</Button>
      </Col>

    </Row>

    <Row>
      <Table border :columns="columns" :data="resultList" :stripe="true" :loading="loading" :no-data-text="noDataText">
        <template slot-scope="{ row }" slot="name">
          <strong>{{ row.name }}</strong>
        </template>
      </Table>
    </Row>
    <br>
    <Row style="margin: 10px 0; text-align: right">
      <Page :total="total"
            :page-size="pageSize"
            :current.sync="page"
            show-sizer
            show-total
            :page-size-opts="[10,20,50,100]"
            @on-change="changePage"
            @on-page-size-change="changePageSize"/>
    </Row>

  </div>
</template>

<script>

  import {
    getProjectSnapshotConfig,
  } from "@/api/project"
  import {
    resultList,
  } from "@/api/result"

  export default {
    name: 'task_detail',
    data () {
      return {
        page: 1,
        pageSize: 20,
        total: 0,
        loading: false,
        noDataText: '空空如也',

        autoRefresh: true,
        autoRefreshTimer: true,

        exportFormats: ['SQL', 'JSON', 'CSV'],
        exportFormat: 'SQL',

        resultList: [],

        pageParams: {
          page: 1,
          pageSize: 20,
        },
        configSnapshot: null,
        columns: [],
      }
    },
    props: {
      taskId: {
        type: [Number, String],
        default: 0,
      },
    },
    methods: {

      async listResult() {
        if (this.taskId === 0) {
          return
        }
        this.loading = true
        this.noDataText = '加载中...'

        const { data: data } = await resultList(this.taskId, this.page, this.pageSize)
        console.log(data)
        if (data && data.code == 0) {
          this.total = data.data.total

          let resultArr = new Array()
          let srcDataArr = data.data.data == null ? [] : data.data.data
          if (this.total > 0) {
            for (let i = 0; i< srcDataArr.length; i++) {
              let item = JSON.parse(srcDataArr[i].result)
              item.id = srcDataArr[i].id
              resultArr.push(item)
            }
          }
          this.resultList = resultArr

          if (srcDataArr.length == 0) {
            this.loading = false
            this.noDataText = '空空如也'
            return
          }

          let obj = JSON.parse(data.data.data[0].result)
          this.columns = []
          let fields = []
          for (let n in obj) {
            console.log(n)
            fields.push(n)
          }

          fields = fields.sort()
          this.columns.push({
            title: 'ID',
            width: 100,
            key: 'id',
            align: 'center',
            render: (h, params) => {
              return h('div', [
                h('span', params.row.id)
              ]);
            }
          })
          for (let i = 0; i < fields.length; i++) {
            if (fields[i] == 'id') {
              continue
            }
            this.columns.push({
              title: fields[i],
              key: fields[i],
              ellipsis: true,
              tooltip: true,
              attrs: {
                ellipsis: true,
                tooltip: true,
              },
              render: (h, params) => {
                console.log(typeof params.row[fields[i]])
                let fd = params.row[fields[i]]
                if (fd) {
                  if ((typeof fd) == 'object' && fd.length > 0) {
                    let arr = new Array()
                    for (let k = 0; k < fd.length; k++) {
                      let m = h('Tag', fd[k])
                      arr.push(m)
                    }
                    return arr
                  } else {
                    let m = params.row[fields[i]].startsWith('http://') || params.row[fields[i]].startsWith('https://') ?
                      h('a',
                        {
                          attrs: {
                            href: params.row[fields[i]],
                            target: '_blank',
                            title: params.row[fields[i]],
                          },
                        },
                        params.row[fields[i]])
                      :
                      h('span', {
                        attrs: {
                          title: params.row[fields[i]],
                        },
                      }, params.row[fields[i]])
                    return h('div', [
                      m
                    ]);
                  }
                }
              }
            })
          }
        } else {
          this.resultList = []
          this.$Message.error('加载失败：' + data.msg)
          this.noDataText = '加载失败'
        }

        this.loading = false
      },


      async getProjectSnapshotConfig() {
        const { data: data } = await getProjectSnapshotConfig(this.taskId)
        console.log(data)
        if (data && data.code == 0) {
          this.configSnapshot = data.data
        } else {
          this.$Message.error('加载失败：' + data.msg)
        }
      },

/*
      changeAutoRefresh (status) {
        let that = this
        if (status) {
          this.autoRefreshTimer = setInterval(function () {
            that.listResult()
          }, 3000)
        } else {
          if (this.autoRefreshTimer != null) {
            this.autoRefreshTimer.cancel()
          }
          this.autoRefreshTimer = null
        }
      },*/


      changePageSize(newSize) {
        this.page = 1
        this.pageSize = newSize
        this.listResult()
      },

      changePage(page) {
        this.page = page
        this.listResult()
      },

      exportData(){
        window.open('/api/v1/results/export' +
          '?taskId=' + this.taskId +
          '&format=' + this.exportFormat +
          '&token=' + window.localStorage.getItem('token'))
      },

      exportDML(type){

        if (this.configSnapshot == null) {
          this.$Message.error('没有配置快照')
          return
        }

        console.log(this.configSnapshot)

        if (this.columns.length === 0) {
          this.$Message.error('没有数据')
          return
        }

        let d = ''
        if (type === 1) {
          d = `
DROP SEQUENCE IF EXISTS "seq_` + this.configSnapshot.name.toLowerCase() + `";
CREATE SEQUENCE "seq_` + this.configSnapshot.name.toLowerCase() + `"
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;
`

          d += `DROP TABLE IF EXISTS "t_` + this.configSnapshot.name.toLowerCase() + `";
CREATE TABLE "t_` + this.configSnapshot.name.toLowerCase() + `" (
  "id" int4 NOT NULL DEFAULT nextval('seq_` + this.configSnapshot.name.toLowerCase() + `'::regclass),
`
          this.columns.forEach(res=>{
            if (res.title == 'ID') {
              return
            }
            d+=`  "`+ res.title +`" text COLLATE "pg_catalog"."default",\n`
          })
          d += `  CONSTRAINT "t_` + this.configSnapshot.name.toLowerCase() + `_pkey" PRIMARY KEY ("id")
);`
        } else {
          d += `DROP TABLE IF EXISTS \`t_` + this.configSnapshot.name.toLowerCase() + `\`;
CREATE TABLE t_` + this.configSnapshot.name.toLowerCase() + ` (
  \`id\` bigint(20) NOT NULL AUTO_INCREMENT,
`
          this.columns.forEach(res=>{
            if (res.title == 'ID') {
              return
            }
            d+=`  \``+ res.title +`\` text COLLATE utf8mb4_general_ci,\n`
          })
          d += `  PRIMARY KEY (\`id\`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC;`
        }

        const downloadUrl = window.URL.createObjectURL(new Blob([d]))

        const link = document.createElement('a')

        link.href = downloadUrl

        link.setAttribute('download', 't_' + this.configSnapshot.name.toLowerCase() + '.sql') // any other extension

        document.body.appendChild(link)
        link.click()
        link.remove()
      },

    },

    watch: {
      taskId(val) {
        this.configSnapshot = null
        this.getProjectSnapshotConfig()
        this.listResult()
      },
    },

    mounted() {
      this.listResult()
    }
  }
</script>
<style>

</style>
