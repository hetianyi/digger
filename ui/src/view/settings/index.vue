<template>
  <div>

    <Form :model="configs" label-position="left" :label-width="150" style="width: 500px;">

      <FormItem label="管理员用户名">
        <Row>
          <Col span="18">
            <Input type="text" v-model="configs.admin_user" placeholder=""></Input>
          </Col>
          <Col span="4" offset="1">
            <Button type="primary" @click="saveConfig('admin_user', configs.admin_user)">保存</Button>
          </Col>
        </Row>
      </FormItem>

      <FormItem label="管理员登录密码">
        <Row>
          <Col span="18">
            <Input type="text" v-model="configs.admin_password" placeholder=""></Input>
          </Col>
          <Col span="4" offset="1">
            <Button type="primary" @click="saveConfig('admin_password', configs.admin_password)">保存</Button>
          </Col>
        </Row>
      </FormItem>

      <FormItem label="通知邮箱">
        <Row>
          <Col span="18">
            <Input type="text" v-model="configs.email_config" placeholder="格式：user:pass@host:port"></Input>
          </Col>
          <Col span="4" offset="1">
            <Button type="primary" @click="saveConfig('email_config', configs.email_config)">保存</Button>
          </Col>
        </Row>
      </FormItem>

      <FormItem label="钉钉机器人webhook">
        <Row>
          <Col span="18">
            <Input type="text" v-model="configs.dingtalk" placeholder="Enter something..."></Input>
          </Col>
          <Col span="4" offset="1">
            <Button type="primary" @click="">保存</Button>
          </Col>
        </Row>
      </FormItem>

    </Form>

  </div>
</template>

<script>
  import {
    configList,
    updateList,
  } from "@/api/settings"

  export default {
    name: 'node_list',
    data () {
      return {
        loading: false,
        configs: {
          admin_user: '',
          admin_password: '',
          email_config: "",
          dingtalk: "",
        },
      }
    },

    computed: {

    },

    methods: {

      async listConfigs() {
        const { data: data } = await configList()
        if (data && data.code == 0) {
          this.configs = data.data == null ? {} : data.data
        } else {
          this.$Message.error('加载失败：' + data.msg)
        }
      },

      async saveConfig(key, value) {
        let d = {
          key: key,
          value: value,
        }
        const { data: data } = await updateList(d)
        if (data && data.code == 0) {
          this.$Message.success('保存成功')
        } else {
          this.$Message.error('保存失败：' + data.msg)
        }
      },

    },

    mounted() {
      this.listConfigs()
    }
  }
</script>
<style>

</style>
