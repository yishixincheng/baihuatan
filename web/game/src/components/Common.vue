<template>
    <div class="alert-demo-wrapper">
        <mu-alert color="success">
            <mu-icon left value="check_circle"></mu-icon> this is success alert
        </mu-alert>
        <mu-alert color="warning">
            <mu-icon left value="priority_high"></mu-icon> this is warning alert
        </mu-alert>
        <mu-alert color="info" >
            <mu-icon left value="info"></mu-icon> this is info alert
        </mu-alert>
        <mu-alert color="error" delete v-if="alert" @delete="closeAlert()">
            <mu-icon left value="warning"></mu-icon> this is error alert
        </mu-alert>
        <mu-button color="primary" v-if="!alert" @click="showAlert()">SHOW ALERT</mu-button>
    </div>
</template>

<script>
export default {
  name: 'Common',  
  data () {
    return {
      alert: true
    };
  },
  methods: {
    showAlert () {
      this.alert = true;
    },
    closeAlert () {
      this.alert = false;
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

.alert-demo-wrapper {
  width: 100%;
}
.mu-scale-transition-enter-active,
.mu-scale-transition-leave-active {
  transition: transform .45s cubic-bezier(0.23, 1, 0.32, 1), opacity .45s cubic-bezier(0.23, 1, 0.32, 1);
  backface-visibility: hidden;
}

.mu-scale-transition-enter,
.mu-scale-transition-leave-active {
  transform: scale(0);
  opacity: 0;
}
</style>