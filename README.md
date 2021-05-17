# TwitchChatBot

This bot will connect to a channel at Twitch and lookup Magic card details
when people type specific commands into the chat.

## Usage:
    TwitchChatBot -config "app.config"

!card cardName
- This will lookup the card with card name and reply to the chat with that card's information. If there are multiple cards that match your cardName, only the first card will be displayed

!changechannel channelName
- This command is only accepted when sent as a whisper to magiccardbot. It will cause the bot to leave the current channel and switch to the specified channel. If channelName is not a real channel, then the bot will not be listening on any channel except for whispers.

!shutdown
- This command is only accepted when sent as a whisper to magiccardbot. This will cause the bot to shutdown.
