<template>
  <div class="page">
    <div class="page-header">
      <div>
        <h2 class="page-title">寄件下单</h2>
        <p class="page-desc">填写寄件和收件信息，提交后可立即生成支付单。</p>
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
        <el-form-item label="物品信息" prop="item_info">
          <el-input v-model.trim="form.item_info" type="textarea" :rows="3" />
        </el-form-item>
        <el-form-item label="重量(kg)" prop="weight">
          <el-input-number
            v-model="form.weight"
            :min="0.01"
            :precision="2"
            :step="0.1"
            controls-position="right"
            style="width: 100%"
          />
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
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { createSendOrder } from '../../api/send'
import { createPayment } from '../../api/pay'
import { getMyCoupons, useCoupon } from '../../api/coupon'

const RECENT_SEND_ORDER_KEY = 'package_station_recent_send_orders'
const formRef = ref()
const loading = ref(false)
const router = useRouter()

const form = reactive({
  sender_name: '',
  sender_phone: '',
  sender_address: '',
  receiver_name: '',
  receiver_phone: '',
  receiver_address: '',
  item_info: '',
  weight: null
})

const rules = {
  sender_name: [{ required: true, message: '请输入寄件人', trigger: 'blur' }],
  sender_phone: [{ required: true, message: '请输入寄件电话', trigger: 'blur' }],
  sender_address: [{ required: true, message: '请输入寄件地址', trigger: 'blur' }],
  receiver_name: [{ required: true, message: '请输入收件人', trigger: 'blur' }],
  receiver_phone: [{ required: true, message: '请输入收件电话', trigger: 'blur' }],
  receiver_address: [{ required: true, message: '请输入收件地址', trigger: 'blur' }],
  item_info: [{ required: true, message: '请输入物品信息', trigger: 'blur' }],
  weight: [
    { required: true, message: '请输入重量', trigger: 'change' },
    {
      validator: (_, value, callback) => {
        if (typeof value !== 'number' || value <= 0) {
          callback(new Error('重量必须大于 0'))
          return
        }
        callback()
      },
      trigger: 'change'
    }
  ]
}

function readRecentOrders() {
  try {
    const parsed = JSON.parse(localStorage.getItem(RECENT_SEND_ORDER_KEY) || '[]')
    return Array.isArray(parsed) ? parsed : []
  } catch {
    return []
  }
}

function rememberRecentOrder(order) {
  const list = readRecentOrders().filter((item) => item.order_no !== order.order_no)
  list.unshift(order)
  localStorage.setItem(RECENT_SEND_ORDER_KEY, JSON.stringify(list.slice(0, 20)))
}

function normalizeCoupons(data) {
  const list = Array.isArray(data) ? data : data?.list || []
  return list
    .map((item) => ({
      user_coupon_id: Number(item.user_coupon_id || item.userCouponId || item.id),
      code: item.code || item.coupon_code || item.couponCode || '',
      amount: Number(item.amount || 0),
      threshold: Number(item.threshold || 0),
      status: (item.status || '').toLowerCase()
    }))
    .filter((item) => Number.isFinite(item.user_coupon_id))
}

function chooseBestCoupon(coupons, orderAmount) {
  const available = coupons.filter((item) => item.status === 'unused')
  const matched = available.filter((item) => orderAmount >= item.threshold)
  if (matched.length === 0) {
    return null
  }
  return matched.sort((a, b) => b.amount - a.amount)[0]
}

async function tryUseCoupon(order) {
  const orderNo = order?.order_no || order?.orderNo
  const estimatedFee = Number(order?.estimated_fee ?? order?.estimatedFee ?? 0)
  if (!orderNo || !Number.isFinite(estimatedFee)) {
    return order
  }

  try {
    const data = await getMyCoupons()
    const coupons = normalizeCoupons(data)
    const bestCoupon = chooseBestCoupon(coupons, estimatedFee)

    if (!bestCoupon) {
      const hasUnused = coupons.some((item) => item.status === 'unused')
      if (hasUnused) {
        ElMessage.info(`当前订单金额 ¥${estimatedFee.toFixed(2)} 未达到红包使用门槛`)
      }
      return order
    }

    const shouldUse = await ElMessageBox.confirm(
      `检测到可用红包 ${bestCoupon.code || ''}，可减 ¥${bestCoupon.amount.toFixed(2)}，门槛 ¥${bestCoupon.threshold.toFixed(2)}。\n是否使用红包？`,
      '红包抵扣',
      {
        type: 'warning',
        confirmButtonText: '使用红包',
        cancelButtonText: '不使用',
        distinguishCancelAndClose: true
      }
    )
      .then(() => true)
      .catch((action) => {
        if (action === 'cancel') {
          return false
        }
        throw new Error(action)
      })

    if (!shouldUse) {
      return order
    }

    const updated = await useCoupon({
      order_no: orderNo,
      user_coupon_id: bestCoupon.user_coupon_id
    })
    ElMessage.success(`红包已使用，减免 ¥${bestCoupon.amount.toFixed(2)}`)
    return updated || order
  } catch (error) {
    const message = error?.response?.data?.msg || error?.message || '红包使用失败'
    if (message !== 'close') {
      ElMessage.warning(message)
    }
    return order
  }
}

async function submit() {
  await formRef.value.validate()
  loading.value = true
  try {
    const payload = {
      sender_name: form.sender_name,
      sender_phone: form.sender_phone,
      sender_address: form.sender_address,
      receiver_name: form.receiver_name,
      receiver_phone: form.receiver_phone,
      receiver_address: form.receiver_address,
      item_info: form.item_info,
      itemInfo: form.item_info,
      weight: form.weight,
      item_desc: form.item_info
    }
    const created = await createSendOrder(payload)
    const orderAfterCoupon = await tryUseCoupon(created)
    const orderNo = orderAfterCoupon?.order_no || orderAfterCoupon?.orderNo
    const estimatedFee = orderAfterCoupon?.estimated_fee ?? orderAfterCoupon?.estimatedFee ?? null

    if (orderNo) {
      rememberRecentOrder({
        order_no: orderNo,
        estimated_fee: estimatedFee,
        pay_status: orderAfterCoupon?.pay_status || orderAfterCoupon?.payStatus || 'unpaid',
        created_at: orderAfterCoupon?.created_at || orderAfterCoupon?.createdAt || ''
      })

      try {
        const pay = await createPayment({
          related_type: 'send_order',
          order_no: orderNo
        })
        const payNo = pay?.pay_no || pay?.payNo || ''
        ElMessage.success(
          payNo ? `寄件订单已提交，支付单已创建：${payNo}` : '寄件订单已提交，支付单已创建'
        )
        router.push({
          path: '/app/user/bills',
          query: { order_no: orderNo, pay_no: payNo }
        })
      } catch {
        ElMessage.warning('寄件订单已提交，支付单创建失败，请稍后在账单页重试')
        router.push({
          path: '/app/user/bills',
          query: { order_no: orderNo }
        })
      }
    } else {
      ElMessage.success('寄件订单已提交')
    }

    formRef.value.resetFields()
    form.weight = null
  } finally {
    loading.value = false
  }
}
</script>
