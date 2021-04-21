# iris.

Pomodoro bot. Built with  :heart: and `discordgo`

### fyi.

Stack: mongoDB - discordgo - Heroku 

cmd line version named after herself `tensrose` will start Iris

### install.

to install `tensrose` binary run 

```sh 
$ go install github.com/aarnphm/iris/cmd/tensrose
```

### todo.
- pkg to construct message
- fix time difference 

--> full on microservices
services
- messageBuilder
- databaseHandler
  - DAO
  - CRUD
  - Repo
- commandHandler
  - this contains error handling
  - given prefix setup different commands
- cache
- metrics
- leaderboard
