<template>
  <div class="container" style="font-family: StandardBlack">
    <div class="row" style="margin-top:140px">
      <div class="col-2" />
      <div class="col-8">
        <h2>欢迎使用 Toaster</h2>
        <h4>自定义信息聚合工具</h4>
        <p style="margin-top:40px">Powered by Quited</p>
        <p>8/19/2019</p>
      </div>
      <div class="col-2" />
    </div>
    <div class="row" style="margin-top:70px">
      <div class="col-3" />
      <div class="col-3">
        <button
          id="start"
          class="btn btn-success"
          style="-webkit-app-region: no-drag; font-size:125%"
        >开始使用</button>
      </div>
      <div class="col-1" />
      <div class="col-3">
        <button
          id="exit"
          class="btn btn-warning"
          style="-webkit-app-region: no-drag; font-size:125%"
          @click="appExit"
        >退出</button>
      </div>
      <div class="col-2" />
    </div>
  </div>
</template>

<script>
export default {
  name: "index",
  data() {
    return {
      placeholder: "null",
    };
  },
  methods: {
    appExit: function() {
      window.opener = null;
      window.open("", "_self");
      window.close();
    }
  }
};
</script>
