package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

const (
	ERR_TIME_BELOW_MIN        = 1
	ERR_TIME_ABOVE_MAX        = 2
	ERR_TIME_INVALID_FORMAT   = 3
	ERR_USER_ALREADY_IN_JAIL  = 4
	ERRR_COLOR_INVALID_FORMAT = 5
)

// var (...): Declaring several vars at once (works also with const (yippee) and type (still not sure what tht is outside of struct, Update: also defining other data types duh))
var (
	annoucementsChannel *discordgo.Channel = nil
	arrestRole          *discordgo.Role    = nil
	defaultTime         string             = "1m"
	minTime             string             = "30s"
	maxTime             string             = "5m"
	embedColor          Color              = Color{hexValue: "#000000", intValue: 0}

	s     *discordgo.Session
	token string
)

type Color struct {
	hexValue string
	intValue int
}

type Embed struct{}

func (e *Embed) arrest(i *discordgo.InteractionCreate, a Arrest, t string) discordgo.MessageEmbed {
	fields := []*discordgo.MessageEmbedField{
		{
			Name:  "What happened?",
			Value: fmt.Sprintf("**<@!%s> was arrested by <@!%s>! 😱**", a.user.ID, i.Interaction.Member.User.ID),
		},
		{
			Name:  "Reason",
			Value: a.reason,
		},
	}
	if t != "0s" {
		timeField := discordgo.MessageEmbedField{
			Name:  "You will see them again in ",
			Value: t,
		}
		fields = append(fields, &timeField)
	}
	return discordgo.MessageEmbed{
		Title:  "Sheriff report",
		Fields: fields,
		Color:  embedColor.intValue,
	}
}
func (e *Embed) release(a Arrest) discordgo.MessageEmbed {
	return discordgo.MessageEmbed{
		Title:       "Sheriff report",
		Description: fmt.Sprintf("**<@!%s> was released from jail!**", a.user.ID),
		Color:       embedColor.intValue,
	}
}

type Arrest struct {
	user               *discordgo.User
	success            bool
	timeString, reason string
	embed              Embed
}

// Acts as a method to Arrest
func (a Arrest) ValidateTime() (int, error) {
	t, err := time.ParseDuration(a.timeString)
	if err != nil {
		return ERR_TIME_INVALID_FORMAT, err
	}

	min, _ := time.ParseDuration(minTime)
	max, _ := time.ParseDuration(maxTime)

	if t.Seconds() < min.Seconds() {
		return ERR_TIME_BELOW_MIN, errors.New("time below minimum")
	}
	if t.Seconds() > max.Seconds() {
		return ERR_TIME_ABOVE_MAX, errors.New("time above maximum")
	}

	return 0, nil
}

func (a *Arrest) makeArrest(s *discordgo.Session, i *discordgo.InteractionCreate) {

	s.GuildMemberRoleAdd(i.GuildID, a.user.ID, arrestRole.ID)

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintln("Succesful arrest"),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		panic(err)
	}

	// embed := a.embed(i, EMBED_NEW_ARREST)
	embed := a.embed.arrest(i, *a, a.timeString)
	msg, err := s.ChannelMessageSendEmbed(annoucementsChannel.ID, &embed)
	if err != nil {
		panic(err)
	}

	// embed := a.embed(i, EMBED_NEW_ARREST)
	// arrestEmbed := []*discordgo.MessageEmbed{&embed}
	a.success = true
	log.Println("Role Added succesfully")

	t, _ := time.ParseDuration(a.timeString)
	seconds := int(t.Seconds())
	// minutes := int(t.Minutes())
	dur := t

	for range time.Tick(time.Second) {
		// Update time in channel?
		if seconds == 0 {
			break
		}
		seconds--
		dur -= time.Second

		embed := a.embed.arrest(i, *a, dur.String())
		_, err = s.ChannelMessageEditEmbed(annoucementsChannel.ID, msg.ID, &embed)
		if err != nil {
			panic(err)
		}

		log.Printf("%v", dur)
	}
	defer a.breakFree(s, i)
}

func (a *Arrest) breakFree(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.GuildMemberRoleRemove(i.GuildID, a.user.ID, arrestRole.ID)
	delete(arrests, a.user.ID)

	// embed := a.embed(i, EMBED_RELEASE)
	embed := a.embed.release(*a)
	_, err := s.ChannelMessageSendEmbed(annoucementsChannel.ID, &embed)
	if err != nil {
		panic(err)
	}

	log.Println("Role Removed succesfully")
}

var arrests = make(map[string]*Arrest) // key == ArrestedUserID

// init() functions should run after all global variable declarations are initialized and packages are initialized
// and befoe main function

// Gets flag's values
func init() { flag.Parse() }

// "Load will read your env file(s) and load them into ENV for this process."
func init() { // Idk if I'm using this correctly
	godotenv.Load()
	token = os.Getenv("BOT_TOKEN")
}

func init() {
	var err error
	s, err = discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
}

var (
	// defaultMemberPermissions int64 = discordgo.PermissionAdministrator

	commands = []*discordgo.ApplicationCommand{ // Array of pointers to ApplicationCommands
		{
			// Type: ,
			Name:        "arrest",
			Description: "Put someone to jail for some time",
			// DefaultMemberPermissions: &defaultMemberPermissions,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user",
					Description: "Who are we arresting bestie?",
					Required:    true,
				},
			},
		},
		{
			Name:        "arrest-config-set",
			Description: "Arrest settings",
			// DefaultMemberPermissions: &defaultMemberPermissions,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionChannel,
					Name:        "annoucement-channel",
					Description: "A channel for arrest/release announcements.",
					ChannelTypes: []discordgo.ChannelType{
						discordgo.ChannelTypeGuildText,
					},
					Required: false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionRole,
					Name:        "arrest-role",
					Description: "A role that will be added to the arrested person.",
					Required:    false,
				},
				{
					Type: discordgo.ApplicationCommandOptionString,
					Name: "min-time",

					Description: "Minimum time for imprisonment (Format [min]m[sec]s).",
					Required:    false,
				},
				{
					Type: discordgo.ApplicationCommandOptionString,
					Name: "max-time",

					Description: "Maximum time for imprisonment (Format [min]m[sec]s).",
					Required:    false,
				},
				{
					Type: discordgo.ApplicationCommandOptionString,
					Name: "default-time",

					Description: "Default time for imprisonment (Format [min]m[sec]s).",
					Required:    false,
				},
				{
					Type: discordgo.ApplicationCommandOptionString,
					Name: "embed-color",

					Description: "HEX value for embed's color stripe (defaults to black).",
					Required:    false,
				},
			},
		},
		{
			Name:        "arrest-config-get",
			Description: "Retrieves current config for /arrest command.",
			// DefaultMemberPermissions: &defaultMemberPermissions,
		},
		{
			Name:        "arrest-unset-channel",
			Description: "Unsets channel for annoucements.",
			// DefaultMemberPermissions: &defaultMemberPermissions,
		},
	}

	componentsHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"unset-channel-yes": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			annoucementsChannel = nil
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseUpdateMessage,
				Data: &discordgo.InteractionResponseData{
					Content: "As you wish!",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			if err != nil {
				panic(err)
			}
		},
		"unset-channel-no": func(s *discordgo.Session, i *discordgo.InteractionCreate) {

			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseUpdateMessage,
				Data: &discordgo.InteractionResponseData{
					Content: "Okie Dokie, the channel stays.",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			if err != nil {
				panic(err)
			}
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, selectedUserID *string){
		"arrest": func(s *discordgo.Session, i *discordgo.InteractionCreate, selectedUserID *string) {
			selectedUser := i.ApplicationCommandData().Options[0].UserValue(s)
			*selectedUserID = selectedUser.ID

			if arrestRole == nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Flags:   discordgo.MessageFlagsEphemeral,
						Content: fmt.Sprintln("Please set an arrest role. You can do it by using **/arrest-config-set**") + fmt.Sprintln("If you don't have access to this command, you can kindly ask your Admin people to do it :)"),
					},
				})
				return
			}

			if annoucementsChannel == nil {
				ch, err := s.Channel(i.ChannelID)
				if err != nil {
					panic(err)
				}
				annoucementsChannel = ch
			}

			if errCode, err := createNewArrest(s, i, selectedUser); errCode != 0 {
				log.Printf("%v", err)
				sendErrorResponse(s, i, selectedUser, errCode)
				return
			}

			reasonValue := ""
			log.Println(arrests[*selectedUserID].reason)
			if (!arrests[*selectedUserID].success) && (arrests[*selectedUserID].reason != "") {
				reasonValue = arrests[*selectedUserID].reason
			}

			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseModal,
				Data: &discordgo.InteractionResponseData{
					CustomID: "arrest_" + i.Interaction.Member.User.ID,
					Title:    "New Arrest",
					Components: []discordgo.MessageComponent{
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.TextInput{
									CustomID:  "reason",
									Label:     "Why is this person being arrested?",
									Style:     discordgo.TextInputParagraph,
									Required:  true,
									MaxLength: 2000,
									Value:     reasonValue,
								},
							},
						},
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.TextInput{
									CustomID:    "time",
									Label:       "For how long? (default 1min)",
									Style:       discordgo.TextInputShort,
									Placeholder: "[min]m[sec]s (e.g. 1m30s / 3m / 30s)",
									Required:    false,
									MinLength:   2,
									MaxLength:   6,
								},
							},
						},
					},
				},
			})
			if err != nil {
				log.Println("Something went wrong during the arrest")
				panic(err)
			}
		},
		"arrest-config-set": func(s *discordgo.Session, i *discordgo.InteractionCreate, selectedUserID *string) {
			options := i.ApplicationCommandData().Options

			// Or convert the slice into a map
			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}

			margs := make([]interface{}, 0, len(options))
			msgformat := "Values set:\n"

			if opt, ok := optionMap["annoucement-channel"]; ok {
				annoucementsChannel = opt.ChannelValue(nil)
				margs = append(margs, opt.ChannelValue(nil).ID)
				msgformat += "> annoucement-channel: <#%s>\n"
			}
			if opt, ok := optionMap["arrest-role"]; ok {
				arrestRole = opt.RoleValue(nil, i.GuildID)
				margs = append(margs, opt.RoleValue(nil, i.GuildID).ID)
				msgformat += "> arrest-role: <@&%s>\n"
			}
			if opt, ok := optionMap["min-time"]; ok {
				minTime = opt.StringValue()
				margs = append(margs, opt.StringValue())
				msgformat += "> min-time: %s\n"
			}
			if opt, ok := optionMap["max-time"]; ok {
				maxTime = opt.StringValue()
				margs = append(margs, opt.StringValue())
				msgformat += "> max-time: %s\n"
			}
			if opt, ok := optionMap["default-time"]; ok {
				defaultTime = opt.StringValue()
				margs = append(margs, opt.StringValue())
				msgformat += "> default-time: %s\n"
			}
			if opt, ok := optionMap["embed-color"]; ok {
				c, err := ParseHexColor(opt.StringValue())
				if err != nil {
					fmt.Println(err)
					sendErrorResponse(s, i, nil, ERRR_COLOR_INVALID_FORMAT)

					margs = append(margs, embedColor.hexValue)
					msgformat += "> embed-color: `%s`\n"
					s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
						Flags: discordgo.MessageFlagsEphemeral,
						Content: fmt.Sprintf(
							msgformat,
							margs...,
						),
					})
					return
				} else {
					embedColor.hexValue = opt.StringValue()
					embedColor.intValue = c
					margs = append(margs, embedColor.hexValue)
					msgformat += "> embed-color: `%s`\n"
				}
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				// Ignore type for now, they will be discussed in "responses"
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Flags: discordgo.MessageFlagsEphemeral,
					Content: fmt.Sprintf(
						msgformat,
						margs...,
					),
				},
			})
		},
		"arrest-config-get": func(s *discordgo.Session, i *discordgo.InteractionCreate, selectedUserID *string) {
			margs := make([]interface{}, 0)
			msgformat := "Values set:\n"
			if annoucementsChannel == nil {
				msgformat += "> annoucement-channel: unset\n"
			} else {
				margs = append(margs, annoucementsChannel.ID)
				msgformat += "> annoucement-channel: <#%s>\n"
			}
			if arrestRole == nil {
				msgformat += "> arrest-role: unset\n"
			} else {
				margs = append(margs, arrestRole.ID)
				msgformat += "> arrest-role: <@&%s>\n"
			}
			margs = append(margs, minTime)
			msgformat += "> min-time: %s\n"
			margs = append(margs, maxTime)
			msgformat += "> max-time: %s\n"
			margs = append(margs, defaultTime)
			msgformat += "> default-time: %s\n"
			margs = append(margs, embedColor.hexValue)
			msgformat += "> embed-color: `%s`\n"

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Flags: discordgo.MessageFlagsEphemeral,
					Content: fmt.Sprintf(
						msgformat,
						margs...,
					),
				},
			})
		},
		"arrest-unset-channel": func(s *discordgo.Session, i *discordgo.InteractionCreate, selectedUserID *string) {
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Are you sure you want to unset the announcement channel? If not set, announcements will show up in the channel where /arrest is called.",
					Flags:   discordgo.MessageFlagsEphemeral,
					Components: []discordgo.MessageComponent{
						// ActionRow is a container of all buttons within the same row.
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.Button{
									// Label is what the user will see on the button.
									Label: "Yes",
									// Style provides coloring of the button. There are not so many styles tho.
									Style: discordgo.SuccessButton,
									// Disabled allows bot to disable some buttons for users.
									Disabled: false,
									// CustomID is a thing telling Discord which data to send when this button will be pressed.
									CustomID: "unset-channel-yes",
								},
								discordgo.Button{
									Label:    "No",
									Style:    discordgo.DangerButton,
									Disabled: false,
									CustomID: "unset-channel-no",
								},
							},
						},
					},
				},
			})
			if err != nil {
				panic(err)
			}
		},
	}
)

func ready(s *discordgo.Session, r *discordgo.Ready) {
	log.Printf("ready")
	log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)

	mem, err := s.GuildMembers(r.Guilds[0].ID, "", 1000)
	if err != nil {
		log.Println(err)
	}
	usernames := []string{}

	for _, m := range mem {
		usernames = append(usernames, m.User.Username)
	}

	log.Println(usernames)

	err = s.RequestGuildMembers(r.Guilds[0].ID, "", 0, "", true) // Doesn't work

	// err := s.RequestGuildMembers(r.Application.GuildID, "", 0, "", false)

	// mem, err := s.GuildMembers(i.GuildID, "", 0)
	if err != nil {
		log.Println(err)
	}
}

func main() {
	s.AddHandler(ready)

	var currArrestUserID string

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand: // Registers slash commands
			if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
				err := s.RequestGuildMembers(i.GuildID, "", 0, "", false)
				// mem, err := s.GuildMembers(i.GuildID, "", 0)
				if err != nil {
					log.Println(err)
				}
				// log.Println(mem)
				defer h(s, i, &currArrestUserID)
			}
		case discordgo.InteractionMessageComponent:
			if h, ok := componentsHandlers[i.MessageComponentData().CustomID]; ok {
				defer h(s, i)
			}
		case discordgo.InteractionModalSubmit:
			data := i.ModalSubmitData()

			if !strings.HasPrefix(data.CustomID, "arrest") {
				return
			}

			a, ok := arrests[currArrestUserID]
			if !ok {
				panic("Something went wrong, user not added")
			}

			a.reason = data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
			timeSubmitted := data.Components[1].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
			if timeSubmitted != "" {
				a.timeString = timeSubmitted
			}

			if errCode, err := a.ValidateTime(); errCode != 0 {
				log.Printf("%v", err)
				a.success = false
				sendErrorResponse(s, i, a.user, errCode)
				return
			}

			a.makeArrest(s, i)

		}
	})

	s.Identify.Intents = discordgo.IntentsAllWithoutPrivileged | discordgo.IntentGuildMembers | discordgo.IntentGuildPresences
	err := s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	log.Println("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, "", v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	s.StateEnabled = true
	log.Println(s.LastHeartbeatSent)

	defer s.Close() // "A defer statement defers the execution of a function until the surrounding function returns."

	// Prevents bot to stop after all is executed
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	log.Println("Removing commands...")
	// I am not sure what is happening here or if it works
	for _, v := range registeredCommands {
		err := s.ApplicationCommandDelete(s.State.User.ID, "", v.ID)
		if err != nil {
			log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
		}
	}

	log.Println("Gracefully shutting down.")
}

func createNewArrest(s *discordgo.Session, i *discordgo.InteractionCreate, selectedUser *discordgo.User) (int, error) {
	// selectedUser.
	m, err := s.GuildMember(i.GuildID, selectedUser.ID)
	if err != nil {
		panic(err)
	}
	if slices.Contains(m.Roles, arrestRole.ID) {
		return ERR_USER_ALREADY_IN_JAIL, errors.New("user already arrested")
	}
	if _, ok := arrests[selectedUser.ID]; ok {
		// User already exists but there was a wrong value passed
		return 0, nil
	}
	// UserID doesn't exist -> create a new one
	arrests[selectedUser.ID] = &Arrest{user: selectedUser, reason: "", timeString: defaultTime, success: false, embed: Embed{}}
	return 0, nil
}

func sendErrorResponse(s *discordgo.Session, i *discordgo.InteractionCreate, selectedUser *discordgo.User, errCode int) {
	responseMsg := ""

	switch errCode {
	case ERR_TIME_INVALID_FORMAT:
		responseMsg = fmt.Sprintln("You entered invalid time.\nThe format should look like this: [minutes]m[seconds]s.\nFor example these are valid time formats: 5m, 30s, 1m20s, 0m50s, 2m0s..")
	case ERR_TIME_BELOW_MIN:
		responseMsg = fmt.Sprintf("Okay that's too short, you gotta be a bit more strict.\nMinimum time is %v.\n", minTime)
	case ERR_TIME_ABOVE_MAX:
		responseMsg = fmt.Sprintf("Way too long!\nMaximum time is %v.\n", maxTime)
	case ERR_USER_ALREADY_IN_JAIL:
		responseMsg = fmt.Sprintf("User <@!%s> is already in jail! Let's calm down.\n", selectedUser.Username)
	case ERRR_COLOR_INVALID_FORMAT:
		responseMsg = fmt.Sprintf("Not a valid color format. Defaulting to `%s`.\n", embedColor.hexValue)
	default:
		fmt.Sprintln("In full honesty I am not sure what happened, sorry :(")
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: responseMsg,
		},
	})
}

func ParseHexColor(s string) (int, error) {
	allowedChars := []string{"A", "B", "C", "D", "E", "F", "1", "2", "3", "4", "5", "6", "7", "8", "9", "0"}

	if s[0] != '#' {
		return 0, errors.New("missing # in color")
	}
	s = strings.Replace(s, "#", "", -1)
	if len(s) != 3 && len(s) != 6 {
		return 0, errors.New("missing # in color")
	}

	for _, char := range s {
		if !slices.Contains(allowedChars, strings.ToUpper(string(char))) {
			return 0, errors.New("invalid color input")
		}
	}

	if len(s) == 3 {
		s += s // FFF -> FFFFFF
	}
	decimal_num, err := strconv.ParseInt(s, 16, 64)
	if err != nil {
		return 0, err
	}

	return int(decimal_num), nil
}
