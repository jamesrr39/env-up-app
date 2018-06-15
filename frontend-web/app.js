Vue.use(VueMaterial.default)

var app = new Vue({
    el: '#app',
    data: {
        environmentsFilePaths: [],
        promptActive: false,
        value: "",
    },
    mounted() {
      this.getConfig();
    },
    methods: {
      getConfig() {
        fetch('/api/config').then(response => {
          if (response.status !== 200) {
            throw new Error(response.text());
          }
          response.json().then(body => {
            this.config = body
          });
        });
      },
      addEnvironment() {
        this.promptActive = true;
      },
    },
});
