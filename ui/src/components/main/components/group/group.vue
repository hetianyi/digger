<template>
  <div>
    <Dropdown trigger="hover" style="margin-right: 20px;">
      <a href="javascript:void(0)">
        加入讨论组
      </a>
      <DropdownMenu slot="list" style="width: 200px; height: 220px; background: #CEBB3B; padding: 0">
        <p style="width: 100%; text-align: center;">微信扫一扫加入讨论群</p>
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

<style lang="less">

</style>
