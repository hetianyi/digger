<template>
    <div ref="dom"></div>
</template>

<script>
import echarts from 'echarts'
import { on, off } from '@/libs/tools'
import tdTheme from '../../components/charts/theme'

echarts.registerTheme('tdTheme', tdTheme)

export default {
  name: 'serviceRequests',
  data () {
    return {
      dom: null,
      color: ['#2d8cf0', '#10A6FF', '#0C17A6'],
      option: {
        backgroundColor: 'transparent',
        title: {
          trigger: 'axis',
          text: '服务统计',
          subtext: this.subtext,
          x: 'right',
          y: 'top',
          bottom: '10'
        },
        tooltip: {
          trigger: 'axis',
        },
        legend: {
          orient: 'horizontal',
          icon: 'roundRect',
          left: '1.2%',
          top: '0',
          data: []
        },
        grid: {
          top: '15%',
          left: '1.2%',
          right: '1%',
          bottom: '3%',
          borderColor: '#F55F5F',
          containLabel: true
        },
        dataZoom: {
          type: 'inside',
          start: 0,
          end: 100,
          handleSize: 8
        },
        xAxis: [],
        yAxis: {
          type: 'value',
          axisLine: {
            lineStyle: {
              color: "#656565"
            }
          },
        },
        series: []
      },
    }
  },
  props: {
    series: Array,
  },
  methods: {
    resize () {
      this.dom.resize()
    }
  },

  watch: {
    series() {
      let that = this
      this.$nextTick(() => {
        that.option.series = []
        that.option.legend.data = []
        that.option.xAxis = []
        if (this.series && this.series.length > 0) {
          for (let i = 0; i < this.series.length; i++) {
            let res = this.series[i]
            that.option.series.push({
              type: 'line',
              symbol: "none",
              name: res.name,
              color: res.color,
              lineStyle: {
                type: "solid"
              },
              smooth: true,
              data: res.data.map(function (item) {
                return Number(item[1]);
              })
            })
            that.option.xAxis.push({
              color: res.color,
              position: 'bottom',
              axisLine: {
                lineStyle: {
                  color: "#656565"
                }
              },
              data: res.data.map(function (item) {
                return item[0];
              })
            })
            that.option.legend.data.push(res.name)
          }
          console.log(that.option.legend.data)
        }
        this.dom = echarts.init(this.$refs.dom)
        this.dom.setOption(that.option)
        on(window, 'resize', this.resize)
      })
    },
  },

  beforeDestroy () {
    off(window, 'resize', this.resize)
  }
}
</script>
