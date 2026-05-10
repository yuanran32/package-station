<template>
  <div class="page">
    <div class="page-header">
      <div>
        <h2 class="page-title">支付账单</h2>
        <p class="page-desc">查看支付账单，并对待支付寄件单执行模拟支付。</p>
      </div>
      <el-button type="primary" :icon="Refresh" :loading="loading" @click="load">刷新</el-button>
    </div>

    <el-alert
      v-if="selectedOrderNo"
      :title="`当前寄件单：${selectedOrderNo}`"
      type="info"
      :closable="false"
      show-icon
    />

    <section class="section-panel">
      <el-table v-loading="loading" :data="rows" border stripe empty-text="暂无账单">
        <el-table-column prop="pay_no" label="支付单号" min-width="180" show-overflow-tooltip />
        <el-table-column prop="related_type" label="业务类型" min-width="120" />
        <el-table-column label="关联订单号" min-width="190" show-overflow-tooltip>
          <template #default="{ row }">
            {{ row.related_order_no || '-' }}
          </template>
        </el-table-column>
        <el-table-column label="金额" min-width="100">
          <template #default="{ row }">
            ¥{{ formatAmount(row.amount) }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="支付状态" min-width="120" />
        <el-table-column prop="created_at" label="创建时间" min-width="180" show-overflow-tooltip />
        <el-table-column label="操作" width="210" fixed="right">
          <template #default="{ row }">
            <el-button
              size="small"
              type="success"
              :loading="payingPayNo === getPayNo(row)"
              :disabled="!canPay(row)"
              @click="mockPay(row)"
            >
              立即支付
            </el-button>
            <el-button size="small" text @click="copyPayNo(row)">复制支付单号</el-button>
          </template>
        </el-table-column>
      </el-table>
    </section>
  </div>
</template>

<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { Refresh } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { getAdminSendOrders } from '../../api/send'
import { getBills, payCallback } from '../../api/pay'
import { useAuthStore } from '../../stores/auth'

const RECENT_SEND_ORDER_KEY = 'package_station_recent_send_orders'
const loading = ref(false)
const payingPayNo = ref('')
const rows = ref([])
const route = useRoute()
const auth = useAuthStore()

const isAdmin = computed(() => auth.role === 'admin')
const selectedOrderNo = computed(() => (route.query.order_no || '').toString())

function getPayNo(row) {
  return row.pay_no || row.payNo || ''
}

function getOrderNo(row) {
  return row.order_no || row.orderNo || row.id || ''
}

function formatAmount(value) {
  const amount = Number(value)
  return Number.isFinite(amount) ? amount.toFixed(2) : '--'
}

function canPay(row) {
  const status = (row.status || '').toLowerCase()
  return status === 'pending' || status === 'created' || status === 'unpaid'
}

function readRecentOrderMap() {
  try {
    const list = JSON.parse(localStorage.getItem(RECENT_SEND_ORDER_KEY) || '[]')
    if (!Array.isArray(list)) {
      return new Map()
    }
    return new Map(list.map((item) => [item.order_no, item]))
  } catch {
    return new Map()
  }
}

async function fetchOrderNoByIdMap() {
  if (!isAdmin.value) {
    return new Map()
  }
  try {
    const data = await getAdminSendOrders('')
    const orders = Array.isArray(data) ? data : data?.list || []
    return new Map(orders.map((order) => [order.id, getOrderNo(order)]))
  } catch {
    return new Map()
  }
}

async function load() {
  loading.value = true
  try {
    const [billData, orderNoById] = await Promise.all([getBills(), fetchOrderNoByIdMap()])
    const billRows = Array.isArray(billData) ? billData : billData?.list || []
    const recentMap = readRecentOrderMap()

    rows.value = billRows.map((row) => {
      const payNo = getPayNo(row)
      const relatedType = row.related_type || row.relatedType || ''
      const relatedId = row.related_id ?? row.relatedId
      const relatedOrderNo =
        row.related_order_no ||
        row.relatedOrderNo ||
        orderNoById.get(relatedId) ||
        recentMap.get(row.order_no || '')?.order_no ||
        ''
      return {
        ...row,
        pay_no: payNo,
        related_type: relatedType,
        related_order_no: relatedOrderNo
      }
    })
  } finally {
    loading.value = false
  }
}

async function mockPay(row) {
  const payNo = getPayNo(row)
  if (!payNo) {
    ElMessage.warning('缺少支付单号，无法支付')
    return
  }
  payingPayNo.value = payNo
  try {
    await payCallback({
      pay_no: payNo,
      status: 'paid'
    })
    ElMessage.success(`支付成功：${payNo}`)
    await load()
  } finally {
    payingPayNo.value = ''
  }
}

async function copyPayNo(row) {
  const payNo = getPayNo(row)
  if (!payNo) {
    ElMessage.warning('当前记录没有支付单号')
    return
  }
  try {
    await navigator.clipboard.writeText(payNo)
    ElMessage.success('支付单号已复制')
  } catch {
    ElMessage.warning('复制失败，请手动复制')
  }
}

watch(
  () => route.query.pay_no,
  (payNo, oldPayNo) => {
    if (payNo && payNo !== oldPayNo) {
      ElMessage.success(`已定位支付单：${payNo}`)
    }
  }
)

onMounted(load)
</script>
