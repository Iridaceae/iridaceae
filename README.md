<div align="center">
    <h1 align="center">~ 虹彩 ~</h1>
    <hr>
    <p align="center">
        <a href="https://pkg.go.dev/github.com/TensRoses/iris"><img alt="goref" src="https://pkg.go.dev/badge/github.com/TensRoses/iris.svg"></a>
        <a href="https://goreportcard.com/report/github.com/Iridaceae/iris"><img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/Iridaceae/iris"></a>
        <a href="https://github.com/Iridaceae/iris/issues"><img alt="GitHub issues" src="https://img.shields.io/github/issues/Iridaceae/iris?style=flat-square"></a>
        <img alt="GitHub commit checks state" src="https://img.shields.io/github/checks-status/Iridaceae/iris/a29703a1367977d2867167fda8c4146aea6cd58e?style=flat-square">
        <img alt="Codecov" src="https://img.shields.io/codecov/c/gh/Iridaceae/iris?style=flat-square">
    </p>
    <strong>
        A general purpose discord bot written in Go (discord.go)<br>
    </strong>
</div>

## install. <img alt="VimL" src="https://img.shields.io/badge/-Atlas-66d124?style=flat-square&logo=mongoDB&logoColor=white" />&nbsp;<img alt="Go" src="https://img.shields.io/badge/-discordgo-46a2f1?style=flat-square&logo=go&logoColor=white" />&nbsp;<img alt="Heroku" src="https://img.shields.io/badge/-Heroku-430098?style=flat-square&logo=heroku&logoColor=white" />&nbsp;<img alt="git" src="https://img.shields.io/badge/-Github Actions-000000?style=flat-square&logo=GitHub&logoColor=white" />
```sh 
$ go install github.com/Iridaceae/iris/cmd/tensroses-server
```

## folder structures.
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

## inspirations.
> <div align="left"><i>why <strong>Go</strong> not <strong>TS/JS</strong>?</i></div>
- TypeScript and JavaScripts are amazing, and the fact that web elements for Iris will use `React` and `Next.js`, refers to [here](https://github.com/TensRoses/dashboard)
- Why Go then? Go is an awesome language with its interface and concurrency I just want to build a bot in Go :smile:
- Also, front-facing API will be written in Rust as a way to learn the language.

> <div align="left"><i>why <strong>iris</strong> when we already had </i><a href="https://github.com/jonas747/yagpdb"><strong>yagpdb</strong></a>?</div>
- `yagpdb` are pretty clunky and bloated in the sense that the maintainer uses a lot of boilerplate, whereas `iris` goals are to stay minimal and bloat-free as possible
- `yagpdb` structures are pretty messy and harder for development and thus `iris` aims to make it more developer-friendly as well as user-friendly
- `iris` got a lot of inspirations from `yagpdb` but takes a more lightweight approaches considering structures and deployments, refers to [here](pkg/README.md) for more details
- `iris` is pretty much <strong>WIP</strong> so any <strong>PR</strong> and help are appreciated
