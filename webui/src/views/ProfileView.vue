<script>
export default {
  data: function() {
    return {
      errmsg: null,
      username: null,
      self_profile: false,
      is_banned: false,
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

      if (this.$user.username == null) {
        this.$router.push("/");
        return;
      }

      if (this.username == this.$user.username) {
        this.self_profile = true;
      }

      try {
        let response = await this.$axios.get("/users/" + this.$user.username + "/bans/", {
          headers: {
            "Authorization": this.$user.token
          }
        });

        this.is_banned = response.data.map(x => x["username_string"]).includes(this.username);

        response = await this.$axios.get("/users/" + this.username + "/followers", {
          headers: {
            "Authorization": this.$user.token
          }
        });

        this.followers = response.data;

        this.is_followed = this.followers.map(x => x["username_string"]).includes(this.$user.username);

        response = await this.$axios.get("/users/" + this.username + "/following/", {
          headers: {
            "Authorization": this.$user.token
          }
        });

        this.following = response.data;

        response = await this.$axios.get("/users/" + this.username + "/profile/posts/", {
          headers: {
            "Authorization": this.$user.token
          }
        });

        this.posts = response.data;

      } catch (e) {
        if (e.response && e.response.status) {
          switch (e.response.status) {
            case 400:
              this.errmsg = "Bad request";
              break;
            case 401:
              this.errmsg = "Unauthorized";
              break;
            case 403:
              this.forbidden = true;
              break;
            case 404:
              this.errmsg = "Not found";
              break;
            case 500:
              this.errmsg = "Internal server error";
              break;
            default:
              this.errmsg = "Unhandled response code";
          }
        } else {
          this.errmsg = "An unexpected error occurred. Please try again later."
        }
      }
    },

    async follow() {
      let response = await this.$axios.put("/users/" + this.$user.username + "/following/" + this.username, {}, {
        headers: {
          Authorization: this.$user.token
        }
      });

      switch (response.status) {
        case 200:
          this.is_followed = true;
          this.followers.push({ "username_string": this.$user.username });
          break;
        case 400:
          this.errmsg = "Bad request";
          break;
        case 401:
          this.errmsg = "Unauthorized";
          break;
        case 403:
          this.forbidden = true;
          break;
        case 404:
          this.errmsg = "Not found";
          break;
        case 500:
          this.errmsg = "Internal server error";
          break;
        default:
          this.errmsg = "Unhandled response code";
      }
    },

    async unfollow() {
      const response = await this.$axios.delete("/users/" + this.$user.username + "/following/" + this.username, {
        headers: {
          Authorization: this.$user.token
        }
      });

      switch (response.status) {
        case 204:
          this.is_followed = false;
          this.followers = this.followers.filter(follower => {
            return follower.username_string !== this.$user.username;
          });
          break;
        case 400:
          this.errmsg = "Bad request";
          break;
        case 401:
          this.errmsg = "Unauthorized";
          break;
        case 403:
          this.forbidden = true;
          break;
        case 404:
          this.errmsg = "Not found";
          break;
        case 500:
          this.errmsg = "Internal server error";
          break;
        default:
          this.errmsg = "Unhandled response code";
      }
    },

    async ban() {
      const response = await this.$axios.put("/users/" + this.$user.username + "/bans/" + this.username, {}, {
        headers: {
          Authorization: this.$user.token
        }
      });

      switch (response.status) {
        case 200:
          this.is_banned = true;
          break;
        case 400:
          this.errmsg = "Bad request";
          break;
        case 401:
          this.errmsg = "Unauthorized";
          break;
        case 403:
          this.forbidden = true;
          break;
        case 404:
          this.errmsg = "Not found";
          break;
        case 500:
          this.errmsg = "Internal server error";
          break;
        default:
          this.errmsg = "Unhandled response code";
      }
    },

    async unban() {
      const response = await this.$axios.delete("/users/" + this.$user.username + "/bans/" + this.username, {
        headers: {
          Authorization: this.$user.token
        }
      });

      switch (response.status) {
        case 200:
          this.is_banned = false;
          break;
        case 400:
          this.errmsg = "Bad request";
          break;
        case 401:
          this.errmsg = "Unauthorized";
          break;
        case 403:
          this.forbidden = true;
          break;
        case 404:
          this.errmsg = "Not found";
          break;
        case 500:
          this.errmsg = "Internal server error";
          break;
        default:
          this.errmsg = "Unhandled response code";
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
        const response = await this.$axios.put("/users/" + this.$user.username + "/set_username", request_body, {
          headers: {
            Authorization: this.$user.token
          }
        });
        this.$user.username = new_name;
        this.username = new_name;

        this.$router.push("/profile/" + new_name);
      } catch (e) {
        if (e && e.response) {
          switch (e.response.status) {
            case 400:
              this.errmsg = "Bad request";
              break;
            case 401:
              this.errmsg = "Unauthorized";
              break;
            case 403:
              this.forbidden = true;
              break;
            case 404:
              this.errmsg = "Not found";
              break;
            case 409:
              this.errmsg = "Already exist a user with that username";
              break;
            case 500:
              this.errmsg = "Internal server error";
              break;
            default:
              this.errmsg = "Unhandled response code";
          }
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
          <div class="count">{{ this.followers == null ? 0 : this.followers.length }}</div>
          <div class="label">followers</div>
        </div>
        <div class="count-container">
          <div class="count">{{ this.following == null ? 0 : this.following.length }}</div>
          <div class="label">following</div>
        </div>
      </div>
      <div class="btn-toolbar mb-2 mb-md-0">
        <div v-if="username == $user.username" class="btn-group me-2">
          <button type="button" class="btn btn-sm btn-outline-primary" @click="ChangeName">
            Change Name
          </button>
        </div>
        <div v-if="username !== $user.username" class="btn-group me-2">
          <button v-if="is_banned" type="button" class="btn btn-sm btn-outline-primary" @click="unban">
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
    <ErrorMsg v-if="errmsg" :msg="errmsg"></ErrorMsg>
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