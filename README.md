# [GoThrough](https://gothrough.dev)

## Development

```bash
$ go run main.go
```

_Requires a least `go1.11` (for go modules support) and `git` to be installed._

_Replace `go` with [`gin`](https://github.com/codegangsta/gin) for auto-restarts._

##

Tests are run using the standard go test command (`go test ./...`).

## Deployment

This project is hosted on [Google Cloud Platform](https://cloud.google.com/)'s [Cloud Run](https://cloud.google.com/run/) platform.

The [deployment workflow](./.github/workflows/push.yaml) uses [GitHub Actions](https://developer.github.com/actions/) to publish a new image and update the running service on push.

## License

[MIT](./LICENSE)