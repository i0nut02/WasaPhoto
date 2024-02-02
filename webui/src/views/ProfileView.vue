<script>
export default {
  data: function() {
    return {
      errormsg: null,
      username: null,
      self_profile: false,
      is_banished: false,
      followers: null,
      following: null,
      is_followed: false,
      forbidden: false,
      posts: [],
    };
  },
  methods: {
    async refresh() {
      this.username = this.$route.params.username;
      this.posts = [];
      if (this.$user.username == null) {
        this.$router.push("/");
        return;
      }

      if (this.username == this.$user.username) {
        this.self_profile = true;
      }

      try {
        let response = await this.$axios.get("/users/" + this.username + "/profile/", {
          headers: {
            "Authorization": this.$user.token
          }
        });
        this.is_banished = response.data.is_banished

        this.followers = response.data.num_followers;

        this.is_followed = response.data.following;

        this.following = response.data.num_following;

        response = await this.$axios.get("/users/" + this.username + "/profile/posts/", {
          headers: {
            "Authorization": this.$user.token
          }
        });
        this.posts = response.data;
      } catch(e) {
          if (e.response) {
            this.errormsg = e.response.data.response;
            if (e.response.status == 403) {
              this.forbidden = true;
            }
          } else {
            this.errormsg = e.toString();
          }
      }
    },

    async follow() {
      try {
        let response = await this.$axios.put("/users/" + this.$user.username + "/following/" + this.username, {}, {
          headers: {
            Authorization: this.$user.token
          }
        });
        this.is_followed = true;
        this.followers += 1;
      } catch(e) {
        if (e.response) {
          this.errormsg = e.response.data.response;
        } else {
          this.errormsg = e.toString();
        }
      }
    },

    async unfollow() {
      try {
        let response = await this.$axios.delete("/users/" + this.$user.username + "/following/" + this.username, {
          headers: {
            Authorization: this.$user.token
          }
        });
        this.is_followed = false;
        this.followers -= 1;
      } catch(e) {
        if (e.response) {
          this.errormsg = e.response.data.response;
        } else {
          this.errormsg = e.toString();
        }
      }
    },

    async ban() {
      try {
        let response = await this.$axios.put("/users/" + this.$user.username + "/bans/" + this.username, {}, {
          headers: {
            Authorization: this.$user.token
          }
        });
        this.is_banished = true;
      } catch(e) {
        if (e.response) {
          this.errormsg = e.response.data.response;
        } else {
          this.errormsg = e.toString();
        }
      }
    },

    async unban() {
      try {
        let response = await this.$axios.delete("/users/" + this.$user.username + "/bans/" + this.username, {
          headers: {
            Authorization: this.$user.token
          }
        });
        this.is_banished = false;
      } catch(e) {
        if (e.response) {
          this.errormsg = e.response.data.response;
        } else {
          this.errormsg = e.toString();
        }
      }
    },

    async deletePost(post_data) {
      this.refresh();
    },

    async ChangeName() {
      const new_name = prompt("Change name", this.$user.username);

      if (new_name == null || new_name == "") {
        return;
      }

      if (!new_name.match("^[a-zA-Z0-9_.]{3,20}$")) {
        alert("Invalid username");
        return;
      }

      const request_body = {
        "username_string": new_name
      };

      try {
        let response = await this.$axios.put("/users/" + this.$user.username + "/set_username", request_body, {
          headers: {
            Authorization: this.$user.token
          }
        });
        this.$user.username = new_name;
        this.username = new_name;

        this.$router.push("/profile/" + new_name);
        
      } catch (e) {
        if (e.response) {
          this.errormsg = e.response.data.response;
        } else {
          this.errormsg = e.toString();
        }
      }
    },
  },
  mounted() {
    this.refresh();
  }
}
</script>

<template>
  <div v-if="forbidden">
    <br>
    <ErrorMsg msg="Forbidden Page"></ErrorMsg>
  </div>
  <div v-else>
    <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
      <h1 class="h2">{{ this.username }} Page</h1>
      <div class="follower-container">
        <div class="count-container">
          <div class="count">{{ this.posts.length }}</div>
          <div class="label">posts</div>
        </div>
        <div class="count-container">
          <div class="count">{{ this.followers }}</div>
          <div class="label">followers</div>
        </div>
        <div class="count-container">
          <div class="count">{{ this.following }}</div>
          <div class="label">following</div>
        </div>
      </div>
      <div class="btn-group me-2 mb-md-0">
					<button type="button" class="btn btn-sm btn-outline-secondary" @click="refresh">
						Refresh
					</button>
      </div>
      <div class="btn-toolbar mb-2 mb-md-0">
        <div v-if="username == $user.username" class="btn-group me-2">
          <button type="button" class="btn btn-sm btn-outline-primary" @click="ChangeName">
            Change Name
          </button>
        </div>
        <div v-if="username !== $user.username" class="btn-group me-2">
          <button v-if="is_banished" type="button" class="btn btn-sm btn-outline-primary" @click="unban">
            Unban
          </button>
          <button v-else type="button" class="btn btn-sm btn-outline-primary" @click="ban">
            Ban
          </button>
        </div>
        <div v-if="username !== $user.username" class="btn-group me-2">
          <button v-if="is_followed" type="button" class="btn btn-sm btn-outline-primary" @click="unfollow">
            Unfollow
          </button>
          <button v-else type="button" class="btn btn-sm btn-outline-primary" @click="follow">
            Follow
          </button>
        </div>
      </div>
    </div>
    <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
    <PostsList :posts="posts" :key="username" @delete-post="deletePost"></PostsList>
  </div>
</template>

<style>
.follower-container {
  display: flex;
  justify-content: space-around;
  align-items: center;
  background-color: #f0f0f0;
  padding: 10px;
  border-radius: 5px;
  margin: 20px auto;
}

.count-container {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.count {
  font-size: 13px;
  font-weight: bold;
}

.label {
  font-size: 13px;
  color: #555;
}

.count-container + .count-container {
  margin-left: 20px;
}
</style>