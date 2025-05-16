<script setup>
import { ref, onMounted, reactive } from 'vue'
import router from '@/router'
import { useRoute, useRouter } from 'vue-router'
import PostCard from '@/components/singlePost/post.vue'
import sidebar from '@/components/side/sidebar.vue'
import sideUsers from '@/components/chat/sideUsers.vue'

const route = useRoute()
const name = route.query.username
const userData = ref(null)
const myData = ref(null)
const me2 = ref(null)
console.log("🚀 ~ name:", name)

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
  const myDataResp = await fetch('/api/data', {
    method: 'POST',
  })
  if (myDataResp.ok) {
    myData.value = await myDataResp.json()
    console.log('🚀 ~ auth ~ myDataResp:', myData)
       if (!me2.value) {
        me2.value = { me: myData.value.Username, mepfp: myData.value.pfp }
      } else {
        me2.value = null
      }
  }
  try {
    let username = new FormData()
    username.append('username', name)
    const Re = await fetch('/api/getUserPost', {
      method: 'post',
      body: username,
    })
    if (Re.ok) {
      userData.value = await Re.json()
      console.log('🚀 ~ userPosts ~ Userposts:', userData.value.user.Username)
    } else {
      console.error('Error fetching posts')
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
  link.href = '/src/assets/profile.css'
  document.head.appendChild(link)
}
async function closeProfile() {
  if (userData.value.user.public == 'public') {
    const pri = await fetch('/api/alPri', {
      method: 'POST',
    })
    if (pri.ok) {
      userData.value.user.public = 'alPri'
    }
  } else if (userData.value.user.public == 'alPri') {
    const pri = await fetch('/api/pri', {
      method: 'POST',
    })
    if (pri.ok) {
      userData.value.user.public = 'pri'
    }
  } else if (userData.value.user.public == 'pri') {
    const pri = await fetch('/api/pub', {
      method: 'post',
    })
    if (pri.ok) {
      userData.value.user.public = 'public'
    }
  }
  console.log(userData.value.user.public)
}

async function followReq() {
  const follData = new FormData()
  follData.append('username', userData.value.user.Username)
  console.log('🚀 ~ followReq ~ Data:', follData.value)

  const followResp = await fetch('/api/followReq', {
    method: 'POST',
    body: follData
  })
  if (followResp.ok) {
    console.log('🚀 ~ followReq ~ followReq is good')
  } else {
    console.log('error in follow req')
  }
}

onMounted(() => {
  auth()
  changeCss()
})

async function follow() {
  const followData = new FormData()
  //im the fer and the fing is the other guy
  followData.append('fing', name)

  if (!userData.value.user.following) {
    console.log('fol')

    await fetch('api/follower', {
      method: 'POST',
      body: followData,
    })
    userData.value.user.following = true
    userData.value.user.followers++
  } else {
    console.log('unfol')

    await fetch('api/unfollow', {
      method: 'POST',
      body: followData,
    })
    userData.value.user.following = false
    userData.value.user.followers--
  }
  console.log('🚀 ~ follow ~ followers--:', userData.value.user.followers)
}
const fol = async () => {
  router.push({ path: '/followersPage', query: { username: userData.value.Username } })
}
</script>

<template>
  <div class="main-container">
    <sidebar :username="myData?.Username" :pfp="myData?.pfp" />

    <div class="content">
      <div class="welcome">
        <div class="options">
          <img v-if="userData" :src="userData.user?.pfp" id="pfpmain" @click="profile" />
          <br />
        </div>
        <br />
        <p @click="fol" id="username" v-if="userData">
          Name: {{ userData.user.Username }} - Age: {{ userData.user.Age }}
          <i class="fa-solid fa-circle-user"> {{ userData?.user.followers }}</i> <br />
        </p>
      </div>

      <div class="proCont">
        <button
          class="lock"
          v-if="userData?.user.public == 'public' && userData?.user.Me"
          @click="closeProfile"
        >
          <i class="fa-solid fa-lock-open" style="color: #f7f7f7"> Public</i>
        </button>
        <button
          class="lock"
          v-else-if="userData?.user.public == 'alPri' && userData?.user.Me"
          @click="closeProfile"
        >
          <i class="fa-solid fa-lock" style="color: #efeff1"> Followers only</i>
        </button>

        <button
          class="lock"
          v-else-if="userData?.user.public == 'pri' && userData?.user.Me"
          @click="closeProfile"
        >
          <i class="fa-solid fa-lock" style="color: #efeff1"> Only me</i>
        </button>

        <!-- !userData?.user.following" -->
        <button
          v-if="
            !userData?.user.following && !userData?.user.Me && userData?.user.public == 'public'
          "
          @click="follow"
          class="fow"
        >
          <i class="fa-solid fa-user-plus"> follow </i>
        </button>

        <button
          v-else-if="userData?.user.following && !userData?.user.Me"
          @click="follow"
          class="fow"
        >
          <i class="fa-solid fa-circle-user"> following </i>
        </button>

        <button
          v-else-if="
            !userData?.user.following &&
            !userData?.user.Me &&
            (userData?.user.public == 'pri' || userData?.user.public == 'alPri')
          "
          @click="followReq"
          class="fow"
        >
          <i class="fa-solid fa-user-plus"> request to follow </i>
        </button>
      </div>

      <hr />
      <div id="postFeed" class="post-feed">
        <!-- this not workdfsdf -->
        <PostCard
          v-if="userData?.posts"
          v-for="post in userData?.posts"
          :key="post.postID"
          :post="post"
          :user="myData"
        />
        <p v-else-if="!userData?.user.public && !userData?.user.Me">this account is Private</p>
        <p v-else>no posts yet</p>
      </div>
    </div>
    <sideUsers v-if="me2" :me="me2" />

    <!-- end of post feed -->
  </div>
  <!-- end of content -->
  <!-- end of main container -->
</template>
<style></style>
