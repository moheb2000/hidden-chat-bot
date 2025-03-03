# Hidden Chat Bot
![License](https://img.shields.io/github/license/moheb2000/hidden-chat-bot)
![go.mod Go version](https://img.shields.io/github/go-mod/go-version/moheb2000/hidden-chat-bot)

A telegram bot written in go that let users creates anonymous message links for others to recieve other messages without revieling their usernames.

## Installation Guide
> [!NOTE]
> Development take place in main branch, so it may have bugs or doesn't run at all! If you want to build this project from source, change the cloned repository to a version tag before starting to build!
### Prerequisites
Hidden chat bot stores data on postgres. So you need to install it on your server and create a database for connecting to it.

Also you need to get a bot token from telegram bot father.

Then create a `.envrc` files to store environment variables. You can use `.envrc_example` for this:
```
cp .envrc_example .envrc
```
And then add bot token and postgres dsn to it. Your dsn may look like this:
```
postgres://<username>:<password>@localhost/<database_name>
```
### Running migration files
Run database migration files with this command:
```
make migration/up
```
You can clear your database data with this command:
```
make migration/down
```
### Running the bot
For running the bot in development mode, run:
```
make run/bot name="<Bot Name>" locale="<en-us or fa-ir>"
```
And for building the bot and getting an executable file, run:
```
make build/bot
```
And the binary file will be in `bin` directory.
