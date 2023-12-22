<script>
export default {
  data: function () {
    return {
      errormsg: null,
      successmsg: null,
      description: "",
      image: null,
      imagePreviewUrl: null,
      postPreview: null,
    };
  },

  methods: {
    async onImageChange() {
      this.errormsg = null;
      this.successmsg = null;
      const file = this.$refs.imageInput.files[0];
      if (file) {
        const reader = new FileReader();
        reader.onload = (e) => {
          this.imagePreviewUrl = e.target.result;
          this.updatePreview(); // Update preview when image changes
        };
        reader.readAsDataURL(file);
      }
    },

    async updatePreview() {
      // Update post preview based on input values
      this.errormsg = null;
      this.successmsg = null;
      this.postPreview = {
        description: this.description,
        image: this.imagePreviewUrl,
      };
      // You can include other information in the preview as needed
    },

    async uploadPhoto() {
      if (this.imagePreviewUrl == null) {
        this.errormsg = "Please select an image.";
        return;
      }
      try {
        let response = await this.$axios.post("/users/" + this.$user.username + "/profile/posts/", {"file": this.imagePreviewUrl, "description": this.description}, {
          headers: {
            Authorization: this.$user.token
          }
        });
        switch (response.status) {
          case 201:
            this.successmsg = "the post is uploaded";
            this.errormsg = null;
            break;
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
      } catch (error) {
        this.errormsg = 'Failed to upload post. Please try again.';
        this.successmsg = null;
      }
    },
  },
};
</script>

<template>
  <div>
    <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
      <h1 class="h2">New Post</h1>
    </div>

    <form @submit.prevent="uploadPhoto" enctype="multipart/form-data">
      <label for="image" class="form-label">Select Image:</label>
      <input type="file" id="image" name="image" accept="image/*" required class="form-input" ref="imageInput" @change="onImageChange">
      <br>

      <label for="description" class="form-label">Description:</label>
      <textarea v-model="description" id="description" name="description" rows="4" cols="50" class="form-textarea" @input="updatePreview"></textarea>
      <br>

      <input type="submit" value="Post" class="form-submit btn btn-success">
    </form>
    <br>

    <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
    <SuccessMsg v-if="successmsg" :msg="successmsg"></SuccessMsg>

    <!-- Post Preview -->
    <div v-if="postPreview" class="post-preview">
      <h3>Post Preview</h3>
      <div class="post-container">
        <div class="post-header">
          <div>
            <i class="bi-person-circle" style="font-size: 2em;"></i>
            <span class="font-weight-bold h4" style="color: #333;">{{ this.$user.username }}</span>
          </div>
        </div>

        <div class="divider"></div>

        <!-- Post Image -->
        <div class="image-container">
          <img :src="postPreview.image" alt="Post Image" v-if="postPreview.image" class="post-img">
        </div>

        <div class="divider"></div>

        <!-- Post Description -->
        <div class="post-description">
          <i class="bi-person-circle mx-1" style="font-size: 1.5em;"></i>
          <span>{{ postPreview.description }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<style>
.post-form {
  background-color: #fff;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
  width: 300px;
  text-align: center;
}

.form-label {
  display: block;
  margin-bottom: 8px;
}

.form-input,
.form-textarea {
  width: 100%;
  padding: 8px;
  margin-bottom: 12px;
  box-sizing: border-box;
}

.form-submit {
  background-color: #4caf50;
  color: #fff;
  cursor: pointer;
}

.form-submit:hover {
  background-color: #45a049;
}

.preview-image {
  max-width: 100%;
  margin-top: 10px;
}

.post-container {
  border: 2px solid #ffcc00;
  border-radius: 15px;
  background-color: #f9f9f9;
  background-color: #000000;
  margin: 20px;
  padding: 20px;
  width: 50%;
  height: 50%;
  margin-left: auto;
  margin-right: auto;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.divider {
  border-top: 1px solid #ccc;
  margin: 10px 0; /* Adjust the margin as needed */
}

.image-container {
  background-color: #000000;
  width: 300px;
  height: 300px;
  margin-left: auto;
  margin-right: auto;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden; /* Ensure the image doesn't overflow the container */
}

.post-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}

.post-header button {
  background: #ff6666;
  color: #fff;
  border: none;
  border-radius: 5px;
  padding: 5px 10px;
  cursor: pointer;
  font-weight: bold;
}

.post-img {
  max-width: 100%;
  max-height: 100%;
  object-fit: contain; /* Maintain image proportions */
}

.post-footer {
  margin-top: 15px;
  color: #666;
}
</style>
