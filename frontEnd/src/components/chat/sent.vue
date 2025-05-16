<script setup>
import { defineProps } from 'vue'

const props = defineProps({
  text: {
    type: String,
    required: true,
  },
  pfp: {
    type: String,
    required: true,
  },
})
</script>

<template>
  <div class="sent-message">
    <div class="message-text">{{ text }}</div>
    <img :src="pfp" alt="User Avatar" class="avatar" />
  </div>
</template>

<style scoped>
.sent-message {
  display: flex;
  align-items: center;
  background-color: #333333; /* darker grey background */
  color: white;
  padding: 8px;
  border-radius: 8px;
  margin: 5px 0;
  width: fit-content;
  margin-left: auto; /* push message to right */
}

.avatar {
  width: 30px;
  height: 30px;
  border-radius: 50%;
  object-fit: cover;
  margin-left: 8px;
}

.message-text {
  display: inline-block;
  text-align: right;
}
</style>
