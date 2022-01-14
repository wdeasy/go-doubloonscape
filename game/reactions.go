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
            game.ReactionIncrement()
        case CAPTAIN_REACTION:
            game.ReactionCaptain(m.UserID, m.GuildID)
        case PRESTIGE_REACTION:
            game.ReactionPrestige(m.UserID)
        case TREASURE_REACTION:           
            game.ReactionTreasure(m.UserID)
        case PICKPOCKET_REACTION:             
            game.ReactionPickPocket(m.UserID)
    }
    
    if m.Emoji.Name == TREASURE_REACTION || m.Emoji.Name == PICKPOCKET_REACTION {
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

    game.checkEventReactions(message)
}

//see if event reactions can be added
func (game *Game) checkEventReactions(message *discordgo.Message) { 
    if treasureChance() && !game.isReactionInReactions(TREASURE_REACTION, message.Reactions) {
        game.treasure.Up = true
        game.addReaction(message.ID, TREASURE_REACTION)		
    }

    for _, e := range game.events {
        eventReaction := game.getReaction(e.Name)

        if game.isReactionInReactions(eventReaction, message.Reactions) {
            continue
        }

        if e.Ready(game.getCooldown(e.Name)) {
            e.Up = true
            game.addCurrentEvent(e.Name)
            game.addReaction(message.ID, eventReaction)
        }        
    }

    for _, e := range game.currentEvents {
        if !game.events[e].Up {
            game.removeReaction(message.ID, game.getReaction(e), game.currentBotID)
            game.removeCurrentEvent(e)
        }
    }
}

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

//function for when the increment reaction is used
func (game *Game) ReactionIncrement() {
    game.incrementCaptain()
}

//function for when the captain reaction is used
func (game *Game) ReactionCaptain(UserID string, GuildID string) {
    if UserID == game.currentCaptainID {
        return
    }

    game.changeCaptainsInGameAndServer(GuildID, UserID)
}

//function for when the prestige reaction is used
func (game *Game) ReactionPrestige(UserID string) {
    if UserID != game.currentCaptainID {
        return
    }

    game.captains[UserID].AddPrestige(PRESTIGE_CONVERSION)
    game.setMessage()		
}

//function for when the treasure reaction is used
func (game *Game) ReactionTreasure(UserID string) {
    if !game.treasure.Up {
        return 
    } 

    game.treasure.Up = false
    game.giveTreasure(UserID)
}

//function for when the pickpocket reaction is used
func (game *Game) ReactionPickPocket(UserID string) {
    if !game.events[PICKPOCKET_NAME].Up {
        return
    }

    game.resetEvent(PICKPOCKET_NAME)

    if UserID == game.currentCaptainID {
        return
    }
   
    game.executePickPocket(UserID) 
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

//get the appropriate reaction from constants
func (game *Game) getReaction(name string) (string) {
    switch name {
    case PICKPOCKET_NAME:
        return PICKPOCKET_REACTION
    default:
        return DEFAULT_REACTION
    }
}