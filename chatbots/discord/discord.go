package discord

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"strings"
	"github.com/bwmarrin/discordgo"
	"github.com/plasmakatt/bdobot/utils/timeconversion"
	"github.com/plasmakatt/bdobot/gametimers/nighttimer"
	"github.com/plasmakatt/bdobot/gametimers/energytimer"
	"github.com/plasmakatt/bdobot/gametimers/imperialtimer"
)

type DiscordBot struct {
	Token string
}

func (discord DiscordBot) Run() {
	bot, err := discordgo.New("Bot " + discord.Token)
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		os.Exit(1)
	}
	bot.AddHandler(MessageCreate)
	err = bot.Open()
	if err != nil {
		fmt.Println("Error opening connection: ", err)
		os.Exit(1)
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	bot.Close()
	fmt.Println("Ctrl+C detected, shutting down")
}

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Do nothing if it's our msg
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m != nil {
		s.ChannelMessageSend(m.ChannelID, GetMessageToSend(m.Content))
	}
}

func GetMessageToSend(m string) string {
	var message string
	switch x := strings.Split(m, " "); {
		case x[0] == "!night":
			timer := nighttimer.New()
			if timer.IsDay {
				message = "Night in " + timeconversion.GetHMSFromSeconds(timer.SecondsUntilNightStart)
			} else {
				message = "Currently Night, will be day in " + timeconversion.GetHMSFromSeconds(timer.SecondsUntilNightEnd)
			}
		case x[0] == "!cooking":
			remSecs := imperialtimer.NewImperialCooking().SecondsUntilReset
			message = "Imperial cooking will reset in " + timeconversion.GetHMSFromSeconds(remSecs) + "\nWill reset at " + timeconversion.GetDateAfterSeconds(remSecs)
		case x[0] == "!trading":
			remSecs := imperialtimer.NewImperialTrading().SecondsUntilReset
			message = "Imperial trading will reset in " + timeconversion.GetHMSFromSeconds(remSecs) + "\nWill reset at " + timeconversion.GetDateAfterSeconds(remSecs)
		case x[0] == "!energy":
			if len(x) == 3 {
				entimer := energytimer.EnergyTimer{CurrentEnergy: x[1], MaxEnergy: x[2]}
				remSecs := entimer.GetRemainingSeconds()
				message = "Time to fill energy: " + timeconversion.GetHMSFromSeconds(remSecs) + "\nWill be full at " + timeconversion.GetDateAfterSeconds(remSecs)
			} else {
				message = "Usage: !energy <current energy> <max energy>"
			}
		case x[0] == "!bosses":
			message = "Not yet implemented!"
		case x[0] == "!hemligt":
			message = "Väldigt hemligt meddelande, shhh!"
		case x[0] == "!commands":
			message = "!night !cooking !trading !energy"
		default:
		}
	return message
}