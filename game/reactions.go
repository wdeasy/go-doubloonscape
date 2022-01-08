package game

import (
	"fmt"
	"math/rand"
	"net/url"

	"github.com/bwmarrin/discordgo"
)

//handler for new reactions in discord
func (game *Game) messageReactionAdd(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
    if m.UserID == s.State.User.ID || m.ChannelID != Channel {
        return
    }

	err := game.checkIfCaptainExists(m.UserID, m.GuildID)
    if err != nil {
		fmt.Printf("error while checking if captain %s exists: %s\n", m.UserID, err)
		return
    }   
	
    switch m.Emoji.Name {
        case "🪙":
            game.coinEmoji()
        case "🏴‍☠️":
            game.pirateEmoji(m.UserID, m.GuildID)
        case "🔱":
            game.tridentEmoji(m.UserID)
        case "👑":
            game.crownEmoji(m.UserID)
    }	
    
    err = game.dg.MessageReactionRemove(Channel, m.MessageID, url.QueryEscape(m.Emoji.Name), m.UserID)
    if err != nil {
		fmt.Printf("could not remove reaction %s from %s: %s\n", m.Emoji.Name, game.captains[m.UserID].Name, err)
    }    

}

//add reactions to a message
func (game *Game) addReactions(message *discordgo.Message) {
    err := game.dg.MessageReactionAdd(Channel, message.ID, url.QueryEscape("🪙"))
    if err != nil {
        fmt.Printf("could not add 🪙 reaction to new message %s: %s\n", message.ID, err)
    }	

    err = game.dg.MessageReactionAdd(Channel, message.ID, url.QueryEscape("🔱"))
    if err != nil {
        fmt.Printf("could not add 🔱 reaction to new message %s: %s\n", message.ID, err)
    }	

    err = game.dg.MessageReactionAdd(Channel, message.ID, url.QueryEscape("🏴‍☠️"))
    if err != nil {
        fmt.Printf("could not add 🏴‍☠️ reaction to new message %s: %s\n", message.ID, err)
    }	

    if rand.Intn(100) == 2 {
        treasure = true
        err = game.dg.MessageReactionAdd(Channel, message.ID, url.QueryEscape("👑"))
		if err != nil {
			fmt.Printf("could not add 👑 reaction to new message %s: %s\n", message.ID, err)
		}			
    }
}

//function for when the coin(🪙) reaction is used
func (game *Game) coinEmoji() {
	game.incrementCaptain()
}

//function for when the pirate(🏴‍☠️) reaction is used
func (game *Game) pirateEmoji(UserID string, GuildID string) {
	if UserID == game.currentCaptainID {
		return
	}

	game.newCaptain(GuildID, UserID)
}

//function for when the trident(🔱) reaction is used
func (game *Game) tridentEmoji(UserID string) {
	if UserID != game.currentCaptainID {
		return
	}

	game.addPrestige(UserID)
	game.setMessage()		
}

//function for when the crown(👑) reaction is used
func (game *Game) crownEmoji(UserID string) {
	if !treasure {
		fmt.Printf("%s clicked on 👑 but the treasure is not up yet\n", UserID)
		return 
	} 

	treasure = false

	game.giveTreasure(UserID)
	game.setMessage()

	messages, err := game.dg.ChannelMessages(Channel, 100, "", "", "")
	if err != nil {
		fmt.Printf("could not retrieve messages for channel %s: %s\n", Channel, err)
		return
	}

	for _, s := range messages {
		if s.Author.ID == game.dg.State.User.ID {
			err := game.dg.MessageReactionRemove(Channel, s.ID, url.QueryEscape("👑"), game.dg.State.User.ID)

			if err != nil {
				fmt.Printf("could not remove 👑 reaction from bot message %s: %s\n", s.ID, err)
				return
			}
		}
	}
}

//create a captain in the map if it does not exist
func (game *Game) checkIfCaptainExists(UserID string, GuildID string) (error) {
	if _, ok := game.captains[UserID]; !ok {
		err := game.addCaptainFromDiscordReaction(GuildID, UserID)

		if err != nil {
			return fmt.Errorf("could not add captain %s from discord reaction: %w", UserID, err)
		}		
	}
	
	return nil
}

//create a new captain with info from the discord reaction
func (game *Game) addCaptainFromDiscordReaction(GuildID string, UserID string) (error) {
	m, err := game.dg.GuildMember(GuildID, UserID)
	if err != nil {
		return fmt.Errorf("could not get guild member info %s: %w", UserID, err)
	}
		
	name := getName(m.Nick, m.User.Username)
	
	if m.User.Bot {
		return fmt.Errorf("%s is a bot", name)
	}

	game.createCaptain(UserID, name)	

	return nil
}