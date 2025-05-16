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
  <div class="received-message">
    <img :src="pfp" alt="User Avatar" class="avatar" />
    <div class="message-text">{{ text }}</div>
  </div>
</template>

<style scoped>
.received-message {
  display: flex;
  align-items: center;
  background-color: #808080; /* gray background */
  color: white;
  padding: 8px;
  border-radius: 8px;
  margin: 5px 0;
  width: fit-content;
  margin-right: auto;
}

.avatar {
  width: 30px;
  height: 30px;
  border-radius: 50%;
  object-fit: cover;
  margin-right: 8px;
}

.message-text {
  display: inline-block;
  text-align: right;
}
</style>
