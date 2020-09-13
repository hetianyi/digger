<template>
  <div>
    <Row style="margin: 15px 0">

      <Col span="24">
        <span style="line-height: 32px;">项目</span>
        <i-select v-model="currentTaskPageSelectedProjectId" clearable filterable
                  style="width: 200px; padding: 0 10px;"
                  @on-change="statusChange"
        >
          <Option :value="p.id" :key="p.id" v-for="p in projectList">{{p.name}}</Option>
        </i-select>
        <span style="line-height: 32px;margin-left: 10px;">状态</span>
        <Select v-model="selectedStatus" style="width: 90px; padding: 0 10px;" clearable
          @on-change="statusChange">
          <Option :value="p.status" :key="p.status" v-for="p in standardTaskStatus">{{p.name}}</Option>
        </Select>
        <Button type="primary" style="margin-left: 10px;" @click="listTasks">搜索</Button>
      </Col>

    </Row>

    <Row>
      <Table border
             :columns="columns"
             :data="taskList"
             :border="false"
             :loading="loading"
             :no-data-text="noDataText">
        <template slot-scope="{ row }" slot="name">
          <strong>{{ row.name }}</strong>
        </template>
        <template slot-scope="{ row, index }" slot="action">
          <Button type="info"
                  shape="circle"
                  ghost
                  icon="md-eye"
                  style="margin-right: 5px"
                  @click="showTaskDetail(row)"
                  title="查看结果"></Button>
          <Button type="info"
                  ghost
                  shape="circle"
                  icon="md-document"
                  style="margin-right: 5px"
                  @click="checkLogs(row)"
                  title="查看日志"></Button>

          <Button type="warning"
                  shape="circle"
                  :disabled="row.status != 1"
                  ghost
                  icon="ios-pause"
                  style="margin-right: 5px"
                  @click="updateTaskStatus(row.id, 0)"
                  title="暂停"></Button>
          <Button type="primary"
                  shape="circle"
                  :disabled="row.status != 0"
                  ghost
                  icon="md-play"
                  style="margin-right: 5px"
                  @click="updateTaskStatus(row.id, 1)"
                  title="继续"></Button>
          <Button type="error"
                  shape="circle"
                  :disabled="row.status != 0 && row.status != 1"
                  ghost
                  icon="ios-square"
                  style="margin-right: 5px"
                  @click="updateTaskStatus(row.id, 2)"
                  title="停止"></Button>
          <Button type="error"
                  shape="circle"
                  :disabled="row.status !== 2 && row.status !== 3"
                  ghost
                  icon="md-trash"
                  style="margin-right: 5px"
                  @click="deleteTask(row.id)"
                  title="删除"></Button>
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
            :page-size-opts="[10,20,50]"
            @on-change="changePage"
            @on-page-size-change="changePageSize"/>
    </Row>

    <Modal v-model="showTaskDetailModal" fullscreen :title="projectDetailModalTitle" :footer-hide="true">
      <TaskDetail :taskId="taskDetailId"></TaskDetail>
      <div slot="footer">
      </div>
    </Modal>

    <Modal v-model="showTaskLogsModal"
           fullscreen
           :title="projectDetailModalTitle"
           @on-visible-change="closeLogModal"
           :footer-hide="true">
      <TaskLogs :taskId="checkLogsTaskId"></TaskLogs>
      <div slot="footer">
      </div>
    </Modal>


  </div>
</template>

<script>
  import {
    taskList,
    pauseTask,
    continueTask,
    stopTask,
    deleteTask,
  } from "@/api/task"
  import {
    dateFormat,
  } from "@/libs/util"
  import {listProjects} from "@/api/project"
  import TaskDetail from "./task-detail"
  import TaskLogs from "./task-logs"

  export default {
    name: 'task_list',
    components: {
      TaskDetail,
      TaskLogs,
    },
    data () {
      return {
        page: 1,
        pageSize: 10,
        total: 0,
        loading: false,
        noDataText: '空空如也',
        showTaskDetailModal: false,
        taskDetailId: 0,
        searchParams: {
          projectId: 0,
          taskStatus: 0,
        },

        showTaskLogsModal: false,
        projectDetailModalTitle: '日志',
        checkLogsTaskId: 0,

        standardTaskStatus: [
          {status: -1, name: '无', color: 'default'},
          {status: 1, name: '进行中', color: 'primary'},
          {status: 0, name: '暂停', color: 'warning'},
          {status: 2, name: '已停止', color: 'error'},
          {status: 3, name: '已完成', color: 'success'},
        ],

        projectList: [
          {id: 1, name: '树莓派实验室'}
        ],
        taskList: [],

        columns: [
          {
            title: 'ID',
            key: 'id',
            width: 80,
            align: 'center',
          },
          {
            title: '项目',
            key: 'display_name',
            render: (h, params) => {
              let pName = ''
              this.projectList.forEach(v => {
                if (v.id == params.row.project_id) {
                  pName = v.name
                }
              })
              return h('div', [
                h('span', pName)
              ]);
            }
          },
          {
            title: '结果',
            width: 120,
            align: 'center',
            render: (h, params) => {
              return h('div', [
                h('strong', params.row.result_count)
              ]);
            }
          },
          {
            title: '错误',
            width: 100,
            align: 'center',
            render: (h, params) => {
              return h('strong', params.row.error_request);
            }
          },
          {
            title: '状态',
            width: 150,
            align: 'center',
            render: (h, params) => {
              let statusStr = ''
              let color = ''
              this.standardTaskStatus.forEach(v => {
                if (v.status == params.row.status) {
                  statusStr = v.name
                  color = v.color
                }
              })
              return h('div', [
                h('Tag', {
                  props: {
                    type: 'dot',
                    color: color,
                  }
                }, statusStr),
              ]);
            }
          },
          {
            title: '启动时间',
            align: 'center',
            width: 150,
            render: (h, params) => {
              return h('div', dateFormat(params.row.create_time, 'yyyy-MM-dd hh:mm:ss'));
            }
          },
          {
            title: '操作',
            slot: 'action',
            width: 280,
            align: 'center'
          }
        ],

        split: 0.5,
      }
    },

    computed: {
      currentTaskPageSelectedProjectId: {
        get() {
          return this.$store.state.currentTaskPageSelectedProjectId
        },
        set(value) {
          this.$store.commit('selectProjectChange', value)
        }
      },
      selectedStatus: {
        get() {
          return this.$store.state.currentTaskPageSelectedStatus
        },
        set(value) {
          this.$store.commit('selectStatusChange', value)
        }
      }
    },

    methods: {

      async listProjects() {
        let that = this

        const { data: data } = await listProjects(1, 2147483647, 1)
        if (data && data.code == 0) {
          this.total = data.data.total
          this.projectList = data.data.data == null ? [] : data.data.data
          if (this.projectList.length == 0) {
          }
          that.loading = false
        } else {
          this.$Message.error('加载失败：' + data.msg)
        }
      },

      async listTasks() {
        if (this.projectId == 0) {
          return
        }
        this.loading = true
        this.noDataText = '加载中...'
        const { data: data } = await taskList(this.currentTaskPageSelectedProjectId, this.page, this.pageSize, this.selectedStatus)
        console.log(data)
        if (data && data.code == 0) {
          this.total = data.data.total
          this.taskList = null == data.data.data ? [] : data.data.data
          if (this.taskList.length == 0) {
            this.noDataText = '空空如也'
          }
        } else {
          this.$Message.error('加载失败：' + data.msg)
          this.noDataText = '加载失败'
        }
        this.loading = false
      },

      async updateTaskStatus(taskId, status) {
        this.$Spin.show({
          render: (h) => {
            return h('div', [
              h('Icon', {
                'class': 'demo-spin-icon-load',
                props: {
                  type: 'ios-loading',
                  size: 18
                }
              }),
              h('div', '请稍等...')
            ])
          }
        });

        if (status === 0) {
          const { data: data } = await pauseTask(taskId)
          console.log(data)

          if (data && data.code == 0) {
            for (let i = 0; i < this.taskList.length; i++) {
              if (this.taskList[i].id === taskId) {
                this.taskList[i].status = 0
              }
            }
          } else {
            this.$Message.error('操作失败：' + data.msg)
          }
        } else if (status === 1) {
          const { data: data } = await continueTask(taskId)
          console.log(data)
          if (data && data.code == 0) {
            for (let i = 0; i < this.taskList.length; i++) {
              if (this.taskList[i].id === taskId) {
                this.taskList[i].status = 1
              }
            }
          } else {
            this.$Message.error('操作失败：' + data.msg)
          }
        } else if (status === 2) {
          const { data: data } = await stopTask(taskId)
          console.log(data)
          if (data && data.code == 0) {
            for (let i = 0; i < this.taskList.length; i++) {
              if (this.taskList[i].id === taskId) {
                this.taskList[i].status = 2
              }
            }
          } else {
            this.$Message.error('操作失败：' + data.msg)
          }
        }
        this.$Spin.hide();
      },

      showTaskDetail(row) {
        this.showTaskDetailModal = true
        this.taskDetailId = row.id
        let pName = ''
        this.projectList.forEach(v => {
          if (v.id == row.project_id) {
            pName = v.name
          }
        })
        this.projectDetailModalTitle = '任务结果（项目：' + pName + '，任务ID：' + row.id + '）'
      },

      changePageSize(newSize) {
        this.page = 1
        this.pageSize = newSize
        this.listTasks()
      },

      changePage(page) {
        this.page = page
        this.listTasks()
      },

      statusChange(status) {
        this.listTasks()
      },
      projectChange(status) {
        this.listTasks()
      },

      checkLogs(row) {
        this.showTaskLogsModal = true
        this.checkLogsTaskId = row.id
      },

      closeLogModal() {
        if (!this.showTaskLogsModal) {
          console.log('关闭日志窗口')
          this.checkLogsTaskId = 0
        }
      },

      deleteTask(taskId) {
        let that = this
        this.$Modal.confirm({
          title: '确定删除此任务吗？',
          loading: true,
          onOk: async () => {
            const { data: data } = await deleteTask(taskId)
            if (data && data.code == 0) {
              that.$Message.success('删除成功');
              that.listTasks()
            } else {
              this.$Message.error('删除失败：' + data.msg)
            }
            that.$Modal.remove();
          }
        });
      },

    },

    mounted() {
      this.listProjects()
      this.listTasks()
    }
  }
</script>
<style>

</style>
