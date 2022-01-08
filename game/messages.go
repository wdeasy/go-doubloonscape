package game

import (
    "fmt"
    "regexp"

    "github.com/bwmarrin/discordgo"
)

//handler for new messages in discord
func (game *Game) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
    if m.Author.ID == s.State.User.ID || m.ChannelID != Channel || m.Author.Bot {
        return
    }

    matched, _ := regexp.MatchString(CAPTAIN_REGEX, m.Content)
    if matched {
        if _, ok := game.captains[m.Author.ID]; !ok {
            game.addCaptainFromDiscordMessage(m.Author.ID, m.Member.Nick, m.Author.Username)
        }

        game.newCaptain(m.GuildID, m.Author.ID)
    }
}

//update the bot's message
func (game *Game) setMessage(){ 
    messages, err := game.dg.ChannelMessages(Channel, 100, "", "", "")
    if err != nil {
        fmt.Printf("could not set message. error while getting messages from channel %s: %s\n", Channel, err)
        return
    }
    
    embed := game.generateEmbed()

    if (messages[0].Author.ID == game.dg.State.User.ID) {
        game.editMessage(&embed, messages[0].ID)
        return
    } 

    game.newMessage(&embed)
    for _, s := range messages {
        if (s.Author.ID == game.dg.State.User.ID) {
            err = game.dg.ChannelMessageDelete(Channel, s.ID)
            if err != nil {
                fmt.Printf("error while setting message. could not delete message %s: %s\n", s.ID, err)
            }			
        }
    }
}

//generate the embed for the bot's message
func (game *Game) generateEmbed() (discordgo.MessageEmbed) {
    embed := discordgo.MessageEmbed{
        Color: 0xf1c40f,
        Title: "𝔏𝔢𝔞𝔡𝔢𝔯𝔅𝔬𝔞𝔯𝔡",
        Description: stats.Leaderboard,
        Fields: []*discordgo.MessageEmbedField{
            {
                Name:   "ℭ𝔞𝔭𝔱𝔞𝔦𝔫",
                Value:  "` " + firstN(game.captains[game.currentCaptainID].Name, 31) + " `",
                Inline: true,
            },
            {
                Name:   "𝔇𝔬𝔲𝔟𝔩𝔬𝔬𝔫𝔰",
                Value:  "` " + fmt.Sprintf("%-7d",game.captains[game.currentCaptainID].Gold) + " `",
                Inline: true,
            },
            {
                Name:   "𝔓𝔯𝔢𝔰𝔱𝔦𝔤𝔢",
                Value:  "` " + fmt.Sprintf("%-4.3f",game.captains[game.currentCaptainID].Prestige) + " `",
                Inline: true,
            },
        },
        Footer: &discordgo.MessageEmbedFooter{
            Text:   stats.Event,
        },		
    }

    return embed
}

//edit the existing bot message
func (game *Game) editMessage(embed *discordgo.MessageEmbed, messageID string) { 
    _, err := game.dg.ChannelMessageEditEmbed(Channel, messageID, embed)	
    if err != nil {
        fmt.Printf("could not edit message: %s\n", err)
    }
}

//create a new bot message
func (game *Game) newMessage(embed *discordgo.MessageEmbed) { 
    treasure = false

    msg, err := game.dg.ChannelMessageSendEmbed(Channel, embed)	
    if err != nil {
        fmt.Printf("could not create new message: %s\n", err)
        return
    }

    game.addReactions(msg)
}

//create a new captain with info from the discord message
func (game *Game) addCaptainFromDiscordMessage(UserID string, Nick string, UserName string) {
    name := getName(Nick, UserName)
    game.createCaptain(UserID, name)	
}

//return the correct discord name
func getName(nick string, user string) (string) {
    if (nick != "") {
        return nick
    } 

    return user
}