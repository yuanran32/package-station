import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const routes = [
  { path: '/', redirect: '/login' },
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/auth/LoginView.vue'),
    meta: { public: true }
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('../views/auth/RegisterView.vue'),
    meta: { public: true }
  },
  {
    path: '/app',
    component: () => import('../components/AppLayout.vue'),
    meta: { requiresAuth: true },
    children: [
      { path: '', redirect: '/app/user/dashboard' },
      {
        path: 'user/dashboard',
        name: 'UserDashboard',
        component: () => import('../views/user/UserDashboard.vue'),
        meta: { role: 'user', title: '用户首页' }
      },
      {
        path: 'user/profile',
        name: 'UserProfile',
        component: () => import('../views/user/ProfileView.vue'),
        meta: { role: 'user', title: '个人资料' }
      },
      {
        path: 'user/parcel',
        name: 'UserParcel',
        component: () => import('../views/user/ParcelStatusView.vue'),
        meta: { role: 'user', title: '快递查询' }
      },
      {
        path: 'user/history',
        name: 'UserPickupHistory',
        component: () => import('../views/user/PickupHistoryView.vue'),
        meta: { role: 'user', title: '取件历史' }
      },
      {
        path: 'user/identity',
        name: 'UserIdentity',
        component: () => import('../views/user/IdentityCodeView.vue'),
        meta: { role: 'user', title: '身份码' }
      },
      {
        path: 'user/send-order',
        name: 'UserSendOrder',
        component: () => import('../views/user/SendOrderView.vue'),
        meta: { role: 'user', title: '寄件下单' }
      },
      {
        path: 'user/pay-send',
        name: 'UserPaySend',
        component: () => import('../views/user/BillsView.vue'),
        meta: { role: 'user', title: '寄件支付' }
      },
      {
        path: 'user/coupons',
        name: 'UserCoupons',
        component: () => import('../views/user/CouponsView.vue'),
        meta: { role: 'user', title: '红包礼券' }
      },
      {
        path: 'user/bills',
        name: 'UserBills',
        component: () => import('../views/user/BillsView.vue'),
        meta: { role: 'user', title: '支付账单' }
      },
      {
        path: 'user/notices',
        name: 'UserNotices',
        component: () => import('../views/user/NotificationsView.vue'),
        meta: { role: 'user', title: '通知提醒' }
      },
      {
        path: 'admin/dashboard',
        name: 'AdminDashboard',
        component: () => import('../views/admin/AdminDashboard.vue'),
        meta: { role: 'admin', title: '管理首页' }
      },
      {
        path: 'admin/inbound',
        name: 'AdminInbound',
        component: () => import('../views/admin/ParcelInboundView.vue'),
        meta: { role: 'admin', title: '快递入库' }
      },
      {
        path: 'admin/parcels',
        name: 'AdminParcels',
        component: () => import('../views/admin/ParcelListView.vue'),
        meta: { role: 'admin', title: '快递列表' }
      },
      {
        path: 'admin/outbound',
        name: 'AdminOutbound',
        component: () => import('../views/admin/ParcelOutboundView.vue'),
        meta: { role: 'admin', title: '出库处理' }
      },
      {
        path: 'admin/pickup',
        name: 'AdminPickup',
        component: () => import('../views/admin/PickupManageView.vue'),
        meta: { role: 'admin', title: '取派件记录' }
      },
      {
        path: 'admin/send-orders',
        name: 'AdminSendOrders',
        component: () => import('../views/admin/SendOrdersView.vue'),
        meta: { role: 'admin', title: '寄件订单' }
      },
      {
        path: 'admin/coupons',
        name: 'AdminCoupons',
        component: () => import('../views/admin/AdminCouponCreateView.vue'),
        meta: { role: 'admin', title: '红包礼券设置' }
      },
      {
        path: 'admin/notices',
        name: 'AdminNotices',
        component: () => import('../views/admin/NoticeSendView.vue'),
        meta: { role: 'admin', title: '通知发送' }
      }
    ]
  },
  { path: '/:pathMatch(.*)*', redirect: '/login' }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to) => {
  const auth = useAuthStore()
  if (to.meta.public) {
    if (auth.isLoggedIn) {
      return auth.isAdmin ? '/app/admin/dashboard' : '/app/user/dashboard'
    }
    return true
  }
  if (to.meta.requiresAuth || to.meta.role) {
    if (!auth.isLoggedIn) {
      return '/login'
    }
    if (to.meta.role && to.meta.role !== auth.role) {
      return auth.isAdmin ? '/app/admin/dashboard' : '/app/user/dashboard'
    }
  }
  return true
})

export default router
