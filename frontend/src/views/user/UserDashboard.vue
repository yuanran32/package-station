<template>
  <div class="page">
    <div class="page-header">
      <div>
        <h2 class="page-title">用户首页</h2>
        <p class="page-desc">集中处理取件查询、寄件下单、礼券和账单。</p>
      </div>
    </div>
    <div class="metric-grid">
      <div v-for="item in metrics" :key="item.label" class="metric-card">
        <div class="metric-label">{{ item.label }}</div>
        <div class="metric-value">{{ item.value }}</div>
      </div>
    </div>
    <section class="section-panel">
      <h3 class="section-title">快捷操作</h3>
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
import { Box, CreditCard, Location, Tickets, Van } from '@element-plus/icons-vue'

const metrics = [
  { label: '待查询快递', value: '实时' },
  { label: '身份码状态', value: '可用' },
  { label: '寄件服务', value: '在线' },
  { label: '通知通道', value: 'WebSocket' }
]

const quickLinks = [
  { path: '/app/user/parcel', label: '查询快递', icon: Location },
  { path: '/app/user/history', label: '取件历史', icon: Box },
  { path: '/app/user/send-order', label: '寄件下单', icon: Van },
  { path: '/app/user/coupons', label: '我的礼券', icon: Tickets },
  { path: '/app/user/bills', label: '支付账单', icon: CreditCard }
]
</script>

<style scoped>
.quick-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(160px, 1fr));
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
  color: #111827;
  font-weight: 700;
}

.quick-item .el-icon {
  color: var(--app-primary);
  font-size: 22px;
}
</style>
