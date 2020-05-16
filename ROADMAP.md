# Roadmap

The road is very long. It is only the beginning of the journey.

Beginning might not be that smooth. Especially when development is not my primary job and writing in Golang is something relatively new for me. But even so, achieved progress already looks fantastic!

So please be patient. But most important, I'm looking for your issue reports, requests and PRs.

The primary focus is on operations stability, robustness and feedback. Especially when `ktempl` tends to be backing some crucial front-facing, edge gateways and load balancers.

- Proper os signals and in-app messages handling, retries, reloads, quitting the app.
- Polish logging
- Add support for Prometheus metrics
- Support complex values in set rather than just `string=string`
- In general enhance templating and make it look like templating in other apps out there. Functions and stuff.
- Support multiple Kubernetes clusters
- Tests :P
