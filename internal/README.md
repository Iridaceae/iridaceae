## Internal

this directory includes core functionality of `iris`. This shouldn't be exposed or used by other parties

### structure.

| folder | description |
| ------ | ----------- |
|[`/internal/configparser`](./configparser) | handles configuration for `iris`. Included a master parser|
|[`/internal/datastore`](./datastore) | handles database services and logics (_current_ `mongodb`, maybe `sql` _in the future_)|
|[`/internal/jog`](./jog) | `Iris` command handlers |
|[`/internal/logging.go`](./logging.go) | contains `iris` structured logger wrapped around `rs/zerolog`|

### refactoring and testing.
- should be more structured as refactoring goes so that it is easier to manage and test
- rewrite tests to reduce boilerplate code
- implements token-bucket `ratelimiter` (now it is a simple `timedmap` to check executions)