package main

import (
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"regexp"
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
		} else {
			fmt.Printf("%s is the captain now.\n", m.Author.Username)
		}
		
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

		env.changeCaptains(m.Author.ID, m.Author.Username)
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
	env.addCaptain(ID)
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

func (env *Env) removeCaptains(){
	updateStmt := `UPDATE captains SET captain = $1`
	_, e := env.db.Exec(updateStmt, false)
	CheckError(e)
}

func (env *Env) addCaptain(ID string) {
	updateStmt := `UPDATE captains SET captain = $1 WHERE id = $2`
	_, e := env.db.Exec(updateStmt, true, ID)
	CheckError(e)
}

func (env *Env) incrementCaptain() {
	updateStmt := `UPDATE captains SET gold = gold + 1 WHERE captain = true`
	_, e := env.db.Exec(updateStmt)
	CheckError(e)
}

func (env *Env) setTopic() {
	rows, err := env.db.Query(`SELECT * FROM captains WHERE captain = true LIMIT 1`)
	CheckError(err)

	for rows.Next(){
		var captain Captain 

		err := rows.Scan(&captain.ID, &captain.Name, &captain.Gold, &captain.Captain)
		if err != nil {
			fmt.Println(err)
			return
		}	

		topic := fmt.Sprintf("Captain: %s. Gold: %d.", captain.Name, captain.Gold)
		
		_, err = env.dg.ChannelEditComplex(Channel, &discordgo.ChannelEdit{
			Topic: topic,
		})		

		if err != nil {
			fmt.Println(err)
			return
		}
		
		fmt.Println("Set Topic.")
	}
}

func (env *Env) GameTimer() {
	i := 1
	ticker := time.NewTicker(60 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
		   select {
			case <- ticker.C:
				//1 minute timer
				env.incrementCaptain()
				i++
				
				//5 minute timer
				if (i % 5 == 0) {
					env.setTopic()
				}
			case <- quit:
				ticker.Stop()
				return
			}
		}
	}()
}

