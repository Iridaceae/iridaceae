## Internal

this directory includes core functionality of `iris`. This shouldn't be exposed or used by other parties

### structure.

| folder | description |
| ------ | ----------- |
|[`/internal/configparser`](./configparser) | handles configuration for `iris`. Included a master parser|
|[`/internal/database`](./database) | handles database services and logics (_current_ `mongodb`, maybe `sql` _in the future_)|
|[`/internal/testutils`](./testutils) | handles all core test dependencies (separated from core dependencies)|

- [`/internal/testutils/cbor`](testutils/cbor) and [`/internal/testutils/json`](testutils/json) are ported from [`rs/zerlog`](git@github.com:rs/zerolog.git) to reduce test dependencies.

### refactoring and testing.
- should be more structured as refactoring goes so that it is easier to manage and test
- rewrite tests to reduce boilerplate code
- implements token-bucket `ratelimiter` (now it is a simple `timedmap` to check executions)

### todos.
- [ ] fixes docker images with envars limiter
- [ ] added `di` container for our dependency issues.
- [ ] fixes ``rosetta`` rate limiter being to aggressive (panic when running into limit cap)
- [ ]  ```sclog``` should have options to persistently add global variables throughout different context. However, the current implementation is to be expected since the use of `goid`
