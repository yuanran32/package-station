<template>
  <div class="page">
    <div class="page-header">
      <div>
        <h2 class="page-title">寄件下单</h2>
        <p class="page-desc">填写寄件和收件信息，提交后由管理员处理。</p>
      </div>
    </div>
    <section class="section-panel form-narrow">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
        <el-form-item label="寄件人" prop="sender_name">
          <el-input v-model.trim="form.sender_name" />
        </el-form-item>
        <el-form-item label="寄件电话" prop="sender_phone">
          <el-input v-model.trim="form.sender_phone" />
        </el-form-item>
        <el-form-item label="寄件地址" prop="sender_address">
          <el-input v-model.trim="form.sender_address" />
        </el-form-item>
        <el-form-item label="收件人" prop="receiver_name">
          <el-input v-model.trim="form.receiver_name" />
        </el-form-item>
        <el-form-item label="收件电话" prop="receiver_phone">
          <el-input v-model.trim="form.receiver_phone" />
        </el-form-item>
        <el-form-item label="收件地址" prop="receiver_address">
          <el-input v-model.trim="form.receiver_address" />
        </el-form-item>
        <el-form-item label="物品描述">
          <el-input v-model.trim="form.item_desc" type="textarea" :rows="3" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="loading" @click="submit">提交订单</el-button>
        </el-form-item>
      </el-form>
    </section>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { createSendOrder } from '../../api/send'

const formRef = ref()
const loading = ref(false)
const form = reactive({
  sender_name: '',
  sender_phone: '',
  sender_address: '',
  receiver_name: '',
  receiver_phone: '',
  receiver_address: '',
  item_desc: ''
})

const rules = {
  sender_name: [{ required: true, message: '请输入寄件人', trigger: 'blur' }],
  sender_phone: [{ required: true, message: '请输入寄件电话', trigger: 'blur' }],
  sender_address: [{ required: true, message: '请输入寄件地址', trigger: 'blur' }],
  receiver_name: [{ required: true, message: '请输入收件人', trigger: 'blur' }],
  receiver_phone: [{ required: true, message: '请输入收件电话', trigger: 'blur' }],
  receiver_address: [{ required: true, message: '请输入收件地址', trigger: 'blur' }]
}

async function submit() {
  await formRef.value.validate()
  loading.value = true
  try {
    await createSendOrder(form)
    ElMessage.success('寄件订单已提交')
    formRef.value.resetFields()
    form.item_desc = ''
  } finally {
    loading.value = false
  }
}
</script>
