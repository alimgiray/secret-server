<template>
  <div>
    <form @submit="getSecret">
      <input type="text" placeholder="Hash" v-model="hash" required>
      <button type="submit">Get</button>
    </form>
    <pre>{{ jsonstr | pretty }}</pre>
  </div>
</template>

<script>
export default {
  data() {
    return {
      hash: '',
      jsonstr: '{}'
    }
  },
  methods: {
    getSecret() {
      this.jsonstr = '{}'
      this.$http.get('secret/'+ this.hash).then(response => {
        console.log(response);
        this.jsonstr = response.bodyText;
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

