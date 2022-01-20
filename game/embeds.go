package game

import (
    "fmt"
    "strings"

    "github.com/bwmarrin/discordgo"
)

//generate the embed for the bot's message
func (game *Game) generateEmbed() (discordgo.MessageEmbed) {
    var Name string      = ""
    var Gold int64       = 0
    var Prestige float64 = 1

    if game.currentCaptainID != "" {
        Name     = game.captains[game.currentCaptainID].Name
        Gold     = game.captains[game.currentCaptainID].Gold
        Prestige = game.captains[game.currentCaptainID].Prestige
    }

    embed := discordgo.MessageEmbed{
        Color: EMBED_COLOR,
        //Title: "𝔏𝔢𝔞𝔡𝔢𝔯𝔅𝔬𝔞𝔯𝔡",
        Description: game.generateDescription(),
        Fields: []*discordgo.MessageEmbedField{
            {
                Name:   "ℭ𝔞𝔭𝔱𝔞𝔦𝔫",
                Value:  "` " + firstN(Name, 31) + " `",
                Inline: true,
            },
            {
                Name:   "𝔇𝔬𝔲𝔟𝔩𝔬𝔬𝔫𝔰",
                Value:  "` " + fmt.Sprintf("%-7d",Gold) + " `",
                Inline: true,
            },
            {
                Name:   "𝔓𝔯𝔢𝔰𝔱𝔦𝔤𝔢",
                Value:  "` " + fmt.Sprintf("%-4.3f",Prestige) + " `",
                Inline: true,
            },
        },
        // Footer: &discordgo.MessageEmbedFooter{
        //     Text:   game.stats.Event,
        // },		
    }

    return embed
}

//generate embed description
func (game *Game) generateDescription() (string) {
    //line :=  "~~ᵤₔᵤₑᵤₔᵤₑᵤₔᵤₑᵤₔᵤₑᵤₔᵤₑᵤₔᵤₑᵤₔᵤₑᵤₔᵤₑᵤₔᵤₑᵤₔᵤₑᵤₔᵤₑᵤₔᵤₑᵤₔᵤ~~\n"
    return *game.stats.Leaderboard + *game.stats.Destinations
}

//truncate names
func firstN(s string, n int) string {
    if len(s) > n {
         return strings.TrimSpace(s[:n])
    }
    return s
}