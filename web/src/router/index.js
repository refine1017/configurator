import Vue from 'vue'
import Router from 'vue-router'
import store from '../store'

Vue.use(Router)

/* Layout */
import Layout from '@/layout'

/**
 * Note: sub-menu only appear when route children.length >= 1
 * Detail see: https://panjiachen.github.io/vue-element-admin-site/guide/essentials/router-and-nav.html
 *
 * hidden: true                   if set true, item will not show in the sidebar(default is false)
 * alwaysShow: true               if set true, will always show the root menu
 *                                if not set alwaysShow, when item has more than one children route,
 *                                it will becomes nested mode, otherwise not show the root menu
 * redirect: noRedirect           if set noRedirect will no redirect in the breadcrumb
 * name:'router-name'             the name is used by <keep-alive> (must set!!!)
 * meta : {
    roles: ['admin','editor']    control the page roles (you can set multiple roles)
    title: 'title'               the name show in sidebar and breadcrumb (recommend set)
    icon: 'svg-name'             the icon show in the sidebar
    noCache: true                if set true, the page will no be cached(default is false)
    affix: true                  if set true, the tag will affix in the tags-view
    breadcrumb: false            if set false, the item will hidden in breadcrumb(default is true)
    activeMenu: '/example/list'  if set path, the sidebar will highlight the path you set
  }
 */

/**
 * constantRoutes
 * a base page that does not have permission requirements
 * all roles can be accessed
 */
export const constantRoutes = [
  {
    path: '/login',
    component: () => import('@/views/login/index'),
    hidden: true
  },
  {
    path: '/auth-redirect',
    component: () => import('@/views/login/auth-redirect'),
    hidden: true
  },
  {
    path: '/404',
    component: () => import('@/views/error-page/404'),
    hidden: true
  },
  {
    path: '/401',
    component: () => import('@/views/error-page/401'),
    hidden: true
  },
  {
    path: '',
    component: Layout,
    redirect: 'dashboard',
    children: [
      {
        path: 'dashboard',
        component: () => import('@/views/dashboard/index'),
        name: 'Dashboard',
        meta: { title: 'Dashboard', icon: 'component' }
      },
      {
        path: 'projects',
        component: () => import('@/views/project/index'),
        name: 'SelectProjects',
        hidden: true,
        meta: { title: 'SelectProjects', icon: 'component' }
      }
    ]
  },
  {
    path: '/system',
    component: Layout,
    name: 'System',
    meta: {
      title: 'System',
      icon: 'table'
    },
    children: [
      {
        path: 'admins',
        component: () => import('@/views/system/admins'),
        name: 'Admins',
        meta: { title: 'Admins', icon: 'component' }
      },
      {
        path: 'AdminLogs',
        component: () => import('@/views/system/admin_logs'),
        name: 'AdminLogs',
        meta: { title: 'AdminLogs', icon: 'component' }
      }
    ]
  },
  {
    path: '/management',
    component: Layout,
    name: 'Management',
    meta: {
      title: 'Management',
      icon: 'table'
    },
    children: [
      {
        path: 'projects',
        component: () => import('@/views/management/projects'),
        name: 'Projects',
        meta: { title: 'Projects', icon: 'component' }
      },
      {
        path: 'environments',
        component: () => import('@/views/management/environments'),
        name: 'Environments',
        meta: { title: 'Environments', icon: 'component' }
      },
      {
        path: 'configs',
        component: () => import('@/views/management/configs'),
        name: 'Configs',
        meta: { title: 'Configs', icon: 'component' }
      },
      {
        path: 'servers',
        component: () => import('@/views/management/servers'),
        name: 'Servers',
        meta: { title: 'Servers', icon: 'component' }
      },
      {
        path: 'fields',
        component: () => import('@/views/management/fields'),
        name: 'Fields',
        meta: { title: 'Fields Editor', noCache: true },
        hidden: true
      }
    ]
  }
]

/**
 * asyncRoutes
 * the routes that need to be dynamically loaded based on user roles
 */
export const asyncRoutes = [
  /** when your routing map is too long, you can split it into small modules **/
  // tableRouter,
  // componentsRouter,
  { path: '*', redirect: '/404', hidden: true }
]

const createRouter = () => new Router({
  // mode: 'history', // require service support
  scrollBehavior: () => ({ y: 0 }),
  routes: constantRoutes
})

const router = createRouter()

export function configRoutes() {
  const env = store.getters.userSelectEnvironment

  const configRouters = {
    path: '/config',
    component: Layout,
    name: 'Configs',
    meta: {
      title: 'Configs',
      icon: 'table'
    },
    children: [
      {
        path: 'json',
        component: () => import('@/views/config/json_editor'),
        name: 'Json',
        meta: { title: 'Json Editor', noCache: true },
        hidden: true
      }
    ]
  }

  const configViews = {
    'table': import('@/views/config/table'),
    'json': import('@/views/config/json'),
    'kv': import('@/views/config/kv')
  }

  for (let j = 0; j < env.configs.length; j++) {
    const config = env.configs[j]

    configRouters.children.push({
      path: config.name,
      component: () => configViews[config.format],
      name: config.name,
      meta: { title: config.name, icon: 'component' },
      params: { envId: env.id }
    })
  }

  return [configRouters]
}

// Detail see: https://github.com/vuejs/vue-router/issues/1234#issuecomment-357941465
export function resetRouter() {
  const newRouter = createRouter()
  router.matcher = newRouter.matcher // reset router
}

export default router
