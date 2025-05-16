<script setup>
const emit = defineEmits(['click'])
</script>

<template>
  <button class="reject-btn" @click="$emit('click')">Reject</button>
</template>

<style scoped>
.reject-btn {
  background-color: #f8d7da;
  border: none;
  padding: 10px 18px;
  color: #721c24;
  border-radius: 12px; /* Increased border-radius for a more modern look */
  cursor: pointer;
  font-size: 14px;
  font-weight: bold;
  transition:
    background-color 0.3s ease,
    transform 0.2s ease,
    box-shadow 0.2s ease;
}

.reject-btn:hover {
  background-color: #f1b0b7;
  transform: scale(1.07); /* Slight increase for better interaction */
  box-shadow: 0px 0px 10px rgba(245, 54, 92, 0.5); /* Subtle red glow effect */
}

.reject-btn:active {
  transform: scale(1); /* Prevent excessive scaling */
  background-color: #d9534f; /* A deeper red when clicked */
  color: white; /* Ensures readability */
}
</style>
