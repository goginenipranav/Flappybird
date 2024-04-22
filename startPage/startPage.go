//pending: no rows...warning

package main

import (
	"database/sql"
	"fmt"
	"os/exec"

	rl "github.com/gen2brain/raylib-go/raylib"
	_ "github.com/go-sql-driver/mysql"
)

const (
	DBUsername = "sql5700809"
	DBPassword = "JcHxlM3Hp9"
	DBHost     = "sql5.freemysqlhosting.net"
	DBPort     = "3306"
	DBName     = "sql5700809"
)

type User struct {
	Name  string
	Lives int
	Score int
}

// TextInput represents a text input box.
type TextInput struct {
	rect       rl.Rectangle
	text       string
	active     bool
	fontSize   int32
	fontColor  rl.Color
	normalText string // Text to display when not active
	label      string // Label for the text input box
}

// NewTextInput creates a new text input box.
func NewTextInput(rect rl.Rectangle, label string, normalText string, fontSize int32, fontColor rl.Color) TextInput {
	return TextInput{
		rect:       rect,
		label:      label,
		normalText: normalText,
		fontSize:   fontSize,
		fontColor:  fontColor,
	}
}

// Draw draws the text input box.
func (t *TextInput) Draw() {
	// Draw label
	rl.DrawText(t.label, int32(t.rect.X), int32(t.rect.Y-25), t.fontSize, t.fontColor)

	if t.active {
		rl.DrawRectangleRec(t.rect, rl.White) // Background color
		rl.DrawRectangleLinesEx(t.rect, 1, rl.Black)
		rl.DrawText(t.text, int32(t.rect.X+5), int32(t.rect.Y+10), t.fontSize, t.fontColor)
	} else {
		rl.DrawRectangleRec(t.rect, rl.LightGray)
		rl.DrawRectangleLinesEx(t.rect, 1, rl.Gray)
		rl.DrawText(t.normalText, int32(t.rect.X+5), int32(t.rect.Y+10), t.fontSize, rl.Gray)
	}
}

// Update handles text input for the text input box.
func (t *TextInput) Update() {
	if rl.CheckCollisionPointRec(rl.GetMousePosition(), t.rect) && rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		t.active = true
	} else if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		t.active = false
	}

	if t.active {
		key := rl.GetKeyPressed()
		if key >= 32 && key <= 125 && len(t.text) < 16 {
			t.text += string(key)
		} else if key == rl.KeyBackspace && len(t.text) > 0 {
			t.text = t.text[:len(t.text)-1]
		}
	}
}

var gameNotFoundErrorDisplayed bool

func FetchUserDetails(username string) (string, int, int, error) {
	// Check if the error message has already been displayed
	if gameNotFoundErrorDisplayed {
		// Return nil error to indicate that no error occurred
		return "", 0, 0, nil
	}

	// Database connection parameters
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", DBUsername, DBPassword, DBHost, DBPort, DBName))
	if err != nil {
		fmt.Println("Error opening database connection:", err)
		return "", 0, 0, err
	}
	defer db.Close()

	fmt.Println("Database connection established successfully")

	// Query to fetch user details
	query := fmt.Sprintf("SELECT username, lives, score FROM User WHERE username='%s'", username)
	row := db.QueryRow(query)

	// Variables to store user details
	var name string
	var lives, score int

	// Scan user details from the database
	err = row.Scan(&name, &lives, &score)
	if err != nil {
		if err == sql.ErrNoRows {
			// Set the flag to indicate that the error message has been displayed
			gameNotFoundErrorDisplayed = true
			// Return a custom error indicating no game exists
			fmt.Println("No rows found for username:", username)
			return "", 0, 0, fmt.Errorf("no game exists")
		}
		fmt.Println("Error scanning row:", err)
		return "", 0, 0, err
	}

	fmt.Println("User details retrieved successfully:", name, lives, score)

	return name, lives, score, nil
}

func main() {
	screenWidth := int32(800)
	screenHeight := int32(600)

	rl.InitWindow(screenWidth, screenHeight, "Flappy Bird")

	// Load texture
	texture := rl.LoadTexture("icon2.png")
	if texture.ID == 0 {
		rl.TraceLog(rl.LogError, "Failed to load texture")
		return
	}

	// Padding for the texture from the left edge
	padding := int32(50)

	// Button position and size
	buttonWidth := int32(200)
	buttonHeight := int32(50)
	paddingY := int32(20)
	buttonX := screenWidth - buttonWidth - 50 // 50 pixels padding from the right edge
	buttonY := screenHeight/2 - (buttonHeight+paddingY*2+buttonHeight)/2

	// Username text box
	usernameBox := NewTextInput(rl.NewRectangle(float32(buttonX), float32(buttonY-buttonHeight-paddingY), float32(buttonWidth), float32(buttonHeight)), "Username:", "", 15, rl.Black)

	// Main loop
	for !rl.WindowShouldClose() {
		// Check if the button is clicked
		mousePosition := rl.GetMousePosition()
		buttonRect := rl.NewRectangle(float32(buttonX), float32(buttonY), float32(buttonWidth), float32(buttonHeight))

		if rl.CheckCollisionPointRec(mousePosition, buttonRect) && rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			// Launch new game
			if usernameBox.text == "" {
				usernameBox.text = "user"
			}
			cmd := exec.Command("go", "run", "../gamePage/gamePage.go", usernameBox.text, "3", "0")
			err := cmd.Run()
			if err != nil {
				rl.TraceLog(rl.LogError, "Failed to open gamePage.go:", err)
			}
			usernameBox.text = "" // Reset the username after starting the game
		}

		// Check if the "Continue Game" button is clicked
		continueButtonRect := rl.NewRectangle(float32(buttonX), float32(buttonY+buttonHeight+paddingY), float32(buttonWidth), float32(buttonHeight))
		if rl.CheckCollisionPointRec(mousePosition, continueButtonRect) && rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			// Fetch user details from the database
			// Launch new game
			if usernameBox.text == "" {
				usernameBox.text = "user"
			}
			username, lives, score, err := FetchUserDetails(usernameBox.text)
			if err != nil {
				rl.TraceLog(rl.LogError, "Failed to fetch user details:", err)
				continue
			}

			// Check if the user exists
			if username == "" || lives == 0 {
				// Inform the user that the user does not exist
				// Start a new game
				cmd := exec.Command("go", "run", "../gamePage/gamePage.go", usernameBox.text, "3", "0")
				err := cmd.Run()
				if err != nil {
					rl.TraceLog(rl.LogError, "Failed to open gamePage.go:", err)
				}
			} else {
				// Check if the user has lives greater than 0
				if lives > 0 {
					// Launch the game with user details
					cmd := exec.Command("go", "run", "../gamePage/gamePage.go", username, fmt.Sprintf("%d", lives), fmt.Sprintf("%d", score))
					err := cmd.Run()
					if err != nil {
						rl.TraceLog(rl.LogError, "Failed to open gamePage.go:", err)
					}
				}
			}
			usernameBox.text = ""
		}

		// Check if the "Leadership" button is clicked
		leadershipButtonRect := rl.NewRectangle(float32(buttonX), float32(buttonY+buttonHeight+paddingY*2+buttonHeight), float32(buttonWidth), float32(buttonHeight))
		if rl.CheckCollisionPointRec(mousePosition, leadershipButtonRect) && rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			// Navigate to the leadership page
			// Implement navigation logic here
			cmd := exec.Command("go", "run", "../leaderboardPage.go")
			err := cmd.Run()
			if err != nil {
				rl.TraceLog(rl.LogError, "Failed to open leaderboardPage.go:", err)
			}
		}

		// Update username text box
		usernameBox.Update()

		// Draw
		rl.BeginDrawing()

		rl.ClearBackground(getBiscuitColor())

		// Draw texture with left padding
		textureX := padding
		textureY := screenHeight/2 - texture.Height/2
		rl.DrawTexture(texture, textureX, textureY, rl.White)

		// Draw buttons
		rl.DrawRectangle(buttonX, buttonY, buttonWidth, buttonHeight, rl.Green)
		rl.DrawText("New Game", buttonX+50, buttonY+15, 20, rl.White)

		rl.DrawRectangle(buttonX, buttonY+buttonHeight+paddingY, buttonWidth, buttonHeight, rl.Green)
		rl.DrawText("Continue Game", buttonX+25, buttonY+buttonHeight+paddingY+15, 20, rl.White)

		rl.DrawRectangle(buttonX, buttonY+buttonHeight+paddingY*2+buttonHeight, buttonWidth, buttonHeight, rl.Blue)
		rl.DrawText("Leaderboard", buttonX+50, buttonY+buttonHeight+paddingY*2+buttonHeight+15, 20, rl.White)

		// Draw username text box
		usernameBox.Draw()

		// Set cursor appropriately
		if rl.CheckCollisionPointRec(mousePosition, buttonRect) || rl.CheckCollisionPointRec(mousePosition, continueButtonRect) || rl.CheckCollisionPointRec(mousePosition, leadershipButtonRect) {
			rl.SetMouseCursor(rl.MouseCursorPointingHand)
		} else {
			rl.SetMouseCursor(rl.MouseCursorDefault)
		}

		rl.EndDrawing()
	}

	// Unload texture
	rl.UnloadTexture(texture)

	rl.CloseWindow()
}

// getBiscuitColor returns the "biscuit" color (RGB: 255, 228, 196).
func getBiscuitColor() rl.Color {
	return rl.NewColor(255, 228, 196, 255)
}
