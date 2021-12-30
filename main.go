package main

import (
    "fmt"
    "os"
    "os/signal"
    "syscall"
    "regexp"

    "github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters
var (
    Token string = os.Getenv("BOT_TOKEN")
    Role string = os.Getenv("ROLE")
    Channel string = os.Getenv("CHANNEL")
)

func main() {

    // Create a new Discord session using the provided bot token.
    dg, err := discordgo.New("Bot " + Token)
    if err != nil {
        fmt.Println("error creating Discord session,", err)
        return
    }

    // Register the messageCreate func as a callback for MessageCreate events.
    dg.AddHandler(messageCreate)

    // In this example, we only care about receiving message events.
    dg.Identify.Intents = discordgo.IntentsGuildMessages

    // Open a websocket connection to Discord and begin listening.
    err = dg.Open()
    if err != nil {
        fmt.Println("error opening connection,", err)
        return
    }

    dg.UserUpdateStatus(discordgo.StatusIdle)

    // Wait here until CTRL-C or other term signal is received.
    fmt.Println("Bot is now running. Press CTRL-C to exit.")
    sc := make(chan os.Signal, 1)
    signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
    <-sc

    // Cleanly close down the Discord session.
    dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

    // Ignore all messages created by the bot itself
    // This isn't required in this specific example but it's a good practice.
    if m.Author.ID == s.State.User.ID {
        return
    }

    if m.ChannelID != Channel {
        return
    }	

	matched, _ := regexp.MatchString(`\b[Ii][’']?[Mm][ \t]+[Tt][Hh][Ee][ \t]+[Cc][Aa][Pp][Tt][Aa][Ii][Nn][ \t]+[Nn][Oo][Ww].?\b`, m.Content)
	if matched {
		err := s.GuildMemberRoleAdd(m.GuildID, m.Author.ID, Role)

		if err != nil {
			fmt.Println(err)
			return
		} else {
			fmt.Printf("%s is the captain now.\n", m.Member.User.Username)
		}
		
		members, err := s.GuildMembers(m.GuildID, "", 1000)
		if err != nil {
			fmt.Println(err)
			return
		} else {
			fmt.Printf("Checking %d Guild Members for Captain Role.\n", len(members))
		}	

		for _, member := range members {
			for _, role := range member.Roles {
				if (role == Role && member.User.ID != m.Author.ID) {
					fmt.Printf("Removing the Captain role from %s.\n", member.User.Username)

					err := s.GuildMemberRoleRemove(m.GuildID, member.User.ID, Role)

					if err != nil {
						fmt.Println(err)
					}						
				}
			}
		}
	}
}