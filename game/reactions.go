package game

import (
    "fmt"
    "net/url"

    "github.com/bwmarrin/discordgo"
)

type Reaction struct {
    Event     string
    Emoji     string
}

//list of valid reactions
func (game *Game) getReactions() ([]string) {
    return []string{CAPTAIN_REACTION, INCREMENT_REACTION, PICKPOCKET_REACTION, PRESTIGE_REACTION, TREASURE_REACTION}
}

//handler for new reactions in discord
func (game *Game) messageReactionAdd(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
    if m.UserID == s.State.User.ID || m.ChannelID != Channel {
        return
    }

    err := game.checkIfCaptainExists(m.UserID, m.GuildID)
    if err != nil {
        printLog(fmt.Sprintf("error while checking if captain %s exists: %s\n", m.UserID, err))
        return
    }   
    
    for _, reaction := range game.reactions {
        if m.Emoji.Name == reaction.Emoji {
            game.EventReactionReceived(game.reactions[m.Emoji.Name].Event, m.UserID, m.GuildID)
            game.removeReaction(m.MessageID, m.Emoji.Name, game.currentBotID)
        }
    }

    game.removeReaction(m.MessageID, m.Emoji.Name, m.UserID)
}

//logic to see if event reactions should be up and add reactions to message
func (game *Game) addReactions(message *discordgo.Message) {
    for _, reaction := range game.reactions {
        if game.isReactionInReactions(reaction.Emoji, message.Reactions) && !game.events[reaction.Event].Up {
            game.removeReaction(message.ID, reaction.Emoji, game.currentBotID)
            continue
        }

        if !game.isReactionInReactions(reaction.Emoji, message.Reactions) && game.events[reaction.Event].Up {
            game.addReaction(message.ID, reaction.Emoji)
            continue
        }
    }
}

//check if reaction is already in message
func (game *Game) isReactionInReactions(reaction string, reactions []*discordgo.MessageReactions) (bool){
    for _, r := range reactions {
        if r.Emoji.Name == reaction {
            return true
        }
    }
    
    return false
}

//add reaction to message
func (game *Game) addReaction(messageID string, reaction string) {
    err := game.dg.MessageReactionAdd(Channel, messageID, url.QueryEscape(reaction))
    if err != nil {
        printLog(fmt.Sprintf("could not add %s reaction to new message %s: %s\n", reaction, messageID, err))
    }
}

//remove reaction to message
func (game *Game) removeReaction(messageID string, reaction string, userID string) {
    err := game.dg.MessageReactionRemove(Channel, messageID, url.QueryEscape(reaction), userID)
    if err != nil {
        printLog(fmt.Sprintf("could not remove %s reaction from bot message %s: %s\n", reaction, messageID, err))
        return
    }	    
}

//load reactions into structs
func (game *Game) loadReactions() (map[string]*Reaction){
    reactions := make(map[string]*Reaction)
    for _, r := range game.getReactions() {
        reaction := new(Reaction)
        event := game.getReactionEvent(r)
        reaction.Event = event
        reaction.Emoji = r
        reactions[r] = reaction
    }  
    
    return reactions
}

//get the appropriate reaction from constants
func (game *Game) getReactionEvent(reaction string) (string) {
    switch reaction{
    case PICKPOCKET_REACTION:
        return PICKPOCKET_NAME
    case CAPTAIN_REACTION:
        return CAPTAIN_NAME
    case INCREMENT_REACTION:
        return INCREMENT_NAME
    case PRESTIGE_REACTION:
        return PRESTIGE_NAME
    case TREASURE_REACTION:
        return TREASURE_NAME      
    default:
        return DEFAULT_NAME
    }   
}