Vue.use(VueMaterial.default)

var app = new Vue({
    el: '#app',
    data: {
        name: "",
        components: [],
    },
    mounted() {
      this.getConfig();
    },
    methods: {
      getConfig() {
        fetch('/api/environment').then(response => {
          if (response.status !== 200) {
            throw new Error(response.text());
          }
          response.json().then(body => {
            this.name = body.name;
            this.components = body.components;
          });
        });
      },
      addEnvironment() {
        this.promptActive = true;
      },
    },
});
