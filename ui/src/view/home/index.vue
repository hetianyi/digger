<template>
  <div>

    <Row :gutter="20">
      <i-col :xs="12" :md="8" :lg="4" v-for="(infor, i) in inforCardData" :key="`infor-${i}`" style="height: 120px;padding-bottom: 10px;">
        <infor-card shadow :color="infor.color" :icon="infor.icon" :icon-size="36">
          <count-to :end="infor.count" count-class="count-style"/>
          <p>{{ infor.title }}</p>
        </infor-card>
      </i-col>
    </Row>

    <Row>
      <Col span="24" align="left" style="margin: 10px 0;">
        <DatePicker type="datetimerange"
                    placeholder="选择日期范围"
                    style="width: 300px"
                    format="yyyy-MM-dd HH:mm:ss"
                    v-model="dateRange"
                    @on-ok="dateChange"
        ></DatePicker>

        <Button type="success" style="margin-left: 10px;" @click="resetDate">现在</Button>
        <br>
      </Col>
    </Row>
    <Row>
      <Col span="24">
        <Card shadow>
          <example style="height: 350px;" :xAxis="xAxis" :series="series"/>
        </Card>
      </Col>
    </Row>
  </div>
</template>

<script>
  import {
    getStatistic,
  } from "@/api/statistic"
  import InforCard from '_c/info-card'
  import CountTo from '_c/count-to'
  import Example from './example.vue'


  export default {
    name: 'task_list',
    components: {
      InforCard,
      CountTo,
      Example,
    },
    data () {
      return {
        dateRange: [new Date(new Date().getTime()-86400000), new Date()],

        inforCardData: [/*
          { title: '项目', icon: 'md-person-add', count: 803, color: '#2d8cf0' },
          { title: '任务', icon: 'md-locate', count: 232, color: '#19be6b' },
          { title: '工作节点', icon: 'md-help-circle', count: 142, color: '#ff9900' },
          { title: '累计请求', icon: 'md-share', count: 657, color: '#ed3f14' },
          { title: '累计结果', icon: 'md-chatbubbles', count: 12, color: '#E46CBB' },
        */],
        xAxis: [],
        series: [],
      }
    },

    computed: {
    },

    watch: {
      /*dateRange(newRange) {
        console.log(newRange)
      }*/
    },

    methods: {

      async listStatistic() {

        let start = this.formatDate(this.dateRange[0])
        let end = this.formatDate(this.dateRange[1])
        console.log(start + ' --- ' + end)

        const { data: data } = await getStatistic({
          start: start,
          end: end,
        })
        if (data && data.code == 0) {
            this.xAxis = data.data.xAxis
            this.series = data.data.series
            this.inforCardData = data.data.inforCardData
        } else {
          this.$Message.error('加载失败：' + data.msg)
        }
      },

      dateChange() {
        console.log(this.formatDate(this.dateRange[0]))
        console.log(this.formatDate(this.dateRange[1]))
        this.listStatistic()
      },

      resetDate() {
        this.dateRange = [new Date(new Date().getTime()-86400000), new Date()]
        this.listStatistic()
      },
      formatDate (date) {
        let yyyy = date.getFullYear() // 获取四位数字表示的年份

        /* 获取月份，
                +1是因为getMonth方法返回的0(一月份)-11（十二月份） */
        let mm = date.getMonth() + 1

        let dd = date.getDate() // 获取日期
        let hh = date.getHours() // 获取小时数
        let min = date.getMinutes() // 获取分钟数
        let ss = date.getSeconds() // 获取秒
        // 分隔符
        let sep1 = '-'
        let sep2 = ':'
        // 用"0"补位不足两位数的时间
        mm = (mm < 10) ? ('0' + mm) : mm
        dd = (dd < 10) ? ('0' + dd) : dd
        hh = (hh < 10) ? ('0' + hh) : hh
        min = (min < 10) ? ('0' + min) : min
        ss = (ss < 10) ? ('0' + ss) : ss
        return yyyy + sep1 + mm + sep1 + dd + ' ' + hh + sep2 + min + sep2 + ss
      },

    },

    mounted() {
      let that = this
      this.listStatistic()
      /*setInterval(function () {
        that.listStatistic()
      }, 10000)*/
    }
  }
</script>
<style>
  .count-style{
    font-size: 50px;
  }
</style>
