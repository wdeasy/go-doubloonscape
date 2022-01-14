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

        game.changeCaptainsInGameAndServer(m.GuildID, m.Author.ID)
    }
}

//update the bot's message
func (game *Game) setMessage(){ 
    messages := game.getMessages()
    if messages == nil {
        return
    }
    
    embed := game.generateEmbed()

    if (messages[0].Author.ID == game.currentBotID) {
        game.editMessage(&embed, messages[0])
        return
    } 
    
    if game.currentMessageID != "" {
        game.deleteMessage(game.currentMessageID)
    }

    tempMessageId := game.currentMessageID
    game.currentMessageID = game.newMessage(&embed)

    for _, s := range messages {
        if (s.ID == tempMessageId) {
            continue
        }

        if (s.Author.ID == game.currentBotID) {
            game.deleteMessage(s.ID)		
        }
    }
}

//delete a discord message
func (game *Game) deleteMessage(messageID string) {
    err := game.dg.ChannelMessageDelete(Channel, messageID)
    if err != nil {
        printLog(fmt.Sprintf("error while setting message. could not delete message %s: %s\n", messageID, err))
    }	    
}

//get previous discord messages
func (game *Game) getMessages() ([]*discordgo.Message){
    messages, err := game.dg.ChannelMessages(Channel, MESSAGE_MAX, "", "", "")
    if err != nil {
        printLog(fmt.Sprintf("could not get messages from channel %s: %s\n", Channel, err))
    }
    
    return messages
} 

//edit the existing bot message
func (game *Game) editMessage(embed *discordgo.MessageEmbed, message *discordgo.Message) { 
    _, err := game.dg.ChannelMessageEditEmbed(Channel, message.ID, embed)	
    if err != nil {
        printLog(fmt.Sprintf("could not edit message: %s\n", err))
    }

    game.checkEventReactions(message)
}

//create a new bot message
func (game *Game) newMessage(embed *discordgo.MessageEmbed) (string) { 
    game.treasure.Up = false

    msg, err := game.dg.ChannelMessageSendEmbed(Channel, embed)	
    if err != nil {
        printLog(fmt.Sprintf("could not create new message: %s\n", err))
        return ""
    }

    game.addReactions(msg)

    return msg.ID
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