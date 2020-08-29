<template>
  <div>
    <Row style="margin: 10px 0; text-align: right">
      <div style="float: left">
        <Tag type="dot" color="primary">运行中</Tag>
        <Tag type="dot" color="warning">暂停</Tag>
        <Tag type="dot" color="error">已停止</Tag>
        <Tag type="dot" color="success">已完成</Tag>
      </div>

      <Button type="primary" shape="circle" icon="md-add" @click="showCreateProjectModal = true">新建项目</Button>
    </Row>

    <Row>
      <Table border :columns="columns" :data="projectList" :border="false" :loading="loading" :no-data-text="noDataText">
        <template slot-scope="{ row }" slot="name">
          <strong>{{ row.name }}</strong>
        </template>
        <template slot-scope="{ row, index }" slot="action">
          <Button type="success"shape="circle" ghost icon="md-add" style="margin-right: 5px" @click="startNewTask(row.id)" title="创建任务"></Button>
          <Button type="primary"shape="circle" ghost icon="md-create" style="margin-right: 5px" @click="showModal(row)" title="编辑"></Button>
          <Button type="primary"shape="circle" ghost icon="md-cloud-download" style="margin-right: 5px" @click="exportProjectConfig(row)" title="导出"></Button>
          <Button type="error"shape="circle" ghost icon="md-trash" @click="deleteProject(row)" title="删除"></Button>
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
            :page-size-opts="[1,2,5,10,20,50]"
            @on-change="changePage"
            @on-page-size-change="changePageSize"/>
    </Row>

    <Modal v-model="showProjectDetailModal"
           fullscreen
           :title="projectDetailModalTitle"
           @on-visible-change="closeProjectEditModal"
           :footer-hide="true">
      <ProjectEdit :projectId="projectDetailId"></ProjectEdit>
      <div slot="footer">
      </div>
    </Modal>


    <Modal v-model="showCreateProjectModal"
           @on-ok="doAddProject"
           title="新建项目">
      <Form :model="newProject"
            label-position="right"
            :label-width="100">
        <FormItem label="项目名">
          <Input v-model="newProject.name"></Input>
        </FormItem>
        <FormItem label="项目显示名称">
          <Input v-model="newProject.display_name"></Input>
        </FormItem>
        <Row>
          <Col>
            <FormItem label="标签">
              <p style="font-style: italic" v-if="newProjectTags.length == 0">无标签</p>
              <Tag color="default" closable @on-close="removeTag(t)"  v-for="t in newProjectTags" :value="t" :key="t">{{t}}</Tag>
            </FormItem>
          </Col>
          <Col>
            <FormItem label="">
              <Input v-model="tempTag" @on-enter="addTag" placeholder="键入标签按回车添加"></Input>
            </FormItem>
          </Col>
        </Row>
        <FormItem label="备注">
          <Input v-model="newProject.remark"></Input>
        </FormItem>
      </Form>


      <div slot="footer">
        <Button type="success" :loading="createProjectModalLoading" @click="doAddProject">确认</Button>
      </div>

    </Modal>

  </div>
</template>

<script>

  import ProjectEdit from './project-edit'
  import {
    listProjects,
    deleteProject,
    createProject,
    startNewTask,
  } from "@/api/project"

  export default {
    name: 'Project',
    components: {
      ProjectEdit,
    },
    data () {
      return {
        page: 1,
        pageSize: 10,
        total: 0,
        loading: false,
        noDataText: '空空如也',
        showProjectDetailModal: false,
        showTaskListModal: false,
        showCreateProjectModal: false,
        createProjectModalLoading: false,
        projectDetailModalTitle: '',
        projectDetailId: 0,
        projectList: [],
        newProject: {
          name: '',
          display_name: '',
          tags: '',
          remark: ''
        },
        newProjectTags: [],

        tempTag: '',
        columns: [
          {
            title: 'ID',
            key: 'id',
            width: 80,
            align: 'center',
          },
          {
            title: '显示名称',
            width: 200,
            key: 'display_name',
            ellipsis: true,
            tooltip: true,
          },
          {
            title: '名称',
            width: 150,
            key: 'name',
            ellipsis: true,
            tooltip: true,
          },
          {
            title: '任务',
            width: 320,
            align: 'center',
            render: (h, params) => {
              return h('div', [
                h('Tag', {
                  props: {
                    type: 'dot',
                    color: 'primary',
                  }
                }, params.row.extras.tasks.active_count),
                h('Tag', {
                  props: {
                    type: 'dot',
                    color: 'warning',
                  }
                }, params.row.extras.tasks.pause_count),
                h('Tag', {
                  props: {
                    type: 'dot',
                    color: 'error',
                  }
                }, params.row.extras.tasks.stop_count),
                h('Tag', {
                  props: {
                    type: 'dot',
                    color: 'success',
                  }
                }, params.row.extras.tasks.finish_count),
                h('Button', {
                  props: {
                    icon: 'ios-search',
                    shape: 'circle',
                  },
                  on: { // 操作事件
                    click: () => {
                      this.$store.commit('selectProjectChange', params.row.id)
                      this.$router.push({ name: 'task-list'})
                    }
                  }
                }),
              ]);
            }
          },
          {
            title: '备注',
            key: 'remark',
            ellipsis: true,
            tooltip: true,
          },
          {
            title: '操作',
            slot: 'action',
            width: 200,
            align: 'center'
          }
        ],
      }
    },
    methods: {

      showModal(row) {
        this.showProjectDetailModal = true
        this.projectDetailId = row.id
        this.projectDetailModalTitle = "配置项目：" + row.name
      },

      async listProjects() {
        let that = this
        that.loading = true
        that.noDataText = "加载中..."

        const { data: data } = await listProjects(that.page, this.pageSize, 1)
        if (data && data.code == 0) {
          this.total = data.data.total
          this.projectList = data.data.data == null ? [] : data.data.data
          if (this.projectList.length == 0) {
            that.noDataText = "空空如也"
          }
          that.loading = false
        } else {
          this.$Message.error('加载失败：' + data.msg)
          that.loading = false
          that.noDataText = "加载失败"
        }
      },

      deleteProject(row) {
        console.log(row)
        let that = this
        this.$Modal.confirm({
          title: '确定删除此项目吗？',
          loading: true,
          onOk: () => {
            deleteProject(row.id).then(res=>{
              console.log(res)
              if (res.status == 200) {
                this.$Message.success('删除成功');
                that.listProjects()
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

      addTag() {
        if (this.tempTag.trim() == '') {
          return
        }
        let exists = false
        this.newProjectTags.forEach(v => {
          if (v == this.tempTag.trim()) {
            exists = true
          }
        })
        if (exists) {
          this.tempTag = ''
          return
        }
        this.newProjectTags.push(this.tempTag.trim())
        this.tempTag = ''
      },

      removeTag(tag) {
        let newTags = new Array()
        this.newProjectTags.forEach(v => {
          if (v != tag) {
            newTags.push(v)
          }
        })
        this.newProjectTags = newTags
      },

      async doAddProject() {

        this.createProjectModalLoading = true

        let that = this
        this.newProject.tags = JSON.stringify(this.newProjectTags)
        const { data: data } = await createProject(this.newProject)
        this.createProjectModalLoading = false
        if (data && data.code == 0) {
          console.log(data.data)
          this.$Message.success('保存成功')
          this.showCreateProjectModal = false
          this.newProjectTags = []
          this.newProject = {
            name: '',
            display_name: '',
            tags: '',
            remark: ''
          }
          await that.listProjects()
        } else {
          this.$Message.error('保存失败：' + data.msg)
        }
      },

      startNewTask(projectId) {
        let that = this
        this.$Modal.confirm({
          title: '为此项目启动新的任务吗？',
          onOk: () => {
            that.doCreateTask(projectId)
          }
        });
      },

      async doCreateTask(projectId) {
        let that = this
        const { data: data } = await startNewTask(projectId)
        if (data && data.code == 0) {
          console.log(data.data)
          this.$Message.success('启动成功')
          //await that.listProjects()
          this.$store.commit('selectProjectChange', projectId)
          this.$router.push({ name: 'task-list'})
        } else {
          this.$Message.error('启动失败：' + data.msg)
        }
      },

      exportProjectConfig(row) {
        window.open('/api/v1/projects/' + row.id + '/export?token=' + window.localStorage.getItem('token'))
      },

      go2Task(projectId) {

      },

      closeProjectEditModal() {
        if (!this.showProjectDetailModal) {
          this.projectDetailId = 0
        }
      },

      changePageSize(newSize) {
        this.page = 1
        this.pageSize = newSize
        this.listProjects()
      },

      changePage(page) {
        this.page = page
        //this.$router.push({ name: 'project-list', query: { page: page }})
        this.listProjects()
      }
    },

    mounted() {

      /*let page = this.$route.query.page
      let pageSize = this.$route.query.pageSize
      if (page == null || page == undefined) {
        page = 1
      }
      if (pageSize == null || pageSize == undefined) {
        pageSize = 10
      }
      this.page = new Number(page)
      this.pageSize = new Number(pageSize)*/

      this.listProjects()
    }
  }
</script>
<style>

</style>
