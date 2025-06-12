package main

import (
	"bot-discord/controller"
	"bot-discord/helper"
	"bot-discord/helper/bot"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/bwmarrin/discordgo"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.DebugMode)

	r := gin.New()
	files, err := loadTemplates("./website/html/")
	if err != nil {
		log.Println(err)
	}
	r.LoadHTMLFiles(files...)

	cfg := helper.ReadConfig()
	b := bot.GetBot()

	b.AddHandlerOnce(func(s *discordgo.Session, r *discordgo.Ready) {
		for _, guild := range s.State.Guilds {
			_, err := s.ApplicationCommandCreate(cfg.Application, guild.ID, &discordgo.ApplicationCommand{
				Name:        "ping",
				Description: "Responde com Pong!",
			})
			if err != nil {
				fmt.Printf("Erro ao criar comando em guild %s: %v\n", guild.ID, err)
			} else {
				fmt.Printf("Comando /ping registrado em guild %s\n", guild.ID)
			}
		}

		s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if i.ApplicationCommandData().Name == "ping" {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Pong!",
					},
				})
				return
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("Ticket criado: <#%s>", channel.ID),
				},
			})
			s.ChannelMessageSend(channel.ID, fmt.Sprintf("Olá %s! Descreva seu problema.", user.Mention()))
		}
	})

	r.GET("/", controller.Index)
	r.GET("/user-list", controller.GetUserList)
	r.GET("/message", controller.GetMessageList)
	r.POST("message/:channel_id", controller.PostMessage)

	r.Static("/static", "./website/static/")
	err = r.Run(":80")
	if err != nil {
		fmt.Println(err)
	}
}

func loadTemplates(root string) (files []string, err error) {
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		fileInfo, err := os.Stat(path)
		if err != nil {
			return err
		}
		if fileInfo.IsDir() {
			if path != root {
				if _, err := loadTemplates(path); err != nil {
					fmt.Println("Erro ao carregar templates em", path, ":", err)
					return err
				}
			}
		} else {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}
