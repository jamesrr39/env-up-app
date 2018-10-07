Vue.use(VueMaterial.default)

Vue.component('cluster-component', {
  template: `
  <div id="component" style="margin: 10px; border: 1px dashed black; padding: 10px;">
    <h2>{{ component.name }}</h2>
    <div v-if="isRunning">
      <button type="button" @click="stop">Stop</button>
    </div>
    <div v-else>
      <button type="button" @click="start">Start</button>
    </div>
    <p>{{ component.description }}</p>
    <code>{{ component.runCmd }}</code>
    <ul>
      <li v-for="logMessage in component.logMessages">
        {{ logMessage }}
      </li>
    </ul>
  </div>`,
  props: ['component'],
  methods: {
    start() {
      this.isRunning = true;
      fetch(`/api/environment/${encodeURIComponent(this.component.name)}/start`, {
        method: 'POST'
      }).then(response => {
        if (response.code >= 300 || response.code < 200) {
          this.isRunning = false;
          console.log("error");
        }
      });
    },
    stop() {
      this.isRunning = false;
      fetch(`/api/environment/${encodeURIComponent(this.component.name)}/stop`, {
      // fetch(`/api/app%201/start`, {
        method: 'POST'
      }).then(response => {
        if (response.code >= 300 || response.code < 200) {
          this.isRunning = true;
          console.log("error");
        }
      });
    }
  },
  data: () => ({
    isRunning: false,
  })
});

new Vue({
    el: '#app',
    data: {
        name: "",
        components: [],
    },
    mounted() {
      this.getConfig();
    },
    methods: {
      getConfig: function() {
        const self = this;

        fetch('/api/environment').then(response => {
          if (response.status !== 200) {
            throw new Error(response.text());
          }
          response.json().then(body => {
            self.name = body.name;
            self.components = body.components.map(component => ({
              ...component,
              logMessages: [],
            }));
            self.openWebsocket();
          });
        });
      },
      openWebsocket () {
        const ws = new WebSocket(`ws://${window.location.host}/api/environment/logs`)
        ws.onmessage = (event) => {
          const data = JSON.parse(event.data);
          const component = this.components.find(component => component.name === data.componentName);
          component.logMessages.push(data.message);
        };

        ws.onerror = (event) => {
          alert(`websocket error: ${JSON.stringify(event.data)}`);
        };

        window.addEventListener('beforeunload', () => {
          ws.close();
        });
      },
    },
});
