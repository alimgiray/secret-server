<template>
  <div>
    <form @submit="createSecret">
      <input type="text" placeholder="Secret" v-model="secret" required>
      <input type="number" placeholder="View count" v-model="count" required>
      <input type="number" placeholder="Expire time (in minutes)" v-model="expire" required>
      <button type="submit">Create</button>
    </form>
    <pre>{{ responseStr }}</pre>
  </div>
</template>

<script>
export default {
  data() {
    return {
      secret: '',
      count: '',
      expire: '',
      responseStr: ''
    }
  },
  methods: {
    createSecret() {
      this.jsonstr = ''
      let data = {}
      data["secret"] = this.secret
      data["expireAfterViews"] = this.count
      data["expireAfter"] = this.expire
      this.$http.post('secret', data, {emulateJSON: true})
      .then(response => {
        this.responseStr = response.bodyText
      }, error => {
        this.responseStr = error.statusText
      });
    }
  }
}
</script>

<style>
  input {
    display: block;
    margin-left: auto;
    margin-right: auto;
    margin-bottom: 20px;
  }
</style>