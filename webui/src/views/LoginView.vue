<script>
export default {
  components: {},
  data: function () {
    return {
      errormsg: null,
      username: null,
    };
  },
  methods: {
    async checkLoged() {
      if (this.$user.username != null) {
        this.$router.push({ path: '/stream/' + this.$user.username });
        return;
      }
    },

    async doLogin() {
      this.username = this.username.trim();
      if (!this.username.match("^[a-zA-Z0-9_.]{3,20}$")) {
        this.errormsg = "Username must contain: from 3 to 20 characters and just alphanumeric characters or \"_.\"";
      } else {
        try {
          let response = await this.$axios.post("/session", { username_string: this.username });
          this.$user.token = response.data["user_id"];
          this.$user.username = this.username;

          this.$router.push('/stream/' + this.$user.username);
        } catch(e) {
          if (e.response) {
            this.errormsg = e.response.data.response;
          } else {
            this.errormsg = e.toString();
          }
        }
      }
    },
  },
  mounted() {
    this.checkLoged();
  },
};
</script>

<template>
  <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
    <h1 class="h2">Welcome to WASAPhoto</h1>
  </div>
  <div class="input-group mb-3">
    <input
      type="text"
      id="username"
      v-model="username"
      class="form-control"
      placeholder="Insert a username to log in WASAPhoto."
      aria-label="Recipient's username"
      aria-describedby="basic-addon2"
    />
    <div class="input-group-append">
      <button class="btn btn-success" type="button" @click="doLogin">Login</button>
    </div>
  </div>
  <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
</template>

<style>
</style>
