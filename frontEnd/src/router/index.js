import { createRouter, createWebHistory } from 'vue-router'
import login from '../components/login/login.vue'
import signUP from '../components/login/signUp.vue'
import home from '../components/home.vue'
import viewPost from '../components/viewPost.vue'
import profile from '@/components/profile.vue'
import explore from '@/components/explore.vue'
import groups from '../components/groups.vue'
import GroupPage from '@/components/GroupPage.vue'
import MemsInvs from '@/components/MemsInvs.vue'
import Notif from '@/components/notification.vue'
import followersPage from '@/components/followersPage.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: login,
    },
    {
      path: '/sign-up',
      name: 'sign-up',
      component: signUP,
    },
    {
      path: '/',
      name: 'home',
      component: home,
    },
    {
      path: '/Explore',
      name: 'explore',
      component: explore,
    },
    {
      path: '/post',
      name: 'viewPost',
      component: viewPost,
    },
    {
      path: '/profile',
      name: 'profile',
      component: profile,
    },
    {
      path: '/groups',
      name: 'groups',
      component: groups,
    },
    {
      path: '/group',
      name: '/group',
      component: GroupPage,
    },
    {
      path: '/groupMembers',
      name: '/groupMembers',
      component: MemsInvs,
    },
    {
      path: '/Notif',
      name: 'Notif',
      component: Notif,
    },
    {
      path: '/followersPage',
      name: 'followersPage',
      component: followersPage,
    },
  ],
})

export default router
