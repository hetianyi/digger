<template>
    <div style="height: 100%;">
      <Layout style="height: 100%;">
        <Sider hide-trigger style="box-shadow: 4px 4px 15px #EBEBEB; background: none; " width="350">
          <Layout style="height: 100%;">
            <Content :style="{padding: '10px 10px'}">
              <Form :model="project" label-position="left" :label-width="100">
                <FormItem label="项目名">
                  <Input v-model="project.name"></Input>
                </FormItem>
                <FormItem label="项目显示名称">
                  <Input v-model="project.display_name"></Input>
                </FormItem>
                <FormItem label="定时启动">
                  <i-switch v-model="project.enable_cron">
                    <span slot="open">开</span>
                    <span slot="close">关</span>
                  </i-switch>
                  <a href="https://crontab.guru/"></a>
                </FormItem>
                <FormItem label="cron表达式" v-if="project.enable_cron">
                  <Input v-model="project.cron">
                    <Button slot="append" icon="md-help" title="帮助" to="https://crontab.guru/" target="_blank"></Button>
                  </Input>
                </FormItem>
                  <Row>
                    <Col>
                      <FormItem label="标签">
                        <p style="font-style: italic" v-if="tags.length == 0">无标签</p>
                        <Tag color="default" closable @on-close="removeTag(t)"  v-for="t in tags" :value="t" :key="t">{{t}}</Tag>
                      </FormItem>
                    </Col>
                    <Col>
                      <FormItem label="">
                        <Input v-model="tempTag" @on-enter="addTag" placeholder="键入标签回车添加"></Input>
                      </FormItem>
                    </Col>
                  </Row>
                <FormItem label="备注">
                  <Input v-model="project.remark"></Input>
                </FormItem>
              </Form>
            </Content>
            <Footer style="text-align: right; padding: 24px 10px;">
              <Button icon="md-checkbox-outline" type="info" :loading="saveProjectLoading" @click="updateProject">保存</Button>
            </Footer>
          </Layout>


        </Sider>
        <Content>

          <Tabs style="height: 100%" class="editor-tabs">
            <TabPane label="基本配置" icon="md-cog">

              <Layout style="height: 98%">
                <Content style="height: 100%; overflow-y: auto">
                  <textarea ref="config_textarea"/>
                </Content>
                <Footer style="text-align: right; padding: 24px 10px;">
                  <Button icon="md-help-circle" @click="showConfigHelp = true"></Button>
                  &nbsp;&nbsp;&nbsp;
                  <Upload :action="'/api/v1/projects/' + projectId + '/import'"
                          :multiple="false"
                          :paste="false"
                          :show-upload-list="false"
                          name="config"
                          type="select"
                          :format="['json']"
                          accept=".json"
                          :on-format-error="function() {
                            this.$Message.error('文件格式错误')
                          }"
                          :before-upload="showLoading"
                          :on-error="hideLoading"
                          :on-success="uploadConfigSuccess"
                          :headers="uploadHeaders"
                          style="display: inline-block">
                    <Button icon="md-cloud-upload" :loading="importLoading">导入</Button>
                  </Upload>
                  &nbsp;&nbsp;&nbsp;
                  <Button icon="md-cloud-download" @click="exportProjectConfig">导出</Button>
                  &nbsp;&nbsp;&nbsp;
                  <Button icon="md-checkbox-outline" :loading="saveProjectConfigLoading" type="info" @click="saveProjectConfig">保存</Button>
                  &nbsp;&nbsp;&nbsp;
                  <Button icon="md-bug" type="success" @click="parseConfig">调试</Button>
                </Footer>
              </Layout>

            </TabPane>
            <TabPane label="插件" icon="logo-dropbox">
              <Layout style="height: 100%">

                <Sider hide-trigger style="background: none; height: 100%" width="240">
                  <!--<Menu width="auto" @on-select="activePluginEditor" :active-name="activeEditPluginName" ref="pluginMenu">
                    <MenuItem :name="p.name" v-for="p in plugins" :key="p.name" style="text-align: left">
                      <Icon type="md-flash" />
                      {{p.name}}
                      <Icon type="md-remove-circle"
                            title="删除"
                            style="float: right; color: red;"
                            @click="removePlugin(p.name)"/>
                    </MenuItem>
                  </Menu>-->


                  <Card title="插件列表">
                    <CellGroup @on-click="selectPluginMenuItem">
                      <Cell :title="p.name"
                            :name="p.name"
                            v-for="p in plugins" :key="p.name"
                            :selected="activeEditPluginName == p.name">
                        <Icon type="md-remove-circle"
                              title="删除"
                              slot="extra"
                              style="color: red;"
                              @click="removePlugin(p.name)"/>
                      </Cell>
                    </CellGroup>
                  </Card>

                  <Button icon="md-add" long type="dashed"
                          dashed ghost
                          style="width: 98%; font-size: 18px; font-weight: bold; color: blue; margin: 5px 0;"
                          @click="newPluginModal.showPluginModal = true"
                  ></Button>
                </Sider>



                <Content style="height: 100%; overflow-y: auto">
                  <Layout style="height: 98%">
                    <Content style="height: 100%; overflow-y: auto">
                      <textarea ref="plugin_textarea"/>
                    </Content>
                    <Footer style="text-align: right; padding: 24px 10px;">
                      <Button icon="md-help-circle" @click="showConfigHelp = true"></Button>
                      &nbsp;&nbsp;&nbsp;
                      <Button icon="md-checkbox-outline" :loading="savePluginLoading" type="success" @click="savePlugins">保存</Button>
                      &nbsp;&nbsp;&nbsp;
                      <Button icon="md-bug" type="success" @click="parseConfig">调试</Button>
                    </Footer>
                  </Layout>
                </Content>



              </Layout>
            </TabPane>
          </Tabs>



        </Content>
      </Layout>


      <Modal
        title="调试选项"
        v-model="showDebugModal"
        :mask-closable="false"
        :loading="debugLoading"
        width="650px">

        <Form :model="debugParams"
              label-position="left"
              :label-width="120">
          <Form-item label="输入调试参数">
            <Input v-model="debugParams.input" placeholder="请输入" clearable/>
          </Form-item>
          <Form-item label="选择调试的stage">
            <Select v-model="debugParams.stageName" style="width:200px">
              <Option v-for="item in debugParams.stages" :value="item.name" :key="item.name">{{ item.name }}</Option>
            </Select>
          </Form-item>
        </Form>

        <Divider dashed>输出结果</Divider>

        <Input v-model="debugOutput" type="textarea" :autosize="{ minRows: 8, maxRows: 8 }" :readonly="true" placeholder="输出结果" />

        <div slot="footer">
          <Button type="default" icon="md-refresh" @click="debugOutput=''">清空</Button>
          <Button type="success" icon="md-bug" :loading="debugProcessLoading" @click="startDebug">调试</Button>
        </div>
      </Modal>


      <Modal title="输入插件名称"
             :footer-hide="true"
             v-model="newPluginModal.showPluginModal">

        <Form :model="newPluginModal"
              :rules="newPluginModal.ruleValidate"
              label-position="left"
              :label-width="120">
          <Form-item label="插件名称" prop="tempPluginName">
            <Input v-model="newPluginModal.tempPluginName" placeholder="字母下划线..." clearable/>
          </Form-item>
          <FormItem style="text-align: right">
            <Button type="success" @click="addNewPlugin">添加</Button>
          </FormItem>
        </Form>
        <div slot="footer">
        </div>
      </Modal>



      <Drawer title="配置指南"
              placement="left"
              :closable="true"
              :transfer="false"
              :mask="false"
              v-model="showConfigHelp">
        <ConfigHelp/>
      </Drawer>



    </div>
</template>

<script>

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

  import {
    getProject,
    parseProjectConfigFile,
    debugStage,
    updateProjectBaseInfo,
    getProjectYamlConfig,
    saveProjectConfig,
    savePlugins,
  } from "@/api/project"
  import ConfigHelp from '../help/config-help'

  export default {
    name: 'project_edit',
    components: {
      ConfigHelp
    },
    data () {
      return {
        split: 0.5,
        crawlerFileEditor: false,
        pluginFileEditor: false,
        showDebugModal: false,
        debugLoading: false,
        debugProcessLoading: false,
        saveProjectLoading: false,
        saveProjectConfigLoading: false,
        savePluginLoading: false,
        importLoading: false,

        showConfigHelp: false,

        uploadHeaders: {
          Authorization: 'Bearer ' + window.localStorage.getItem('token')
        },

        newPluginModal: {
          showPluginModal: false,
          tempPluginName: '',
          ruleValidate: {
            tempPluginName: [
              { required: true, message: '请输入插件名称', trigger: 'blur' },
              { type: 'string', max: 50, message: '最多50字符', trigger: 'blur' }
            ],
          },
        },
        activeEditPluginName: '',

        debugParams: {
          stages: [],
          stageName: '',
          input: '',
        },
        plugins: [],
        debugOutput: '',
        yamlValue: '',
        project: {},
        tags: [],
        tempTag: '',
      }
    },
    props: {
      projectId: {
        type: [Number, String],
        default: 0,
      }
    },
    methods: {
      async getProjectDetail() {
        if (this.projectId == 0) {
          return
        }
        let that = this
        const { data: data } = await getProject(this.projectId)

        if (data && data.code == 0) {

          this.saveProjectLoading = false
          this.saveProjectConfigLoading = false
          this.savePluginLoading = false

          this.project = data.data.project
          let _tags = data.data.project.tags
          if (_tags == '') {
            _tags = '[]'
          }
          this.tags = JSON.parse(_tags)
          this.plugins = data.data.plugins == null ? [] : data.data.plugins
          if (this.plugins.length > 0) {
            that.activePluginEditor(this.plugins[0].name)
          }
          this.crawlerFileEditor.setValue(data.data.yaml)
          this.crawlerFileEditor.refresh()
          this.crawlerFileEditor.focus()
        } else {
          this.$Message.error('加载失败：' + data.msg)
        }
      },

      async parseConfig() {
        let that = this
        const { data: data } = await parseProjectConfigFile({
          project: that.crawlerFileEditor.getValue()
        })
        if (data && data.code == 0) {
          console.log(data.data)
          this.debugParams.stages = data.data.stages
          if (this.debugParams.input == '' && data.data.start_urls != null && data.data.start_urls.length > 0) {
            this.debugParams.input = data.data.start_urls[0]
          }
          this.showDebugModal = true
        } else {
          this.$Message.error('配置解析失败：' + data.msg)
        }
      },

      async startDebug() {
        if (this.debugParams.input == "") {
          this.$Message.error('请输入调试参数')
          return
        }
        if (this.debugParams.stageName == "") {
          this.$Message.error('请选择stage')
          return
        }
        this.debugProcessLoading = true
        let that = this
        debugStage({
          stage_name: that.debugParams.stageName,
          url: that.debugParams.input,
          project: that.crawlerFileEditor.getValue(),
          project_id: that.projectId
        }).then(data=>{
          data = data.data
          if (data && data.code == 0) {
            that.debugOutput = JSON.stringify(data.data, null, "\t")
          } else {
            that.debugOutput = '调试错误：' + data.msg
          }
          this.debugProcessLoading = false
        }).catch(err=>{
          this.debugProcessLoading = false
          this.$Message.error('调试失败：' + err)
        })
      },

      async updateProject() {

        if (this.project.name == "") {
          this.$Message.error('项目名称不能为空')
          return
        }
        this.saveProjectLoading = true

        let that = this
        const { data: data } = await updateProjectBaseInfo({
          id: that.project.id,
          name: that.project.name,
          display_name: that.project.display_name,
          remark: that.project.remark,
          cron: that.project.cron,
          enable_cron: that.project.enable_cron,
          tags: JSON.stringify(that.tags)
        })
        this.saveProjectLoading = false
        if (data && data.code == 0) {
          console.log(data.data)
          this.$Message.success('保存成功')
        } else {
          this.$Message.error('保存失败：' + data.msg)
        }
      },

      async saveProjectConfig() {

        this.saveProjectConfigLoading = true

        let that = this
        const { data: data } = await saveProjectConfig({
          id: that.projectId,
          project: that.crawlerFileEditor.getValue()
        })
        this.saveProjectConfigLoading = false
        if (data && data.code == 0) {
          console.log(data.data)
          this.$Message.success('保存成功')
        } else {
          this.$Message.error('保存失败：' + data.msg)
        }
      },

      addTag() {
        if (this.tempTag.trim() == '') {
          return
        }
        let exists = false
        this.tags.forEach(v => {
          if (v == this.tempTag.trim()) {
            exists = true
          }
        })
        if (exists) {
          this.tempTag = ''
          return
        }
        this.tags.push(this.tempTag.trim())
        this.tempTag = ''
      },

      removeTag(tag) {
        let newTags = new Array()
        this.tags.forEach(v => {
          if (v != tag) {
            newTags.push(v)
          }
        })
        this.tags = newTags
      },

      addNewPlugin() {

        if (this.newPluginModal.tempPluginName.trim() == '' || this.newPluginModal.tempPluginName.trim().length > 50) {
          this.$Message.error('插件名称非法');
          return
        }

        let exists = false
        this.plugins.forEach(v => {
          if (v.name == this.newPluginModal.tempPluginName.trim()) {
            exists = true
          }
        })
        if (exists) {
          this.$Message.error('插件名称已存在');
          return
        }
        this.plugins.push({
          name: this.newPluginModal.tempPluginName,
          script: `(function(){
  // Plugin Name: `+this.newPluginModal.tempPluginName+`
  // Start here...
})()`,
        })
        this.selectPluginMenuItem(this.newPluginModal.tempPluginName)
        this.newPluginModal.showPluginModal = false
        this.newPluginModal.tempPluginName = ''
      },

      activePluginEditor(name) {
        this.selectPluginMenuItem(name)
      },

      selectPluginMenuItem(name) {
        let newSc = ''
        if (name == '') {
          this.pluginFileEditor.setValue(newSc)
          this.pluginFileEditor.getWrapperElement().style.height='100%'
          this.pluginFileEditor.refresh();
          this.pluginFileEditor.focus();
          this.activeEditPluginName = name
          console.log('current plugin：' + this.activeEditPluginName)
          return
        }
        let exist = false

        this.plugins.forEach(v => {
          if (v.name == name) {
            newSc = v.script
            exist = true
          }
        })

        // 拦截
        if (!exist) {
          return
        }

        if (this.activeEditPluginName != '') {
          this.plugins.forEach(v => {
            if (v.name == this.activeEditPluginName) {
              v.script = this.pluginFileEditor.getValue()
            }
          })
        }
        this.plugins.forEach(v => {
          if (v.name == name) {
            newSc = v.script
          }
        })
        this.pluginFileEditor.setValue(newSc)
        this.pluginFileEditor.getWrapperElement().style.height='100%'
        this.pluginFileEditor.refresh();
        this.pluginFileEditor.focus();
        this.activeEditPluginName = name

        console.log('current plugin：' + this.activeEditPluginName)

      },

      removePlugin(name) {

        console.log('remove:' + name)

        let newArr = new Array()
        let hasNext = true
        let hasBefore = true
        let nextName = ''
        let beforeName = ''

        for (let i = 0; i < this.plugins.length; i++) {
          let v = this.plugins[i]
          if (v.name != name) {
            newArr.push(v)
          } else {
            if (i == this.plugins.length - 1) {
              hasNext = false
            } else {
              nextName = this.plugins[i+1].name
            }
            if (i == 0) {
              hasBefore = false
            } else {
              beforeName = this.plugins[i-1].name
            }
          }
        }

        console.log('nextName=' + nextName)
        console.log('beforeName=' + beforeName)

        this.plugins = newArr
        if (hasNext) {
          this.selectPluginMenuItem(nextName)
        } else if (hasBefore) {
          this.selectPluginMenuItem(beforeName)
        } else {
          this.selectPluginMenuItem('')
        }
      },

      async savePlugins() {
        this.savePluginLoading = true
        // 先同步编辑器中的数据到队列中
        if (this.activeEditPluginName != '') {
          this.plugins.forEach(v => {
            if (v.name == this.activeEditPluginName) {
              v.script = this.pluginFileEditor.getValue()
            }
          })
        }
        let that = this
        const { data: data } = await savePlugins({
          projectId: that.projectId,
          plugins: that.plugins
        })
        this.savePluginLoading = false
        if (data && data.code == 0) {
          console.log(data.data)
          this.$Message.success({
            content: '保存成功',
            duration: 5
          })
        } else {
          this.$Message.error({
            content: '保存失败：' + data.msg,
            duration: 5
          })
        }
      },

      uploadConfigSuccess(res, file) {
        this.importLoading = false
        console.log(res)
        if (res && res.code == 0) {
          this.$Message.success('导入成功');

          this.saveProjectLoading = true
          this.saveProjectConfigLoading = true
          this.savePluginLoading = true
          this.getProjectDetail()
        } else {
          this.$Message.error('导入失败');
        }
      },

      exportProjectConfig() {
        window.open('/api/v1/projects/' + this.projectId + '/export?token=' + window.localStorage.getItem('token'))
      },

      showLoading() {
        this.importLoading = true
      },
      hideLoading() {
        this.importLoading = false
      },

    },

    watch: {
      projectId(val) {
        this.plugins = []
        this.project = {}
        this.tags = []
        this.activeEditPluginName = ''
        this.crawlerFileEditor && this.crawlerFileEditor.setValue('')
        this.pluginFileEditor && this.pluginFileEditor.setValue('')
        this.saveProjectLoading = true
        this.saveProjectConfigLoading = true
        this.savePluginLoading = true
        this.getProjectDetail()
      },
    },

    mounted() {
      let that = this
      // 初始化配置文件编辑器
      this.crawlerFileEditor = CodeMirror.fromTextArea(this.$refs.config_textarea, {
        lineNumbers: true, // 显示行号
        mode: 'text/x-yaml', // 语法model
        gutters: ['CodeMirror-lint-markers'],  // 语法检查器
        theme: 'darcula', // 编辑器主题
        tabSize: 2,
        autofocus: true,
        lint: false // 开启语法检查
      })
      this.crawlerFileEditor.setValue('')
      this.crawlerFileEditor.getWrapperElement().style.fontSize='15px'
      this.crawlerFileEditor.getWrapperElement().style.fontFamily='Consolas'
      this.crawlerFileEditor.getWrapperElement().style.height='100%'
      this.crawlerFileEditor.refresh();

      // 初始化插件编辑器
      this.pluginFileEditor = CodeMirror.fromTextArea(this.$refs.plugin_textarea, {
        lineNumbers: true, // 显示行号
        mode: 'text/javascript', // 语法model
        gutters: ['CodeMirror-lint-markers'],  // 语法检查器
        theme: 'darcula', // 编辑器主题
        tabSize: 2,
        autofocus: true,
        lint: false // 开启语法检查
      })
      this.pluginFileEditor.setValue('')
      this.pluginFileEditor.getWrapperElement().style.fontSize='15px'
      this.pluginFileEditor.getWrapperElement().style.fontFamily='Consolas'
      this.pluginFileEditor.getWrapperElement().style.height='100%'
      this.pluginFileEditor.refresh();


      CodeMirror.commands.saveProject = function (cm) {
        that.saveProjectConfig()
      }
      CodeMirror.commands.savePlugin = function (cm) {
        that.savePlugins()
      }
      // 判断是否为Mac
      let mac = CodeMirror.keyMap.default == CodeMirror.keyMap.macDefault
      let runKey = (mac ? "Cmd" : "Ctrl") + "-S"
      let projectExtraKeys = {}
      projectExtraKeys[runKey] = "saveProject"
      let pluginExtraKeys = {}
      pluginExtraKeys[runKey] = "savePlugin"

      this.crawlerFileEditor.setOption("extraKeys", projectExtraKeys)
      this.pluginFileEditor.setOption("extraKeys", pluginExtraKeys)
    }
  }
</script>
<style>
  .editor-tabs{
    height: 100%;
  }
  .editor-tabs .ivu-tabs-tabpane, .editor-tabs .ivu-tabs-content{
    height: 98%;
  }
</style>
