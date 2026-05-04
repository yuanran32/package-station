<template>
  <div class="page">
    <div class="page-header">
      <div>
        <h2 class="page-title">红包礼券</h2>
        <p class="page-desc">领取礼券并查看当前账号可用优惠。</p>
      </div>
      <div class="toolbar">
        <el-input v-model.trim="couponCode" placeholder="礼券码" clearable />
        <el-button type="primary" :loading="receiving" @click="receive">领取</el-button>
      </div>
    </div>
    <section class="section-panel">
      <DataTable :rows="rows" :columns="columns" :loading="loading" />
    </section>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { getMyCoupons, receiveCoupon } from '../../api/coupon'
import DataTable from '../../components/DataTable.vue'

const couponCode = ref('')
const loading = ref(false)
const receiving = ref(false)
const rows = ref([])
const columns = [
  { prop: 'coupon_code', label: '礼券码' },
  { prop: 'amount', label: '金额' },
  { prop: 'status', label: '状态' },
  { prop: 'expire_time', label: '过期时间', minWidth: 180 }
]

async function load() {
  loading.value = true
  try {
    const data = await getMyCoupons()
    rows.value = Array.isArray(data) ? data : data?.list || []
  } finally {
    loading.value = false
  }
}

async function receive() {
  if (!couponCode.value) {
    ElMessage.warning('请输入礼券码')
    return
  }
  receiving.value = true
  try {
    await receiveCoupon({ coupon_code: couponCode.value })
    ElMessage.success('领取成功')
    couponCode.value = ''
    await load()
  } finally {
    receiving.value = false
  }
}

onMounted(load)
</script>
