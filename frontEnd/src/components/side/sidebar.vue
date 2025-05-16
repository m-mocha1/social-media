<script setup>
import { ref, computed, defineProps, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import router from '@/router'
import { connectWebSocket, addListener } from '@/wepSock.js'


const allUsers = ref([])

// which  user im chating

const props = defineProps({
  username: String,
  pfp: String,
})

const notiNum = ref()
async function notiNumF() {
  const nc = await fetch('/api/countNoti', {
    method: 'post',
  })
  if (nc.ok) {
    notiNum.value = await nc.json()
    console.log('🚀 ~  notiNum.value:', notiNum.value)
  }
}

async function getUsers() {
  try {
    const usersResp = await fetch('/api/allUsers', {
      method: 'POSt',
    })
    if (usersResp.ok) {
      allUsers.value = await usersResp.json()
      console.log('🚀 ~ getUsers ~ allUsers:', allUsers.value)
    } else {
      console.log('no users')
    }
  } catch {
    console.log('Error fetching users')
  }
}

onMounted(() => {
  getUsers()
  notiNumF()
})

const route = useRoute()
const isHomePage = computed(() => route.path === '/')
const isExplorePage = computed(() => route.path === '/Explore')
const isGroupsPage = computed(() => route.path === '/Groups')
const isNotPage = computed(() => route.path === '/Notif')

async function logout() {
  await fetch('/api/out', {
    method: 'POST',
  })
  router.push('/login')
}
const myProfile = async () => {
  router.push({ path: '/profile', query: { username: props.username } })
}
</script>

<template>
  <div id="side" class="sidebar">
    <div class="options">
      <img v-if="pfp" :src="pfp" id="pfpPro" @click="myProfile" />
      <p @click="myProfile" id="username" v-if="username">{{ username }}</p>
    </div>
    <br />
    <router-link to="/" :class="{ glow: isHomePage }">
      <i class="fa-solid fa-house fa-2xs" style="color: #f5f5f5"> Home</i>
    </router-link>
    <br />
    <router-link to="/Groups" :class="{ glow: isGroupsPage }">
      <i class="fa-solid fa-user-group fa-2xs" style="color: #fafafa"> Groups</i>
    </router-link>
    <br />
    <router-link to="/Explore" :class="{ glow: isExplorePage }">
      <i class="fa-solid fa-globe fa-2xs" style="color: #f2f2f2"> Explore</i>
    </router-link>
    <br />
    <router-link to="/Notif" :class="{ glow: isNotPage }">
      <i class="fa-solid fa-inbox fa-2xs" style="color: #f7f7f7">
        notifications&nbsp;&nbsp; {{ notiNum }}</i
      >
    </router-link>
    <form @submit.prevent="logout">
      <button id="out">
        <i class="fa-solid fa-arrow-right-from-bracket fa-2xs">log-out</i>
      </button>
      <br />
    </form>
  </div>
</template>

<style scoped>
.noi {
  margin-left: 35px;
}
</style>
