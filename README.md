# discord-reminder-bot

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
