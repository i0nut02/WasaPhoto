<script>
  import { Axios } from 'axios';

  export default {
    data: function() {
      return {
        errormsg: null,
        loading: false,
        stream: [],
        stream_top: 0,
        new_posts: true,
      }
    },
    methods: {
      async initialize() {
        if (this.$user.username == null) {
          this.$router.push("/");
          return;
        }

        this.loading = true;
        this.errormsg = null;

        let new_posts = await this.getStream(this.stream_top);

        if (new_posts == null) { 
          this.loading = false;
          return;
        }

        if (new_posts.length == 0) {
          this.new_posts = false;
        }

        this.stream.push(...new_posts);
        this.stream_top += new_posts.length;
        this.loading = false;

        window.onscroll = async () => {
          if (window.innerHeight + window.scrollY >= document.body.offsetHeight) {
            let new_posts = await this.getStream(this.stream_top);

            if (this.there_are_more_posts == false) {
              return;
            }

            if (new_posts.length == 0) {
              this.there_are_more_posts = false;
            }

            this.stream.push(...new_posts);
            this.stream_top += new_posts.length;
          }
        };
      },

      async getStream(start) {
        let response = await this.$axios.get("/users/" + this.$user.username + "/stream?from=" + start + "&max_quantity=" + "10", {
          headers: {
            "Authorization": this.$user.token
          }
        });

        switch (response.status) {
          case 200:
            return response.data;
          case 400:
            this.errmsg = "Bad request";
            break;
          case 401:
            this.errmsg = "Unauthorized";
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
        return null;
      },

      async deletePost(post_data) {
        this.initializes();
      },

      async refresh() {
        this.initialize();
      },
    },

    mounted() {
      this.initialize();
    }
  }
</script>

<template>
  <div>
    <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
      <h1 class="h2">Home page</h1>
      <div class="btn-group me-2 mb-md-0">
					<button type="button" class="btn btn-sm btn-outline-secondary" @click="refresh">
						Refresh
					</button>
      </div>
    </div>
    <PostsList :posts="stream" :key="stream.length" @delete-post="deletePost"></PostsList>
    <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
  </div>
</template>

<style>
</style>
