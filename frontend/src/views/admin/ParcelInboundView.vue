<template>
  <div class="page">
    <div class="page-header">
      <div>
        <h2 class="page-title">快递入库</h2>
        <p class="page-desc">录入快递单号、货架位置和收件人手机号。</p>
      </div>
    </div>
    <section class="section-panel form-narrow">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="120px">
        <el-form-item label="快递单号" prop="tracking_no">
          <el-input v-model.trim="form.tracking_no" placeholder="SF1234567890" />
        </el-form-item>
        <el-form-item label="存储位置" prop="location">
          <el-input v-model.trim="form.location" placeholder="A-01-03" />
        </el-form-item>
        <el-form-item label="收件人手机号" prop="recipient_phone">
          <el-input v-model.trim="form.recipient_phone" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="loading" @click="submit">确认入库</el-button>
        </el-form-item>
      </el-form>
    </section>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { inboundParcel } from '../../api/parcel'

const formRef = ref()
const loading = ref(false)
const form = reactive({
  tracking_no: '',
  location: '',
  recipient_phone: ''
})

const rules = {
  tracking_no: [{ required: true, message: '请输入快递单号', trigger: 'blur' }],
  location: [{ required: true, message: '请输入存储位置', trigger: 'blur' }],
  recipient_phone: [{ required: true, message: '请输入收件人手机号', trigger: 'blur' }]
}

async function submit() {
  await formRef.value.validate()
  loading.value = true
  try {
    await inboundParcel(form)
    ElMessage.success('快递已入库')
    formRef.value.resetFields()
  } finally {
    loading.value = false
  }
}
</script>
