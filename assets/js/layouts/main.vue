<template>
  <div id="page-container">
    <nav class="navbar">
      爱的小屋
    </nav>
    <div id="video-screen" v-show="authed">
      <img id="frame-image" :src="frameSrc" alt="">
      <div class="status" v-if="online&&showFps">
        {{fps}}fps
      </div>
    </div>

    <el-form v-show="!authed" class="login-box" label-width="60px" label-position="left">
      <el-form-item label="用户名">
        <el-input v-model="loginForm.user"></el-input>
      </el-form-item>
      <el-form-item label="密码">
        <el-input v-model="loginForm.pwd" type="password"></el-input>
      </el-form-item>

      <el-button @click="onLogin" style="width:100%;" type="primary">登录</el-button>

    </el-form>

    <el-form v-show="authed" class="verify-success" label-position="left" label-width="80px">
      <el-form-item label="状态">
        <template v-if="online">
          连接成功
        </template>
        <template v-else>
          离线
        </template>
      </el-form-item>
      <el-form-item label="显示FPS">
        <el-switch
          v-model="showFps"
          active-color="#13ce66"
          inactive-color="#ff4949">
        </el-switch>
      </el-form-item>

      <el-button @click="onLogout" style="width:100%;">注销</el-button>
    </el-form>

  </div>
</template>

<script>
export default {
  data: function() {
    return {
      loginForm: {
        user: "",
        pwd: ""
      },
      authed: false,
      socket: null,
      wsuri: "ws://" + document.location.host + "/socket/video_stream",
      fps: 0,
      fpsCounter: 0,
      lastFpsUpdatedAt: 0,
      online: false,
      showFps: true,
      frameSrc: null
    };
  },
  methods: {
    onMessage: function(e) {
      if (e.data == 'heartbeat') {
        return;
      }
      if (typeof(e.data) == 'string') {
        let data = JSON.parse(e.data)
        if (typeof(data) == 'object') {
          console.log(data)
          switch(data.option) {
            case 'auth':
            if (data.message == 'ok') {
              this.authed = true
            }
            break
          }
        }
        return;
      }

      if (this.lastFpsUpdatedAt == 0) {
        this.lastFpsUpdatedAt = new Date().getTime();
      } else if (new Date().getTime() - this.lastFpsUpdatedAt > 1000) {
        this.fps = this.fpsCounter;
        this.lastFpsUpdatedAt = new Date().getTime();
        this.fpsCounter = 1;
      }
      this.fpsCounter += 1;
      this.frameSrc = window.URL.createObjectURL(e.data);
    },
    onConnOpen: function() {
      console.log("connected to " + this.wsuri);
      this.online = true;
    },
    onConnClose: function(e) {
      console.log("connection closed (" + e.code + ")");
      this.online = false;
    },
    onLogin: function() {
      let param = {
        user: this.loginForm.user,
        password: this.loginForm.pwd
      }
      switch (this.socket.readyState) {
        case WebSocket.OPEN:
        // ready 
        this.socket.send(JSON.stringify(param))
        break;
        case WebSocket.CLOSED:
        case WebSocket.CLOSING:
        // error
        this.$message.error("网络出错")
        break;
        case WebSocket.CONNECTING:
        // waiting
        this.$message({
          type: 'warning',
          message: '网络连接中...稍后再试'
        })
        break;
      }
    },
    onLogout: function() {
      window.location.href = "/";
    }
  },
  created: function() {
    this.socket = new WebSocket(this.wsuri);
    this.socket.onopen = this.onConnOpen;
    this.socket.onclose = this.onConnClose;
    this.socket.onmessage = this.onMessage;
    
  }
};
</script>

<style lang="scss">
body {
  margin: 0px;
  padding: 0px;
}
</style>

<style lang="scss" scoped>
.navbar {
  height: 60px;
  text-align:center;
  line-height: 60px;
  color: #409EFF;
  font-size: 24px;
}
#video-screen {
  font-size: 14px;
  position: relative;
  #frame-image{
    width: 100vw;
    height: 75vw;
    background-color: #ccc;
  }
  .status {
    color: #00cc00;
    position: absolute;
    left: 20px;
    top: 20px;
    border: 1px solid #fff;
    padding: 0px 5px;
    line-height: 20px;
    height:20px;
  }
}
.login-box {
  padding: 10px;
  margin-top: 20px;
}
.verify-success {
  padding: 10px;
}
</style>

