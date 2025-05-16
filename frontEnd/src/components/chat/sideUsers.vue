<script setup>
import { ref, defineProps, onMounted } from 'vue'
import users from '@/components/chat/users.vue'
import ChatBox from '@/components/chat/chatBox.vue' // make sure ChatBox.vue exists

const props = defineProps({
  me: Object,
})
console.log('🚀 ~ user:', props.user)

const allUsers = ref([])

const selectedUser = ref(null)

async function getUsers() {
  try {
    const usersResp = await fetch('/api/allUsers', {
      method: 'POSt',
    })
    if (usersResp.ok) {
      allUsers.value = await usersResp.json()
      console.log('🚀 ~ getUsers ~ allUsers:', allUsers)
    } else {
      console.log('no users')
    }
  } catch {}
}
onMounted(() => {
  getUsers()
})

function openChatBox(user) {
  if (!selectedUser.value) {
    selectedUser.value = { Username: user.Username, pfp: user.pfp, me: props.me }
  } else {
    selectedUser.value = null
  }

  console.log('🚀 ~ openChatBox ~ selectedUser:', selectedUser.value.Username)
}
</script>

<template>
  <div id="side" class="sideUser">
    <users v-for="u in allUsers" :key="u.id" :user="u" @click="openChatBox(u)" />
  </div>
  <transition-group name="fade" tag="div">
    <ChatBox v-if="selectedUser" :user="selectedUser" @close="selectedUser = null" />
  </transition-group>
</template>

<style scoped></style>
