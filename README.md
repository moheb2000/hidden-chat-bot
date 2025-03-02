# Hidden Chat Bot
**Notice: This bot is not ready to use in production**

A telegram bot written in go that let users creates anonymous message links for others to recieve other messages without revieling their usernames.

## Installation Guide
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
make run/bot
```
Or this for persian locale:
```
make run/bot/fa
```
And for building the bot and getting an executable file, run:
```
make build/bot
```
And the binary file will be in `bin` directory.
