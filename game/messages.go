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
    if !matched {
        return
    }

    if _, ok := game.captains[m.Author.ID]; !ok {
        game.addCaptainFromDiscordMessage(m.Author.ID, m.Member.Nick, m.Author.Username)
    }

    game.changeCaptainsInGameAndServer(m.GuildID, m.Author.ID)
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

//return the correct discord name
func getName(nick string, user string) (string) {
    var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9 ]+`)

    if (nick != "") {
        return nonAlphanumericRegex.ReplaceAllString(nick, "")
    } 

    return nonAlphanumericRegex.ReplaceAllString(user, "")
}

//create a new bot message
func (game *Game) newMessage(embed *discordgo.MessageEmbed) (string) { 
    //msg, err := game.dg.ChannelMessageSendEmbed(Channel, embed)	
    msg, err := game.dg.ChannelMessageSendComplex(Channel, game.generateMessageSend(embed))
    
    
    if err != nil {
        printLog(fmt.Sprintf("could not create new message: %s\n", err))
        return ""
    }

    game.addReactions(msg)

    return msg.ID
}

//edit the existing bot message
func (game *Game) editMessage(embed *discordgo.MessageEmbed, message *discordgo.Message) { 
    msg, err := game.dg.ChannelMessageEditComplex(game.generateMessageEdit(message.ID, embed))
    if err != nil {
        printLog(fmt.Sprintf("could not edit message: %s\n", err))
        return
    }

    game.addReactions(msg)
}

//generate MessageSend object for complex message send
func (game *Game) generateMessageSend(embed *discordgo.MessageEmbed) (*discordgo.MessageSend) {
    message := discordgo.MessageSend{
        Content: *game.stats.Log,
        Embed: embed,
    }

    return &message
}

//generate MessageEdit object for complex message edit
func (game *Game) generateMessageEdit(messageID string, embed *discordgo.MessageEmbed) (*discordgo.MessageEdit) {
    message := discordgo.MessageEdit{
        Content: game.stats.Log,
        Embed: embed,
        ID: messageID,
        Channel: Channel,
    }

    return &message
}