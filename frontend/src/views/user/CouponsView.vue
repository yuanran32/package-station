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
      <el-table v-loading="loading" :data="rows" border stripe empty-text="暂无礼券">
        <el-table-column prop="code" label="礼券码" min-width="130" show-overflow-tooltip />
        <el-table-column prop="name" label="礼券名称" min-width="160" show-overflow-tooltip />
        <el-table-column label="面额" min-width="90">
          <template #default="{ row }">
            ¥{{ formatAmount(row.amount) }}
          </template>
        </el-table-column>
        <el-table-column label="门槛" min-width="90">
          <template #default="{ row }">
            ¥{{ formatAmount(row.threshold) }}
          </template>
        </el-table-column>
        <el-table-column prop="status_text" label="状态" min-width="100" />
        <el-table-column prop="received_at_text" label="领取时间" min-width="180" show-overflow-tooltip />
        <el-table-column prop="expire_time_text" label="过期时间" min-width="180" show-overflow-tooltip />
      </el-table>
    </section>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { getMyCoupons, receiveCoupon } from '../../api/coupon'

const couponCode = ref('')
const loading = ref(false)
const receiving = ref(false)
const rows = ref([])

function statusToText(status) {
  const value = (status || '').toLowerCase()
  if (value === 'unused') return '未使用'
  if (value === 'used') return '已使用'
  if (value === 'expired') return '已过期'
  if (value === 'invalid') return '已失效'
  return status || '-'
}

function formatAmount(value) {
  const amount = Number(value)
  return Number.isFinite(amount) ? amount.toFixed(2) : '0.00'
}

function formatTime(value) {
  if (!value) {
    return '-'
  }
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) {
    return value
  }
  const pad = (n) => String(n).padStart(2, '0')
  const y = date.getFullYear()
  const m = pad(date.getMonth() + 1)
  const d = pad(date.getDate())
  const hh = pad(date.getHours())
  const mm = pad(date.getMinutes())
  const ss = pad(date.getSeconds())
  return `${y}-${m}-${d} ${hh}:${mm}:${ss}`
}

function normalizeRows(data) {
  const list = Array.isArray(data) ? data : data?.list || []
  return list.map((item) => {
    const usedAt = item.used_at || item.usedAt || ''
    const receivedAt = item.received_at || item.receivedAt || ''
    const expireAt = item.expire_time || item.expireTime || item.end_at || item.endAt || ''
    const status = item.status || ''

    let expireText = formatTime(expireAt)
    if (status.toLowerCase() === 'used' && usedAt) {
      expireText = '已使用'
    }

    return {
      ...item,
      code: item.code || item.coupon_code || item.couponCode || '-',
      name: item.name || item.coupon_name || item.couponName || '-',
      threshold: item.threshold ?? 0,
      status_text: statusToText(status),
      received_at_text: formatTime(receivedAt),
      expire_time_text: expireText
    }
  })
}

async function load() {
  loading.value = true
  try {
    const data = await getMyCoupons()
    rows.value = normalizeRows(data)
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
