# HAProxy Access Log Parser

[![Download from Artifactory](https://img.shields.io/badge/download-artifactory-brightgreen)](https://artifactory.ubisoft.org/generic/uks/pdappollonio/haproxy-parser/)

`haproxy-parser` is a tiny app that allows you to parse in a more human fashion the output of the HAProxy Controller access logs.

This program supports feeding the data through `stdin`:

```bash
kubectl logs -n ingress-haproxy haproxy-ingress-controller-abcde -c access-log | haproxy-parser
```

And it will error out if accessed directly:

```bash
$ haproxy-parser
This program doesn't accept input from stdin. You must pipe the output of the HAProxy Access Logs:
Example: kubectl logs -n ingress-haproxy haproxy-ingress-controller-abcde -c access-log | haproxy-parser
```

You can also use `less` after the parser app to have a paginated output.

### `(truncated)` output?

Some IP addresses might not be able to be parsed off the log. In fact, this program does not parse the log output as JSON, but instead, uses regular expressions to capture match the output needed. The reason being that when you're pulling the logs using `kubectl`, there's a maximum length of a string you can capture. If the IP was partially or completely there, then there's no way to know what IP address originated it.

You'll see those as `(truncated)` in the output of the program.
