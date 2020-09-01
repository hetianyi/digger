<template>
  <div
    id="__scroll_container">
      <p v-for="log in logs" v-html="log"></p>

      <Spin>
        <Icon type="ios-loading" size=18 class="demo-spin-icon-load"></Icon>
      </Spin>

  </div>
</template>

<script>
  export default {
    name: 'task_logs',
    data () {
      return {
        ws: null,
        logs: []
      }
    },
    props: {
      taskId: {
        type: [Number, String],
        default: 0,
      }
    },
    methods: {
      startLogs() {
        if (this.ws) {
          this.ws.close()
        }
        let that = this

        let pro = window.location.protocol

        let _ws = new WebSocket(
          (pro.startsWith('https:') ? 'wss://' : 'ws://')
          + window.location.host + '/api/tasks/'+ this.taskId +'/logs/ws');

        // let _ws = new WebSocket(
        //   (pro.startsWith('https:') ? 'wss://' : 'ws://')
        //   + 'localhost:9012/api/tasks/'+ this.taskId +'/logs/ws');

        _ws.onopen = function(evt) {
          that.logs.push('开始监听日志')
          console.log('开始监听日志')
        }
        _ws.onclose = function(evt) {
          that.logs.push('日志监听断开')
          console.log('日志监听断开')
        }
        _ws.onmessage = function(evt) {
          // console.log(decodeURIComponent(window.atob(evt.data)))
          // console.log(decodeURIComponent(decodeURIComponent(window.atob(evt.data))))
          that.logs.push(decodeURIComponent(window.atob(evt.data).replace(/\+/g, " ")))
          setTimeout(function () {
            let div = document.getElementById('__scroll_container')
            div.parentElement.scrollTop = div.parentElement.scrollHeight
          }, 200)
        }
        _ws.onerror = function(evt) {
          console.log('日志监听错误')
          that.logs.push(evt.data)
        }
        this.ws = _ws
      },
    },

    watch: {
      taskId(val) {
        if (val == 0) {
          console.log('关闭websocket')
          this.logs = []
          if (this.ws) {
            console.log('关闭')
            this.ws.close()
          }
          return
        } else {
          this.startLogs()
        }
      },
    },

    mounted() {
    }
  }
</script>
<style>
  .demo-spin-icon-load{
    animation: ani-demo-spin 1s linear infinite;
  }
  @keyframes ani-demo-spin {
    from { transform: rotate(0deg);}
    50%  { transform: rotate(180deg);}
    to   { transform: rotate(360deg);}
  }
  .demo-spin-col{
    height: 100px;
    position: relative;
    border: 1px solid #eee;
  }
  #__scroll_container{
    width: 100%;
    background: #001529;
    color: #c3c3c3;
    padding: 10px 10px;
  }
</style>
