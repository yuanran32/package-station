<template>
  <div class="page">
    <div class="page-header">
      <div>
        <h2 class="page-title">管理首页</h2>
        <p class="page-desc">处理入库、出库、寄件订单和通知发送。</p>
      </div>
    </div>
    <div class="metric-grid">
      <div v-for="item in metrics" :key="item.label" class="metric-card">
        <div class="metric-label">{{ item.label }}</div>
        <div class="metric-value">{{ item.value }}</div>
      </div>
    </div>
    <section class="section-panel">
      <h3 class="section-title">管理入口</h3>
      <div class="quick-grid">
        <RouterLink v-for="item in quickLinks" :key="item.path" :to="item.path" class="quick-item">
          <el-icon><component :is="item.icon" /></el-icon>
          <span>{{ item.label }}</span>
        </RouterLink>
      </div>
    </section>
  </div>
</template>

<script setup>
import { Box, Collection, Document, Message, Van } from '@element-plus/icons-vue'

const metrics = [
  { label: '快递入库', value: '可处理' },
  { label: '出库核销', value: '可处理' },
  { label: '寄件订单', value: '待维护' },
  { label: '通知发送', value: '在线' }
]

const quickLinks = [
  { path: '/app/admin/inbound', label: '快递入库', icon: Box },
  { path: '/app/admin/parcels', label: '快递列表', icon: Document },
  { path: '/app/admin/pickup', label: '取派件记录', icon: Collection },
  { path: '/app/admin/send-orders', label: '寄件订单', icon: Van },
  { path: '/app/admin/notices', label: '通知发送', icon: Message }
]
</script>

<style scoped>
.quick-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(170px, 1fr));
  gap: 12px;
}

.quick-item {
  display: flex;
  align-items: center;
  gap: 10px;
  min-height: 64px;
  padding: 14px;
  border: 1px solid var(--app-border);
  border-radius: 8px;
  background: #ffffff;
  font-weight: 700;
}

.quick-item .el-icon {
  color: var(--app-accent);
  font-size: 22px;
}
</style>
