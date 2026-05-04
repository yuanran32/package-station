<template>
  <div class="page">
    <div class="page-header">
      <div>
        <h2 class="page-title">个人身份码</h2>
        <p class="page-desc">用于驿站核验身份和辅助取件。</p>
      </div>
      <el-button type="primary" :loading="loading" @click="load">获取身份码</el-button>
    </div>
    <section class="section-panel identity-wrap">
      <div class="fake-qr" :style="{ '--seed': code.length || 1 }">
        <span v-for="n in 121" :key="n" :class="{ dark: isDark(n) }" />
      </div>
      <div>
        <h3 class="section-title">身份码</h3>
        <el-text class="identity-code">{{ code || '点击获取身份码' }}</el-text>
      </div>
    </section>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { getIdentityCode } from '../../api/user'

const loading = ref(false)
const code = ref('')

function isDark(n) {
  if (!code.value) return n % 5 === 0 || n % 7 === 0
  const sum = code.value.split('').reduce((acc, char) => acc + char.charCodeAt(0), 0)
  return (n * 13 + sum) % 4 === 0 || n <= 3 || n >= 119
}

async function load() {
  loading.value = true
  try {
    const data = await getIdentityCode()
    code.value = data?.identity_code || data?.identityCode || data?.code || String(data || '')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.identity-wrap {
  display: flex;
  align-items: center;
  gap: 24px;
  flex-wrap: wrap;
}

.fake-qr {
  display: grid;
  grid-template-columns: repeat(11, 1fr);
  gap: 4px;
  width: 220px;
  height: 220px;
  padding: 14px;
  border: 1px solid var(--app-border);
  border-radius: 8px;
  background: #ffffff;
}

.fake-qr span {
  border-radius: 2px;
  background: #eef2f7;
}

.fake-qr .dark {
  background: #111827;
}

.identity-code {
  display: block;
  max-width: 520px;
  padding: 14px 16px;
  border-radius: 8px;
  background: #f3f7f8;
  color: var(--app-primary);
  font-size: 22px;
  font-weight: 800;
  overflow-wrap: anywhere;
}
</style>
