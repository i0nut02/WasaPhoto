<script>
export default {
  props: {
    posts: {
      type: Object
    }
  },

  watch: {
    posts: function(newPosts) {
      this.posts_ = newPosts;
    }
  },

  data: function () {
    return {
      posts_: null // Initialize as null
    };
  },

  emits: ["delete-post"],

  methods: {
    async removePost(postData) {
      this.$emit("delete-post", postData);
    }
  },

  mounted() {
    // Initialize posts_ here if needed
    this.posts_ = this.posts;
  }
}
</script>

<template>
  <div>
    <div v-for="post in posts_" :key="post.id">
      <Post :post_data="post" @delete-post="removePost"></Post>
    </div>
  </div>
</template>

<style>
/* Add your styles if needed */
</style>
