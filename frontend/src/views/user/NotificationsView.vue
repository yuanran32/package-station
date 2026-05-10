<template>
  <div class="page">
    <div class="page-header">
      <div>
        <h2 class="page-title">通知提醒</h2>
        <p class="page-desc">通过 WebSocket 接收驿站实时通知。</p>
      </div>
      <el-button :type="connected ? 'success' : 'primary'" :loading="connecting" @click="connect">
        {{ connected ? '已连接' : '连接通知' }}
      </el-button>
    </div>
    <section class="section-panel">
      <el-timeline>
        <el-timeline-item v-for="item in notices" :key="item.id" :timestamp="item.time">
          {{ item.content }}
        </el-timeline-item>
      </el-timeline>
      <el-empty v-if="!notices.length" description="暂无通知" />
    </section>
  </div>
</template>

<script setup>
import { onBeforeUnmount, ref } from 'vue'
import { storeToRefs } from 'pinia'
import { ElMessage } from 'element-plus'
import { createNotifySocket } from '../../api/notice'
import { getPickupHistory } from '../../api/user'
import { useAuthStore } from '../../stores/auth'
import { useNoticeStore } from '../../stores/notices'

const auth = useAuthStore()
const noticeStore = useNoticeStore()
const { notices } = storeToRefs(noticeStore)
const connected = ref(false)
const connecting = ref(false)
let socket = null
const pickupHistoryCache = {
  loadedAt: 0,
  list: []
}

function addNotice(content, time = new Date().toLocaleString()) {
  noticeStore.addNotice(content, time)
}

function hasValue(value) {
  return value !== undefined && value !== null && value !== ''
}

function hasChinese(text) {
  return /[\u4e00-\u9fa5]/.test(String(text || ''))
}

function pickValue(source, keys) {
  for (const key of keys) {
    const value = source?.[key]
    if (hasValue(value)) {
      return value
    }
  }
  return ''
}

function formatTime(value) {
  if (!hasValue(value)) {
    return new Date().toLocaleString()
  }
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) {
    return String(value)
  }
  return date.toLocaleString()
}

function normalizePlainMessage(text) {
  const raw = String(text || '').trim()
  if (!raw) return ''
  if (hasChinese(raw)) return raw

  const lower = raw.toLowerCase()
  if (lower === 'ws notify connected' || lower === 'connected') return '通知通道已连接'
  if (lower.includes('parcel') && lower.includes('arriv')) return '到站通知：您的包裹已到站，请及时取件。'
  if (lower.includes('pickup') && lower.includes('remind')) return '取件提醒：您有包裹待取件。'
  if (lower.includes('notify') || lower.includes('notice') || lower.includes('notification')) {
    return '系统通知：您有一条新通知'
  }
  return raw
}

function extractByRegex(text, patterns) {
  const source = String(text || '')
  for (const pattern of patterns) {
    const match = source.match(pattern)
    if (match && hasValue(match[1])) {
      return match[1].trim()
    }
  }
  return ''
}

function normalizeList(data) {
  if (Array.isArray(data)) return data
  if (Array.isArray(data?.list)) return data.list
  if (Array.isArray(data?.records)) return data.records
  if (Array.isArray(data?.items)) return data.items
  if (Array.isArray(data?.rows)) return data.rows
  if (Array.isArray(data?.data)) return data.data
  return []
}

function normalizeCode(value) {
  return String(value || '').trim().toLowerCase()
}

async function loadPickupHistory() {
  const now = Date.now()
  const cacheValidMs = 30 * 1000
  if (pickupHistoryCache.list.length && now - pickupHistoryCache.loadedAt < cacheValidMs) {
    return pickupHistoryCache.list
  }
  const data = await getPickupHistory()
  const list = normalizeList(data)
  pickupHistoryCache.list = list
  pickupHistoryCache.loadedAt = now
  return list
}

async function resolveTrackingNoByPickupCode(pickupCode) {
  if (!hasValue(pickupCode)) return ''
  const target = normalizeCode(pickupCode)
  if (!target) return ''

  try {
    const list = await loadPickupHistory()
    for (const item of list) {
      const itemPickupCode = pickValue(item, [
        'pickup_code',
        'pickupCode',
        'storage_code',
        'storageCode',
        'code',
        'parcel.pickup_code',
        'parcel.pickupCode'
      ])
      if (normalizeCode(itemPickupCode) !== target) continue

      const trackingNo = pickValue(item, [
        'tracking_no',
        'trackingNo',
        'tracking_number',
        'trackingNumber',
        'waybill_no',
        'waybillNo',
        'parcel.tracking_no',
        'parcel.trackingNo'
      ])
      if (hasValue(trackingNo)) return String(trackingNo).trim()
    }
  } catch {
    // 解析通知失败不影响页面主流程
  }

  return ''
}

function buildNoticeContent(payload) {
  const type = String(pickValue(payload, ['type', 'notice_type', 'noticeType'])).toLowerCase()
  const messageRaw = pickValue(payload, ['message', 'content', 'msg', 'text'])
  const message = normalizePlainMessage(messageRaw)

  if (type === 'connected' || String(messageRaw).toLowerCase() === 'ws notify connected') {
    return { ignore: true }
  }

  const trackingNo = pickValue(payload, [
    'tracking_no',
    'trackingNo',
    'order_no',
    'orderNo',
    'waybill_no',
    'waybillNo'
  ]) || extractByRegex(messageRaw, [/单号[是为:\s]*([A-Za-z0-9-]+)/, /tracking[_\s-]*no[=: \s]*([A-Za-z0-9-]+)/i])
  const pickupCode =
    pickValue(payload, ['pickup_code', 'pickupCode', 'storage_code', 'storageCode', 'code']) ||
    extractByRegex(messageRaw, [/取件码[是为:\s]*([A-Za-z0-9-]+)/, /pickup[_\s-]*code[=: \s]*([A-Za-z0-9-]+)/i])
  const location = pickValue(payload, ['location', 'pickup_location', 'pickupLocation', 'receiver_address', 'receiverAddress'])

  const messageRawLower = String(messageRaw || '').toLowerCase()
  const arrivalLikeByType =
    type === 'station_arrival' || type === 'arrival' || type === 'pickup_ready' || type === 'parcel_arrived'
  const arrivalLikeByText =
    String(message).includes('到站') ||
    messageRawLower.includes('arriv') ||
    (messageRawLower.includes('pickup') && (messageRawLower.includes('ready') || messageRawLower.includes('code')))
  const isArrivalLike = arrivalLikeByType || arrivalLikeByText

  let text = ''
  let appendExtras = true
  if (isArrivalLike) {
    const hasTrackingNo = hasValue(trackingNo)
    const hasPickupCode = hasValue(pickupCode)
    if (hasTrackingNo && hasPickupCode) {
      text = `你的包裹单号是：${trackingNo}，取件码是：${pickupCode}，请及时取件`
    } else if (hasTrackingNo) {
      text = `你的包裹单号是：${trackingNo}，请及时取件`
    } else if (hasPickupCode) {
      text = `你的包裹取件码是：${pickupCode}，请及时取件`
    } else {
      text = message || '到站通知：您的包裹已到站，请及时取件。'
      appendExtras = false
    }
    appendExtras = false
  } else if (type === 'pickup_reminder' || type === 'pickup') {
    text = message || '取件提醒：您有包裹待取件。'
  } else if (type === 'notice' || type === 'system' || type === 'broadcast') {
    text = `系统通知：${message || '您有一条新通知'}`
  } else {
    text = message || '收到一条新通知'
  }

  const extras = []
  if (hasValue(trackingNo)) extras.push(`单号：${trackingNo}`)
  if (hasValue(pickupCode)) extras.push(`取件码：${pickupCode}`)
  if (hasValue(location)) extras.push(`位置：${location}`)
  if (appendExtras && extras.length) {
    text += `（${extras.join('，')}）`
  }

  return {
    ignore: false,
    time: formatTime(pickValue(payload, ['timestamp', 'time', 'created_at', 'createdAt'])),
    content: text,
    meta: {
      isArrivalLike,
      trackingNo: hasValue(trackingNo) ? String(trackingNo).trim() : '',
      pickupCode: hasValue(pickupCode) ? String(pickupCode).trim() : ''
    }
  }
}

async function normalizeMessageData(rawData) {
  let data = rawData
  if (typeof Blob !== 'undefined' && data instanceof Blob) {
    data = await data.text()
  }

  if (typeof data === 'string') {
    const text = data.trim()
    if (!text) {
      return { ignore: true }
    }
    try {
      const payload = JSON.parse(text)
      if (payload && typeof payload === 'object') {
        const notice = buildNoticeContent(payload)
        if (
          !notice.ignore &&
          notice.meta?.isArrivalLike &&
          !hasValue(notice.meta?.trackingNo) &&
          hasValue(notice.meta?.pickupCode)
        ) {
          const resolvedTrackingNo = await resolveTrackingNoByPickupCode(notice.meta.pickupCode)
          if (hasValue(resolvedTrackingNo)) {
            notice.content = `你的包裹单号是：${resolvedTrackingNo}，取件码是：${notice.meta.pickupCode}，请及时取件`
            notice.meta.trackingNo = resolvedTrackingNo
          }
        }
        return notice
      }
    } catch {
      if (text.toLowerCase() === 'ws notify connected') {
        return { ignore: true }
      }
      return {
        ignore: false,
        time: new Date().toLocaleString(),
        content: normalizePlainMessage(text)
      }
    }
  }

  if (data && typeof data === 'object') {
    const notice = buildNoticeContent(data)
    if (!notice.ignore && notice.meta?.isArrivalLike && !hasValue(notice.meta?.trackingNo) && hasValue(notice.meta?.pickupCode)) {
      const resolvedTrackingNo = await resolveTrackingNoByPickupCode(notice.meta.pickupCode)
      if (hasValue(resolvedTrackingNo)) {
        notice.content = `你的包裹单号是：${resolvedTrackingNo}，取件码是：${notice.meta.pickupCode}，请及时取件`
        notice.meta.trackingNo = resolvedTrackingNo
      }
    }
    return notice
  }

  return {
    ignore: false,
    time: new Date().toLocaleString(),
    content: String(data)
  }
}

function connect() {
  if (connected.value || connecting.value) return
  connecting.value = true
  socket = createNotifySocket(auth.token)
  socket.onopen = () => {
    connecting.value = false
    connected.value = true
    ElMessage.success('通知通道已连接')
    addNotice('通知通道已连接，正在等待驿站推送消息。')
  }
  socket.onmessage = async (event) => {
    try {
      const notice = await normalizeMessageData(event.data)
      if (notice.ignore) return
      addNotice(notice.content, notice.time)
    } catch {
      addNotice('收到一条通知，但消息格式无法解析。')
    }
  }
  socket.onerror = () => {
    connecting.value = false
    ElMessage.error('通知通道连接失败')
    addNotice('通知通道连接失败，请检查后端服务与地址配置。')
  }
  socket.onclose = () => {
    connecting.value = false
    connected.value = false
    addNotice('通知通道已断开。')
  }
}

onBeforeUnmount(() => {
  socket?.close()
})
</script>
