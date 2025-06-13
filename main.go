package main

import (
	"botdiscord/controller"
	"botdiscord/helper"
	"botdiscord/helper/bot"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/bwmarrin/discordgo"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

var (
	discordClientID    = os.Getenv("DISCORD_CLIENT_ID")
	discordRedirectURI = "http://localhost:80/callback"

	discordOauthConfig = &oauth2.Config{
		ClientID:     discordClientID,
		ClientSecret: os.Getenv("DISCORD_CLIENT_SECRET"),
		RedirectURL:  discordRedirectURI,
		Scopes:       []string{"identify"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://discord.com/api/oauth2/authorize",
			TokenURL: "https://discord.com/api/oauth2/token",
		},
	}
)

func main() {
	gin.SetMode(gin.DebugMode)

	r := gin.New()

	// Adicione o middleware de sessão aqui:
	store := cookie.NewStore([]byte("super-secret"))
	r.Use(sessions.Sessions("mysession", store))
	r.Use(UserSessionMiddleware())

	files, err := loadTemplates("./website/html/")
	if err != nil {
		log.Println(err)
	}
	r.LoadHTMLFiles(files...)

	cfg := helper.ReadConfig()
	b := bot.GetBot() // GetBot já faz b.Open(), não precisa abrir de novo!
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
	r.GET("/api/users", controller.GetUserList)
	r.GET("/message", controller.GetMessageList)
	r.POST("message/:channel_id", controller.PostMessage)

	// Serve arquivos estáticos (CSS, JS, imagens)
	r.Static("/static", "./website/static")

	// Serve páginas HTML
	r.GET("/user-list", func(c *gin.Context) {
		c.File("./website/html/user-list.html")
	})

	// Rota para iniciar login com Discord
	r.GET("/login-discord", func(c *gin.Context) {
		clientID := os.Getenv("DISCORD_CLIENT_ID")
		redirectURI := "http://localhost:80/callback"
		url := fmt.Sprintf(
			"https://discord.com/oauth2/authorize?client_id=%s&redirect_uri=%s&response_type=code&scope=identify&state=state",
			clientID,
			redirectURI,
		)
		c.Redirect(http.StatusFound, url)
	})

	r.GET("/callback", func(c *gin.Context) {
		session := sessions.Default(c)
		code := c.Query("code")
		token, err := discordOauthConfig.Exchange(c, code)
		if err != nil {
			c.String(500, "Erro ao trocar código: %v", err)
			return
		}

		// Buscar dados do usuário
		client := discordOauthConfig.Client(c, token)
		resp, err := client.Get("https://discord.com/api/users/@me")
		if err != nil {
			c.String(500, "Erro ao buscar usuário: %v", err)
			return
		}
		defer resp.Body.Close()

		var user struct {
			ID            string `json:"id"`
			Username      string `json:"username"`
			Discriminator string `json:"discriminator"`
			Avatar        string `json:"avatar"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
			c.String(500, "Erro ao decodificar usuário: %v", err)
			return
		}

		// Salva dados na sessão
		session.Set("user_id", user.ID)
		session.Set("username", user.Username)
		session.Set("avatar", user.Avatar)
		session.Save()

		c.Redirect(http.StatusFound, "/")
	})

	r.GET("/logout", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Clear()
		session.Save()
		c.Redirect(http.StatusFound, "/")
	})

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

func UserSessionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		c.Set("UserID", session.Get("user_id"))
		c.Set("Username", session.Get("username"))
		c.Set("Avatar", session.Get("avatar"))
		c.Next()
	}
}

func RenderWithUser(c *gin.Context, code int, name string, data gin.H) {
	if data == nil {
		data = gin.H{}
	}
	data["UserID"] = c.MustGet("UserID")
	data["Username"] = c.MustGet("Username")
	data["Avatar"] = c.MustGet("Avatar")
	c.HTML(code, name, data)
}
