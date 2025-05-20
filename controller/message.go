package controller

import (
	"bot-discord/helper/bot"
	"fmt"

	"github.com/gin-gonic/gin"
)

type ChannelsStruct struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func GetMessageList(c *gin.Context) {
	Guilds := bot.GetBot().State.Guilds
	IChannels := []ChannelsStruct{}
	for _, guild := range Guilds {
		channels, err := bot.GetBot().GuildChannels(guild.ID)
		if err != nil {
			continue // Ignora guilds que não puderam ser listadas
		}
		for _, channel := range channels {
			// Apenas canais de texto (Type 0)
			if channel.Type == 0 {
				IChannels = append(IChannels, ChannelsStruct{
					ID:   channel.ID,
					Name: channel.Name,
				})
			}
		}
	}
	context := gin.H{
		"Channels": IChannels,
	}
	c.HTML(200, "message.html", context)
}

func PostMessage(c *gin.Context) {
	channelID := c.Param("channel_id")
	message := c.PostForm("message")
	_, err := bot.GetBot().ChannelMessageSend(channelID, message)
	if err != nil {
		fmt.Println("Error sending message:", err)
		c.JSON(500, gin.H{"error": "Failed to send message"})
		return
	}
	c.JSON(200, gin.H{"status": "ok"})
}
