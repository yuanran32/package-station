<template>
  <div class="page">
    <div class="page-header">
      <div>
        <h2 class="page-title">个人资料</h2>
        <p class="page-desc">查看并维护姓名、手机号等基础信息。</p>
      </div>
    </div>
    <section class="section-panel form-narrow">
      <el-form ref="formRef" :model="form" label-width="90px">
        <el-form-item label="用户名">
          <el-input v-model="form.username" disabled />
        </el-form-item>
        <el-form-item label="姓名">
          <el-input v-model.trim="form.name" />
        </el-form-item>
        <el-form-item label="手机号">
          <el-input v-model.trim="form.phone" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="loading" @click="save">保存资料</el-button>
          <el-button :loading="fetching" @click="load">刷新</el-button>
        </el-form-item>
      </el-form>
    </section>
  </div>
</template>

<script setup>
import { reactive, ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getProfile, updateProfile } from '../../api/user'
import { useAuthStore } from '../../stores/auth'

const auth = useAuthStore()
const loading = ref(false)
const fetching = ref(false)
const form = reactive({
  username: '',
  name: '',
  phone: ''
})

function fill(data = {}) {
  Object.assign(form, {
    username: data.username || '',
    name: data.name || '',
    phone: data.phone || ''
  })
}

async function load() {
  fetching.value = true
  try {
    const data = await getProfile()
    fill(data)
  } finally {
    fetching.value = false
  }
}

async function save() {
  loading.value = true
  try {
    await updateProfile({ name: form.name, phone: form.phone })
    await auth.refreshProfile()
    ElMessage.success('资料已保存')
  } finally {
    loading.value = false
  }
}

onMounted(() => fill(auth.user || {}))
</script>
