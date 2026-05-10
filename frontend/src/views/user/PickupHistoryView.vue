<template>
  <div class="page">
    <div class="page-header">
      <div>
        <h2 class="page-title">取件历史</h2>
        <p class="page-desc">展示账号关联的取件记录。</p>
      </div>
      <el-button type="primary" :icon="Refresh" :loading="loading" @click="load">刷新</el-button>
    </div>
    <section class="section-panel">
      <el-table v-loading="loading" :data="rows" border stripe empty-text="暂无取件记录">
        <el-table-column prop="tracking_no_text" label="快递单号" min-width="160" show-overflow-tooltip />
        <el-table-column prop="pickup_code_text" label="存储码/取件码" min-width="140" show-overflow-tooltip />
        <el-table-column prop="location_text" label="取件位置" min-width="150" show-overflow-tooltip />
        <el-table-column prop="pickup_time_text" label="取件时间" min-width="180" show-overflow-tooltip />
        <el-table-column prop="status_text" label="状态" min-width="100" />
      </el-table>
    </section>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { Refresh } from '@element-plus/icons-vue'
import { getPickupHistory } from '../../api/user'

const loading = ref(false)
const rows = ref([])

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

function normalizeKey(key) {
  return String(key || '')
    .toLowerCase()
    .replace(/[^a-z0-9]/g, '')
}

function flattenEntries(source, prefix = '') {
  if (!source || typeof source !== 'object') {
    return []
  }
  const entries = []
  for (const [key, value] of Object.entries(source)) {
    const path = prefix ? `${prefix}.${key}` : key
    if (value && typeof value === 'object' && !Array.isArray(value)) {
      entries.push(...flattenEntries(value, path))
      continue
    }
    entries.push([path, value])
  }
  return entries
}

function pickByKeywords(source, includeKeywords, excludeKeywords = []) {
  const includes = includeKeywords.map(normalizeKey).filter(Boolean)
  const excludes = excludeKeywords.map(normalizeKey).filter(Boolean)
  if (!includes.length) {
    return ''
  }
  const entries = flattenEntries(source)
  for (const [path, value] of entries) {
    if (!hasValue(value)) {
      continue
    }
    const key = normalizeKey(path)
    const hitInclude = includes.every((token) => key.includes(token))
    const hitExclude = excludes.some((token) => token && key.includes(token))
    if (hitInclude && !hitExclude) {
      return value
    }
  }
  return ''
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
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())} ${pad(date.getHours())}:${pad(
    date.getMinutes()
  )}:${pad(date.getSeconds())}`
}

function statusText(status) {
  const value = String(status ?? '').toLowerCase()
  if (value === '0') return '待取件'
  if (value === '1') return '已取件'
  if (value === '2') return '已派送'
  if (value === '3') return '异常'
  if (value === 'inbound' || value === 'stored') return '已入库'
  if (value === 'pending') return '待取件'
  if (value === 'picked' || value === 'picked_up' || value === 'outbound') return '已取件'
  if (value === 'delivered') return '已派送'
  if (value === 'lost') return '异常'
  return status || '-'
}

function normalizeRows(data) {
  let list = []
  if (Array.isArray(data)) {
    list = data
  } else {
    list = data?.list || data?.records || data?.items || data?.rows || []
    if (!Array.isArray(list) && Array.isArray(data?.data)) {
      list = data.data
    }
  }

  return list.map((item) => {
    const trackingNo =
      pickValue(item, [
        'tracking_no',
        'trackingNo',
        'tracking_number',
        'trackingNumber',
        'track_no',
        'trackNo',
        'express_no',
        'expressNo',
        'express_number',
        'expressNumber',
        'waybill_no',
        'waybillNo',
        'parcel.tracking_no',
        'parcel.trackingNo'
      ]) ||
      pickByKeywords(item, ['tracking']) ||
      pickByKeywords(item, ['waybill']) ||
      pickByKeywords(item, ['express', 'no'])
    const pickupCode =
      pickValue(item, [
        'pickup_code',
        'pickupCode',
        'storage_code',
        'storageCode',
        'fetch_code',
        'fetchCode',
        'identity_code',
        'identityCode',
        'pickup_no',
        'pickupNo',
        'extract_code',
        'extractCode',
        'code',
        'parcel.pickup_code',
        'parcel.pickupCode'
      ]) ||
      pickByKeywords(item, ['pickup', 'code']) ||
      pickByKeywords(item, ['storage', 'code']) ||
      pickByKeywords(item, ['identity', 'code'])
    const location =
      pickValue(item, [
        'location',
        'pickup_location',
        'pickupLocation',
        'storage_location',
        'storageLocation',
        'pickup_address',
        'pickupAddress',
        'storage_address',
        'storageAddress',
        'locker_no',
        'lockerNo',
        'cabinet_no',
        'cabinetNo',
        'cabinet_name',
        'cabinetName',
        'slot',
        'shelf',
        'position',
        'parcel.location',
        'parcel.storage_location',
        'parcel.storageLocation'
      ]) ||
      pickByKeywords(item, ['location']) ||
      pickByKeywords(item, ['locker']) ||
      pickByKeywords(item, ['cabinet']) ||
      pickByKeywords(item, ['position'])
    const pickupTime =
      pickValue(item, [
        'pickup_time',
        'pickupTime',
        'outbound_time',
        'outboundTime',
        'finished_at',
        'finishedAt',
        'picked_at',
        'pickedAt',
        'created_at',
        'createdAt',
        'updated_at',
        'updatedAt'
      ]) ||
      pickByKeywords(item, ['pickup', 'time']) ||
      pickByKeywords(item, ['picked', 'at']) ||
      pickByKeywords(item, ['outbound', 'time']) ||
      pickByKeywords(item, ['created', 'at'])
    const status =
      pickValue(item, [
        'status',
        'pickup_status',
        'pickupStatus',
        'state',
        'parcel.status'
      ]) ||
      pickByKeywords(item, ['status']) ||
      pickByKeywords(item, ['state'])

    return {
      ...item,
      tracking_no_text: trackingNo || '-',
      pickup_code_text: pickupCode || '-',
      location_text: location || '-',
      pickup_time_text: formatTime(pickupTime),
      status_text: statusText(status)
    }
  })
}

async function load() {
  loading.value = true
  try {
    const data = await getPickupHistory()
    rows.value = normalizeRows(data)
  } finally {
    loading.value = false
  }
}

onMounted(load)
</script>
