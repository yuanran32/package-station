<template>
  <el-table v-loading="loading" :data="rows" border stripe empty-text="暂无数据">
    <el-table-column
      v-for="column in columns"
      :key="column.prop"
      :prop="column.prop"
      :label="column.label"
      :min-width="column.minWidth || 120"
      show-overflow-tooltip
    />
    <slot />
  </el-table>
</template>

<script setup>
defineProps({
  rows: {
    type: Array,
    default: () => []
  },
  columns: {
    type: Array,
    default: () => []
  },
  loading: {
    type: Boolean,
    default: false
  }
})
</script>
