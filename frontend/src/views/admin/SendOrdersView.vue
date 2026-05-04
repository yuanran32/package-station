<template>
  <div class="page">
    <div class="page-header">
      <div>
        <h2 class="page-title">寄件订单</h2>
        <p class="page-desc">查看用户寄件申请，并执行接单、分配快递员、完成订单。</p>
      </div>
      <div class="toolbar">
        <el-select v-model="status" clearable placeholder="订单状态" style="width: 160px" @change="load">
          <el-option label="待处理" value="pending" />
          <el-option label="已接单" value="accepted" />
          <el-option label="派件中" value="assigned" />
          <el-option label="已完成" value="completed" />
        </el-select>
        <el-button type="primary" :icon="Refresh" :loading="loading" @click="load">刷新</el-button>
      </div>
    </div>
    <section class="section-panel">
      <el-table v-loading="loading" :data="rows" border stripe empty-text="暂无订单">
        <el-table-column prop="order_no" label="订单号" min-width="160" show-overflow-tooltip />
        <el-table-column prop="sender_name" label="寄件人" min-width="100" />
        <el-table-column prop="receiver_name" label="收件人" min-width="100" />
        <el-table-column prop="receiver_address" label="收件地址" min-width="220" show-overflow-tooltip />
        <el-table-column prop="status" label="状态" min-width="100" />
        <el-table-column label="操作" width="300" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="process(row, 'accept')">接单</el-button>
            <el-button size="small" type="warning" @click="openAssign(row)">分配</el-button>
            <el-button size="small" type="success" @click="process(row, 'complete')">完成</el-button>
          </template>
        </el-table-column>
      </el-table>
    </section>
    <el-dialog v-model="assignVisible" title="分配快递员" width="420px">
      <el-form :model="assignForm" label-width="90px">
        <el-form-item label="订单号">
          <el-input v-model="assignForm.order_no" disabled />
        </el-form-item>
        <el-form-item label="快递员">
          <el-input v-model.trim="assignForm.courier_name" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="assignVisible = false">取消</el-button>
        <el-button type="primary" :loading="processing" @click="assign">确认</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { Refresh } from '@element-plus/icons-vue'
import { getAdminSendOrders, processSendOrder } from '../../api/send'

const loading = ref(false)
const processing = ref(false)
const status = ref('')
const rows = ref([])
const assignVisible = ref(false)
const assignForm = reactive({
  order_no: '',
  courier_name: ''
})

function getOrderNo(row) {
  return row.order_no || row.orderNo || row.id
}

async function load() {
  loading.value = true
  try {
    const data = await getAdminSendOrders(status.value)
    rows.value = Array.isArray(data) ? data : data?.list || []
  } finally {
    loading.value = false
  }
}

async function process(row, action) {
  processing.value = true
  try {
    await processSendOrder({ order_no: getOrderNo(row), action })
    ElMessage.success('订单处理成功')
    await load()
  } finally {
    processing.value = false
  }
}

function openAssign(row) {
  assignForm.order_no = getOrderNo(row)
  assignForm.courier_name = ''
  assignVisible.value = true
}

async function assign() {
  if (!assignForm.courier_name) {
    ElMessage.warning('请输入快递员姓名')
    return
  }
  processing.value = true
  try {
    await processSendOrder({
      order_no: assignForm.order_no,
      action: 'assign_pickup',
      courier_name: assignForm.courier_name
    })
    ElMessage.success('已分配快递员')
    assignVisible.value = false
    await load()
  } finally {
    processing.value = false
  }
}

onMounted(load)
</script>
