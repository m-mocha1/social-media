<script setup>
</script>

<template>
  <button class="accept-btn" @click="acc">Accept</button>
</template>

<style scoped>
.accept-btn {
  background-color: #2e6b30;
  border: none;
  padding: 10px 18px;
  color: #ffffff;
  border-radius: 12px;
  cursor: pointer;
  font-size: 14px;
  font-weight: bold;
  transition:
    background-color 0.3s ease,
    transform 0.2s ease,
    box-shadow 0.2s ease;
}

.accept-btn:hover {
  background-color: #388e3c;
  transform: scale(1.07);
  box-shadow: 0px 0px 10px rgba(72, 239, 128, 0.7);
}

.accept-btn:active {
  transform: scale(1);
  background-color: #2e7d32;
}
</style>
