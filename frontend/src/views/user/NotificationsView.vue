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
import { ElMessage } from 'element-plus'
import { createNotifySocket } from '../../api/notice'
import { useAuthStore } from '../../stores/auth'

const auth = useAuthStore()
const connected = ref(false)
const connecting = ref(false)
const notices = ref([])
let socket = null

function addNotice(content) {
  notices.value.unshift({
    id: Date.now() + Math.random(),
    time: new Date().toLocaleString(),
    content
  })
}

function connect() {
  if (connected.value || connecting.value) return
  connecting.value = true
  socket = createNotifySocket(auth.token)
  socket.onopen = () => {
    connecting.value = false
    connected.value = true
    ElMessage.success('通知通道已连接')
  }
  socket.onmessage = (event) => {
    addNotice(event.data)
  }
  socket.onerror = () => {
    connecting.value = false
    ElMessage.error('通知通道连接失败')
  }
  socket.onclose = () => {
    connecting.value = false
    connected.value = false
  }
}

onBeforeUnmount(() => {
  socket?.close()
})
</script>
