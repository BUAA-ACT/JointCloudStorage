<template>
  <el-card>
    <div class="box">
      <div class="loader" v-if="key !== 3"></div>
      <i
        v-if="key === 3"
        style="text-align: center;
        display: block;color:#0bbd87;
        font-size: 4rem;"
        class="el-icon-success"
      ></i>
    </div>
    <h2>{{ key !== 3 ? activities[key].content + "中，请耐心等待......." : "当前集群已创建完成。" }}</h2>
    <div class="block">
      <el-timeline>
        <el-timeline-item
          v-for="(activity, index) in activities"
          :key="index"
          :icon="activity.icon"
          :type="activity.type"
          :color="activity.color"
          :size="activity.size"
          :timestamp="activity.timestamp"
        >
          <span>{{ activity.content }}</span>
          <span>{{ activity.duration }}</span>
          <span>{{ activity.status }}</span>
        </el-timeline-item>
      </el-timeline>
    </div>
  </el-card>
</template>

<script>
export default {
  data() {
    return {
      activities: [
        {
          content: "创建集群安全测试规则",
          duration: "大约3分钟",
          status: "已完成",
          size: "large",
          type: "primary",
          color: "#0bbd87",
          icon: "el-icon-check"
        },
        {
          content: "创建控制节点网络",
          duration: "大约3分钟",
          status: "创建中",
          icon: "el-icon-loading",
          size: "large"
        },
        {
          content: "创建控制节点",
          duration: "大约3分钟",
          status: "待创建",
          size: "large",
          icon: "el-icon-more"
        },
        {
          content: "创建工作节点",
          duration: "大约3分钟",
          status: "待创建",
          size: "large",
          icon: "el-icon-more"
        }
      ],
      key: 0
    };
  },
  created() {
    const interval = setInterval(() => {
      this.key += 1;
      if (this.key === 3) {
        clearInterval(interval);
      }
      this.activities[this.key].icon = "el-icon-check";
      if (this.key !== 3) this.activities[this.key + 1].icon = "el-icon-loading";
      this.activities[this.key].color = "#0bbd87";
      this.activities[this.key].status = "已完成";
      if (this.key !== 3) this.activities[this.key + 1].status = "构建中";
    }, 2000);
  }
};
</script>
<style lang="scss">
.el-timeline-item__node--large {
  left: -10px;
  width: 30px;
  height: 30px;
}
.el-timeline-item__content {
  height: 40px;
  line-height: 2.3rem;
  text-indent: 4rem;
  span {
    display: block;
    width: 33%;
    float: left;
  }
}
</style>
<style lang="scss" scoped>
.loader {
  display: block;
  margin: 30px auto 100px auto;
  width: 1em;
  height: 1em;
  font-size: 35px;
  color: #409eff;
  pointer-events: none;
  border-radius: 50%;
  box-shadow: 0 1em 0 -0.2em currentcolor;
  position: relative;
  -webkit-animation: loader 0.8s ease-in-out alternate infinite;
  animation: loader 0.8s ease-in-out alternate infinite;
  -webkit-animation-delay: 0.32s;
  animation-delay: 0.32s;
  top: -1em;
}
h2 {
  text-align: center;
  margin-bottom: 70px;
  font-weight: 400;
}
.loader:after,
.loader:before {
  content: "";
  position: absolute;
  width: inherit;
  height: inherit;
  border-radius: inherit;
  box-shadow: inherit;
  -webkit-animation: inherit;
  animation: inherit;
}
.loader:before {
  left: -1em;
  -webkit-animation-delay: 0.48s;
  animation-delay: 0.48s;
}
.loader:after {
  right: -1em;
  -webkit-animation-delay: 0.16s;
  animation-delay: 0.16s;
}
@-webkit-keyframes loader {
  0% {
    box-shadow: 0 2em 0 -0.2em currentcolor;
  }
  100% {
    box-shadow: 0 1em 0 -0.2em currentcolor;
  }
}
@keyframes loader {
  0% {
    box-shadow: 0 2em 0 -0.2em currentcolor;
  }
  100% {
    box-shadow: 0 1em 0 -0.2em currentcolor;
  }
}
</style>
