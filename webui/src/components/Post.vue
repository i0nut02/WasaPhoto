<script>

export default {
    props: {
        post_data: {
            type:Object
        }
    },

    data: function () {
        return {
            upload_time: null,
            is_user_post: false,
            photo_id: null,
            file: "",
            user_liked: false,
            author: null,
            likes: 0,
            comments: [],
            description: "",
            errormsg: null,
            newCommentText: ""
        }
    },

    emits: ["delete-post"],

    methods: {
        async initialize() {
            this.author = this.post_data.author;
            this.photo_id = this.post_data.photo_id;
            this.description = this.post_data.description;
            this.file = this.post_data.photo_file;

            this.upload_time = this.post_data.upload_time.split("T");
            let date = this.upload_time[0].split("-")
            let time = this.upload_time[1].split(":")

            this.upload_time = date[2] + "/" + date[1] + "/" + date[0] + " at " + time[0] + ":" + time[1];

            this.is_user_post = this.author == this.$user.username;

            this.likes = this.post_data.num_likes;
            this.user_liked = this.post_data.liked_photo;

            try {
                let response  = await this.$axios.get("/users/" + this.author + "/profile/posts/" + this.photo_id + "/comments/", {
                    headers: {
                        "Authorization": this.$user.token
                    }
                });
                this.comments = response.data;
            } catch (e) {
                if (e.response) {
                    this.errormsg = e.response.data.response;
                } else {
                    this.errormsg = e.toString();
                }
            }
        },

        async deletePost() {
            try {
                let response = await this.$axios.delete("/users/" + this.author + "/profile/posts/" + this.photo_id + "/", {
                    headers: {
                            "Authorization": this.$user.token
                    }
                });
                this.$emit("delete-post", this.post_data);
                
            } catch(e) {
                if (e.response) {
                    this.errormsg = e.response.data.response;
                } else {
                    this.errormsg = e.toString();
                }
            }
        },

        async like() {
            try {
                let response = await this.$axios.put("/users/" + this.author + "/profile/posts/" + this.photo_id + "/likes/" + this.$user.token, {}, {
                    headers: {
                        Authorization: this.$user.token
                    }
                });
                if (this.user_liked == false) {
                    this.likes += 1;
                }
                this.user_liked = true;
            } catch (e) {
                if (e.response){
                    this.errormsg = e.response.data.response;
                } else {
                    this.errormsg = e.toString();
                }
            }
        },

        async unlike() {
            try {
                let response = await this.$axios.delete("/users/" + this.author + "/profile/posts/" + this.photo_id + "/likes/" + this.$user.token, {
                    headers: {
                        "Authorization": this.$user.token
                    }
                });
                if (this.user_liked == true) {
                    this.likes -= 1;
                }
                this.user_liked = false;
            } catch(e) {
                if (e.response) {
                    this.errormsg = e.response.data.response;
                } else {
                    this.errormsg = e.toString();
                }
            }
        },

        async addComment() {
            if (this.newCommentText.trim() == ""){
                return;
            }
            try {
                let response = await this.$axios.post("/users/" + this.author + "/profile/posts/" + this.photo_id + "/comments/", 
                    { "comment" : this.newCommentText }, { headers: {
                                                "Authorization" : this.$user.token
                }});
                this.comments.push(response.data);
            } catch(e) {
                if (e.response) {
                    this.errormsg = response.data.response;
                } else {
                    this.errormsg = e.toString();
                }
            }
        },

        async deleteComment(comment) {
            try {
                let response = await this.$axios.delete("/users/" + this.author + "/profile/posts/" + this.photo_id + "/comments/" + comment.id, 
                    { headers: {
                                    "Authorization" : this.$user.token
                    }});
                this.comments = this.comments.filter((c) => c.id != comment.id);
            } catch(e) {
                if (e.response) {
                    this.errormsg = e.response.data.response;
                } else {
                    this.errormsg = e.toString();
                }
            }
        },
    },

    mounted() {
        this.initialize();
    }
}

</script>

<template>
    <div class="post-container">
        <div class="post-header">
            <div>
                <i class="bi-person-circle" style="font-size: 2em;"></i>
                <RouterLink class="text-dark text-decoration-none m-0" :to="'/profile/' + author">
                    <span class="font-weight-bold h4" style="color: #333;">{{ this.author }}</span>
                </RouterLink>
            </div>
            <div>
                <button v-if="author == $user.username" @click="deletePost">Delete Post</button>
            </div>
        </div>

        <div class="post-divider"></div>

        <!-- Post Image -->
        <div class="image-container">
          <img :src="file" alt="Post Image" class="post-img">
        </div>

        <div class="post-divider"></div>

        <!-- Post Description -->
        <div class="post-description">
            <i class="bi-person-circle mx-1" style="font-size: 1.5em;"></i>
            <span>{{ this.description }}</span>
        </div>

        <div class="post-divider"></div>

        <!-- Like and Comment Section -->
        <div class="post-footer">
            <div class="d-flex justify-content-between">
                <div>
                    <div v-if="!is_user_post">
                        <button v-if="user_liked" type="button" class="btn btn-success btn-xs" @click="unlike"> liked </button>
                        <button v-else type="button" class="btn btn-danger btn-xs" @click="like"> like it </button>
                    </div>
                    <span><b>{{ this.likes }} likes</b></span>
                </div>
                <div>
                    <i class="bi-chat" style="color: #333;"></i> <span class="font-weight-bold" style="color: #333;">Comments: {{ this.comments.length }}</span>
                </div>
            </div>
        </div>

        <div class="post-divider"></div>

        <!-- Comments Section -->
        <div class="comments-section">
            <div v-if="comments.length == 0" class="col-12 align-content-center w-100"><!-- Center the text -->
                <span class="h5 mx-1 font-weight-bold align-middle text-muted text-center">No comments yet.</span>
            </div>
            <div v-else class="col-12">
                <span class="h4 mx-1 font-weight-bold align-middle mb-2 text-start">Comments: </span>
            </div>
            <div>
                <div v-for="comment in comments" :key="comment.id" class="comment">
                    <div>
                        <RouterLink class="text-dark text-decoration-none m-0" :to="'/profile/' + comment.author">
                            <strong>{{ comment.author }}:</strong> 
                        </RouterLink>
                        <button v-if="is_user_post" class="delete-comment-btn" @click="deleteComment(comment)">Delete</button>
                        <br>
                        <div>{{ comment.content }}</div>
                    </div>
                </div>
            </div>
        </div>

        <div class="post-divider"></div>

        <!-- Comment Input -->
        <div class="comment-input">
            <input type="text" placeholder="Add a comment..." style="flex: 1;" v-model="newCommentText">
            <button class="btn btn-success" @click="addComment()">Post</button>
        </div>
    </div>
    
    <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
</template>

<style>

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

  .post-img {
    max-width: 100%;
    max-height: 100%;
    object-fit: contain; /* Maintain image proportions */
  }
        .post-container {
            border: 2px solid #ffcc00;
            border-radius: 15px;
            background-color: #f9f9f9;
            margin: 20px;
            padding: 20px;
            max-width: 600px;
            margin-left: auto;
            margin-right: auto;
            box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
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
            width: 100%;
            border-radius: 15px;
            margin-bottom: 15px;
        }
        .post-footer {
            margin-top: 15px;
            color: #666;
        }
        .comments-section {
            margin-top: 15px;
        }
        .comment {
            position: relative;
            margin-bottom: 10px;
            padding: 10px;
            background-color: #e6e6e6;
            border-radius: 8px;
            display: flex;
            align-items: center; /* Center items vertically */
        }
        .comment strong {
            color: #333;
            flex: 1; /* Allow the comment text to take the remaining space */
        }
        .delete-comment-btn {
            background: #ff6666;
            color: #fff;
            border: none;
            border-radius: 5px;
            padding: 5px 10px;
            cursor: pointer;
            font-size: 12px;
            margin-left: 10px; /* Add margin for better spacing */
        }
        .comment-input {
            display: flex;
            margin-top: 15px;
        }
    .comment-input input {
        flex: 1;
        padding: 8px;
        border: 1px solid #ddd;
        border-radius: 5px;
        margin-right: 10px;
    }
    .comment-input button {
        background: #4CAF50;
        color: #fff;
        border: none;
        border-radius: 5px;
        padding: 8px 15px;
        cursor: pointer;
    }

    .post-divider {
        border-top: 1px solid #ccc;
        margin: 10px 0; /* Adjust the margin as needed */
    }
</style>