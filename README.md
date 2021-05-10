# discord-reminder-bot
![build](https://github.com/TwinProduction/discord-reminder-bot/workflows/build/badge.svg?branch=master)

This is a simple Discord bot for managing reminders.

If you react to a message with `:timer:`, a message will be sent to you through direct message informing you
that you will be reminded about the message you sent in 8 hours.

You may also use the following syntax:
```
!RemindMe <DURATION> [NOTE]
```
**Where:**
- `<DURATION>` must be use one of the following formats: `30m`, `6h30m`, `48h`
- `[NOTE]` is an optional note to attach to the reminder with less than 240 characters

Note that `!RemindMe` can be replaced by directly pinging the bot (e.g. `@reminder-bot 2h30m meeting about cookies`)


## Usage

| Environment variable | Description | Required | Default |
|:--- |:--- |:--- |:--- |
| DISCORD_BOT_TOKEN | Discord bot token | yes | `""` |
| COMMAND_PREFIX | Character prepending all bot commands. Must be exactly 1 character, or it will default to `!` | no | `!` |


## Getting started

### Discord

1. Create an application
2. Add a bot in your application
3. Save the bot's token and set it as the `DISCORD_BOT_TOKEN` environment variable
4. Go to `https://discordapp.com/oauth2/authorize?client_id=<YOUR_BOT_CLIENT_ID>&scope=bot&permissions=2112`
5. Add the bot to a server of your choice
