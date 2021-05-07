<div align="center">
    <h1 align="center">~ 虹彩 ~</h1>
    <hr>
    <p align="center">
        <a href="https://pkg.go.dev/github.com/Iridaceae/iridaceae"><img alt="golang" src="https://pkg.go.dev/badge/github.com/Iridaceae/iridaceae.svg"></a>
        <a href="https://goreportcard.com/report/github.com/Iridaceae/iridaceae"><img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/Iridaceae/iridaceae"></a>
        <a href="https://codecov.io/gh/Iridaceae/iridaceae"><img alt="codecov" src="https://codecov.io/gh/Iridaceae/iridaceae/branch/main/graph/badge.svg?token=qGdkowt7ki"/></a>
        <a href="https://github.com/Iridaceae/iridaceae/issues"><img alt="GitHub issues" src="https://img.shields.io/github/issues/Iridaceae/iridaceae?style=flat-square"></a>
        <img alt="GitHub commit checks state" src="https://img.shields.io/github/checks-status/Iridaceae/iridaceae/a29703a1367977d2867167fda8c4146aea6cd58e?style=flat-square">
    </p>
    <br>
    <strong>A general purpose discord bot that focuses on readability and developer-friendly<br></strong>
</div>

## install. <img alt="VimL" src="https://img.shields.io/badge/-Atlas-66d124?style=flat-square&logo=mongoDB&logoColor=white">&nbsp;<img alt="Go" src="https://img.shields.io/badge/-discordgo-46a2f1?style=flat-square&logo=go&logoColor=white">&nbsp;<img alt="Heroku" src="https://img.shields.io/badge/-Heroku-430098?style=flat-square&logo=heroku&logoColor=white">&nbsp;<img alt="git" src="https://img.shields.io/badge/-Github Actions-000000?style=flat-square&logo=GitHub&logoColor=white">

```sh 
$ go install github.com/Iridaceae/iridaceae/cmd/iridaceae-server
```

## folder structures.
```bash
.
├── bin
├── cmd           # bot cli, including both iridaceae and concertina        # lg: Go
├── internal      # internal core of iridaceae, logging, config handling    # lg: Go
│   └── jog       # iridaceae' internal command parser                      # lg: Go
├── pkg                                                                     
│   └── core      # bot logics, including message queries with configstore  # lg: Go
├── web           # bots front-facing api,                                  # lg: Rust
└── scripts
```

## inspirations.
> <div align="left"><i>why <strong>Go</strong> not <strong>TS/JS</strong>?</i></div>
- TypeScript and JavaScripts are amazing, and the fact that web elements for iridaceae will use `React` and `Next.js`, refers to [here](https://github.com/TensRoses/dashboard)
- Why Go then? Go is an awesome language with its interface and concurrency I just want to build a bot in Go :smile:
- Also, front-facing API will be written in Rust as a way to learn the language.

> <div align="left"><i>why <strong>iridaceae</strong> when we already had </i><a href="https://github.com/jonas747/yagpdb"><strong>yagpdb</strong></a>?</div>
- `yagpdb` are pretty clunky and bloated in the sense that the maintainer uses a lot of boilerplate, whereas `iridaceae` goals are to stay minimal and bloat-free as possible
- `yagpdb` structures are pretty messy and harder for development and thus `iridaceae` aims to make it more developer-friendly as well as user-friendly
- `iridaceae` got a lot of inspirations from `yagpdb` but takes a more lightweight approaches considering structures and deployments, refers to [here](pkg/README.md) for more details
- `iridaceae` is pretty much <strong>WIP</strong> so any help are appreciated
