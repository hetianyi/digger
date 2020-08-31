
<style lang="less">
  .ivu-select-dropdown {
    margin: 0;
    padding: 0;
  }
</style>


<template>
  <div>
    <Dropdown trigger="hover" style="margin: 0 20px; padding: 0; ">
      <a href="javascript:void(0)">
        加入讨论组
      </a>
      <DropdownMenu slot="list" style="width: 200px; background: #CEBB3B; ">
        <p style="width: 100%; text-align: center; color: #FFF; margin: 0; padding: 0; font-weight: bold; font-size: 14px;">微信扫一扫加入讨论群</p>
        <img :src="require('@/assets/images/author_qr.jpg')"
        style="width: 200px; height: 200px;"
        >
      </DropdownMenu>
    </Dropdown>
  </div>
</template>

<script>

export default {
  name: 'Disz',
  computed: {
  },
  methods: {

  },
  mounted () {

  }
}
</script>
