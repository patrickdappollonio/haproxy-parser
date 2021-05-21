# HAProxy Access Log Parser

[![Download from Artifactory](https://img.shields.io/badge/download-artifactory-brightgreen)](https://artifactory.ubisoft.org/generic/uks/pdappollonio/haproxy-parser/)

`haproxy-parser` is a tiny app that allows you to parse in a more human fashion the output of the HAProxy Controller access logs.

This program supports feeding the data through `stdin`:

```text
kubectl logs -n ingress-haproxy haproxy-ingress-controller-abcde -c access-log | haproxy-parser
```

And it will error out if accessed directly:

```text
$ haproxy-parser
This program doesn't accept input from stdin. You must pipe the output of the HAProxy Access Logs:
Example: kubectl logs -n ingress-haproxy haproxy-ingress-controller-abcde -c access-log | haproxy-parser
```

You can also use `less` after the parser app to have a paginated output.

### `(truncated)` output?

Some IP addresses might not be able to be parsed off the log. In fact, this program does not parse the log output as JSON, but instead, uses regular expressions to capture match the output needed. The reason being that when you're pulling the logs using `kubectl`, there's a maximum length of a string you can capture. If the IP was partially or completely there, then there's no way to know what IP address originated it.

You'll see those as `(truncated)` in the output of the program.

### Example output

The output below has been truncated. The program will output a full log after executing.

```
REMOTE ADDRESS     COUNT    ENDPOINT
216.98.56.75       528      /v3/connect/config
216.98.56.75       381      /version
66.131.169.216     231      /k8s/clusters/c-ttz8z/apis/authorization.k8s.io/v1/selfsubjectaccessreviews
194.2.155.254      173      /version
66.131.169.216     117      /k8s/clusters/c-gjz9t/apis/authorization.k8s.io/v1/selfsubjectaccessreviews
10.172.38.17       117      /v3/connect/config
194.2.155.254      96       /v3/connect/config
66.131.169.216     80       /k8s/clusters/c-ttz8z/api/v1/namespaces/lens-metrics/services/prometheus:80/proxy/api/v1/query_range
216.98.56.75       71       /v3
216.98.56.75       71       /v3/schemas
66.131.169.216     60       /k8s/clusters/c-gjz9t/api/v1/namespaces/lens-metrics/services/prometheus:80/proxy/api/v1/query_range
216.98.60.148      57       /k8s/clusters/c-kmwnx/api/v1/namespaces/projectcontour/pods
216.98.60.148      46       /k8s/clusters/c-l6xbw/apis/authorization.k8s.io/v1/selfsubjectaccessreviews
```
