import {createRouter, createWebHashHistory} from 'vue-router'
import ProfileView from '../views/ProfileView.vue'
import StreamView from '../views/StreamView.vue'
import LoginView from '../views/LoginView.vue'


const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/', component: LoginView},
		{path: '/stream/:username', component: StreamView},
		{path: '/profile/:username/', component: ProfileView},
	]
})

export default router
