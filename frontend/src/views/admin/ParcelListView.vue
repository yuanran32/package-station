<template>
  <div class="page">
    <div class="page-header">
      <div>
        <h2 class="page-title">快递列表</h2>
        <p class="page-desc">查询并浏览所有入库快递。</p>
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
import { getParcelList } from '../../api/parcel'
import DataTable from '../../components/DataTable.vue'

const loading = ref(false)
const rows = ref([])
const columns = [
  { prop: 'tracking_no', label: '快递单号', minWidth: 160 },
  { prop: 'recipient_phone', label: '收件人手机号', minWidth: 140 },
  { prop: 'location', label: '存储位置' },
  { prop: 'pickup_code', label: '取件码' },
  { prop: 'status', label: '状态' },
  { prop: 'created_at', label: '入库时间', minWidth: 180 }
]

async function load() {
  loading.value = true
  try {
    const data = await getParcelList()
    rows.value = Array.isArray(data) ? data : data?.list || []
  } finally {
    loading.value = false
  }
}

onMounted(load)
</script>
