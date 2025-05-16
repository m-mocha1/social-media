<script setup>
import { ref, onMounted } from 'vue'
import router from '@/router'
import notiCard from '@/components/noti/notif.vue'
import { connectWebSocket, addListener } from '../wepSock.js'
import sidebar from '@/components/side/sidebar.vue'
import Select from './postType/select.vue'
import sideUsers from '@/components/chat/sideUsers.vue'

const userData = ref(null)
const notification = ref([])
const me = ref(null)

async function auth() {
  try {
    const AuthResp = await fetch('/api/Auth', {
      method: 'POST',
    })
    if (!AuthResp.ok) {
      router.push('/login')
    }
  } catch (error) {
    console.error('Fetch error:', error)
  }
  try {
    const userResp = await fetch('/api/data', {
      method: 'POST',
    })
    if (userResp.ok) {
      userData.value = await userResp.json()
      console.log('🚀 ~ auth ~ userData:', userData)
      if (!me.value) {
        me.value = { me: userData.value.Username, mepfp: userData.value.pfp }
      } else {
        me.value = null
      }
    } else {
    }
  } catch {}
}

onMounted(() => {
  auth()
  changeCss()
  getNoti()
})

addListener((data) => {
  console.log('🚀 ~ addListener ~ data.type:', data)

  if (
    userData.value &&
    (data.type === 'followreq' ||
      data.type === 'startFollow' ||
      data.type === 'like' ||
      data.type === 'comment')
  ) {
    if (!Array.isArray(notification.value)) {
      notification.value = []
    }
    notification.value.unshift(data)
  }
})

const removeNoti = (notiId) => {
  notification.value = notification.value.filter((noti) => noti.id !== notiId)
}
async function getNoti() {
  try {
    const notiResp = await fetch('/api/getNoti', {
      method: 'POST',
    })
    if (notiResp.ok) {
      notification.value = await notiResp.json()
      console.log('🚀 ~ getNoti ~ notification:', notification.value)
    } else {
      const errorText = await notiResp.text()
      console.error('Error fetching notifications:', notiResp.status, errorText)
    }
  } catch (error) {
    console.error('Fetch error:', error)
  }
}

function changeCss() {
  const oldcss = document.getElementById('dc')

  if (oldcss) {
    oldcss.remove()
  }
  const link = document.createElement('link')
  link.id = 'dc'
  link.rel = 'stylesheet'
  link.href = '/src/assets/mainPage.css'
  document.head.appendChild(link)
}
</script>

<template v-cloak>
  <div class="main-container">
    <sidebar :username="userData?.Username" :pfp="userData?.pfp" />

    <div class="content">
      <transition-group name="fade" tag="div" id="postFeed" class="post-feed">
        <notiCard
          v-for="noti in notification"
          :key="noti.id"
          :noti="noti"
          @notification-deleted="removeNoti"
        />
      </transition-group>
    </div>
    <sideUsers v-if="me" :me="me" />
  </div>

  <!-- end of post feed -->
  <!-- end of content -->
  <!-- end of main container -->
</template>

<style>
body {
  background-color: black;
}
[v-cloak] {
  display: none;
}
</style>
