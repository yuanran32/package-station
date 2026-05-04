<template>
  <div class="page">
    <div class="page-header">
      <div>
        <h2 class="page-title">取派件记录</h2>
        <p class="page-desc">生成取件码，记录取件和派件信息。</p>
      </div>
    </div>
    <el-row :gutter="16">
      <el-col :xs="24" :lg="8">
        <section class="section-panel">
          <h3 class="section-title">生成取件码</h3>
          <el-form :model="codeForm" label-position="top">
            <el-form-item label="快递单号">
              <el-input v-model.trim="codeForm.tracking_no" />
            </el-form-item>
            <el-form-item label="收件人手机号">
              <el-input v-model.trim="codeForm.recipient_phone" />
            </el-form-item>
            <el-button type="primary" :loading="codeLoading" @click="createCode">生成</el-button>
          </el-form>
          <el-alert v-if="pickupCode" :title="`取件码：${pickupCode}`" type="success" show-icon class="result-alert" />
        </section>
      </el-col>
      <el-col :xs="24" :lg="8">
        <section class="section-panel">
          <h3 class="section-title">记录取件</h3>
          <el-form :model="pickupForm" label-position="top">
            <el-form-item label="快递单号">
              <el-input v-model.trim="pickupForm.tracking_no" />
            </el-form-item>
            <el-form-item label="取件码">
              <el-input v-model.trim="pickupForm.pickup_code" />
            </el-form-item>
            <el-form-item label="取件人">
              <el-input v-model.trim="pickupForm.pickup_user" />
            </el-form-item>
            <el-button type="primary" :loading="pickupLoading" @click="savePickup">记录</el-button>
          </el-form>
        </section>
      </el-col>
      <el-col :xs="24" :lg="8">
        <section class="section-panel">
          <h3 class="section-title">记录派件</h3>
          <el-form :model="deliveryForm" label-position="top">
            <el-form-item label="快递单号">
              <el-input v-model.trim="deliveryForm.tracking_no" />
            </el-form-item>
            <el-form-item label="派件员">
              <el-input v-model.trim="deliveryForm.courier_name" />
            </el-form-item>
            <el-form-item label="备注">
              <el-input v-model.trim="deliveryForm.remark" />
            </el-form-item>
            <el-button type="primary" :loading="deliveryLoading" @click="saveDelivery">记录</el-button>
          </el-form>
        </section>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { generatePickupCode, recordDelivery, recordPickup } from '../../api/pickup'

const codeLoading = ref(false)
const pickupLoading = ref(false)
const deliveryLoading = ref(false)
const pickupCode = ref('')
const codeForm = reactive({ tracking_no: '', recipient_phone: '' })
const pickupForm = reactive({ tracking_no: '', pickup_code: '', pickup_user: '' })
const deliveryForm = reactive({ tracking_no: '', courier_name: '', remark: '' })

async function createCode() {
  codeLoading.value = true
  try {
    const data = await generatePickupCode(codeForm)
    pickupCode.value = data?.pickup_code || data?.pickupCode || data?.code || String(data || '')
  } finally {
    codeLoading.value = false
  }
}

async function savePickup() {
  pickupLoading.value = true
  try {
    await recordPickup(pickupForm)
    ElMessage.success('取件记录已保存')
  } finally {
    pickupLoading.value = false
  }
}

async function saveDelivery() {
  deliveryLoading.value = true
  try {
    await recordDelivery(deliveryForm)
    ElMessage.success('派件记录已保存')
  } finally {
    deliveryLoading.value = false
  }
}
</script>

<style scoped>
.result-alert {
  margin-top: 14px;
}

.el-col {
  margin-bottom: 16px;
}
</style>
