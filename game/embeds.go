package game

import (
    "fmt"
    "sort"
    "strings"
    "time"

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
        Footer: &discordgo.MessageEmbedFooter{
            Text:   stats.Event,
        },		
    }

    return embed
}

//generate embed description
func (game *Game) generateDescription() (string) {
    return stats.Leaderboard + game.destinationsString()
}


//leaderboard sorting
type Pair struct {
    Key string
    Value int
  }

//leaderboard sorting  
type PairList []Pair

//leaderboard sorting
func (p PairList) Len() int { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int){ p[i], p[j] = p[j], p[i] }

//format leaderboard to string
func (game *Game) printLeaderboard() (string) {

    pl := make(PairList, len(game.captains))
    i := 0
    for k, v := range game.captains {
        pl[i] = Pair{k, v.Gold}
        i++
    }

    sort.Sort(sort.Reverse(pl))

    var b strings.Builder
    fmt.Fprintf(&b, "**𝔏𝔢𝔞𝔡𝔢𝔯𝔅𝔬𝔞𝔯𝔡**\n") 
    for j, k := range pl {
        fmt.Fprintf(&b, "` %2d ` ` %-27s ` ` %7d `\n", j+1, firstN(game.captains[k.Key].Name,27), k.Value)
    }
    return b.String()
}

//truncate names
func firstN(s string, n int) string {
    if len(s) > n {
         return s[:n]
    }
    return s
}

//generate destinations string for embed
func (game *Game) destinationsString() (string) {
    var b strings.Builder

    if time.Now().Before(game.destinations["atlantis"].End) {
        fmt.Fprintf(&b, "%s%s%d%s\n", "` 𝔄𝔱𝔩𝔞𝔫𝔱𝔦𝔰  ` ", "` 𝔇𝔬𝔲𝔟𝔩𝔬𝔬𝔫𝔰 𝔪𝔲𝔩𝔱𝔦𝔭𝔩𝔦𝔢𝔡 𝔟𝔶 ", game.destinations["atlantis"].Amount, " `")
    }

    if time.Now().Before(game.destinations["bermuda"].End) {
        fmt.Fprintf(&b, "%s%s%d%s\n", "` 𝔅𝔢𝔯𝔪𝔲𝔡𝔞 ` ", "` 𝔗𝔦𝔪𝔢 𝔞𝔩𝔱𝔢𝔯𝔢𝔡 𝔟𝔶 ", game.destinations["bermuda"].Amount, " 𝔭𝔢𝔯𝔠𝔢𝔫𝔱 `")
    }	

    if b.Len() != 0 {
        return "**𝔇𝔢𝔰𝔱𝔦𝔫𝔞𝔱𝔦𝔬𝔫𝔰**\n" + b.String()
    }

    return b.String()
}


