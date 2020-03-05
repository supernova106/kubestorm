<template>
  <div id="app">
    <img alt="Kubernetes logo" src="./assets/master-128.png">
    <Eks msg="Simple kubestorm"/>
  </div>
</template>

<script>
import Kubeconfig from "./components/Kubeconfig.vue";
import Eks from "./components/Eks.vue";

export default {
  name: "app",
  components: {
    Eks
  }
};
</script>

<style>
#app {
  font-family: "Avenir", Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
