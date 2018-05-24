package discord

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/plasmakatt/bdobot/gametimers/energytimer"
	"github.com/plasmakatt/bdobot/gametimers/imperialtimer"
	"github.com/plasmakatt/bdobot/gametimers/nighttimer"
	"github.com/plasmakatt/bdobot/utils/timeconversion"
)

/*
TODOS:
- split out message handling
- split out notifications
- read channelid from config
*/

const CommandPrefix = "!"
const CommandNight = CommandPrefix + "night"
const CommandCooking = CommandPrefix + "cooking"
const CommandTrading = CommandPrefix + "trading"
const CommandEnergy = CommandPrefix + "energy"
const NotificationStart = "start"
const NotificationStop = "stop"

// This is a constant, go please...
var Commands = [...]string{CommandNight, CommandCooking, CommandTrading, CommandEnergy}
var nightcallbacks = make(map[string]string)
var energycallbacks = make(map[string]time.Timer)
var ChannelID = "enterchannelidwhereyouwantyourbot"

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
	go NightNotifier(bot)
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

func NightNotifier(bot *discordgo.Session) {
	timer := nighttimer.New()
	notifySentThisNight := false
	i := 0
	for true {
		for k, _ := range energycallbacks {
			fmt.Println(i)
			fmt.Println(k)
			i++
		}
		if !timer.IsDay && !notifySentThisNight {
			for k, _ := range nightcallbacks {
				SendNotification(bot, k, "\nNight is here bro!")
			}
			notifySentThisNight = true
		} else {
			notifySentThisNight = false
		}
		time.Sleep(1000 * time.Millisecond)
	}
}

func SendNotification(bot *discordgo.Session, userId string, message string) {
	bot.ChannelMessageSend(ChannelID, "<@"+userId+"> "+message)
}

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Dont answer our own messages and ignore other chats
	if m.Author.ID == s.State.User.ID || m.ChannelID != ChannelID {
		return
	}
	if m != nil {
		s.ChannelMessageSend(ChannelID, GetMessageToSend(m.Content, m.Author, s))
	}
}

func GetMessageToSend(m string, user *discordgo.User, bot *discordgo.Session) string {
	var message string
	switch x := strings.Split(m, " "); {
	case x[0] == CommandNight:
		message = HandleNightMessages(x, user)
	case x[0] == CommandCooking:
		message = HandleCookingMessages()
	case x[0] == CommandTrading:
		message = HandleTradingMessages()
	case x[0] == CommandEnergy:
		message = HandleEnergyMessages(x, user, bot)
	default:
		message = "Available commands: " + GetCommands() + "\nNight and energy have notifications!" +
			"\nUsage: !night <" + NotificationStart + "/" + NotificationStop + ">" +
			"\nExample: !night start" +
			"\nUsage: !energy <current energy> <max energy> <" + NotificationStart + "/" + NotificationStop + ">" +
			"\nExample: !energy 150 300 start"
	}
	return message
}

func GetCommands() string {
	commands := ""
	for _, c := range Commands {
		commands += c + " "
	}
	return commands
}

func HandleNightMessages(m []string, user *discordgo.User) string {
	message := ""
	timer := nighttimer.New()
	if len(m) >= 2 {
		message = HandleNightNotifications(m[1], timer, user)
	} else {
		if timer.IsDay {
			message = "Night in " + timeconversion.GetHMSFromSeconds(timer.SecondsUntilNightStart)
		} else {
			message = "Currently Night, will be day in " + timeconversion.GetHMSFromSeconds(timer.SecondsUntilNightEnd)
		}
	}
	return message
}

func HandleNightNotifications(m string, timer nighttimer.NightTimer, user *discordgo.User) string {
	message := ""
	if m == NotificationStart {
		nightcallbacks[user.ID] = user.Username
		if timer.IsDay {
			message = "Notify added for user " + user.Username + ".\nWill alert at night."
		} else {
			message = "Notify added for user " + user.Username + ".\nCurrently night will notify at next."
		}
	} else if m == NotificationStop {
		delete(nightcallbacks, user.ID)
		message = "Notify stopped for user " + user.Username + ".\nWill not alert at nights."
	}
	return message
}

func HandleCookingMessages() string {
	remSecs := imperialtimer.NewImperialCooking().SecondsUntilReset
	return "Imperial cooking will reset in " + timeconversion.GetHMSFromSeconds(remSecs) + "\nWill reset at " + timeconversion.GetDateAfterSeconds(remSecs)
}

func HandleTradingMessages() string {
	remSecs := imperialtimer.NewImperialTrading().SecondsUntilReset
	return "Imperial trading will reset in " + timeconversion.GetHMSFromSeconds(remSecs) + "\nWill reset at " + timeconversion.GetDateAfterSeconds(remSecs)
}

func HandleEnergyMessages(m []string, user *discordgo.User, bot *discordgo.Session) string {
	message := ""
	if len(m) == 2 && m[1] == NotificationStop {
		energycallbacks[user.ID] = *time.NewTimer(0)
		<-energycallbacks[user.ID].C
		delete(energycallbacks, user.ID)
		message = "Notify stopped for user " + user.Username + ".\nWill not alert energy status."
	} else if len(m) == 3 {
		entimer := energytimer.EnergyTimer{CurrentEnergy: m[1], MaxEnergy: m[2]}
		remSecs := entimer.GetRemainingSeconds()
		message = "Time to fill energy: " + timeconversion.GetHMSFromSeconds(remSecs) + "\nWill be full at " + timeconversion.GetDateAfterSeconds(remSecs)
	} else if len(m) == 4 {
		if m[3] == NotificationStart {
			entimer := energytimer.EnergyTimer{CurrentEnergy: m[1], MaxEnergy: m[2]}
			remSecs := entimer.GetRemainingSeconds()
			energycallbacks[user.ID] = *time.NewTimer(time.Second * time.Duration(remSecs))
			go func() {
				<-energycallbacks[user.ID].C
				SendNotification(bot, user.ID, "\nEnergy is full bro!")
				delete(energycallbacks, user.ID)
			}()
			message = "Notify added for user " + user.Username + ".\nWill alert when energy is full."
		} else if m[3] == NotificationStop {
			energycallbacks[user.ID] = *time.NewTimer(0)
			<-energycallbacks[user.ID].C
			delete(energycallbacks, user.ID)
			message = "Notify stopped for user " + user.Username + ".\nWill not alert energy status."
		}
	} else {
		message = "Usage: !energy <current energy> <max energy>"
	}
	return message
}
