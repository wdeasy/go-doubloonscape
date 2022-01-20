package game

import (
    "fmt"
    "os"

    "github.com/bwmarrin/discordgo"
)

var (
    Token string = os.Getenv("BOT_TOKEN")
    Role string = os.Getenv("ROLE")
    Channel string = os.Getenv("CHANNEL")
)

//initialize the discord bot
func (game *Game) InitBot() (*discordgo.Session, error) {
    if Token == "" {
        return nil, fmt.Errorf("BOT_TOKEN environment variable is not set")
    }

    if Channel == "" {
        return nil, fmt.Errorf("CHANNEL environment variable is not set")
    }   

    dg, err := discordgo.New("Bot " + Token)
    if err != nil {
        return nil, fmt.Errorf("error creating Discord session: %w", err)
    }

    dg.AddHandler(game.messageCreate)
    dg.AddHandler(game.messageReactionAdd)

    dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuildMessageReactions

    err = dg.Open()
    if err != nil {
        return nil, fmt.Errorf("error opening connection: %w", err)
    }

    return dg, nil
}

//stop the bot
func (game *Game) CloseBot() {
    err := game.dg.Close()
    if err != nil {
        printLog(fmt.Sprintf("error while closing discord bot: %s\n", err))
    }
}

//give the captain role to discord user and remove for all others
func (game *Game) changeRoles(GuildID string, UserID string) (error){
    if Role == "" {
        return fmt.Errorf("ROLE environment variable is not set")
    }

    err := game.dg.GuildMemberRoleAdd(GuildID, UserID, Role)

    if err != nil {
        return fmt.Errorf("could not add captain role to user %s: %w", UserID, err)
    }	

    members, err := game.dg.GuildMembers(GuildID, "", MAX_GUILD_MEMBERS)
    if err != nil {
        return fmt.Errorf("could not get guild members for %s: %w", GuildID, err)
    }	

    for _, member := range members {
        for _, role := range member.Roles {
            if (role == Role && member.User.ID != UserID) {
                err := game.dg.GuildMemberRoleRemove(GuildID, member.User.ID, Role)

                if err != nil {
                    return fmt.Errorf("could not remove captain role from %s: %w", member.User.Username, err)
                }						
            }
        }
    }

    return nil
}