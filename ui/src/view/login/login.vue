
<template>
  <div style="text-align: center;" class="loginPage">
    <div class="container" id="container">
      <div class="form-container sign-in-container">
        <form action="#" type="button">
          <h1>登录</h1>
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

  @import url('https://fonts.googleapis.com/css?family=Montserrat:400,800');

  .loginPage * {
    box-sizing: border-box;
  }

  .loginPage h1 {
    font-weight: bold;
    margin: 0;
  }

  .loginPage h2 {
    text-align: center;
  }

  .loginPage p {
    font-size: 14px;
    font-weight: 100;
    line-height: 20px;
    letter-spacing: 0.5px;
    /*margin: 20px 0 30px;*/
  }

  .loginPage span {
    font-size: 12px;
  }

  .loginPage a {
    color: #333;
    font-size: 14px;
    /* 	text-decoration: none; */
    margin: 15px 0;
  }

  .loginPage button {
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

  .loginPage button:active {
    transform: scale(0.95);
  }

  .loginPage button:focus {
    outline: none;
  }

  .loginPage button.ghost {
    background-color: transparent;
    border-color: #FFFFFF;
  }

  .loginPage form {
    background-color: #FFFFFF;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-direction: column;
    padding: 0 50px;
    height: 100%;
    text-align: center;
  }

  .loginPage input {
    background-color: #eee;
    border: none;
    padding: 12px 15px;
    margin: 8px 0;
    width: 100%;
  }

  .loginPage .container {
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

  .loginPage .form-container {
    position: absolute;
    top: 0;
    height: 100%;
    transition: all 0.6s ease-in-out;
  }

  .loginPage .sign-in-container {
    left: 0;
    width: 50%;
    z-index: 2;
  }

  .loginPage .container.right-panel-active .sign-in-container {
    transform: translateX(100%);
  }

  .loginPage .sign-up-container {
    left: 0;
    width: 50%;
    opacity: 0;
    z-index: 1;
  }

  .loginPage .container.right-panel-active .sign-up-container {
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

  .loginPage .overlay-container {
    position: absolute;
    top: 0;
    left: 50%;
    width: 50%;
    height: 100%;
    overflow: hidden;
    transition: transform 0.6s ease-in-out;
    z-index: 100;
  }

  .loginPage .container.right-panel-active .overlay-container{
    transform: translateX(-100%);
  }

  .loginPage .overlay {
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

  .loginPage .container.right-panel-active .overlay {
    transform: translateX(50%);
  }

  .loginPage .overlay-panel {
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

  .loginPage .overlay-left {
    transform: translateX(-20%);
  }

  .loginPage .container.right-panel-active .overlay-left {
    transform: translateX(0);
  }

  .loginPage .overlay-right {
    right: 0;
    transform: translateX(0);
  }

  .container.right-panel-active .overlay-right {
    transform: translateX(20%);
  }

  .loginPage .social-container {
    margin: 20px 0;
  }

  .loginPage .social-container a {
    border: 1px solid #DDDDDD;
    border-radius: 50%;
    display: inline-flex;
    justify-content: center;
    align-items: center;
    margin: 0 5px;
    height: 40px;
    width: 40px;
  }

  .loginPage footer {
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

  .loginPage footer p {
    margin: 10px 0;
  }

  .loginPage footer i {
    color: red;
  }

  .loginPage footer a {
    color: #3c97bf;
    text-decoration: none;
  }
</style>
