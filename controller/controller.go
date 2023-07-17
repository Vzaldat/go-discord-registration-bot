package controller

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/Vzaldat/registration-bot/Playermodel"
	"github.com/bwmarrin/discordgo"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

var registrations map[string]Playermodel.Player

var questions = []string{
	"Name of the Player: ",
	"In Game Name (excluding tag):",
	"In Game Tag: ",
	"Rank (Choose from the list): ",
}

var rankOptions = []string{
	"Unranked",
	"Iron1",
	"Iron2",
	"Iron3",
	"Bronze1",
	"Bronze2",
	"Bronze3",
	"Silver1",
	"Silver2",
	"Silver3",
	"Gold1",
	"Gold2",
	"Gold3",
	"Platinum1",
	"Platinum2",
	"Platinum3",
	"Diamond1",
	"Diamond2",
	"Diamond3",
	"Ascendant1",
	"Ascendant2",
	"Ascendant3",
	"Immortal1",
	"Immortal2",
	"Immortal3",
	"Radiant1",
	"Radiant2",
	"Radiant3",
}

var spreadSheetID = "1kx4SKmDrITzVTM1DNStgzgl0YglqVXziMc27Ozx4Uh8"
var sheetName = "Sheet2"

var sheetsService *sheets.Service

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.Contains(m.Content, fmt.Sprintf("<@!%s>", s.State.User.ID)) {
		sendRegistrationQuestions(s, m.ChannelID, m.Author.ID)
	}
}
func ReactionAdd(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	if r.UserID == s.State.User.ID {
		return
	}

	registration, ok := registrations[r.MessageID]
	if !ok {
		return
	}

	if r.Emoji.Name != "" && r.MessageID != "" && r.UserID != "" {
		registration.Rank = r.Emoji.Name

		registrations[r.MessageID] = registration
		err := s.MessageReactionRemove(r.ChannelID, r.MessageID, r.Emoji.Name, r.UserID)
		if err != nil {
			fmt.Println("Error in removing reaction: ", err)
		}
	}
}
func sendRegistrationQuestions(s *discordgo.Session, channelID, userID string) {
	dmChannel, err := s.UserChannelCreate(userID)
	if err != nil {
		fmt.Println("Error creating DM channel:", err)
		return
	}

	registrations[dmChannel.ID] = Playermodel.Player{}

	for i, question := range questions {
		message := &discordgo.MessageSend{
			Content: question,
		}

		msg, err := s.ChannelMessageSendComplex(dmChannel.ID, message)
		if err != nil {
			fmt.Println("Error sending message:", err)
			return
		}

		if i == len(questions)-1 {
			for _, rank := range rankOptions {
				err := s.MessageReactionAdd(dmChannel.ID, msg.ID, rank)
				if err != nil {
					fmt.Println("Error adding reaction:", err)
				}
			}
		}
	}
	err = storeRegistrationInGoogleSheets(dmChannel.ID)
	if err != nil {
		fmt.Println("Error in storing registration in Google Sheets: ", err)
		return
	}
	_, err = s.ChannelMessageSend(dmChannel.ID, "Thank you for registering")
	if err != nil {
		fmt.Println("Error sending message: ", err)
	}
}
func storeRegistrationInGoogleSheets(dmChannelID string) error {
	err := createSheetsService()
	if err != nil {
		fmt.Println("Error in creating sheets: ", err)
	}
	registration, ok := registrations[dmChannelID]
	if !ok {
		return fmt.Errorf("registration not found for dm channel ID: %s", dmChannelID)
	}

	values := [][]interface{}{
		{registration.Name, registration.InGameName, registration.Rank},
	}

	valueRange := &sheets.ValueRange{
		Values: values,
	}

	_, err = sheetsService.Spreadsheets.Values.Append(spreadSheetID, sheetName, valueRange).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		return err
	}

	return nil
}
func createSheetsService() error {
	data, err := os.ReadFile("./creds.json")
	if err != nil {
		return err
	}

	credentials, err := google.CredentialsFromJSON(context.Background(), data, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		return err
	}

	sheetsService, err = sheets.NewService(context.Background(), option.WithTokenSource(credentials.TokenSource))
	if err != nil {
		return err
	}

	return nil
}
