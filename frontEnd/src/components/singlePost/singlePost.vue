<script setup>
import { ref, defineProps, reactive, computed } from 'vue'
import com from '@/components/singlePost/comment.vue'
import { addListener } from '@/wepSock'
// Define props to pass post data dynamically
const props = defineProps({
  post: Object,
  user: String,
})
const icon = computed(() => {
  switch (post.postType) {
    case 'Pub':
      return 'fa-solid fa-earth-americas'
    case 'Fol':
      return 'fa-solid fa-user-group'
    case 'Pri':
      return 'fa-solid fa-lock'
    default:
      return 'fa-solid fa-earth-americas'
  }
})
console.log('🚀 ~ props:', props.user)

const post = reactive({ ...props.post })
let comment = ref('')
const comments = ref(Array.isArray(post.Comment) ? post.Comment : [])
let imgPrv = ref('')
let imgFile = ref('')

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
const handleLike = async () => {
  const formData = new FormData()
  formData.append('postID', props.post.postID)
  if (post.postLiked) {
    const resp = await fetch('/api/unlike', {
      method: 'POST',
      body: formData,
    })
    if (resp.ok) {
      console.log('unlike')
      post.postLiked = false
      post.likes--
    } else {
      console.error('unlike err')
    }
  } else {
    const resp = await fetch('/api/like', {
      method: 'POST',
      body: formData,
    })
    if (resp.ok) {
      console.log('liked')
      post.postLiked = true
      post.likes++
    } else {
      console.error('like err')
    }
  }
}

const handleComment = async () => {
  const id = props.post.postID
  console.log(`Liked post ID: ${id}`)

  const comData = new FormData()
  comData.append('comment', comment.value)
  comData.append('postId', id)
  if (imgFile.value) {
    comData.append('comImg', imgFile.value) // Append the image file
  }
  try {
    const c = await fetch('/api/com', {
      method: 'POST',
      body: comData,
    })
    if (c.ok) {
      comment.value = ''
      imgPrv.value = ''
      imgFile.value = null
    } else {
      console.log('err')
    }
  } catch {
    console.error('Error fetching com:', error)
  }
}
addListener((data) => {
  if (data.type === 'comment') {
    if (!Array.isArray(comments.value)) {
      comments.value = []
    }
    console.log(data.data)
    comments.value.unshift(data.data)
  }
})
function postDel() {
  imgPrv.value = ''
  imgFile.value = null
}
const removeCom = (comId) => {
  console.log('🔹 Deleting c ID:', comId)
  comments.value = comments.value.filter((c) => c.comId !== comId)
}
</script>

<template>
  <div v-if="post" class="post-card-single">
    <p>
      <strong>
        <img id="pfpPo" :src="post.userPfp" alt="User Profile Picture" />
        <span id="spP">{{ post.username }}</span>
      </strong>
      <small>{{ post.time }} <i :class="icon"></i></small>
    </p>

    <br />
    <p id="text">{{ post.text }}</p>
    <br />

    <!-- Post Image (hides if not found) -->
    <img
      v-if="post.img"
      id="postImg"
      :src="post.img"
      @error="(event) => (event.target.style.display = 'none')"
    />

    <br />
    <!-- Like Button -->
    <button id="postlike" :class="['nolike', { liked: post.postLiked }]" @click="handleLike">
      <i class="fa-regular fa-heart">
        <span id="likeNum">{{ post.likes }}</span>
      </i>
    </button>

    <br />
    <div id="com" class="comment">
      <form class="com-form" @submit.prevent="handleComment" enctype="multipart/form-data">
        <textarea
          id="comment"
          v-model="comment"
          required
          placeholder="Comment"
          class="comment"
        ></textarea>
        <br />
        <br />
        <br />
        <label for="Cimg">
          <img :src="imgPrv" v-if="imgPrv" id="compic" />
        </label>
        <br /><br />
        <button type="submit" class="cb" id="coms">
          <i class="fa-regular fa-comment"></i>
        </button>
        <input
          type="file"
          id="Cimg"
          name="Cimg"
          style="display: none"
          accept="image/*"
          @change="previewPostImage"
        />

        <label for="Cimg" style="cursor: pointer">
          <i class="fa-solid fa-image fa-lg" style="color: #ffffff"></i>
        </label>

        <label class="cancel" @click="postDel">
          <i class="fa-solid fa-x" style="color: #ffffff"></i>
        </label>
      </form>

      <br />
      <transition-group name="fade" tag="div" id="comCont">
          <com
          v-for="c in comments"
          :key="c.comId"
          :com="c"
          :user="user"
          @com-deleted="removeCom"
          ></com>
        </transition-group>
    </div>
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
