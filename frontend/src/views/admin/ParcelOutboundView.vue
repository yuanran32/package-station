<template>
  <div class="page">
    <div class="page-header">
      <div>
        <h2 class="page-title">出库处理</h2>
        <p class="page-desc">按快递单号和取件码核销出库。</p>
      </div>
    </div>
    <section class="section-panel form-narrow">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
        <el-form-item label="快递单号" prop="tracking_no">
          <el-input v-model.trim="form.tracking_no" />
        </el-form-item>
        <el-form-item label="取件码" prop="pickup_code">
          <el-input v-model.trim="form.pickup_code" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="loading" @click="submit">确认出库</el-button>
        </el-form-item>
      </el-form>
    </section>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { outboundParcel } from '../../api/parcel'

const formRef = ref()
const loading = ref(false)
const form = reactive({
  tracking_no: '',
  pickup_code: ''
})

const rules = {
  tracking_no: [{ required: true, message: '请输入快递单号', trigger: 'blur' }],
  pickup_code: [{ required: true, message: '请输入取件码', trigger: 'blur' }]
}

async function submit() {
  await formRef.value.validate()
  loading.value = true
  try {
    await outboundParcel(form)
    ElMessage.success('出库处理完成')
    formRef.value.resetFields()
  } finally {
    loading.value = false
  }
}
</script>
