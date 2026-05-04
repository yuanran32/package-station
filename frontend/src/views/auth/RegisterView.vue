<template>
  <div class="register-page">
    <section class="register-box">
      <h1>创建用户账号</h1>
      <p>注册后可查询快递、提交寄件订单、查看身份码和账单。</p>
      <el-form ref="formRef" :model="form" :rules="rules" label-position="top">
        <el-form-item label="用户名" prop="username">
          <el-input v-model.trim="form.username" placeholder="例如 zhangsan" size="large" />
        </el-form-item>
        <el-form-item label="姓名" prop="name">
          <el-input v-model.trim="form.name" placeholder="请输入姓名" size="large" />
        </el-form-item>
        <el-form-item label="手机号" prop="phone">
          <el-input v-model.trim="form.phone" placeholder="请输入手机号" size="large" />
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input v-model="form.password" type="password" show-password placeholder="不少于 6 位" size="large" />
        </el-form-item>
        <div class="actions">
          <el-button size="large" @click="router.push('/login')">返回登录</el-button>
          <el-button type="primary" size="large" :loading="loading" @click="submit">注册</el-button>
        </div>
      </el-form>
    </section>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { registerUser } from '../../api/auth'

const router = useRouter()
const formRef = ref()
const loading = ref(false)
const form = reactive({
  username: '',
  password: '',
  name: '',
  phone: ''
})

const rules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  name: [{ required: true, message: '请输入姓名', trigger: 'blur' }],
  phone: [{ required: true, message: '请输入手机号', trigger: 'blur' }],
  password: [{ required: true, min: 6, message: '密码至少 6 位', trigger: 'blur' }]
}

async function submit() {
  await formRef.value.validate()
  loading.value = true
  try {
    await registerUser(form)
    ElMessage.success('注册成功，请登录')
    router.replace('/login')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.register-page {
  display: grid;
  min-height: 100vh;
  place-items: center;
  padding: 32px;
  background:
    linear-gradient(180deg, rgba(244, 247, 251, 0.86), rgba(244, 247, 251, 0.96)),
    url("https://images.unsplash.com/photo-1578575437130-527eed3abbec?auto=format&fit=crop&w=1400&q=80") center/cover;
}

.register-box {
  width: min(560px, 100%);
  padding: 30px;
  border: 1px solid var(--app-border);
  border-radius: 8px;
  background: #ffffff;
}

.register-box h1 {
  margin: 0;
  font-size: 28px;
}

.register-box p {
  margin: 10px 0 24px;
  color: var(--app-muted);
}

.actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style>
