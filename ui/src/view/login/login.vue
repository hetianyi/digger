<style>

  @import url('https://fonts.googleapis.com/css?family=Montserrat:400,800');

  * {
    box-sizing: border-box;
  }

  h1 {
    font-weight: bold;
    margin: 0;
  }

  h2 {
    text-align: center;
  }

  p {
    font-size: 14px;
    font-weight: 100;
    line-height: 20px;
    letter-spacing: 0.5px;
    margin: 20px 0 30px;
  }

  span {
    font-size: 12px;
  }

  a {
    color: #333;
    font-size: 14px;
    /* 	text-decoration: none; */
    margin: 15px 0;
  }

  button {
    border-radius: 20px;
    border: 1px solid #2D8CF0;
    background-color: #2D8CF0;
    color: #FFFFFF;
    font-size: 12px;
    font-weight: bold;
    padding: 12px 45px;
    letter-spacing: 1px;
    text-transform: uppercase;
    transition: transform 80ms ease-in;
  }

  button:active {
    transform: scale(0.95);
  }

  button:focus {
    outline: none;
  }

  button.ghost {
    background-color: transparent;
    border-color: #FFFFFF;
  }

  form {
    background-color: #FFFFFF;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-direction: column;
    padding: 0 50px;
    height: 100%;
    text-align: center;
  }

  input {
    background-color: #eee;
    border: none;
    padding: 12px 15px;
    margin: 8px 0;
    width: 100%;
  }

  .container {
    background-color: #fff;
    border-radius: 10px;
    box-shadow: 0 14px 28px rgba(0,0,0,0.25),
    0 10px 10px rgba(0,0,0,0.22);
    position: relative;
    overflow: hidden;
    width: 768px;
    max-width: 100%;
    min-height: 480px;

    display: inline-block;
    margin-top: 10%;
  }

  .form-container {
    position: absolute;
    top: 0;
    height: 100%;
    transition: all 0.6s ease-in-out;
  }

  .sign-in-container {
    left: 0;
    width: 50%;
    z-index: 2;
  }

  .container.right-panel-active .sign-in-container {
    transform: translateX(100%);
  }

  .sign-up-container {
    left: 0;
    width: 50%;
    opacity: 0;
    z-index: 1;
  }

  .container.right-panel-active .sign-up-container {
    transform: translateX(100%);
    opacity: 1;
    z-index: 5;
    animation: show 0.6s;
  }

  @keyframes show {
    0%, 49.99% {
      opacity: 0;
      z-index: 1;
    }

    50%, 100% {
      opacity: 1;
      z-index: 5;
    }
  }

  .overlay-container {
    position: absolute;
    top: 0;
    left: 50%;
    width: 50%;
    height: 100%;
    overflow: hidden;
    transition: transform 0.6s ease-in-out;
    z-index: 100;
  }

  .container.right-panel-active .overlay-container{
    transform: translateX(-100%);
  }

  .overlay {
    background: #2D8CF0;
    background: -webkit-linear-gradient(to right, #2DF0EE, #2D8CF0);
    background: linear-gradient(to right, #2DF0EE, #2D8CF0);
    background-repeat: no-repeat;
    background-size: cover;
    background-position: 0 0;
    color: #FFFFFF;
    position: relative;
    left: -100%;
    height: 100%;
    width: 200%;
    transform: translateX(0);
    transition: transform 0.6s ease-in-out;
  }

  .container.right-panel-active .overlay {
    transform: translateX(50%);
  }

  .overlay-panel {
    position: absolute;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-direction: column;
    padding: 0 40px;
    text-align: center;
    top: 0;
    height: 100%;
    width: 50%;
    transform: translateX(0);
    transition: transform 0.6s ease-in-out;
  }

  .overlay-left {
    transform: translateX(-20%);
  }

  .container.right-panel-active .overlay-left {
    transform: translateX(0);
  }

  .overlay-right {
    right: 0;
    transform: translateX(0);
  }

  .container.right-panel-active .overlay-right {
    transform: translateX(20%);
  }

  .social-container {
    margin: 20px 0;
  }

  .social-container a {
    border: 1px solid #DDDDDD;
    border-radius: 50%;
    display: inline-flex;
    justify-content: center;
    align-items: center;
    margin: 0 5px;
    height: 40px;
    width: 40px;
  }

  footer {
    background-color: #222;
    color: #fff;
    font-size: 14px;
    bottom: 0;
    position: fixed;
    left: 0;
    right: 0;
    text-align: center;
    z-index: 999;
  }

  footer p {
    margin: 10px 0;
  }

  footer i {
    color: red;
  }

  footer a {
    color: #3c97bf;
    text-decoration: none;
  }
</style>

<template>
  <div style="text-align: center;">
    <div class="container" id="container">
      <div class="form-container sign-in-container">
        <form action="#" type="button">
          <h1>登录</h1>
          <div class="social-container">
            <a href="#" class="social"><i class="fab fa-facebook-f"></i></a>
            <a href="#" class="social"><i class="fab fa-google-plus-g"></i></a>
            <a href="#" class="social"><i class="fab fa-linkedin-in"></i></a>
          </div>
          <input type="text" v-model="username" placeholder="用户名">
          <input type="password" v-model="password" placeholder="密码">
          <p style="width: 100%; font-size: 12px; color: #c3c3c3; margin: 5px 0">默认用户名密码：admin/admin</p>
          <button type="button" @click="handleSubmit">登录</button>
        </form>
      </div>
      <div class="overlay-container">
        <div class="overlay">
          <div class="overlay-panel overlay-right">
            <img src="">
            <h1>Digger</h1>
            <p>A powerful and flexible web crawler</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
  import LoginForm from '_c/login-form'
  import { mapActions } from 'vuex'
  export default {
    components: {
      LoginForm
    },
    data() {
      return {
        username: '',
        password: '',
      }
    },
    methods: {
      ...mapActions([
        'handleLogin',
        'getUserInfo'
      ]),

      handleSubmit () {
        if (this.username.trim() == '') {
          this.$Message.error('请输入用户名')
          return
        }
        if (this.password.trim() == '') {
          this.$Message.error('请输入密码')
          return
        }

        this.handleLogin({
          username: this.username,
          password: this.password
        }).then(res => {
          this.getUserInfo().then(res => {
            this.$router.push({
              name: this.$config.homeName
            })
          })
        })
      }
    }
  }
</script>

<style>

</style>
