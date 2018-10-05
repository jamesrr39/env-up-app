Vue.use(VueMaterial.default)

Vue.component('cluster-component', {
  template: `
  <div id="component" style="margin: 10px; border: 1px dashed black; padding: 10px;">
    <h2>{{ component.name }}</h2>
    <p>{{ component.description }}</p>
    <code>{{ component.runCmd }}</code>
    <ul>
      <li v-for="logMessage in component.logMessages">
        {{ logMessage }}
      </li>
    </ul>
  </div>`,
  props: ['component'],
});

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
            const components = body.components.map(component => ({
              ...component,
              logMessages: [],
            }));
            this.components = components;

            const ws = new WebSocket(`ws://${window.location.host}/api/environment/logs`)
            ws.onmessage = (event) => {
              console.log(event.data);
              const data = JSON.parse(event.data);
              components = components.filter(component => {
                if (component.name !== data.componentName) {
                  return component;
                }

                return {
                  ...component,
                  logMessages: component.logMessages.concat([data.message]),
                };
              });
              this.components = components;
            }

            ws.onerror = (event) => {
              alert(`websocket error: ${JSON.stringify(event.data)}`);
            }
          });
        });
      },
      addEnvironment() {
        this.promptActive = true;
      },
    },
});
