package game

import (
    "fmt"

    "github.com/bwmarrin/discordgo"
)

//generate the embed for the bot's message
func (game *Game) generateEmbed() (discordgo.MessageEmbed) {
    embed := discordgo.MessageEmbed{
        Color: 0xf1c40f,
        //Title: "𝔏𝔢𝔞𝔡𝔢𝔯𝔅𝔬𝔞𝔯𝔡",
        Description: game.generateDescription(),
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
        // Footer: &discordgo.MessageEmbedFooter{
        //     Text:   game.stats.Event,
        // },		
    }

    return embed
}

//generate embed description
func (game *Game) generateDescription() (string) {
    return *game.stats.Leaderboard + *game.stats.Log + *game.stats.Destinations
}

//truncate names
func firstN(s string, n int) string {
    if len(s) > n {
         return s[:n]
    }
    return s
}