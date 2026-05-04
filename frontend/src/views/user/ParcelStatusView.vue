<template>
  <div class="page">
    <div class="page-header">
      <div>
        <h2 class="page-title">快递查询</h2>
        <p class="page-desc">输入快递单号，查看当前状态和存储位置。</p>
      </div>
    </div>
    <section class="section-panel">
      <div class="toolbar">
        <el-input v-model.trim="trackingNo" placeholder="请输入快递单号" clearable style="max-width: 360px" />
        <el-button type="primary" :icon="Search" :loading="loading" @click="query">查询</el-button>
      </div>
    </section>
    <section v-if="result" class="section-panel">
      <h3 class="section-title">查询结果</h3>
      <el-descriptions :column="2" border>
        <el-descriptions-item label="快递单号">{{ result.tracking_no || result.trackingNo || trackingNo }}</el-descriptions-item>
        <el-descriptions-item label="状态">{{ result.status || '未知' }}</el-descriptions-item>
        <el-descriptions-item label="存储位置">{{ result.location || '-' }}</el-descriptions-item>
        <el-descriptions-item label="取件码">{{ result.pickup_code || result.pickupCode || '-' }}</el-descriptions-item>
      </el-descriptions>
    </section>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import { Search } from '@element-plus/icons-vue'
import { getParcelStatus } from '../../api/parcel'

const trackingNo = ref('')
const loading = ref(false)
const result = ref(null)

async function query() {
  if (!trackingNo.value) {
    ElMessage.warning('请输入快递单号')
    return
  }
  loading.value = true
  try {
    result.value = await getParcelStatus(trackingNo.value)
  } finally {
    loading.value = false
  }
}
</script>
