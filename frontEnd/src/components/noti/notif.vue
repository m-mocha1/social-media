<script setup>
import router from '@/router'
import { ref } from 'vue'
import accept from '@/components/noti/buttons/accept.vue'
import reject from '@/components/noti/buttons/reject.vue'

const emit = defineEmits(['notification-deleted'])
const props = defineProps({
  noti: {
    type: Object,
    required: true,
  },
})

console.log('🚀 ~ clickedProfile ~  props.noti:', props.noti)
const clickedProfile = async () => {
  router.push({ path: '/profile', query: { username: props.noti.sender } })
}
async function acceptRe() {
  const followData = new FormData()
  followData.append('id', props.noti.id)
  followData.append('fing', props.noti.sender)
  const acc = await fetch('/api/followerAcc', {
    method: 'POST',
    body: followData,
  })
  if (acc.ok) {
    console.log('🚀 ~ acceptRe ~ followData.ok:', followData.ok)
    emit('notification-deleted', props.noti.id)
  }
}
async function rejectRe() {
  const followData = new FormData()
  followData.append('id', props.noti.id)
  followData.append('fing', props.noti.sender)
  const acc = await fetch('/api/followerRej', {
    method: 'POST',
    body: followData,
  })
  if (acc.ok) {
    emit('notification-deleted', props.noti.id)
  }
}

async function acceptInvFriend() {
  const inviteData = new FormData()
  inviteData.append('id', props.noti.id)
  inviteData.append('inviter', props.noti.sender)
  inviteData.append('groupId', props.noti.groupId)
  inviteData.append('groupName', props.noti.groupName)

  const acc = await fetch('api/inviteFrAcc', {
    method: 'POST',
    body: inviteData,
  })
  if (acc.ok) {
    emit('notification-deleted', props.noti.id)
  }
}
async function rejectInvFriend() {
  const inviteData = new FormData()
  inviteData.append('id', props.noti.id)
  inviteData.append('inviter', props.noti.sender)
  inviteData.append('groupId', props.noti.groupId)
  inviteData.append('groupName', props.noti.groupName)
  const acc = await fetch('api/inviteFrRej', {
    method: 'POST',
    body: inviteData,
  })
  if (acc.ok) {
    emit('notification-deleted', props.noti.id)
  }
}

const goToGroupMembers = () => {
  router.push({
    path: '/groupMembers',
    query: { id: props.noti.groupId, name: props.noti.groupName, isAdmin: true },
  })
}
function retext() {
  const parts = props.noti.text.split('groupId')
  const groupName = parts[0].replace('requested to join ', '').trim()
  return 'requested to join ' + groupName + ' group'
}

async function deleteNot() {
  const noID = new FormData()
  noID.append('id', props.noti.id)
  const de = await fetch('/api/deleteNot', {
    method: 'POST',
    body: noID,
  })
  if (de.ok) {
    emit('notification-deleted', props.noti.id)
  }
}
const clickedPost = async () => {
  router.push({ path: '/post', query: { id: props.noti.actId } })
}
</script>

<template>
  <span class="post" :pid="noti?.id">
    <div class="post-card">
      <div class="post-header">
        <div class="profile" @click="clickedProfile">
          <p>
            <strong>
              <!-- Use optional chaining for post properties -->
              <img id="pfpPo" :src="noti?.sender_pfp" alt="Profile Picture" />
            </strong>
            <span id="spP">{{ noti?.sender }}</span>
            <small>{{ noti?.time }}</small>
          </p>
        </div>
        <button @click="deleteNot" class="trash">
          <i class="fa-solid fa-trash fa-xs" style="color: #fafafa"></i>
        </button>
      </div>
      <p @click="clickedPost">{{ noti?.text }}</p>
      <br />
      <div v-if="noti?.type == 'followreq'">
        <accept @click="acceptRe" /> <reject @click="rejectRe" />
      </div>
      <div v-if="noti?.type == 'join group'">
        <button @click="goToGroupMembers()">Go to the group to mange it</button>
      </div>
      <div v-if="noti?.type == 'invite group'">
        <accept @click="acceptInvFriend" /> <reject @click="rejectInvFriend" />
      </div>
      <hr />
    </div>
  </span>
</template>

<style>
.post-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.profile {
  cursor: pointer;
}

.trash {
  background: transparent;
  border: none;
  cursor: pointer;
  opacity: 0;
  transition: opacity 0.3s ease-in-out;
}
.post-card:hover .trash {
  opacity: 1;
}
</style>
