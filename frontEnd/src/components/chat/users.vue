<script setup>
const props = defineProps({
  user: {
    type: Object,
    required: true, // Ensure post is always provided
  },
})

</script>

<template>
  <div id="res" class="user">
    <span id="spPs">
      <img id="upfp" :src="user?.pfp" class="on" />
      <p id="name">{{ user?.Username }}</p>
    </span>
  </div>
</template>

<style>

</style>
