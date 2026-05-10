<template>
  <div class="page">
    <div class="page-header">
      <div>
        <h2 class="page-title">红包礼券设置</h2>
        <p class="page-desc">管理员创建红包券，并按状态/关键字分页查看已发放礼券。</p>
      </div>
    </div>

    <section class="section-panel form-narrow">
      <h3 class="section-title">创建红包券</h3>
      <el-form ref="formRef" :model="form" :rules="rules" label-width="110px">
        <el-form-item label="券名称" prop="name">
          <el-input v-model.trim="form.name" placeholder="如：618满减券" />
        </el-form-item>
        <el-form-item label="券面额" prop="amount">
          <el-input-number
            v-model="form.amount"
            :min="0.01"
            :precision="2"
            :step="1"
            controls-position="right"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="券码">
          <el-input v-model.trim="form.code" placeholder="不填则系统自动生成" />
        </el-form-item>
        <el-form-item label="活动规则">
          <el-input
            v-model.trim="form.activity_rule"
            type="textarea"
            :rows="3"
            placeholder="如：618活动券"
          />
        </el-form-item>
        <el-form-item label="使用门槛">
          <el-input-number
            v-model="form.threshold"
            :min="0"
            :precision="2"
            :step="5"
            controls-position="right"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="发放总量">
          <el-input-number
            v-model="form.total"
            :min="1"
            :step="100"
            controls-position="right"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="有效天数">
          <el-input-number
            v-model="form.valid_days"
            :min="1"
            :step="1"
            controls-position="right"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="submitting" @click="submit">创建红包券</el-button>
          <el-button :disabled="submitting" @click="reset">重置</el-button>
        </el-form-item>
      </el-form>
    </section>

    <section class="section-panel">
      <div class="toolbar">
        <el-select
          v-model="query.status"
          clearable
          placeholder="状态筛选"
          style="width: 140px"
          @change="handleFilterChange"
        >
          <el-option label="生效中" value="active" />
          <el-option label="已失效" value="inactive" />
          <el-option label="已结束" value="expired" />
        </el-select>
        <el-input
          v-model.trim="query.keyword"
          placeholder="输入券码或名称搜索"
          clearable
          style="width: 260px"
          @keyup.enter="handleSearch"
        />
        <el-button :loading="listLoading" @click="handleSearch">搜索</el-button>
        <el-button :loading="listLoading" @click="resetFilters">重置</el-button>
      </div>

      <el-table v-loading="listLoading" :data="rows" border stripe empty-text="暂无红包券">
        <el-table-column prop="code" label="券码" min-width="140" show-overflow-tooltip />
        <el-table-column prop="name" label="券名称" min-width="160" show-overflow-tooltip />
        <el-table-column label="面额" min-width="90">
          <template #default="{ row }">
            ¥{{ formatAmount(row.amount) }}
          </template>
        </el-table-column>
        <el-table-column label="门槛" min-width="90">
          <template #default="{ row }">
            ¥{{ formatAmount(row.threshold) }}
          </template>
        </el-table-column>
        <el-table-column prop="total" label="总量" min-width="90" />
        <el-table-column prop="remaining" label="剩余" min-width="90" />
        <el-table-column label="状态" min-width="100">
          <template #default="{ row }">
            {{ statusText(row.status) }}
          </template>
        </el-table-column>
      </el-table>

      <div class="pager-wrap">
        <el-pagination
          v-model:current-page="query.page"
          v-model:page-size="query.page_size"
          :total="total"
          :page-sizes="[20, 50, 100, 200]"
          background
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handlePageSizeChange"
          @current-change="handlePageChange"
        />
      </div>
    </section>
  </div>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { createAdminCoupon, getAdminCouponList } from '../../api/coupon'

const formRef = ref()
const submitting = ref(false)
const listLoading = ref(false)
const rows = ref([])
const total = ref(0)

const form = reactive({
  name: '',
  amount: null,
  code: '',
  activity_rule: '',
  threshold: 0,
  total: 1000,
  valid_days: 30
})

const query = reactive({
  status: '',
  keyword: '',
  page: 1,
  page_size: 20
})

const rules = {
  name: [{ required: true, message: '请输入券名称', trigger: 'blur' }],
  amount: [
    { required: true, message: '请输入券面额', trigger: 'change' },
    {
      validator: (_, value, callback) => {
        if (typeof value !== 'number' || value <= 0) {
          callback(new Error('券面额必须大于 0'))
          return
        }
        callback()
      },
      trigger: 'change'
    }
  ]
}

function formatAmount(value) {
  const amount = Number(value)
  return Number.isFinite(amount) ? amount.toFixed(2) : '0.00'
}

function statusText(status) {
  const value = (status || '').toLowerCase()
  if (value === 'active') return '生效中'
  if (value === 'inactive') return '已失效'
  if (value === 'expired') return '已结束'
  return status || '-'
}

function buildCreatePayload() {
  const payload = {
    name: form.name,
    amount: Number(form.amount),
    threshold: Number(form.threshold ?? 0),
    total: Number(form.total ?? 1000),
    valid_days: Number(form.valid_days ?? 30)
  }
  if (form.code) {
    payload.code = form.code
  }
  if (form.activity_rule) {
    payload.activity_rule = form.activity_rule
  }
  return payload
}

function buildListParams() {
  const params = {
    page: query.page,
    page_size: query.page_size
  }
  if (query.status) {
    params.status = query.status
  }
  if (query.keyword) {
    params.keyword = query.keyword
  }
  return params
}

async function loadList() {
  listLoading.value = true
  try {
    const data = await getAdminCouponList(buildListParams())
    const list = Array.isArray(data) ? data : data?.list || []
    rows.value = list
    total.value = Number(data?.total ?? list.length ?? 0)
    if (data?.page) {
      query.page = Number(data.page)
    }
    if (data?.page_size) {
      query.page_size = Number(data.page_size)
    }
  } finally {
    listLoading.value = false
  }
}

async function submit() {
  await formRef.value.validate()
  submitting.value = true
  try {
    const payload = buildCreatePayload()
    const created = await createAdminCoupon(payload)
    const code = created?.code || payload.code || '自动生成'
    ElMessage.success(`红包券创建成功：${code}`)
    reset()
    query.page = 1
    await loadList()
  } finally {
    submitting.value = false
  }
}

function reset() {
  formRef.value?.resetFields()
  form.name = ''
  form.amount = null
  form.code = ''
  form.activity_rule = ''
  form.threshold = 0
  form.total = 1000
  form.valid_days = 30
}

async function handleSearch() {
  query.page = 1
  await loadList()
}

async function handleFilterChange() {
  query.page = 1
  await loadList()
}

async function handlePageChange(page) {
  query.page = page
  await loadList()
}

async function handlePageSizeChange(size) {
  query.page_size = size
  query.page = 1
  await loadList()
}

async function resetFilters() {
  query.status = ''
  query.keyword = ''
  query.page = 1
  query.page_size = 20
  await loadList()
}

onMounted(loadList)
</script>

<style scoped>
.pager-wrap {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}
</style>
