package main

import (
	"database/sql"
	"fmt"
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

    // In this example, we only care about receiving message events.
    dg.Identify.Intents = discordgo.IntentsGuildMessages

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
		err := s.GuildMemberRoleAdd(m.GuildID, m.Author.ID, Role)

		if err != nil {
			fmt.Println(err)
			return
		}

		name := formatName(m.Member.Nick, m.Author.Username)
		fmt.Printf("%s is the captain now.\n", name)		

		members, err := s.GuildMembers(m.GuildID, "", 1000)
		if err != nil {
			fmt.Println(err)
			return
		}	

		for _, member := range members {
			for _, role := range member.Roles {
				if (role == Role && member.User.ID != m.Author.ID) {
					fmt.Printf("Removing the Captain role from %s.\n", member.User.Username)

					err := s.GuildMemberRoleRemove(m.GuildID, member.User.ID, Role)

					if err != nil {
						fmt.Println(err)
					}						
				}
			}
		}

		env.changeCaptains(m.Author.ID, name)
		env.setMessage()
	}
}

func CheckError(err error) {
    if err != nil {
        panic(err)
    }
}

type Captain struct {
	ID string
	Name string
	Gold int
	Captain bool
}

func (env *Env) changeCaptains(ID string, Name string) {
	env.getCaptain(ID, Name)
	env.removeCaptains()
	env.addCaptain(ID, Name)
}

func (env *Env) getCaptain(ID string, Name string){
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
		err := rows.Scan(&captain.ID, &captain.Name, &captain.Gold, &captain.Captain)

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
	updateStmt := `UPDATE captains SET name = $1, captain = true WHERE id = $2`
	_, e := env.db.Exec(updateStmt, Name, ID)
	CheckError(e)
}

func (env *Env) incrementCaptain() {
	updateStmt := `UPDATE captains SET gold = gold + 1 WHERE captain = true`
	_, e := env.db.Exec(updateStmt)
	CheckError(e)
}	

func (env *Env) GameTimer() {
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
	leaderboard, captain := env.getStats()

	embed := discordgo.MessageEmbed{
		Color: 0xf1c40f,
		Title: "𝔏𝔢𝔞𝔡𝔢𝔯𝔅𝔬𝔞𝔯𝔡",
		Description: leaderboard,
		Fields: []*discordgo.MessageEmbedField{
            {
                Name:   "ℭ𝔞𝔭𝔱𝔞𝔦𝔫",
                Value:  "` " + firstN(captain.Name, 31) + " `",
                Inline: true,
            },
            {
                Name:   "𝔇𝔬𝔲𝔟𝔩𝔬𝔬𝔫𝔰",
                Value:  "` " + fmt.Sprintf("%-7d",captain.Gold) + " `",
                Inline: true,
            },
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
	_, err := env.dg.ChannelMessageSendEmbed(Channel, embed)	

	if err != nil {
		fmt.Println(err)
	}
}

func (env *Env) getStats() (string, Captain){
	captains := env.getCaptains()
	leaderboard := env.printLeaderboard(captains)

	var captain Captain
	for _, s := range captains {
		if s.Captain {
			captain = s
		}
	}
	
	return leaderboard, captain
}

func formatName(nick string, user string) (string) {
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