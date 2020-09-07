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
        grid: {
          top: '15%',
          left: '1.2%',
          right: '1%',
          bottom: '3%',
          containLabel: true
        },
        legend: {
          orient: 'horizontal',
          icon: 'roundRect',
          x: 'left',
          data: []
        },
        xAxis: [
          {
            type: 'time',
          }
        ],
        yAxis: {
          type: 'value',
          boundaryGap: [0, '100%'],
          splitLine: {
            show: false
          }
        },
        /*dataZoom: [{
          startValue: '2014-06-01'
        }, {
          type: 'inside'
        }],*/
        series: []
      },
    }
  },
  props: {
    xAxis: Array,
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
        if (this.series && this.series.length > 0) {
          for (let i = 0; i < this.series.length; i++) {
            let res = this.series[i]
            res.symbol = "none"
            that.option.series.push(res)
            that.option.legend.data[i] = res.name
          }
        }
        this.dom = echarts.init(this.$refs.dom, 'tdTheme')
        this.dom.setOption(that.option)
        on(window, 'resize', this.resize)
      })
    },
  },

  mounted () {
    let that = this
    this.$nextTick(() => {
      that.option.series = []
      that.option.legend.data = []
      if (this.series && this.series.length > 0) {
        for (let i = 0; i < this.series.length; i++) {
          let res = this.series[i]
          res.showSymbol = false
          res.clip = true
          that.option.series.push(res)
          that.option.legend.data[i] = res.name
        }
      }

      this.dom = echarts.init(this.$refs.dom, 'tdTheme')
      this.dom.setOption(that.option)
      on(window, 'resize', this.resize)
    })
  },
  beforeDestroy () {
    off(window, 'resize', this.resize)
  }
}
</script>
