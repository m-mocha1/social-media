<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import sidebar from '@/components/side/sidebar.vue'
import router from '@/router'

const route = useRoute()
const groupID = route.query.id
const groupName = route.query.name
const IsAdmin = route.query.isAdmin === 'true' ? true : false
const groupRequests = ref([])
const groupMembers = ref([])
console.log("🚀 ~ groupMembers:", groupMembers)

const MutualFollowers = ref([])
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
}

async function fetchMutualFollowers() {
  try {
    const formData = new FormData()
    formData.append('groupId', groupID)
    const response = await fetch('/api/getMutualFollowers', {
      method: 'POST',
      body: formData,
    })

    if (!response.ok) throw new Error('Failed to load members')

    MutualFollowers.value = await response.json()
  } catch (error) {
    console.error('Error fetching members:', error.message)
  }
}
async function fetchGroupRequests() {
  try {
    const formData = new FormData()
    formData.append('groupId', groupID)

    const response = await fetch('/api/getGroupRequests', {
      method: 'POST',
      body: formData,
    })

    if (!response.ok) throw new Error('Failed to load requests')

    groupRequests.value = await response.json()
  } catch (error) {
    console.error('Error fetching requests:', error.message)
  }
}

async function fetchGroupMembers() {
  try {
    const formData = new FormData()
    formData.append('groupId', groupID)

    const response = await fetch('/api/getGroupMembers', {
      method: 'POST',
      body: formData,
    })

    if (!response.ok) throw new Error('Failed to load members')

    groupMembers.value = await response.json()
  } catch (error) {
    console.error('Error fetching members:', error.message)
  }
}

async function acceptRequest(requestId, userName) {
  try {
    const formData = new FormData()
    formData.append('groupId', groupID)
    formData.append('userName', userName)
    formData.append('requestId', requestId)
    formData.append('groupName', groupName)
    await fetch('/api/acceptInvite', { method: 'POST', body: formData })

    groupRequests.value = groupRequests.value.filter((request) => request.id !== requestId)

    await fetchGroupMembers()
  } catch (error) {
    console.error('Error accepting request:', error.message)
  }
}

async function rejectRequest(requestId, userName) {
  try {
    const formData = new FormData()
    formData.append('groupId', groupID)
    formData.append('userName', userName)
    formData.append('requestId', requestId)
    formData.append('groupName', groupName)

    await fetch('/api/rejectInvite', { method: 'POST', body: formData })

    groupRequests.value = groupRequests.value.filter((request) => request.id !== requestId)
  } catch (error) {
    console.error('Error rejecting request:', error.message)
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
  link.href = '/src/assets/mems.css'
  document.head.appendChild(link)
}

const GoToProfile = async (username) => {
  router.push({ path: '/profile', query: { username: username } })
}

async function inviteFriend(userName) {
  try {
    const formData = new FormData()
    formData.append('groupId', groupID)
    formData.append('userName', userName)
    formData.append('groupName', groupName)
    await fetch('/api/InviteFriend', { method: 'POST', body: formData })
    MutualFollowers.value = MutualFollowers.value.filter(
      (follower) => follower.username !== userName,
    )
    await fetchMutualFollowers()
  } catch (error) {
    console.error('Error accepting request:', error.message)
  }
}

onMounted(() => {
  auth()
  changeCss()
  fetchGroupRequests()
  fetchGroupMembers()
  fetchMutualFollowers()
})
</script>

<template>
  <div class="main-container">
    <sidebar />
    <div class="content">
      <h1 v-if="IsAdmin" class="group-title">{{ groupName }} Group Members & Requests</h1>
      <h1 v-else class="group-title">{{ groupName }} Group Members & Invites</h1>

      <div class="group-sections">
        <!-- Pending Requests Section -->
        <div v-if="IsAdmin" class="group-card">
          <p class="group-title">Pending Requests</p>
           <ul v-if="groupRequests?.length" class="group-list">
            <li v-for="request in groupRequests" :key="request.id" class="group-item">
             <div @click="GoToProfile(request.username)" style="align-items: center; display: flex; gap: 8px;">
                <strong>
                    <img id="pfpPo" :src="request.pfp" />
                  </strong>
                <span class="group-name"> {{ request.username }}</span>
              </div>
              <div>
                 <button @click="rejectRequest(request.id, request.username)" class="fow4">
                  <i class="fa-solid fa-user-plus">
                    <span> Accept </span>
                  </i>
                </button>
                <button @click="rejectRequest(request.id)" class="fow4">
                  <i class="fa fa-user-times">
                    <span> Reject </span>
                  </i>
                </button>
              </div>
            </li>
          </ul>
          <p v-else class="empty-text">No pending requests.</p>
        </div>
    <div class="group-card">
          <p class="group-title">Invite your friends</p>
          <ul v-if="MutualFollowers?.length" class="group-list">
            <li v-for="follower in MutualFollowers" :key="follower.id" class="group-item">
              <div @click="GoToProfile(follower.username)" style="align-items: center; display: flex; gap: 8px;">
                <strong>
                    <img id="pfpPo" :src="follower.pfp" />
                  </strong>
                <span class="group-name"> {{ follower.username }}</span>
              </div>
              <button @click="inviteFriend(follower.username)" class="fow3">
                  <i class="fa-solid fa-plus">
                    <span>
                       Invite
                    </span>
                   </i>
                </button>
            </li>
          </ul>
          <p v-else class="empty-text">No members yet.</p>
        </div>
        <!-- Group Members Section -->
        <div class="group-card">
          <p class="group-title">Members</p>
       <ul v-if="groupMembers?.length" class="group-list">
            <li v-for="member in groupMembers" :key="member.id" class="group-item" @click="GoToProfile(member.username)">
              <div @click="GoToProfile(member.username)" style="align-items: center; display: flex; gap: 8px;">
                <strong>
                    <img id="pfpPo" :src="member.pfp" />
                  </strong>
                <span class="group-name"> {{ member.username }}</span>
              </div>
            </li>
          </ul>
          <p v-else class="empty-text">No members yet.</p>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* * {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
  font-family: 'Arial', sans-serif;
}

body {
  background-color: #0f0f0f;
  min-height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
}

.main-container {
  display: flex;
  width: 90%;
  max-width: 1200px;
  margin: 20px auto;
}

.content {
  flex-grow: 1;
  background: #1e1e1e;
  border-radius: 10px;
  padding: 20px;
  box-shadow: 0px 4px 8px rgba(0, 0, 0, 0.2);
  color: white;
  text-align: center;
}

.group-sections {
  display: flex;
  flex-direction: column;
  flex-wrap: wrap;
  align-items: center;
}

.group-card {
  background: #282828;
  border-radius: 8px;
  padding: 20px;
  width: 70%;
  min-width: 280px;
  margin-bottom: 20px;

}

.group-title {
  font-size: 10;
  color: #23edc9;
  margin-bottom: 15px;
}

.group-list {
  list-style: none;
  align-items: baseline;
}

.group-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: #2a2a2a;
  padding: 10px;
  border-radius: 5px;
  margin-bottom: 8px;
}

.empty-text {
  color: #888;
  font-style: italic;
} */
</style>
