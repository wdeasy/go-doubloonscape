package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"

	_ "github.com/lib/pq"
)

// Variables used for command line parameters
var (
    Token string = os.Getenv("BOT_TOKEN")
    Role string = os.Getenv("ROLE")
    Channel string = os.Getenv("CHANNEL")
    DatabaseURL string = os.Getenv("DATABASE_URL")

    stats Stats
    treasure bool
)

type Env struct {
    db *sql.DB
    dg *discordgo.Session
}

func main() {

    db, err := sql.Open("postgres", DatabaseURL)
    if err != nil {
        fmt.Println("could not open sql: %w", err)
         return
    }

    if err = db.Ping(); err != nil {
        fmt.Println("could not ping DB: %w", err)
         return
    }

    // Create a new Discord session using the provided bot token.
    dg, err := discordgo.New("Bot " + Token)
    if err != nil {
        fmt.Println("error creating Discord session, %w", err)
        return
    }

    // Create an instance of Env containing the connection pool.
    env := &Env{db: db, dg: dg}	

    // Register the messageCreate func as a callback for MessageCreate events.
    dg.AddHandler(env.messageCreate)
    dg.AddHandler(env.messageReactionAdd)

    // In this example, we only care about receiving message events.
    dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuildMessageReactions

    // Open a websocket connection to Discord and begin listening.
    err = dg.Open()
    if err != nil {
        fmt.Println("error opening connection, %w", err)
        return
    }

    // Timer
    env.GameTimer()

    // Wait here until CTRL-C or other term signal is received.
    fmt.Println("Bot is now running. Press CTRL-C to exit.")
    sc := make(chan os.Signal, 1)
    signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
    <-sc

    // Cleanly close down the Discord session.
    dg.Close()
    db.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func (env *Env) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

    // Ignore all messages created by the bot itself
    // This isn't required in this specific example but it's a good practice.
    if m.Author.ID == s.State.User.ID {
        return
    }

    if m.ChannelID != Channel {
        return
    }	

    matched, _ := regexp.MatchString(`\b[Ii][’']?[Mm][ \t]+[Tt][Hh][Ee][ \t]+[Cc][Aa][Pp][Tt][Aa][Ii][Nn][ \t]+[Nn][Oo][Ww].?\b`, m.Content)
    if matched {
        env.newCaptain(m.GuildID, m.Author.ID, m.Member.Nick, m.Author.Username)
    }
}

func (env *Env) newCaptain(GuildID string, UserID string, MemberNick string, UserName string) {
    err := env.dg.GuildMemberRoleAdd(GuildID, UserID, Role)

    if err != nil {
        fmt.Println(err)
        return
    }

    name := getName(MemberNick, UserName)
    fmt.Printf("%s is the captain now.\n", name)		

    members, err := env.dg.GuildMembers(GuildID, "", 1000)
    if err != nil {
        fmt.Println(err)
        return
    }	

    for _, member := range members {
        for _, role := range member.Roles {
            if (role == Role && member.User.ID != UserID) {
                fmt.Printf("Removing the Captain role from %s.\n", member.User.Username)

                err := env.dg.GuildMemberRoleRemove(GuildID, member.User.ID, Role)

                if err != nil {
                    fmt.Println(err)
                }						
            }
        }
    }

    env.changeCaptains(UserID, name)
    env.setMessage()    
}

func (env *Env) messageReactionAdd(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
    if m.UserID == s.State.User.ID {
        return
    }

    if m.ChannelID != Channel {
        return
    }
    
    switch m.Emoji.Name {
        case "🪙":
            env.coinEmoji()
        case "🏴‍☠️":
            env.pirateEmoji(m.UserID, m.GuildID)
        case "🔱":
            env.tridentEmoji(m.UserID, m.GuildID)
        case "👑":
            env.crownEmoji(m.UserID, m.GuildID)
        default:
            return        
    }	
    
    err := env.dg.MessageReactionRemove(Channel, m.MessageID, url.QueryEscape(m.Emoji.Name), m.UserID)
    
    if err != nil {
        fmt.Println(err)
    }    

}

func (env *Env) coinEmoji() {
    env.incrementCaptain()
}

func (env *Env) pirateEmoji(UserID string, GuildID string) {
    if (UserID == stats.Captain.ID) {
        return
    }

    m, err := env.dg.GuildMember(GuildID, UserID)

    if err != nil {
        fmt.Println(err)
        return
    } 

    env.newCaptain(GuildID, UserID, m.User.Username, m.Nick)
}

func (env *Env) tridentEmoji(UserID string, GuildID string) {
    if (UserID != stats.Captain.ID) {
        return
    }

    m, err := env.dg.GuildMember(GuildID, UserID)

    if err != nil {
        fmt.Println(err)
        return
    } 
    
    env.findCaptain(UserID, getName(m.Nick, m.User.Username))
    env.addPrestige(UserID)
    env.setMessage()
}

func (env *Env) crownEmoji(UserID string, GuildID string) {
    if !treasure {
        return
    } else {
        treasure = false
    }

    m, err := env.dg.GuildMember(GuildID, UserID)

    if err != nil {
        fmt.Println(err)
        return
    } 
    
    name := getName(m.Nick, m.User.Username)
    env.findCaptain(UserID, name)
    treasure := env.giveTreasure(UserID)
    stats.Event = fmt.Sprintf("%s looted Treasure worth %d gold!", name, treasure)
    env.setMessage()

    messages, err := env.dg.ChannelMessages(Channel, 100, "", "", "")
    if err != nil {
        fmt.Println(err)
        return
    }
    
    for _, s := range messages {
        if (s.Author.ID == env.dg.State.User.ID) {
            err := env.dg.MessageReactionRemove(Channel, s.ID, url.QueryEscape("👑"), env.dg.State.User.ID)

            if err != nil {
                fmt.Println(err)
                return
            }
        }
    }   
}

func CheckError(err error) {
    if err != nil {
        panic(err)
    }
}

type Stats struct {
    Captain Captain 
    Leaderboard string
    Event string
}

type Captain struct {
    ID string
    Name string
    Gold int
    Captain bool
    Prestige float32
}

func (env *Env) changeCaptains(ID string, Name string) {
    env.findCaptain(ID, Name)
    env.removeCaptains()
    env.addCaptain(ID, Name)
}

func (env *Env) findCaptain(ID string, Name string){
    rows, err := env.db.Query(`SELECT * FROM captains WHERE id = $1 LIMIT 1`, ID)
    CheckError(err)

    i := 0
    for rows.Next(){
        i++
    }

    if i == 0 {
        fmt.Println("New captain joined!")
        insertStmt := `INSERT INTO captains(id, name, gold, captain) VALUES($1, $2, $3, $4)`
        _, err := env.db.Exec(insertStmt, ID, Name, 0, false)
        CheckError(err)
    }
}	

func (env *Env) getCaptains() ([]Captain){
    rows, err := env.db.Query(`SELECT * FROM captains ORDER BY gold DESC`)
    CheckError(err)

    var captains []Captain
    for rows.Next(){
        var captain Captain 
        err := rows.Scan(&captain.ID, &captain.Name, &captain.Gold, &captain.Captain, &captain.Prestige)

        if err != nil {
            fmt.Println(err)
        }
        
        captains = append(captains, captain)
    }

    return captains
}	

func (env *Env) removeCaptains(){
    updateStmt := `UPDATE captains SET captain = false`
    _, e := env.db.Exec(updateStmt)
    CheckError(e)
}

func (env *Env) addCaptain(ID string, Name string) {
    updateStmt := `UPDATE captains SET name = $1, captain = true WHERE id = $2 RETURNING id, name, gold, captain, prestige`
    e := env.db.QueryRow(updateStmt, Name, ID).Scan(&stats.Captain.ID, &stats.Captain.Name, &stats.Captain.Gold, &stats.Captain.Captain, &stats.Captain.Prestige)
    CheckError(e)
}

func (env *Env) incrementCaptain() {
    if stats.Captain.ID == "" {
        return
    }

    updateStmt := `UPDATE captains SET gold = gold + FLOOR(prestige) WHERE captain = true RETURNING id, name, gold, captain, prestige`
    e := env.db.QueryRow(updateStmt).Scan(&stats.Captain.ID, &stats.Captain.Name, &stats.Captain.Gold, &stats.Captain.Captain, &stats.Captain.Prestige)
    CheckError(e)
}	

func (env *Env) addPrestige(UserID string) {
    updateStmt := `UPDATE captains SET prestige = prestige + gold * 0.001, gold = 0 WHERE id = $1 RETURNING id, name, gold, captain, prestige`
    e := env.db.QueryRow(updateStmt, UserID).Scan(&stats.Captain.ID, &stats.Captain.Name, &stats.Captain.Gold, &stats.Captain.Captain, &stats.Captain.Prestige)
    CheckError(e)
}

func (env *Env) giveTreasure(UserID string) (int) {
    treasure := rand.Intn(1000)

    updateStmt := `UPDATE captains SET gold = gold + $1 WHERE id = $2`
    _, e := env.db.Exec(updateStmt, treasure, UserID)
    CheckError(e)

    return treasure
} 

func (env *Env) GameTimer() {
    env.setStats()

    i := 1
    ticker := time.NewTicker(60 * time.Second)
    quit := make(chan struct{})
    go func() {
        for {
            select {
            case <- ticker.C:
                env.incrementCaptain()
                env.setMessage()
                i++

                if (i % 60 == 0) {
                    env.setStats()
                }
            case <- quit:
                ticker.Stop()
                return
            }
        }
    }()
}

func (env *Env) setMessage() { 
    messages, err := env.dg.ChannelMessages(Channel, 100, "", "", "")
    if err != nil {
        fmt.Println(err)
        return
    }
    
    embed := env.generateEmbed()

    if (messages[0].Author.ID == env.dg.State.User.ID) {
        env.editMessage(&embed, messages[0].ID)
    } else {
        env.newMessage(&embed)
        
        for _, s := range messages {
            if (s.Author.ID == env.dg.State.User.ID) {
                env.dg.ChannelMessageDelete(Channel, s.ID)
            }
        }
    }
}

func (env *Env) printLeaderboard(captains []Captain) (string) {
    var b strings.Builder
    for i, s := range captains {
        fmt.Fprintf(&b, "` %2d ` ` %-27s ` ` %7d `\n", i+1, firstN(s.Name,27), s.Gold)
    }

    return b.String()
}

func (env *Env) generateEmbed() (discordgo.MessageEmbed) {
    embed := discordgo.MessageEmbed{
        Color: 0xf1c40f,
        Title: "𝔏𝔢𝔞𝔡𝔢𝔯𝔅𝔬𝔞𝔯𝔡",
        Description: stats.Leaderboard,
        Fields: []*discordgo.MessageEmbedField{
            {
                Name:   "ℭ𝔞𝔭𝔱𝔞𝔦𝔫",
                Value:  "` " + firstN(stats.Captain.Name, 31) + " `",
                Inline: true,
            },
            {
                Name:   "𝔇𝔬𝔲𝔟𝔩𝔬𝔬𝔫𝔰",
                Value:  "` " + fmt.Sprintf("%-7d",stats.Captain.Gold) + " `",
                Inline: true,
            },
            {
                Name:   "𝔓𝔯𝔢𝔰𝔱𝔦𝔤𝔢",
                Value:  "` " + fmt.Sprintf("%-4.3f",stats.Captain.Prestige) + " `",
                Inline: true,
            },
        },
        Footer: &discordgo.MessageEmbedFooter{
            Text:   stats.Event,
        },		
    }

    return embed
}

func (env *Env) editMessage(embed *discordgo.MessageEmbed, messageID string) { 
    _, err := env.dg.ChannelMessageEditEmbed(Channel, messageID, embed)	

    if err != nil {
        fmt.Println(err)
    }
}

func (env *Env) newMessage(embed *discordgo.MessageEmbed) { 
    msg, err := env.dg.ChannelMessageSendEmbed(Channel, embed)	

    if err != nil {
        fmt.Println(err)
    }

    env.addReactions(msg)
}

func (env *Env) setStats() (){
    captains := env.getCaptains()
    stats.Leaderboard = env.printLeaderboard(captains)

    for _, s := range captains {
        if s.Captain {
            stats.Captain = s
        }
    }

    stats.Event = ""
}

func (env *Env) addReactions(message *discordgo.Message) {
    env.dg.MessageReactionAdd(Channel, message.ID, url.QueryEscape("🪙"))
    env.dg.MessageReactionAdd(Channel, message.ID, url.QueryEscape("🔱"))
    env.dg.MessageReactionAdd(Channel, message.ID, url.QueryEscape("🏴‍☠️"))

    //if rand.Intn(100) == 2 {
        treasure = true
        env.dg.MessageReactionAdd(Channel, message.ID, url.QueryEscape("👑"))
    //}
}

func getName(nick string, user string) (string) {
    if (nick != "") {
        return nick
    } else {
        return user
    }
}
func firstN(s string, n int) string {
    if len(s) > n {
         return s[:n]
    }
    return s
}