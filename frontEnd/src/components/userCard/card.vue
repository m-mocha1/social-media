<script setup>
import router from '@/router'

const props = defineProps({
  user: {
    type: Object,
  },
})
console.log('🚀 ~ props:', props.user)

const clickedProfile = async () => {
  router.push({ path: '/profile', query: { username: props.user.username } })
}
</script>

<template>
  <span class="post">
    <div class="post-card">
        <div class="profile" @click="clickedProfile">
          <p>
            <strong>
              <!-- Use optional chaining for post properties -->
              <img id="pfpPo" :src="user?.pfp" alt="Profile Picture" />
            </strong>
            <span id="spP">{{ user?.username }}</span>
          </p>
        </div>
      </div>
  </span>
</template>

<style scoped>
.post-card{
    width: fit-content;
}
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
