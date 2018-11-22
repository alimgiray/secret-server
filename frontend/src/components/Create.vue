<template>
  <div>
    <form @submit="createSecret">
      <input type="text" placeholder="Secret" v-model="secret" required>
      <input type="number" placeholder="View count" v-model="count" required>
      <input type="number" placeholder="Expire time (in minutes)" v-model="expire" required>
      <button type="submit">Create</button>
    </form>
    <pre>{{ jsonstr | pretty }}</pre>
  </div>
</template>

<script>
export default {
  data() {
    return {
      secret: '',
      count: '',
      expire: '',
      jsonstr: '{}'
    }
  },
  methods: {
    createSecret() {
      this.jsonstr = '{}'
      let data = {}
      data["secret"] = this.secret
      data["expireAfterViews"] = this.count
      data["expireAfter"] = this.expire
      this.$http.post('secret', data, {emulateJSON: true})
      .then(response => {
        this.jsonstr = response.bodyText
      }, error => {
        this.jsonstr = '{"error": "' + error.statusText + '"}'
      });
    }
  },
  filters: {
    pretty: function(value) {
      return JSON.stringify(JSON.parse(value), null, 2);
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
  button {
    
  }
</style>