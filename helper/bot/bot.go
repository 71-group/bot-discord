package bot

import (
	"botdiscord/helper"
	"log"
	"sync"

	"github.com/bwmarrin/discordgo"
)

var (
	bot    *discordgo.Session
	once   sync.Once
	botErr error
)

func GetBot() *discordgo.Session {
	once.Do(func() {
		cfg := helper.ReadConfig()

		bot, botErr = discordgo.New("Bot " + cfg.Token)
		if botErr != nil {
			log.Panicln("Erro ao criar sessão do Discord:", botErr)
		}

		botErr = bot.Open()
		if botErr != nil {
			log.Panicln("Erro ao abrir conexão com o Discord:", botErr)
		}
	})

	if bot == nil {
		log.Panicln("Sessão do bot é nula!")
	}
	return bot
}
