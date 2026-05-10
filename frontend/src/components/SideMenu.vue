<template>
  <el-menu :default-active="route.path" router class="side-menu">
    <template v-for="item in menus" :key="item.path">
      <el-menu-item :index="item.path">
        <el-icon><component :is="item.icon" /></el-icon>
        <span>{{ item.label }}</span>
      </el-menu-item>
    </template>
  </el-menu>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import {
  Box,
  Bell,
  Collection,
  CreditCard,
  Document,
  HomeFilled,
  Location,
  Message,
  Money,
  Postcard,
  Tickets,
  User,
  Van
} from '@element-plus/icons-vue'

const props = defineProps({
  role: {
    type: String,
    required: true
  }
})

const route = useRoute()

const userMenus = [
  { path: '/app/user/dashboard', label: '用户首页', icon: HomeFilled },
  { path: '/app/user/profile', label: '个人资料', icon: User },
  { path: '/app/user/parcel', label: '快递查询', icon: Location },
  { path: '/app/user/history', label: '取件历史', icon: Collection },
  { path: '/app/user/identity', label: '身份码', icon: Postcard },
  { path: '/app/user/send-order', label: '寄件下单', icon: Van },
  { path: '/app/user/pay-send', label: '寄件支付', icon: Money },
  { path: '/app/user/coupons', label: '红包礼券', icon: Tickets },
  { path: '/app/user/bills', label: '支付账单', icon: CreditCard },
  { path: '/app/user/notices', label: '通知提醒', icon: Bell }
]

const adminMenus = [
  { path: '/app/admin/dashboard', label: '管理首页', icon: HomeFilled },
  { path: '/app/admin/inbound', label: '快递入库', icon: Box },
  { path: '/app/admin/parcels', label: '快递列表', icon: Document },
  { path: '/app/admin/outbound', label: '出库处理', icon: Location },
  { path: '/app/admin/pickup', label: '取派件记录', icon: Collection },
  { path: '/app/admin/send-orders', label: '寄件订单', icon: Van },
  { path: '/app/admin/coupons', label: '红包礼券设置', icon: Tickets },
  { path: '/app/admin/notices', label: '通知发送', icon: Message }
]

const menus = computed(() => (props.role === 'admin' ? adminMenus : userMenus))
</script>

<style scoped>
.side-menu {
  padding: 12px;
}

.side-menu :deep(.el-menu-item) {
  height: 44px;
  margin-bottom: 6px;
  border-radius: 8px;
}

.side-menu :deep(.el-menu-item.is-active) {
  color: var(--app-primary);
  background: #e7f4f3;
}
</style>
