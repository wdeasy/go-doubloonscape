package game

import (
    "fmt"
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
        printLog(fmt.Sprintf("error while checking if captain %s exists: %s\n", m.UserID, err))
        return
    }   
    
    switch m.Emoji.Name {
        case INCREMENT_REACTION:
            game.incrementCaptainByOne()
        case CAPTAIN_REACTION:
            game.changeCaptainsInGameAndServer(m.GuildID, m.UserID)
        case PRESTIGE_REACTION:
            game.setPrestige(m.UserID)
        case TREASURE_REACTION:           
            game.giveTreasure(m.UserID)
        case PICKPOCKET_REACTION:             
            game.EventReactionReceived(PICKPOCKET_NAME, m.UserID)
    }
    
    if m.Emoji.Name != INCREMENT_REACTION {
        game.removeReaction(m.MessageID, m.Emoji.Name, game.currentBotID)  
    }
    
    game.removeReaction(m.MessageID, m.Emoji.Name, m.UserID)
}

//add reactions to a message
func (game *Game) addReactions(message *discordgo.Message) {
    reactions := []string{INCREMENT_REACTION, PRESTIGE_REACTION, CAPTAIN_REACTION}
    for _, e := range reactions {
        game.addReaction(message.ID, e)
    } 

    game.checkReactions(message)
}

//see if additional reactions can be added
func (game *Game) checkReactions(message *discordgo.Message) { 
    game.checkTreasureReaction(message)
    game.checkEventReactions(message)
}

//logic to see if event reactions should be up
func (game *Game) checkEventReactions(message *discordgo.Message) {
    for _, e := range game.events {
        eventReaction := game.getReaction(e.Name)
        if game.isReactionInReactions(eventReaction, message.Reactions) {
            if !e.Up {
                game.removeReaction(message.ID, eventReaction, game.currentBotID)
            }
            return
        }

        if e.Up {
            game.addReaction(message.ID, eventReaction)
            return
        }

        if e.Ready(game.getCooldown(e.Name)) {
            e.Up = true        
            game.addReaction(message.ID, eventReaction)
        }
    }
}

//logic to see if treasure reaction should be up
func (game *Game) checkTreasureReaction(message *discordgo.Message) {
    if game.isReactionInReactions(TREASURE_REACTION, message.Reactions) {
        if !game.treasure.Up {
            game.removeReaction(message.ID, TREASURE_REACTION, game.currentBotID)
        }
        return
    }

    if game.treasure.Up {
        game.addReaction(message.ID, TREASURE_REACTION)	
        return
    }

    if treasureChance() {
        game.treasure.Up = true
        game.addReaction(message.ID, TREASURE_REACTION)		
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

//get the appropriate reaction from constants
func (game *Game) getReaction(name string) (string) {
    switch name {
    case PICKPOCKET_NAME:
        return PICKPOCKET_REACTION
    default:
        return DEFAULT_REACTION
    }
}