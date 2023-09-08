# HTMX-TodoApp

This is a simple todo application built with Golang and HTMX.

[Live](https://htmx-todo-23.fly.dev/)

This site is functional as an app with only 2 lines of JS!

Stack:
- HTMX
- Golang
- Gin
- SqlLite

Features:
- Create, delete and update todos.
- Login and signup with JWT authentication.

### Running the app locally.

Head over to the [turso](https://docs.turso.tech/) docs to install the cli tool.

The following is only for Macos.

install the libsql server to your machine 
```bash
brew tap libsql/sqld
brew install sqld-beta
```

```bash
turso dev
```

paste the url into a .env file. See example.env.

```bash
APP_ENV=DEV go run main.go
```




## Deployment to fly

Create a fly.toml config file.

```toml
app = <APP_NAME>
primary_region = <REGION>

[build]
  builder = "paketobuildpacks/builder:base"
  buildpacks = ["gcr.io/paketo-buildpacks/go"]

[build.args]
    BP_KEEP_FILES = "templates/*.html"

[env]
  PORT = "8080"
  URL = <TURSO_DB_URL>
  APP_ENV = "PROD"

[http_service]
  internal_port = 8080
  force_https = true
  auto_stop_machines = true
  auto_start_machines = true
  min_machines_running = 0
  processes = ["app"]
```

Create a database with turso and get the URL and auth token.

You set the url as an env in the fly config.

Set the auth token as a secret through the flycli.

Deploy with fly.