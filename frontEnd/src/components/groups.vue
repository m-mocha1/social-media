<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import sidebar from '@/components/side/sidebar.vue'
import sideUsers from '@/components/chat/sideUsers.vue'

import profileImg from '@/assets/profile.jpg'
const imgPrv = ref(profileImg)
const imgFile = ref(profileImg)

const router = useRouter()
const userData = ref(null)
const groupName = ref('')
const groupDesc = ref('')
const myGroups = ref([])
const allGroups = ref([])
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
  autoResizeTextarea()
  fetchGroups()
  fetchUserData()
})

function changeCss() {
  const oldcss = document.getElementById('dc')

  if (oldcss) {
    oldcss.remove()
  }
  const link = document.createElement('link')
  link.id = 'dc'
  link.rel = 'stylesheet'
  link.href = '/src/assets/groups.css'
  document.head.appendChild(link)
}

const autoResizeTextarea = () => {
  const postArea = ref(null)

  if (postArea.value) {
    postArea.value.style.width = 'auto'
    postArea.value.style.height = 'auto'
    postArea.value.style.height = `${postArea.value.scrollHeight}px`
  }
}

async function fetchGroups() {
  try {
    const response = await fetch('/api/getGroups')
    if (!response.ok) throw new Error('Failed to fetch groups')

    const data = await response.json()

    console.log('Fetched My Groups:', data.myGroups)
    console.log('Fetched All Groups:', data.allGroups)

    myGroups.value = data.myGroups || []
    allGroups.value = data.allGroups || []
  } catch (error) {
    console.error('Error fetching groups:', error.message)
  }
}

async function createGroup() {
  const groupData = new FormData()
  groupData.append('groupName', groupName.value)
  groupData.append('groupDesc', groupDesc.value)

  if (imgFile.value) {
    groupData.append('groupImage', imgFile.value)
  }

  const response = await fetch('/api/createGroup', {
    method: 'POST',
    body: groupData,
  })

  if (response.ok) {
    fetchGroups()
    groupName.value = ''
    groupDesc.value = ''
    imgFile.value = profileImg 
    imgPrv.value = profileImg
  }
}


const goToGroup = (groupId, groupName) => {
  console.log('Navigating to group:', groupName)
  router.push({ path: '/group', query: { id: groupId, name: groupName } })
}

async function fetchUserData() {
  const response = await fetch('/api/data', { method: 'POST' })
  if (response.ok) {
    userData.value = await response.json()
  }
}

function previewPostImage(e) {
  const file = e.target.files[0]
  if (file) {
    imgFile.value = file
    const reader = new FileReader()
    reader.onload = (e) => {
      imgPrv.value = e.target.result
    }
    reader.readAsDataURL(file)
  }
}
</script>

<template>
  <div class="main-container">
    <sidebar :username="userData?.Username" :pfp="userData?.pfp" />

    <div class="content2">
      <!-- Create Group Section -->
      <div class="post-card">
        <p class="post-content">Create a Group</p>
        <form @submit.prevent="createGroup">
          <input type="file" id="img" name="img" accept="image/*" @change="previewPostImage" />
      <label for="img">
        <img name="pfp" id="pfp" :src="imgPrv" alt="Profile Picture" />
      </label>
          <div class="input-box">
            <input
              v-model="groupName"
              type="text"
              class="area-name"
              placeholder="Group Name"
              required
            />
          </div>
          <div class="input-box">
            <textarea
              v-model="groupDesc"
              class="area-des"
              placeholder="Group Description"
              required
            ></textarea>
          </div>
            
          <div class="post-actions">
            <button class="send-btn" type="submit">
              <i class="fa-regular fa-paper-plane"></i>
            </button>
          </div>
        </form>
      </div>

      <!--  Groups Section (My Groups & All Groups) -->
      <div class="groups-container">
        <!--  My Groups -->
        <div class="group-card">
          <p class="group-title">My Groups</p>
          <ul v-if="myGroups.length" class="group-list">
            <li v-for="group in myGroups" :key="group.id" class="group-item"  @click="goToGroup(group.id, group.name)">
              <strong>
                  <img id="pfpPo" :src="group.pic" />
                </strong>
              <span class="group-name">{{
                group.name
              }}</span>
            </li>
          </ul>
          <p v-else class="empty-text">You haven't joined any groups yet.</p>
        </div>

        <!--  All Groups -->
        <div class="group-card">
          <p class="group-title">All Groups</p>
          <ul v-if="allGroups.length" class="group-list">
            <li v-for="group in allGroups" :key="group.id" class="group-item">
              <strong>
                  <img id="pfpPo" :src="group.pic" />
                </strong>
              <span class="group-name" @click="goToGroup(group.id, group.name)">{{
                group.name
              }}</span>
              <!-- <button class="join-btn" v-if="!isMember(group.id)" @click.stop="joinGroup(group.id)">Join</button> -->
            </li>
          </ul>
          <p v-else class="empty-text">No groups available.</p>
        </div>
      </div>
    </div>
    <sideUsers v-if="me" :me="me" />
  </div>
</template>

<style scoped></style>
