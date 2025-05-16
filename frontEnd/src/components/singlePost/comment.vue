<script setup>
import { ref, defineProps, computed } from 'vue'

const props = defineProps({
    user: {
    type: String,
  },
  com: {
    type: Object,
    required: true,
  },
})
console.log("🚀 ~ props c:", props.com)
const emit = defineEmits(["com-deleted"]);

const comLiked = ref(props.com.comLiked)
const likes = ref(props.com.likes)
const buttonClass = computed(() => (comLiked.value ? 'likedC' : 'nolikeC'))
const handleComLike = async () => {
  console.log('in like')

  const formData = new FormData()
  formData.append('comId', props.com.comId)
  console.log('com liked', com.comLiked)

  if (comLiked.value) {
    const resp = await fetch('/api/dislikeCom', {
      method: 'POST',
      body: formData,
    })
    if (resp.ok) {
      console.log('ok')

      comLiked.value = false
      likes.value--
    } else {
      console.error('unlike err')
    }
  } else {
    const resp = await fetch('/api/likeCom', {
      method: 'POST',
      body: formData,
    })
    if (resp.ok) {
      comLiked.value = true
      likes.value++
    } else {
      console.error('like err')
    }
  }
}
const hideImage = (event) => {
  event.target.style.display = 'none' //
}
const deleteCom = async () => {

  const com_id = new FormData()
  com_id.append('id', props.com.comId)
  const delResp = await fetch('/api/deletCom', {
    method: 'post',
    body: com_id,
  })
  if (delResp.ok) {
    console.log('com Deleted')
    emit("com-deleted", props.com.comId);
  } else {
    console.log('err com delet')
  }
}
 const me = props.user == props.com.username
</script>

<template>
  <span>
    <hr />
    <br />
    <div class="com-header">
      <div class="profile" @click="clickedProfile">
        <img :src="com.img" alt="Profile picture" id="pfpCo" />
        <p id="spC">{{ com.username }}</p>
        <small>{{ com.time }}</small>
      </div>
      <button v-if="me" @click="deleteCom" class="trash">
        <i class="fa-solid fa-trash fa-xs" style="color: #fafafa"></i>
      </button>
    </div>
    <br />
    <p class="comm">{{ com.text }}</p>
    <img id="postImg" :src="com?.comPic" @error="hideImage" alt="com Image" v-if="com?.comPic" />
    <button @click="handleComLike" id="comlike" :data-com-id="com.comId" :class="buttonClass">
      <i class="fa-regular fa-heart"></i>
      <span class="likes">{{ likes }}</span>
    </button>
    <br />
  </span>
</template>
<style>
.com-header {
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
.com-header:hover .trash {
  opacity: 1;
}
</style>
