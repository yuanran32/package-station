<template>
  <div class="auth-page">
    <section class="auth-visual">
      <div class="visual-content">
        <h1>菜鸟驿站快递管理系统</h1>
        <p>面向用户取件、寄件、支付和管理员入库、出库、通知的一体化工作台。</p>
      </div>
    </section>
    <section class="auth-panel">
      <div class="auth-box">
        <h2>登录</h2>
        <p>管理员默认账号：admin / admin123</p>
        <el-form ref="formRef" :model="form" :rules="rules" label-position="top" @keyup.enter="submit">
          <el-form-item label="账号" prop="account">
            <el-input v-model.trim="form.account" placeholder="用户名或手机号" size="large" />
          </el-form-item>
          <el-form-item label="密码" prop="password">
            <el-input v-model="form.password" type="password" placeholder="请输入密码" show-password size="large" />
          </el-form-item>
          <el-button type="primary" size="large" :loading="loading" class="wide" @click="submit">登录</el-button>
        </el-form>
        <div class="auth-link">
          还没有账号？
          <RouterLink to="/register">去注册</RouterLink>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useAuthStore } from '../../stores/auth'

const router = useRouter()
const auth = useAuthStore()
const formRef = ref()
const loading = ref(false)
const form = reactive({
  account: '',
  password: ''
})

const rules = {
  account: [{ required: true, message: '请输入账号', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
}

async function submit() {
  await formRef.value.validate()
  loading.value = true
  try {
    const role = await auth.login(form)
    ElMessage.success('登录成功')
    router.replace(role === 'admin' ? '/app/admin/dashboard' : '/app/user/dashboard')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.auth-page {
  display: grid;
  grid-template-columns: minmax(360px, 1.1fr) minmax(360px, 0.9fr);
  min-height: 100vh;
}

.auth-visual {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 64px;
  color: #ffffff;
  background:
    linear-gradient(135deg, rgba(18, 109, 106, 0.92), rgba(33, 41, 54, 0.74)),
    url("https://images.unsplash.com/photo-1586528116311-ad8dd3c8310d?auto=format&fit=crop&w=1600&q=80") center/cover;
}

.visual-content {
  max-width: 650px;
  text-align: center;
}

.visual-content h1 {
  margin: 0;
  font-size: 44px;
  line-height: 1.18;
}

.visual-content p {
  max-width: 560px;
  margin: 18px 0 0;
  font-size: 18px;
  line-height: 1.8;
}

.auth-panel {
  display: grid;
  place-items: center;
  padding: 32px;
  background: #f4f7fb;
}

.auth-box {
  width: min(420px, 100%);
  padding: 28px;
  border: 1px solid var(--app-border);
  border-radius: 8px;
  background: #ffffff;
}

.auth-box h2 {
  margin: 0;
  font-size: 28px;
}

.auth-box p {
  margin: 10px 0 26px;
  color: var(--app-muted);
}

.wide {
  width: 100%;
}

.auth-link {
  margin-top: 18px;
  color: var(--app-muted);
  text-align: center;
}

.auth-link a {
  color: var(--app-primary);
  font-weight: 700;
}

@media (max-width: 860px) {
  .auth-page {
    grid-template-columns: 1fr;
  }

  .auth-visual {
    min-height: 320px;
    padding: 32px;
  }

  .visual-content h1 {
    font-size: 32px;
  }
}
</style>
