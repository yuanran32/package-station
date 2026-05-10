<template>
  <div class="page">
    <div class="page-header">
      <div>
        <h2 class="page-title">寄件订单</h2>
        <p class="page-desc">查看用户寄件申请，并执行接单、分配快递员、完成订单。</p>
      </div>
      <div class="toolbar">
        <el-select v-model="status" clearable placeholder="订单状态" style="width: 160px" @change="load">
          <el-option label="待处理" value="created" />
          <el-option label="已接单" value="accepted" />
          <el-option label="派件中" value="assigned" />
          <el-option label="已完成" value="completed" />
        </el-select>
        <el-button type="primary" :icon="Refresh" :loading="loading" @click="load">刷新</el-button>
      </div>
    </div>
    <section class="section-panel">
      <el-table v-loading="loading" :data="rows" border stripe empty-text="暂无订单">
        <el-table-column prop="order_no" label="订单号" min-width="160" show-overflow-tooltip />
        <el-table-column prop="sender_name" label="寄件人" min-width="100" />
        <el-table-column prop="receiver_name" label="收件人" min-width="100" />
        <el-table-column prop="receiver_address" label="收件地址" min-width="120" show-overflow-tooltip />
        <el-table-column label="预估费用" min-width="110">
          <template #default="{ row }">
            ¥{{ formatAmount(row.estimated_fee ?? row.estimatedFee) }}
          </template>
        </el-table-column>
        <el-table-column label="支付状态" min-width="100">
          <template #default="{ row }">
            {{ payStatusText(row.pay_status || row.payStatus) }}
          </template>
        </el-table-column>
        <el-table-column label="状态" min-width="100">
          <template #default="{ row }">
            {{ orderStatusText(row.status) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="300" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="process(row, 'accept')">接单</el-button>
            <el-button size="small" type="warning" @click="openAssign(row)">分配</el-button>
            <el-button size="small" type="success" @click="process(row, 'complete')">完成</el-button>
          </template>
        </el-table-column>
      </el-table>
    </section>
    <el-dialog v-model="assignVisible" title="分配快递员" width="420px">
      <el-form :model="assignForm" label-width="90px">
        <el-form-item label="订单号">
          <el-input v-model="assignForm.order_no" disabled />
        </el-form-item>
        <el-form-item label="快递员">
          <el-input v-model.trim="assignForm.courier_name" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="assignVisible = false">取消</el-button>
        <el-button type="primary" :loading="processing" @click="assign">确认</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Refresh } from '@element-plus/icons-vue'
import { getAdminSendOrders, processSendOrder } from '../../api/send'
import { getParcelStatus } from '../../api/parcel'

const loading = ref(false)
const processing = ref(false)
const status = ref('')
const rows = ref([])
const assignVisible = ref(false)
const assignForm = reactive({
  order_no: '',
  courier_name: ''
})

function getOrderNo(row) {
  return row.order_no || row.orderNo || row.id
}

function getReceiverAddress(row) {
  return row?.receiver_address || row?.receiverAddress || ''
}

function hasValue(value) {
  return value !== undefined && value !== null && value !== ''
}

function getByPath(source, path) {
  return path.split('.').reduce((acc, key) => acc?.[key], source)
}

function pickValue(source, paths) {
  for (const path of paths) {
    const value = getByPath(source, path)
    if (hasValue(value)) {
      return value
    }
  }
  return ''
}

function toBoolean(value) {
  if (!hasValue(value)) return undefined
  if (typeof value === 'boolean') return value
  const text = String(value).trim().toLowerCase()
  if (text === 'true' || text === '1' || text === 'yes' || text === 'y') return true
  if (text === 'false' || text === '0' || text === 'no' || text === 'n') return false
  return undefined
}

function formatBooleanText(value, trueText, falseText) {
  const bool = toBoolean(value)
  if (bool === true) return trueText
  if (bool === false) return falseText
  return '后端未返回（需后端补字段）'
}

function parseCompleteResult(payload, row) {
  const defaultLocation = getReceiverAddress(row)
  return {
    trackingNo:
      pickValue(payload, [
        'inbound_parcel.tracking_no',
        'inboundParcel.tracking_no',
        'inboundParcel.trackingNo',
        'auto_parcel.tracking_no',
        'autoParcel.tracking_no',
        'autoParcel.trackingNo',
        'parcel.tracking_no',
        'parcel.trackingNo',
        'tracking_no',
        'trackingNo'
      ]) || getOrderNo(row),
    location:
      pickValue(payload, [
        'inbound_parcel.location',
        'inboundParcel.location',
        'auto_parcel.location',
        'autoParcel.location',
        'parcel.location',
        'location'
      ]) || defaultLocation,
    pickupCode: pickValue(payload, [
      'inbound_parcel.pickup_code',
      'inboundParcel.pickup_code',
      'inboundParcel.pickupCode',
      'auto_parcel.pickup_code',
      'autoParcel.pickup_code',
      'autoParcel.pickupCode',
      'parcel.pickup_code',
      'parcel.pickupCode',
      'pickup_code',
      'pickupCode'
    ]),
    receiverBound: pickValue(payload, [
      'receiver_bound',
      'receiverBound',
      'receiver_user_bound',
      'receiverUserBound',
      'bind_user',
      'bindUser'
    ]),
    noticeSent: pickValue(payload, [
      'notice_sent',
      'noticeSent',
      'station_notice_sent',
      'stationNoticeSent',
      'notify_sent',
      'notifySent'
    ]),
    status: pickValue(payload, [
      'inbound_parcel.status',
      'inboundParcel.status',
      'parcel.status',
      'status'
    ])
  }
}

async function enrichCompleteResult(rawData, row) {
  const payload = rawData && typeof rawData === 'object' ? rawData : {}
  const result = parseCompleteResult(payload, row)
  const defaultLocation = getReceiverAddress(row)
  const needsParcelFallback = !hasValue(result.pickupCode) || !hasValue(result.location)

  if (!needsParcelFallback || !hasValue(result.trackingNo)) {
    return result
  }

  try {
    const parcelData = await getParcelStatus(result.trackingNo)
    const parcelLocation = pickValue(parcelData, [
      'location',
      'storage_location',
      'storageLocation',
      'pickup_location',
      'pickupLocation'
    ])
    const parcelPickupCode = pickValue(parcelData, [
      'pickup_code',
      'pickupCode',
      'storage_code',
      'storageCode',
      'code'
    ])
    const parcelStatus = pickValue(parcelData, ['status', 'state'])

    if ((!hasValue(result.location) || result.location === defaultLocation) && hasValue(parcelLocation)) {
      result.location = parcelLocation
    }
    if (!hasValue(result.pickupCode) && hasValue(parcelPickupCode)) {
      result.pickupCode = parcelPickupCode
    }
    if (!hasValue(result.status) && hasValue(parcelStatus)) {
      result.status = parcelStatus
    }
  } catch {
    // 补查失败时保留主接口返回，避免影响“完成订单”主流程
  }

  return result
}

function buildCompleteSummary(result) {
  const statusText = hasValue(result.status) ? orderStatusText(result.status) : '后端未返回（可忽略）'

  return [
    '系统已执行“完成订单”。',
    '若收件侧尚未入库，系统会自动创建入库包裹。',
    `入库单号：${result.trackingNo || '-'}`,
    `入库位置：${result.location || '-'}`,
    `取件码：${result.pickupCode || '后端未返回（已尝试自动补查）'}`,
    `包裹状态：${statusText}`,
    `收件人手机号绑定状态：${formatBooleanText(result.receiverBound, '已绑定账号', '未绑定账号')}`,
    `到站通知发送状态：${formatBooleanText(result.noticeSent, '已发送', '未发送')}`
  ].join('\n')
}

function formatAmount(value) {
  const amount = Number(value)
  return Number.isFinite(amount) ? amount.toFixed(2) : '--'
}

function payStatusText(status) {
  const value = (status || '').toLowerCase()
  if (value === 'unpaid') return '未支付'
  if (value === 'pending') return '待支付'
  if (value === 'paid') return '已支付'
  if (value === 'failed') return '支付失败'
  return status || '-'
}

function orderStatusText(status) {
  const value = (status || '').toLowerCase()
  if (value === 'created') return '待处理'
  if (value === 'accepted') return '已接单'
  if (value === 'assigned') return '派件中'
  if (value === 'completed') return '已完成'
  if (value === 'canceled' || value === 'cancelled') return '已取消'
  return status || '-'
}

async function load() {
  loading.value = true
  try {
    const data = await getAdminSendOrders(status.value)
    rows.value = Array.isArray(data) ? data : data?.list || []
  } finally {
    loading.value = false
  }
}

async function process(row, action) {
  if (action === 'complete') {
    try {
      await ElMessageBox.confirm(
        '完成后将自动为收件侧尝试入库并生成取件码（如未入库，位置默认使用收件地址），是否继续？',
        '确认完成订单',
        {
          confirmButtonText: '继续',
          cancelButtonText: '取消',
          type: 'warning'
        }
      )
    } catch {
      return
    }
  }

  processing.value = true
  try {
    const data = await processSendOrder({ order_no: getOrderNo(row), action })
    if (action === 'complete') {
      const completeResult = await enrichCompleteResult(data, row)
      ElMessage.success('订单已完成')
      await ElMessageBox.alert(buildCompleteSummary(completeResult), '完成结果', {
        confirmButtonText: '知道了'
      })
    } else {
      ElMessage.success('订单处理成功')
    }
    await load()
  } finally {
    processing.value = false
  }
}

function openAssign(row) {
  assignForm.order_no = getOrderNo(row)
  assignForm.courier_name = ''
  assignVisible.value = true
}

async function assign() {
  if (!assignForm.courier_name) {
    ElMessage.warning('请输入快递员姓名')
    return
  }
  processing.value = true
  try {
    await processSendOrder({
      order_no: assignForm.order_no,
      action: 'assign_pickup',
      courier_name: assignForm.courier_name
    })
    ElMessage.success('已分配快递员')
    assignVisible.value = false
    await load()
  } finally {
    processing.value = false
  }
}

onMounted(load)
</script>
