# iris.

[![Go Reference](https://pkg.go.dev/badge/github.com/TensRoses/iris.svg)](https://pkg.go.dev/github.com/TensRoses/iris)

Pomodoro bot. Self-develop into a general purpose bot

### fyi.

Stack: mongoDB - discordgo - Heroku - Github Action


### install.

to install `tensrose` binary run

```sh 
$ go install github.com/TensRoses/iris/cmd/tensrose
```

### folder structures.

```bash
.
├── bin
├── web           # bots front-facing api,                                  # lg: Rust
├── internal      # internal core of iris, logging, config handling         # lg: Go
├── pkg                                                                     
│   ├── core      # bot logics, including message queries with configstore  # lg: Go
│   ├── cmd       # bot cli                                                 # lg: Go
│   └── commands  # handles commands building and future commands plugins   # lg: Go
└── scripts
```

### inspirations.
> <div align="center"><i>why <strong>Go</strong> not <strong>TS/JS</strong>?</i></div>
- TypeScript and JavaScripts are amazing, and the fact that web elements for Iris will use `React` and `Next.js`, refers to [here](https://github.com/TensRoses/dashboard)
- Why Go then? Go is an awesome language with its interface and concurrency I just want to build a bot in Go :smile:
- Also, front-facing API will be written in Rust as a way to learn the language.

> <div align="center"><i>why <strong>iris</strong> when we already had </i><a href="https://github.com/jonas747/yagpdb"><strong>yagpdb</strong></a>?</div>
- `yagpdb` are pretty clunky and bloated in the sense that the maintainer uses a lot of boilerplate, whereas `iris` goals are to stay minimal and bloat-free as possible
- `yagpdb` structures are pretty messy and harder for development and thus `iris` aims to make it more developer-friendly as well as user-friendly
- `iris` got a lot of inspirations from `yagpdb` but takes a more lightweight approaches considering structures and deployments, refers to [here](pkg/README) for more details
- `iris` is pretty much <strong>WIP</strong> so any PR and help are appreciated
