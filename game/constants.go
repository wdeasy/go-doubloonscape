package game

const (

    ATLANTIS_CHANCE       = 100
    ATLANTIS_DURATION     = 2
    ATLANTIS_LOW          = 2
    ATLANTIS_MAX          = 100
    ATLANTIS_MOD_MAX      = 5
    ATLANTIS_NAME         = "atlantis"
    
    BERMUDA_CHANCE        = 100
    BERMUDA_DURATION      = 2
    BERMUDA_LOW           = (BERMUDA_MOD_MAX * -1)
    BERMUDA_MAX           = 100
    BERMUDA_MOD_MAX       = 50    
    BERMUDA_NAME          = "bermuda"

    CAPTAIN_CHANCE        = 100
    CAPTAIN_COOLDOWN      = 1
    CAPTAIN_MAX           = 100
    CAPTAIN_NAME          = "captain"
    CAPTAIN_REACTION      = "🏴‍☠️"
    CAPTAIN_REGEX         = `\b[Ii][’']?[Mm][ \t]+[Tt][Hh][Ee][ \t]+[Cc][Aa][Pp][Tt][Aa][Ii][Nn][ \t]+[Nn][Oo][Ww].?\b`

    DEFAULT_CHANCE        = 1	
    DEFAULT_COOLDOWN      = 60     
    DEFAULT_GOLD          = 0    
    DEFAULT_LOW           = 1
    DEFAULT_MAX           = 100   
    DEFAULT_MOD_MAX       = 1
    DEFAULT_NAME          = "default"           
    DEFAULT_PRESTIGE      = 1
    DEFAULT_REACTION      = "⚠️"    

    EMBED_COLOR           = 0xb58900 //gold // 0xf1c40f //yellow 0x2f3136 //grey

    IDLE_EVENT            = "idle"

    INCREMENT_CHANCE      = 100
    INCREMENT_COOLDOWN    = 1      
    INCREMENT_NAME        = "increment"
    INCREMENT_REACTION    = "🪙"
    INCREMENT_MAX         = 100    

    LEADERBOARD_RESET     = 300

    LOG_LINE_LENGTH       = 44    

    MAX_GUILD_MEMBERS     = 1000

    MAX_LOG_LINES         = 10

    MESSAGE_MAX           = 10

    PICKPOCKET_CHANCE     = 1
    PICKPOCKET_COOLDOWN   = 60    
    PICKPOCKET_MAX        = 100    
    PICKPOCKET_NAME       = "pickpocket"
    PICKPOCKET_REACTION   = "💰"

    PRESTIGE_CHANCE       = 100
    PRESTIGE_CONVERSION   = 0.001
    PRESTIGE_COOLDOWN     = 1    
    PRESTIGE_MAX          = 100        
    PRESTIGE_MULTIPLIER   = 1000    
    PRESTIGE_NAME         = "prestige"
    PRESTIGE_REACTION     = "🔱"

    TREASURE_CHANCE       = 1    
    TREASURE_COOLDOWN     = 1    
    TREASURE_MAX          = 1000
    TREASURE_NAME         = "treasure"    
    TREASURE_REACTION     = "👑"    
    TREASURE_START        = 1

)