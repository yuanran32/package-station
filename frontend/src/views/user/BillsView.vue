<template>
  <div class="page">
    <div class="page-header">
      <div>
        <h2 class="page-title">支付账单</h2>
        <p class="page-desc">查看寄件订单和仓储费用相关支付记录。</p>
      </div>
      <el-button type="primary" :icon="Refresh" :loading="loading" @click="load">刷新</el-button>
    </div>
    <section class="section-panel">
      <DataTable :rows="rows" :columns="columns" :loading="loading" />
    </section>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { Refresh } from '@element-plus/icons-vue'
import { getBills } from '../../api/pay'
import DataTable from '../../components/DataTable.vue'

const loading = ref(false)
const rows = ref([])
const columns = [
  { prop: 'bill_no', label: '账单号', minWidth: 160 },
  { prop: 'related_type', label: '业务类型' },
  { prop: 'amount', label: '金额' },
  { prop: 'status', label: '支付状态' },
  { prop: 'created_at', label: '创建时间', minWidth: 180 }
]

async function load() {
  loading.value = true
  try {
    const data = await getBills()
    rows.value = Array.isArray(data) ? data : data?.list || []
  } finally {
    loading.value = false
  }
}

onMounted(load)
</script>
