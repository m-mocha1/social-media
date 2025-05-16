<script>
export default {
  props: {
    modelValue: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      isOpen: false,
      options: [
        { value: 'Pub', icon: 'fa-solid fa-earth-americas' },
        { value: 'Fol', icon: 'fa-solid fa-user-group' },
        { value: 'Pri', icon: 'fa-solid fa-lock' },
      ],
    }
  },
  methods: {
    toggleDropdown() {
      this.isOpen = !this.isOpen
    },
    selectOption(option) {
      this.$emit('update:modelValue', option)
      this.isOpen = false
    },
  },
}
</script>

<template>
  <div class="custom-select">
    <div class="selected-option" @click="toggleDropdown">
      <i :class="modelValue.icon"></i>
    </div>
    <div class="dropdown" v-if="isOpen">
      <div
        v-for="option in options"
        :key="option.value"
        @click="selectOption(option)"
        class="dropdown-item"
      >
        <i :class="option.icon"></i>
      </div>
    </div>
  </div>
</template>

<style scoped>
.custom-select {
  position: relative;
  display: inline-block;
  cursor: pointer;
}
.selected-option {
  background: transparent;
  padding: 5px;
  border-radius: 5px;
}
.dropdown {
  position: absolute;
  background: black;
  border-radius: 5px;
  padding: 5px;
}
.dropdown-item {
  padding: 5px;
  cursor: pointer;
}
</style>
