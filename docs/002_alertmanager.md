# Alertmanager


<img width="1440" alt="image" src="https://github.com/alextanhongpin/go-prometheus/assets/6033638/f928ee60-32ce-4349-8ead-319d39cceb1e">

View from loki:

<img width="1440" alt="image" src="https://github.com/alextanhongpin/go-prometheus/assets/6033638/a595d267-d213-4d44-ae56-1bae3ae8a55c">


Example logs:
```bash
app-1  | {"time":"2024-03-31T08:29:45.632209408Z","level":"INFO","msg":"get handler","release":"canary"}
app-1  | /
app-1  | /
app-1  | {"time":"2024-03-31T08:30:40.483606332Z","level":"INFO","msg":"post handler","body":"{\"receiver\":\"webhook\",\"status\":\"firing\",\"alerts\":[{\"status\":\"firing\",\"labels\":{\"alertname\":\"HighRequestErrorRate\",\"severity\":\"critical\"},\"annotations\":{\"description\":\"The error rate for HTTP requests has exceeded 5% for 5 seconds.\",\"title\":\"High request error rate\"},\"startsAt\":\"2024-03-31T08:30:10.412Z\",\"endsAt\":\"0001-01-01T00:00:00Z\",\"generatorURL\":\"http://0811a999f7bf:9090/graph?g0.expr=sum%28rate%28request_duration_seconds_count%7Bstatus%21~%222..%22%7D%5B1m%5D%29+or+vector%280%29%29+%2F+sum%28rate%28request_duration_seconds_count%5B1m%5D%29+or+vector%280%29%29+%3E+0.05\\u0026g0.tab=1\",\"fingerprint\":\"36b365830b032b13\"}],\"groupLabels\":{\"alertname\":\"HighRequestErrorRate\"},\"commonLabels\":{\"alertname\":\"HighRequestErrorRate\",\"severity\":\"critical\"},\"commonAnnotations\":{\"description\":\"The error rate for HTTP requests has exceeded 5% for 5 seconds.\",\"title\":\"High request error rate\"},\"externalURL\":\"http://6b217cb74cba:9093\",\"version\":\"4\",\"groupKey\":\"{}:{alertname=\\\"HighRequestErrorRate\\\"}\",\"truncatedAlerts\":0}\n"}
app-1  | {"time":"2024-03-31T08:30:40.483615435Z","level":"INFO","msg":"post handler","body":"{\"receiver\":\"webhook\",\"status\":\"firing\",\"alerts\":[{\"status\":\"firing\",\"labels\":{\"alertname\":\"Test default recipient\"},\"annotations\":{\"message\":\"Testalert default recipient\"},\"startsAt\":\"2024-03-31T08:30:10.412Z\",\"endsAt\":\"0001-01-01T00:00:00Z\",\"generatorURL\":\"http://0811a999f7bf:9090/graph?g0.expr=vector%281%29\\u0026g0.tab=1\",\"fingerprint\":\"4d6544dff935c7c6\"}],\"groupLabels\":{\"alertname\":\"Test default recipient\"},\"commonLabels\":{\"alertname\":\"Test default recipient\"},\"commonAnnotations\":{\"message\":\"Testalert default recipient\"},\"externalURL\":\"http://6b217cb74cba:9093\",\"version\":\"4\",\"groupKey\":\"{}:{alertname=\\\"Test default recipient\\\"}\",\"truncatedAlerts\":0}\n"}

```
