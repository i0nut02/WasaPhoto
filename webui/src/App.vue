<script setup>
import { RouterLink, RouterView } from 'vue-router'
</script>
<script>
export default {

	data: function () {
		return {
			search_results: null,
		}
	},

	methods: {

		async logout(){
			this.$user.username = null;
			this.$user.token = null;
			this.$router.push("/");
		},

		async performSearch() {
			let search = document.querySelector("input").value;

			search = search.trim();

			if (search.length > 0) {

				const searcher_id = this.$user.token;

				if (searcher_id == null) {
					this.$user.username = null;
					this.$router.push("/");
					return
				}

				const header = {
					"Authorization": searcher_id
				}
				let response = await this.$axios.get("/users", {
					params: {
						"search_term": search
					},
					headers: header
				});

				if (response.status == 200) {
					this.search_results = response.data;
				} else {
					this.search_results = null;
				}
			}
			else {
				this.search_results = null;
			}
		},

		async refresh() {
			if (this.$user.username == null) {
				this.$router.push("/");
			}
		},
	},

	mounted(){
		this.refresh();
	}
}
</script>

<template>

	<header class="navbar navbar-dark sticky-top bg-dark flex-md-nowrap p-0 shadow">
		<a class="navbar-brand col-md-3 col-lg-2 me-0 px-3 fs-6" href="#/">WasaPhoto</a>
		<button class="navbar-toggler position-absolute d-md-none collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#sidebarMenu" aria-controls="sidebarMenu" aria-expanded="false" aria-label="Toggle navigation">
			<span class="navbar-toggler-icon"></span>
		</button>
	</header>

	<div class="container-fluid">
		<div class="row">
			<nav id="sidebarMenu" class="col-md-3 col-lg-2 d-md-block bg-light sidebar collapse">
				<div class="position-sticky pt-3 sidebar-sticky">
					<h3 class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted">
						<span>{{ $user.username == null ? "Not logged in" : "Hi " + $user.username }}</span>
					</h3>
					<ul class="nav flex-column">
						<li class="nav-item" v-if="$user.username !== null">
						<RouterLink :to="`/profile/${$user.username}`" class="nav-link">
							<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#home"/></svg>
							Profile
						</RouterLink>
						</li>

						<li class="nav-item" v-if="$user.username !== null">
						<RouterLink :to="`/stream/${$user.username}`" class="nav-link">
							<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#layout"/></svg>
							Stream
						</RouterLink>
						</li>


						<li class="nav-item" v-if="$user.username !== null">
							<button @click="logout" class="nav-link btn">
								<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#key"/></svg>
								Logout
							</button>
						</li>

						<li class="nav-item" v-if="$user.username !== null">
							<form class="navbar-form my-2 my-md-0 d-flex justify-content-end">
							<div class="input-group">
								<input class="form-control rounded-pill shorter-search" id="SearchBox" type="text" placeholder="Search" aria-label="Search"
								@input="performSearch()">
							</div>

							<datalist class="list-group custom-select w-10 dropdown mt-5 position-absolute">
								<option class="list-group-item align-middle" v-for="user in search_results" :key="user['username-string']">
								<i class="bi bi-person-circle m-2 fa-lg" style="font-size: 1.5rem;"></i>
								<router-link class="text-dark text-decoration-none m-0" style="font-size: 1.0rem;"
									:to="'/profile/' + user['username_string']">
									{{ user['username_string'] }}
								</router-link>
								</option>
							</datalist>
							</form>
						</li>
					</ul>
				</div>
			</nav>

			<main class="col-md-9 ms-sm-auto col-lg-10 px-md-4">
				<RouterView />
			</main>
		</div>
	</div>
</template>

<style>
	.small-search {
	height: 10px;
	font-size: 5px;
	width: 10px;
	}

  .nav-link.btn {
    border: none !important; /* Use !important to ensure this style is applied */
  }

</style>