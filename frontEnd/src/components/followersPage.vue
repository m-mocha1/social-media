<script setup>
import { ref, onMounted } from 'vue'
import router from '@/router'
import sidebar from '@/components/side/sidebar.vue'
import Select from './postType/select.vue'
import sideUsers from '@/components/chat/sideUsers.vue'
import userCard from '@/components/userCard/card.vue'

const userData = ref(null)
const me = ref(null)
const fing = ref([])
const fer = ref([])

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
async function fetchFing() {
  const getFing = await fetch('/api/get-following', {
    method: 'POST',
  })
  if (getFing.ok) {
    fing.value = await getFing.json()
    console.log('🚀 ~ fetchFing ~ fing:', fing.value)
  }
}
async function fetchFer() {
  const getFer = await fetch('/api/get-followers', {
    method: 'POST',
  })
  if (getFer.ok) {
    fer.value = await getFer.json()
    console.log('🚀 ~ fetchFing ~ fier:', fer.value)
  }
}
onMounted(() => {
  auth()
  fetchFer()
  fetchFing()
  changeCss()
})
const myProfile = async () => {
  router.push({ path: '/profile', query: { username: userData.value.Username } })
}
</script>

<template v-cloak>
  <div class="main-container">
    <sidebar :username="userData?.Username" :pfp="userData?.pfp" />
    <div class="content">

      <transition-group name="fade" tag="div"id="postFeed" class="post-feed" >
        <p> <i id="fi" class="fa-solid fa-user-plus fa-xs" style="color: #ffffff;"></i> Following</p>
       <br>
        <userCard class="fing3" v-for="user in fing" :user="user" />
      </transition-group>

      <transition-group name="fade" tag="div" id="postFeed" class="post-feed" >
        <p><i id="fer" class="fa-solid fa-user-group fa-xs" style="color: #ffffff;"></i> Followers</p>
        <br>
        <userCard class="fer3" v-for="user in fer" :user="user" />
      </transition-group>
    
    </div>

    <sideUsers v-if="me" :me="me" />
  </div>
  <!-- end of post feed -->
  <!-- end of content -->
  <!-- end of main container -->
</template>
<style scoped>
#fer{
  margin:10px; 
}
#fi{
  margin:10px; 
}
.content {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  grid-template-rows: repeat(5, 1fr);
  gap: 8px;
}
.fing3 {
  grid-row: span 5 / span 5;
}
.fer3 {
  grid-row: span 5 / span 5;
  grid-column-start: 3;
}

body {
  background-color: black;
}
[v-cloak] {
  display: none;
}

.fade-enter-active,
.fade-leave-active {
  transition:
    opacity 0.5s ease,
    transform 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}
</style>
