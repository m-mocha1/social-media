<script setup>
import { ref, defineProps, onMounted, watch, nextTick } from 'vue'
import { connectWebSocket, addListener } from '@/wepSock.js'
import sent from './sent.vue'
import received from './received.vue'
import router from '@/router'


const props = defineProps({
  user: Object,
})
console.log('🚀 ~ user:', props.user)
const text = ref('')
const messages = ref([])
const oldmessages = ref({})
const emit = defineEmits(['close'])

const chatMessagesRef = ref(null)

function closeChatBox() {
  emit('close')
}

async function getMessegs() {
  const formData = new FormData()
  formData.append('receiver', props.user.Username)
  const msgResp = await fetch('/api/getMsg', {
    method: 'POST',
    body: formData,
  })
  if (msgResp.ok) {
    oldmessages.value = await msgResp.json()
    console.log('🚀 ~ getMessegs ~ messages:', oldmessages.value)
    await nextTick() // Wait for Vue to update the DOM
    scrollToBottom()
  }
}

async function sendMsg() {
  if (text.value.length > 0) {
    const msgData = new FormData()
    msgData.append('rec', props.user.Username)
    msgData.append('text', text.value)
    const msgResp = await fetch('/api/msg', {
      method: 'POST',
      body: msgData,
    })
    if (msgResp.ok) {
      text.value = ''
    } else {
      console.log('bad resp')
    }
  }
}

onMounted(() => {
  getMessegs()
})
console.log('🚀 ~ addListener ~ props.user:', props.user)

addListener(async (data) => {
  if (data.type === 'msg') {
    console.log('🚀 ~ addListener ~ props.user:', props.user)
    console.log('🚀 ~ addListener ~ data.pfp:', data)
    if (data.to == props.user.me.me && data.Sender == props.user.Username) {
      messages.value.push({
        type: 'received',
        text: data.text,
        pfp: data.pfp,
      })
    } else if (data.Sender == props.user.me.me && data.to == props.user.Username) {
      messages.value.push({
        type: 'sent',
        text: data.text,
        pfp: props.user.me.mepfp,
      })
    }

    await nextTick()
    scrollToBottom()
  }
})
function scrollToBottom() {
  const chatBox = chatMessagesRef.value
  if (chatBox) {
    chatBox.scrollTop = chatBox.scrollHeight
  }
}
const myProfile = async () => {
  console.log("🚀 ~ myProfile ~ userData.value.Username:", props.user.Username)
  router.push({ path: '/profile', query: { username: props.user.Username} })
}
</script>

<template>
  <div class="chat-box">
    <div class="chat-header">
      <div id="res1" class="user">
        <span @click="myProfile" id="spP2">
          <img id="upfp" :src="user?.pfp" class="on2" />
          <p id="name">{{ user?.Username }}</p>
        </span>
      </div>
      <button @click="closeChatBox">
        <i class="fa-solid fa-xmark fa-sm" style="color: #ffffff"></i>
      </button>
    </div>
    <div class="chat-messages" ref="chatMessagesRef">
      <component
        v-for="(msg, index) in oldmessages.msgs"
        :key="index"
        :is="msg.Sender === user.me.me ? sent : received"
        :text="msg.text"
        :pfp="msg.Sender === user.me.me ? user.me.mepfp : msg.RecPfp"
      />
      <transition-group name="fade" tag="div">
        <component
          v-for="(msg, index) in messages"
          :key="'live-' + index"
          :is="msg.type === 'sent' ? sent : received"
          :text="msg.text"
          :pfp="msg.pfp"
        />
      </transition-group>
    </div>
    <div class="chat-input">
      <input v-model="text" type="text" placeholder="Type a message..." />
      <button @click="sendMsg">Send</button>
    </div>
  </div>
</template>

<style scoped>
/* Hide scrollbar for Webkit browsers */
.chat-messages::-webkit-scrollbar {
  display: none;
}

/* Hide scrollbar for IE, Edge and Firefox */
.chat-messages {
  -ms-overflow-style: none; /* IE and Edge */
  scrollbar-width: none; /* Firefox */
}

.on2 {
  transition: all 1s;
  border: 2px solid #4caf50;
  width: 40px;
  height: 40px;
  border-radius: 70%;
  object-fit: cover;
  cursor: pointer;
  margin-bottom: 5px;
}

#spP2 {
  align-items: center;
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  grid-template-rows: repeat(1, 1fr);
  gap: 8px;
}

#res1 {
  display: flex;
  flex-direction: row-reverse;
  align-content: center;
  justify-content: center;
  padding: 4%;
}

.chat-box {
  grid-row: span 2 / span 2;
  grid-column-start: 5;
  grid-row-start: 4;
  position: fixed;
  bottom: 10px;
  right: 10px;
  width: 300px;
  background-color: #1e1e1e;
  color: white;
  border: 1px solid #444;
  border-radius: 8px;
  display: flex;
  flex-direction: column;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
  height: 50vh;
}

.chat-header {
  background-color: #282828;
  padding: 1px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 1px solid #444;
}

.chat-messages {
  flex-grow: 1;
  padding: 10px;
  overflow-y: auto;
}

.chat-input {
  display: flex;
  padding: 10px;
  border-top: 1px solid #444;
}

.chat-input input {
  flex-grow: 1;
  padding: 5px;
  border: none;
  border-radius: 4px;
  margin-right: 5px;
}

.chat-input button {
  padding: 5px 10px;
  border: none;
  border-radius: 4px;
  background-color: #444;
  color: white;
  cursor: pointer;
}
</style>
