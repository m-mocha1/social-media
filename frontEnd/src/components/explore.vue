<script setup>
import { ref, onMounted } from 'vue'
import router from '@/router'
import PostCard from '@/components/singlePost/post.vue'
import { connectWebSocket, addListener } from '../wepSock.js'
import sidebar from '@/components/side/sidebar.vue'
import Select from './postType/select.vue'
import sideUsers from '@/components/chat/sideUsers.vue'

const userData = ref(null)
const postArea = ref(null)
const posts = ref([])
const selectedOption = ref({ value: 'Pub', icon: 'fa-solid fa-earth-americas' })
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
  if (postArea.value) {
    postArea.value.addEventListener('input', autoResizeTextarea)
  }
  fetchPosts()
})

// addListener((data) => {
//   if (data.type === 'post') {
//     if (!Array.isArray(posts.value)) {
//       posts.value = []
//     }
//     posts.value.unshift(data)
//   }
// })

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
  postData.append('groupId', null)
  postData.append('postType', selectedOption.value.value)
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
  const postResp = await fetch('/api/getAllPost')
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
const removePost = (postId) => {
  posts.value = posts.value.filter((post) => post.postID !== postId)
}

</script>

<template v-cloak>
  <div class="main-container">
    <sidebar :username="userData?.Username" :pfp="userData?.pfp" />

    <div class="content">
      <div class="welcome">
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
                <label>
                  <Select v-model="selectedOption" />
                </label>
              </li>
            </ul>
          </form>
        </div>
      </div>
      <hr />
          <transition-group name="fade" tag="div"  id="postFeed" class="post-feed">
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
    <sideUsers v-if="me" :me="me" />

  </div>
  <!-- end of post feed -->
  <!-- end of content -->
  <!-- end of main container -->
</template>

<style>
body {
  background-color: black;
}
[v-cloak] {
  display: none;
}
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.5s ease, transform 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}
</style>
