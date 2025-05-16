<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import sidebar from '@/components/side/sidebar.vue'
import router from '@/router'
import PostCard from '@/components/singlePost/post.vue'
import { connectWebSocket, addListener } from '@/wepSock.js'
import sideUsers from '@/components/chat/sideUsers.vue'

const route = useRoute()
const groupID = route.query.id
const groupName = route.query.name
const groupDetails = ref(null)
const userData = ref(null)
const postArea = ref(null)
const posts = ref([])
const events = ref([])
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

addListener((data) => {
  if (data.type === 'post') {
    if (!Array.isArray(posts.value)) {
      posts.value = []
    }
    posts.value.unshift(data)
  }
})

const fetchGroupDetails = async () => {
  if (!groupID) {
    console.error('Error: Group ID is missing from query parameters.')
    return
  }

  console.log('Fetching details for group ID:', groupID)

  try {
    let formData = new FormData()
    formData.append('groupId', groupID)
    formData.append('groupName', groupName)

    const response = await fetch('/api/getGroup', {
      method: 'POST',
      body: formData,
    })

    if (!response.ok) throw new Error('Failed to load group details')

    groupDetails.value = await response.json()
    console.log('Fetched group details:', groupDetails.value)
  } catch (error) {
    console.error('Error fetching group details:', error.message)
  }
}

async function joinGroup() {
  const joinData = new FormData()
  joinData.append('groupId', groupID)
  joinData.append('groupName', groupName)
  if (!groupDetails.value.IsMem) {
    if (!groupDetails.value.IsInv) {
      console.log('join')

      await fetch('api/joinGroup', {
        method: 'POST',
        body: joinData,
      })
      groupDetails.value.IsInv = true
    } else {
      console.log('not join')

      await fetch('api/notJoinGroup', {
        method: 'POST',
        body: joinData,
      })
      groupDetails.value.IsInv = false
    }
  } else {
    console.log('groupid: ', groupID)
    await fetch('api/outOfGroup', {
      method: 'POST',
      body: joinData,
    })
    groupDetails.value.IsMem = false
    groupDetails.value.IsInv = false
    groupDetails.value.NumOfMems--
    console.log('out of group')
  }
}

const goToGroupMembers = (groupId, groupName, IsAdmin) => {
  console.log('Navigating to group:', groupName)
  router.push({ path: '/groupMembers', query: { id: groupId, name: groupName, isAdmin: IsAdmin } })
}

let imgPrv = ref('')
let imgFile = ref('')
let postText = ref('')

const autoResizeTextarea = () => {
  if (postArea.value) {
    postArea.value.style.width = 'auto'
    postArea.value.style.height = 'auto'
    postArea.value.style.height = `${postArea.value.scrollHeight}px`
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
async function submitPost() {
  const postData = new FormData()
  postData.append('postText', postText.value)
  if (imgFile.value != null) {
    postData.append('postImg', imgFile.value)
  }
  postData.append('groupId', groupID)
  const resp = await fetch('/api/posting', {
    method: 'POST',
    body: postData,
  })
  if (resp.ok) {
    console.log('ok')
    postText.value = ''
    imgPrv.value = ''
    imgFile.value = null
  } else {
    console.error('faild')
  }
}
function postDel() {
  imgPrv.value = ''
  imgFile.value = null
}
async function fetchPosts() {
  const Data = new FormData()
  Data.append('groupId', groupID)
  const postResp = await fetch('/api/getGroupPosts', {
    method: 'POST',
    body: Data,
  })
  if (postResp.ok) {
    posts.value = await postResp.json()
    console.log('Fetched posts:', posts.value)
  } else {
    console.error('Error fetching posts')
  }
}

const myProfile = async () => {
  router.push({ path: '/profile', query: { username: userData.value.Username } })
}
const myProfile2 = async (user) => {
  router.push({ path: '/profile', query: { username: user } })
}
const isPosts = ref(true)
const isEvents = ref(false)

function toggle(showPosts) {
  if (showPosts) {
    isPosts.value = true
    isEvents.value = false
  } else {
    isPosts.value = false
    isEvents.value = true
  }
}

const eventTitle = ref('')
const eventDateTime = ref('')
const eventDescription = ref('')

async function createEvent() {
  const formData = new FormData()
  formData.append('title', eventTitle.value)
  formData.append('dateTime', eventDateTime.value)
  formData.append('description', eventDescription.value)
  formData.append('groupId', groupID)

  try {
    const response = await fetch('/api/createEvent', {
      method: 'POST',
      body: formData,
    })

    if (response.ok) {
      console.log('Event created successfully!')
      eventTitle.value = ''
      eventDateTime.value = ''
      eventDescription.value = ''
    } else {
      console.error('Failed to create event')
    }
  } catch (error) {
    console.error('Error:', error)
  }
  fetchEvents()
}
async function fetchEvents() {
  const response = await fetch(`/api/getEvents?groupId=${groupID}`)
  if (response.ok) {
    events.value = await response.json()
    console.log('1', events.value)
    console.log('2', events)
  } else {
    console.error('Failed to fetch events')
  }
}

async function respondToEvent(eventId, response) {
  const formData = new FormData()
  formData.append('eventId', eventId)
  formData.append('response', response)

  const resp = await fetch('/api/respondToEvent', {
    method: 'POST',
    body: formData,
  })

  if (resp.ok) {
    console.log('Response recorded successfully')
  } else {
    console.error('Failed to respond to event')
  }
  fetchEvents()
}
const deleteEvent = async (eventID) => {
  const event_id = new FormData()
  event_id.append('eventID', eventID)
  const delResp = await fetch('/api/deletEvent', {
    method: 'post',
    body: event_id,
  })
  if (delResp.ok) {
    fetchEvents()
  } else {
    console.log('err')
  }
}

const formatDateTime = (isoString) => {
  console.log('is the time', isoString)
  if (!isoString) return ''
  const date = new Date(isoString)
  let month = date.getMonth()
  let day = date.getDate()
  let year = date.getFullYear()
  let hours = date.getHours() - 2
  let minutes = date.getMinutes()
  minutes = minutes < 10 ? '0' + minutes : minutes
  let ampm = hours >= 12 ? 'PM' : 'AM'
  hours = hours % 12 || 12
  console.log(`${month}/${day}/${year}, ${hours}:${minutes} ${ampm}`)
  return `${month}/${day}/${year}, ${hours}:${minutes} ${ampm}`
}

onMounted(() => {
  auth()
  changeCss()
  autoResizeTextarea()
  if (postArea.value) {
    postArea.value.addEventListener('input', autoResizeTextarea)
  }
  fetchGroupDetails()
  fetchPosts()
  fetchEvents()
})

const removePost = (postId) => {
  posts.value = posts.value.filter((post) => post.postID !== postId)
}
</script>

<template>
  <div class="main-container">
    <sidebar />
    <div class="content">
      <div class="welcome2">
        <div v-if="groupDetails" class="options2">
          <strong>
            <img id="pfpPo" :src="groupDetails.pic" />
          </strong>
          <div class="group-header">
            <h1 class="group-title">{{ groupDetails.name }}</h1>
            <p class="group-admin"><strong>Admin:</strong> {{ groupDetails.admin }}</p>
          </div>
          <div class="group-content">
            <p class="group-description">
              <strong>Description:</strong> {{ groupDetails.description }}
            </p>
            <p class="group-meta">
              <strong>Created at:</strong> {{ groupDetails.created_at.split('T')[0] }}
            </p>
          </div>
          <br />
          <div class="proCont">
            <div>
              <button
                @click="goToGroupMembers(groupID, groupDetails.name, groupDetails.IsAdmin)"
                :disabled="!groupDetails?.IsMem"
                class="fow2"
              >
                <i class="fa fa-users">
                  <span>{{ ' ' + groupDetails.NumOfMems }}</span>
                </i>
              </button>
            </div>

            <button
              v-if="!groupDetails?.IsMem && !groupDetails?.IsInv"
              @click="joinGroup"
              class="fow"
            >
              <i class="fa-solid fa-user-plus"></i>
            </button>

            <button
              v-else-if="!groupDetails?.IsMem && groupDetails?.IsInv"
              @click="joinGroup"
              class="fow"
            >
              <i class="fa fa-user-times"></i>
            </button>

            <button v-else-if="groupDetails?.IsMem" @click="joinGroup" class="fow">
              <i class="fa-solid fa-circle-user"></i>
            </button>
          </div>
        </div>
      </div>
      <div class="welcome3">
        <div>
          <button @click="toggle(true)" :disabled="!groupDetails?.IsMem" class="fow2">
            <i class="fa-solid fa-signs-post">
              <span> Posts</span>
            </i>
          </button>
        </div>
        <div>
          <button @click="toggle(false)" :disabled="!groupDetails?.IsMem" class="fow2">
            <i class="fa fa-calendar">
              <span> Events</span>
            </i>
          </button>
        </div>
      </div>
      <div v-if="groupDetails?.IsMem && isPosts" class="welcome">
        <div class="options">
          <img v-if="userData" :src="userData?.pfp" id="pfpPro" @click="myProfile" />

          <p @click="myProfile" id="username" v-if="userData">{{ userData.Username }}</p>
          <form @submit.prevent="submitPost" enctype="multipart/form-data">
            <textarea
              ref="postArea"
              v-model="postText"
              id="postArea"
              name="text"
              required
              placeholder="Create a post"
              class="txt"
            ></textarea>
            <br />
            <br />
            <br />
            <label for="Postimg">
              <img :src="imgPrv" v-if="imgPrv" id="postpic" />
            </label>
            <p id="err"><br /></p>

            <ul>
              <li>
                <button type="submit" id="post"><i class="fa-regular fa-paper-plane"></i></button>

                <input
                  type="file"
                  id="Pimg"
                  name="Postimg"
                  accept="image/*"
                  style="display: none"
                  @change="previewPostImage"
                />

                <label for="Pimg" style="cursor: pointer">
                  <i class="fa-solid fa-image"></i>
                </label>

                <label class="cancel" @click="postDel">
                  <i class="fa-solid fa-x"></i>
                </label>
              </li>
            </ul>
          </form>
        </div>
      </div>
      <div v-else-if="groupDetails?.IsMem && isEvents" class="welcome">
        <div class="post-card">
          <p class="post-content">Create a Event</p>
          <form @submit.prevent="createEvent">
            <div class="input-box">
              <input
                class="area-tit2"
                type="text"
                v-model="eventTitle"
                placeholder="Event Title"
                required
              />
              <input class="area-tit3" type="datetime-local" v-model="eventDateTime" required />
            </div>
            <div class="input-box">
              <textarea
                class="area-des2"
                v-model="eventDescription"
                placeholder="Event Description"
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
      </div>
      <div v-else class="welcome">
        <p class="group-admin"><strong>Join our group to share with us and see posts.</strong></p>
      </div>

      <div
        v-if="groupDetails?.IsMem && !isPosts"
        v-for="event in events"
        :key="event.id"
        class="event-card"
      >
        <span class="post" :eid="event.id">
          <div class="post-card">
            <div @click="myProfile2(event.creator)" class="add">
              <p class="space">
                <strong>
                  <img
                    id="pfpPo"
                    :src="event?.pic || '/default-avatar.png'"
                    alt="Creator Profile"
                  />
                </strong>
                <span id="spP">{{ event?.creator }}</span>
                <small>{{ new Date(event.createdAt).toLocaleDateString() }}</small>
              </p>
            </div>
            <button v-if="event.me" @click="deleteEvent(event.id)" class="trash">
              <i class="fa-solid fa-trash fa-xs" style="color: #fafafa"></i>
            </button>

            <p class="group-admin"><strong>Title:&nbsp;</strong> {{ event.title }}</p>
            <p class="group-admin"><strong>Description:&nbsp;</strong> {{ event.description }}</p>
            <p class="group-admin">
              <strong>Date & Time:&nbsp;</strong>
              {{ formatDateTime(event.eventDatetime) }}
            </p>

            <div class="post-actions">
              <button
                id="postlike"
                :eid="event.id"
                :class="['nolike', { liked: event.userResponse === 'Going' }]"
                @click="respondToEvent(event.id, 'Going') && event.goingCount++"
              >
                <!-- <i class="fa-solid fa-check-circle"></i>
        <span>Going ({{ event.goingCount }})</span> -->
                <i class="fa fa-calendar-check-o">
                  <span>&nbsp;{{ event.goingCount }}&nbsp;Going</span>
                </i>
              </button>

              <button
                id="postlike"
                :eid="event.id"
                :class="['nolike', { liked: event.userResponse === 'Not going' }]"
                @click="respondToEvent(event.id, 'Not going')"
              >
                <i class="fa fa-calendar-times-o"></i>
                <span>&nbsp;{{ event.notGoingCount }}&nbsp;Not Going</span>
              </button>
            </div>

            <br />
            <hr />
          </div>
        </span>
      </div>
      <div id="postFeed" class="post-feed">
        <!-- Loop over posts and display each -->
        <transition-group
          v-if="groupDetails?.IsMem && isPosts"
          name="fade"
          tag="div"
          id="postFeed"
          class="post-feed"
        >
          <PostCard
            v-for="post in posts"
            :key="post.postID"
            v-if="userData"
            :user="userData"
            :post="post"
            @post-deleted="removePost"
          />
        </transition-group>
      </div>
      <!-- <div v-if="!groupDetails?.IsMem" id="postFeed" class="post-feed"></div> -->
    </div>
    <sideUsers v-if="me" :me="me" />
  </div>
</template>

<style scoped>
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
