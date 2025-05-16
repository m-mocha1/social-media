<script setup>
// Define the post prop
import router from '@/router'
import { watchEffect, ref, onMounted, reactive, computed } from 'vue'
const props = defineProps({
  user: {
    type: Object,
  },
  post: {
    type: Object,
    required: true,
  },
})
const emit = defineEmits(['post-deleted'])
const post = reactive({ ...props.post })
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

const hideImage = (event) => {
  event.target.style.display = 'none' //
}

const toggleLike = async () => {
  const formData = new FormData()
  console.log(post.postLiked)
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

const clickedPost = async () => {
  router.push({ path: '/post', query: { id: props.post.postID, user: props.user.Username } })
}
const clickedProfile = async () => {
  router.push({ path: '/profile', query: { username: props.post.username } })
}
const deletePost = async () => {
  const post_id = new FormData()
  post_id.append('id', post.postID)
  const delResp = await fetch('/api/deletPost', {
    method: 'post',
    body: post_id,
  })
  if (delResp.ok) {
    console.log('post Deleted')
    emit('post-deleted', post.postID)
  } else {
    console.log('err post delet')
  }
}

const me = props.user.Username == post.username
</script>

<template>
  <span class="post" :pid="post?.postID">
    <div class="post-card">
      <div class="post-header">
        <div class="profile" @click="clickedProfile">
          <p>
            <strong>
              <!-- Use optional chaining for post properties -->
              <img id="pfpPo" :src="post?.userPfp" alt="Profile Picture" />
            </strong>
            <span id="spP">{{ post?.username }}</span>
            <small>{{ post?.time }} <i :class="icon"></i></small>
          </p>
        </div>
        <button v-if="me" @click="deletePost" class="trash">
          <i class="fa-solid fa-trash fa-xs" style="color: #fafafa"></i>
        </button>
      </div>

      <p @click="clickedPost">{{ post?.text }}</p>

      <img id="postImg" :src="post?.img" @error="hideImage" alt="Post Image" v-if="post?.img" />

      <button
        id="postlike"
        :postId="post?.postID"
        :class="['nolike', { liked: post.postLiked }]"
        @click="toggleLike"
      >
        <i class="fa-regular fa-heart">
          <span id="likeNum">{{ post?.likes }}</span>
        </i>
      </button>
      <br />
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
