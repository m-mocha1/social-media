<script setup>
import router from '@/router'
import { ref } from 'vue'
import { onMounted } from 'vue'
const isError = ref(false)

async function auth() {
  try {
    const AuthResp = await fetch('/api/Auth', {
      method: 'POST',
    })
    if (AuthResp.ok) {
      router.push('/')
    }
  } catch (error) {
    console.error('Fetch error:', error)
  }
}
function changeCss() {
  const oldcss = document.getElementById('dc')

  if (oldcss) {
    console.log(oldcss)

    oldcss.remove()
  }
  const link = document.createElement('link')
  link.id = 'dc'
  link.rel = 'stylesheet'
  link.href = '/src/assets/enteryCss.css'
  document.head.appendChild(link)
}
onMounted(() => {
  changeCss()
  auth()
})
const username = ref('')
const password = ref('')
// send username and pass to the backend
async function login() {
  const formData2 = new FormData()
  formData2.append('username', username.value)
  formData2.append('password', password.value)
  const resp = await fetch('/api/login', {
    method: 'POST',
    body: formData2,
  })
  if (resp.ok) {
    console.log('login stauts')
    router.push('/') //should be home
    isError.value = false
  } else {
    console.error('login faild')
    isError.value = true
  }
}
</script>

<template v-cloak>
  <div class="wrapper" :class="{ error: isError }">
    <form id="login" @submit.prevent="login">
      <h1>Login</h1>

      <div class="input-box">
        <input
          id="user"
          type="text"
          autocomplete="off"
          v-model="username"
          placeholder="UserName or Email"
          required
          :class="{ error: isError }"
        />
      </div>
      <div class="input-box">
        <input
          id="pass"
          type="password"
          autocomplete="off"
          v-model="password"
          placeholder="Password"
          required
        />
      </div>
      <button type="submit" class="btn">Login</button>
      <div class="register">
        <router-link to="/sign-up">
          <p id="register"><a>Don't have an account? Register</a></p>
        </router-link>
      </div>
    </form>
  </div>
</template>

<style>
.error {
  border: 2px solid red;
  background-color: #ffe5e5;
}
</style>
