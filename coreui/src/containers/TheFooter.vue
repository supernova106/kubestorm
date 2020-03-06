<template>
  <CFooter :fixed="false">
    <div>
      <a href="https://kubestorm.io" target="_blank">CoreUI</a>
      <span class="ml-1">&copy; {{new Date().getFullYear()}} KubeStorm.</span>
    </div>
    <div class="ml-auto">
      <span class="mr-1">Powered by</span>
      <a href="https://kubestorm.io/ui" target="_blank">CoreUI for Vue</a>
    </div>
  </CFooter>
</template>

<script>
export default {
  name: "TheFooter"
};
</script>
