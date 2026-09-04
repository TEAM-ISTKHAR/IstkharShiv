# IstkharShiv

Production-ready Telegram music bot written in Go. The repository includes one
runtime, one Docker image, a Heroku worker definition, and a Railway deployment
definition so the same build runs consistently on both platforms.

[![Deploy to Heroku](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy?template=https://github.com/TEAM-ISTKHAR/IstkharShiv)
[![Deploy on Railway](https://railway.com/button.svg)](https://railway.com/new/template)

**Heroku:** click the button above to open the one-click deploy flow using
`app.json`. **Railway:** click the button above, choose **Deploy from GitHub
repo**, select `TEAM-ISTKHAR/IstkharShiv`, and add the variables from
`sample.env`.

## Runtime requirements

- Go 1.25+ for local builds
- `ffmpeg` and `ffprobe`
- `yt-dlp`
- A Telegram bot token
- MongoDB for persistent state (recommended)

The Docker image installs the media dependencies automatically.

## Required environment

Copy `sample.env` to `.env` for local development. At minimum set:

```text
BOT_TOKEN=...
MONGO_DB_URI=mongodb+srv://...
OWNER_ID=123456789
```

`API_ID` and `API_HASH` are kept in the template for Telegram integrations.
The bot can start without MongoDB in in-memory mode, but data will be lost on
restart.

## Local run

```bash
go mod download
go run .
```

## Docker

```bash
docker build -t istkharshiv .
docker run --env-file .env --restart unless-stopped istkharshiv
```

## Heroku worker

This repository uses a **worker dyno**, not a web dyno. It does not sleep, and
the `heroku.yml`/`Procfile` configuration starts `/app/istkharshiv`.

```bash
heroku create
heroku stack:set container -a YOUR_APP
heroku config:set BOT_TOKEN=... MONGO_DB_URI=... OWNER_ID=... -a YOUR_APP
heroku container:push worker -a YOUR_APP
heroku container:release worker -a YOUR_APP
heroku ps:scale worker=1 -a YOUR_APP
```

Heroku dyno scaling is controlled by the Heroku account, not by application
code. The final `ps:scale` command turns the worker on and keeps it running.

## Railway VPS

Create a Railway service from this repository. Railway detects `Dockerfile`
and `railway.json`; it builds the same image and restarts the worker on failure.
Add the variables from `sample.env` in the Railway Variables panel, then deploy.

## Safety

Never commit `.env`, bot tokens, MongoDB URLs, API keys, or Telegram string
sessions. All provider credentials are read from environment variables.

## License

Use this project only where your use complies with Telegram, YouTube, and
third-party service terms.