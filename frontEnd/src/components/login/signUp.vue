<script setup>
import { ref, onMounted } from 'vue'
import router from '@/router'
const isError = ref(false)
const selected = ref('♂️ Male')
const emit = defineEmits(['update:modelValue'])
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
onMounted(() => {
  auth()
  changeCss()
})
//defult pfp
import profileImg from '@/assets/profile.jpg'
const imgPrv = ref(profileImg)
const imgFile = ref(profileImg)
//this  to hold the img file

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

//userInfo
const firstname = ref('')
const lastname = ref('')
const age = ref('')
const email = ref('')
const username = ref('')
const password = ref('')
const aboutMe = ref('')
const gender = ref('♂️ Male')
const updateValue = () => {
  emit('update:modelValue', selected.value)
}

// send info to the backend
async function signUp() {
  // Append form data
  const formData = new FormData()
  formData.append('img', imgFile.value)
  formData.append('firstname', firstname.value)
  formData.append('lastname', lastname.value)
  formData.append('age', age.value)
  formData.append('email', email.value)
  formData.append('username', username.value)
  formData.append('password', password.value)
  formData.append('gender', gender.value)
  formData.append('aboutMe', aboutMe.value)

  //sending
  const resp = await fetch('/api/reg', {
    method: 'POST',
    body: formData,
  })
  if (resp.ok) {
    router.push('/')
    isError.value = false
  } else {
    console.error('login faild')
    isError.value = true
  }
}

//func to change the pfp
function previewPostImage(e) {
  const file = e.target.files[0]
  if (file) {
    imgFile.value = file
    const reader = new FileReader()
    reader.onload = (e) => {
      imgPrv.value = e.target.result
    }
    reader.readAsDataURL(file)
  }
}
</script>
<template v-cloak>
  <div class="wrapper" :class="{ error: isError }">
    <form id="signUP" @submit.prevent="signUp" enctype="multipart/form-data">
      <h1>Register</h1>
      <input type="file" id="img" name="img" accept="image/*" @change="previewPostImage" />
      <label for="img">
        <img name="pfp" id="pfp" :src="imgPrv" alt="Profile Picture" />
      </label>
      <p v-if="isError" style="color: red">username or email taken</p>
      <div class="input-box">
        <input
          type="text"
          autocomplete="off"
          v-model="firstname"
          placeholder="first name"
          required
        />
      </div>
      <div class="input-box">
        <input type="text" autocomplete="off" v-model="lastname" placeholder="last name" required />
      </div>
      <div class="input-box">
        <input
          min="0"
          max="99"
          id="age-num"
          type="number"
          autocomplete="off"
          v-model="age"
          placeholder="age"
          required
        />
      </div>
      <div class="input-box">
        <input type="email" autocomplete="off" v-model="email" placeholder="Email" required />
      </div>
      <div class="input-box">
        <input type="text" autocomplete="off" v-model="username" placeholder="UserName" required />
      </div>
      <div class="input-box">
        <input
          type="password"
          autocomplete="off"
          v-model="password"
          placeholder="Password"
          required
        />
      </div>
      <div class="input-box">
        <input type="text" autocomplete="off" v-model="aboutMe" placeholder="about me" />
      </div>
      <label class="ge_box" for="gender">
        <select @change="updateValue" v-model="gender" class="gender">
          <option value="♂️ Male">♂️ Male</option>
          <option value="♀️ Female">♀️ Female</option>
        </select>
      </label>
      <br />
      <br />
      <button type="submit" class="btn">register</button>
      <div class="register">
        <router-link to="/login">
          <p id="register"><a>Don't have an account?Register</a></p>
        </router-link>
      </div>
    </form>
  </div>
</template>

<style>
body {
  background-color: black;
}

select {
  -webkit-appearance: none;
  -moz-appearance: none;
  appearance: none;
}

/* Add a custom transparent effect */
select {
  color: white;
  background: transparent;
  border: none;
  font-size: 1.2rem;
}
option {
  background: black;
  color: white;
}
</style>
