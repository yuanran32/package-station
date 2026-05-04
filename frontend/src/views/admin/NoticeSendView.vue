<template>
  <div class="page">
    <div class="page-header">
      <div>
        <h2 class="page-title">通知发送</h2>
        <p class="page-desc">向指定用户或全体用户发送提醒消息。</p>
      </div>
    </div>
    <section class="section-panel form-narrow">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
        <el-form-item label="目标用户">
          <el-input v-model.trim="form.target_user" placeholder="可为空，表示全体用户" />
        </el-form-item>
        <el-form-item label="通知标题" prop="title">
          <el-input v-model.trim="form.title" />
        </el-form-item>
        <el-form-item label="通知内容" prop="content">
          <el-input v-model.trim="form.content" type="textarea" :rows="5" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="loading" @click="submit">发送通知</el-button>
        </el-form-item>
      </el-form>
    </section>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { sendNotice } from '../../api/notice'

const formRef = ref()
const loading = ref(false)
const form = reactive({
  target_user: '',
  title: '',
  content: ''
})

const rules = {
  title: [{ required: true, message: '请输入标题', trigger: 'blur' }],
  content: [{ required: true, message: '请输入内容', trigger: 'blur' }]
}

async function submit() {
  await formRef.value.validate()
  loading.value = true
  try {
    await sendNotice(form)
    ElMessage.success('通知已发送')
    formRef.value.resetFields()
  } finally {
    loading.value = false
  }
}
</script>
