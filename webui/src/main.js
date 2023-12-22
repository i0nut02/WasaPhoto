import {createApp, reactive} from 'vue'
import App from './App.vue'
import router from './router'
import axios from './services/axios.js';

import ErrorMsg from './components/ErrorMsg.vue'
import SuccessMsg from './components/SuccessMsg.vue'
import LoadingSpinner from './components/LoadingSpinner.vue'
import PostsList from './components/PostsList.vue'
import Post from './components/Post.vue'

import './assets/dashboard.css'
import './assets/main.css'

const app = createApp(App)

var user = {
    token: null,
    username: null,
}

app.config.globalProperties.$axios = axios;

app.component("ErrorMsg", ErrorMsg);
app.component("SuccessMsg", SuccessMsg);
app.component("LoadingSpinner", LoadingSpinner);
app.component("PostsList", PostsList);
app.component("Post", Post);

app.config.globalProperties.$user = reactive(user);

app.use(router)
app.mount('#app')
