<template>
  <div>
    <form @submit="getSecret">
      <input type="text" placeholder="Hash" v-model="hash" required>
      <button type="submit">Get</button>
    </form>
    <pre>{{ responseStr }}</pre>
  </div>
</template>

<script>
export default {
  data() {
    return {
      hash: '',
      responseStr: ''
    }
  },
  methods: {
    getSecret() {
      this.jsonstr = ''
      this.$http.get('secret/'+ this.hash).then(response => {
        this.responseStr = response.bodyText
      }, error => {
        this.responseStr = error.statusText
      });
    }
  }
}
</script>

