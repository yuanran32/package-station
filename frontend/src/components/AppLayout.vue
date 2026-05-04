<template>
  <el-container class="layout">
    <el-aside class="layout-aside" width="248px">
      <div class="brand">
        <div class="brand-mark">驿</div>
        <div>
          <strong>菜鸟驿站</strong>
          <span>快递管理系统</span>
        </div>
      </div>
      <SideMenu :role="auth.role" />
    </el-aside>

    <el-container>
      <el-header class="layout-header">
        <div>
          <h1>{{ route.meta.title || '工作台' }}</h1>
          <span>{{ auth.isAdmin ? '管理员后台' : '用户服务台' }}</span>
        </div>
        <div class="header-actions">
          <el-tag :type="auth.isAdmin ? 'danger' : 'success'" effect="light">
            {{ auth.isAdmin ? '管理员' : '普通用户' }}
          </el-tag>
          <span class="user-name">{{ auth.user?.name || auth.user?.username || '未命名用户' }}</span>
          <el-button :icon="SwitchButton" @click="handleLogout">退出</el-button>
        </div>
      </el-header>
      <el-main class="layout-main">
        <RouterView />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { useRoute, useRouter } from 'vue-router'
import { SwitchButton } from '@element-plus/icons-vue'
import { useAuthStore } from '../stores/auth'
import SideMenu from './SideMenu.vue'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()

function handleLogout() {
  auth.logout()
  router.replace('/login')
}
</script>

<style scoped>
.layout {
  min-height: 100vh;
}

.layout-aside {
  position: sticky;
  top: 0;
  height: 100vh;
  border-right: 1px solid var(--app-border);
  background: #ffffff;
}

.brand {
  display: flex;
  align-items: center;
  gap: 12px;
  height: 72px;
  padding: 0 18px;
  border-bottom: 1px solid var(--app-border);
}

.brand-mark {
  display: grid;
  width: 40px;
  height: 40px;
  place-items: center;
  border-radius: 8px;
  background: var(--app-primary);
  color: #ffffff;
  font-weight: 800;
}

.brand strong,
.brand span {
  display: block;
}

.brand span {
  margin-top: 4px;
  color: var(--app-muted);
  font-size: 12px;
}

.layout-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  height: 72px;
  padding: 0 24px;
  border-bottom: 1px solid var(--app-border);
  background: rgba(255, 255, 255, 0.88);
  backdrop-filter: blur(12px);
}

.layout-header h1 {
  margin: 0;
  font-size: 20px;
  color: #111827;
}

.layout-header span {
  color: var(--app-muted);
  font-size: 13px;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.user-name {
  max-width: 160px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.layout-main {
  padding: 24px;
}

@media (max-width: 920px) {
  .layout {
    display: block;
  }

  .layout-aside {
    position: static;
    width: 100% !important;
    height: auto;
  }

  .layout-header {
    align-items: flex-start;
    height: auto;
    padding: 16px;
    flex-direction: column;
  }

  .layout-main {
    padding: 16px;
  }
}
</style>
