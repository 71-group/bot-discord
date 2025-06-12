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
	b := bot.GetBot() // GetBot já faz b.Open(), não precisa abrir de novo!

	// Remova ou comente esta linha:
	// err = b.Open()
	// if err != nil {
	// 	fmt.Println("Erro ao conectar o bot:", err)
	// 	return
	// }
	defer b.Close()
	fmt.Println("Bot do Discord está online!")

	// Registra comandos quando o bot fica pronto
	b.AddHandlerOnce(func(s *discordgo.Session, r *discordgo.Ready) {
		for _, guild := range s.State.Guilds {
			cmd, err := s.ApplicationCommandCreate(cfg.Application, guild.ID, &discordgo.ApplicationCommand{
				Name:        "ping",
				Description: "Responde com Pong!",
			})
			if err != nil {
				fmt.Printf("Erro ao criar comando em guild %s: %v\n", guild.ID, err)
			} else {
				fmt.Printf("Comando /ping registrado em guild %s: %+v\n", guild.ID, cmd)
			}
		}

		// Adiciona o comando /ticket
		for _, guild := range s.State.Guilds {
			_, err := s.ApplicationCommandCreate(cfg.Application, guild.ID, &discordgo.ApplicationCommand{
				Name:        "ticket71", // <-- novo nome aqui
				Description: "Abre um ticket privado para suporte.",
			})
			if err != nil {
				fmt.Printf("Erro ao criar comando /ticket em guild %s: %v\n", guild.ID, err)
			} else {
				fmt.Printf("Comando /ticket registrado em guild %s\n", guild.ID)
			}
		}
	})

	// REGISTRE O HANDLER DE INTERAÇÃO AQUI, FORA DO READY!
	b.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.ApplicationCommandData().Name {
		case "ping":
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Pong!",
				},
			})
		case "ticket71": // <-- novo nome aqui
			user := i.Member.User
			guildID := i.GuildID
			channelName := fmt.Sprintf("ticket-%s", user.Username)

			// Cria canal privado
			channel, err := s.GuildChannelCreateComplex(guildID, discordgo.GuildChannelCreateData{
				Name: channelName,
				Type: discordgo.ChannelTypeGuildText,
				PermissionOverwrites: []*discordgo.PermissionOverwrite{
					{
						ID:   guildID,
						Type: discordgo.PermissionOverwriteTypeRole,
						Deny: discordgo.PermissionViewChannel,
					},
					{
						ID:    user.ID,
						Type:  discordgo.PermissionOverwriteTypeMember,
						Allow: discordgo.PermissionViewChannel,
					},
				},
			})
			if err != nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Erro ao criar o ticket.",
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
