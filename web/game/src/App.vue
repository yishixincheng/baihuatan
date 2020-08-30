<template>
  <div id="app">
    <router-view/>
  </div>
</template>

<script>
import Common from '@/components/Common'
export default {
  name: 'App',
}
</script>

<style>
body,h3,h2,ul,li,dl,dd,ol,dt {
  margin: 0;
  padding: 0;
  list-style: none;
}
#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
}

</style>
