<script setup>
import { ref, onMounted } from 'vue'
import router from '@/router'
import { useRoute } from 'vue-router'
import { connectWebSocket, addListener } from '../wepSock.js'
import { triggerRef } from 'vue'
import sidebar from '@/components/side/sidebar.vue'
import sideUsers from '@/components/chat/sideUsers.vue'


import postTemp from '@/components/singlePost/singlePost.vue'
const route = useRoute()
const postData = ref(null)
const userData = ref(null)
const name = route.query.user

function changeCss() {
  const oldcss = document.getElementById('dc')

  if (oldcss) {
    console.log(oldcss)

    oldcss.remove()
  }
  const link = document.createElement('link')
  link.id = 'dc'
  link.rel = 'stylesheet'
  link.href = '../src/assets/mainPage.css'
  document.head.appendChild(link)
}

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
      console.log('🚀 ~ auth ~ userData:', userData.value)
    } else {
      console.log('wat');
      
    }
  } catch {}
}



async function onePost() {
  const postId = route.query.id
  try {
    const formData = new FormData()
    formData.append('postid', postId)

    const response = await fetch('/api/onepost', {
      method: 'POST',
      body: formData,
    })
    if (response.ok) {
      postData.value = await response.json()
      console.log('🚀 ~ onePost ~ postData:', postData.value)
    } else {
      console.error('Failed to fetch post details')
    }
  } catch (error) {
    console.error('Error fetching post:', error)
  }
}
// addListener((data) => {
//   if (data.type === 'comment') {
//     comments.value.unshift(data)
//     console.log(comments.value)
//   }
// })
triggerRef(postData)
onMounted(() => {
  changeCss()
  onePost()
  auth()
})
</script>

<template>
  <sidebar :username="userData?.Username" :pfp="userData?.pfp" />
  
  <div id="postFeed" class="post-feed">
          
    <postTemp v-if="postData" :post="postData" :user="name"></postTemp>
  </div>
</template>

<style></style>
