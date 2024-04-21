package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

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

type Apple struct {
	posX   int32
	posY   int32
	width  int32
	height int32
	Color  rl.Color
}

type User struct {
	ID    int
	Name  string
	Lives int
	Score int
}

var appleTexture rl.Texture2D // Declaring the apple texture globally

func main() {
	screenWidth := int32(1000)
	screenHeight := int32(700)
	rl.InitAudioDevice()
	eatNoise := rl.LoadSound("../sound/eat.wav")
	rl.InitWindow(screenWidth, screenHeight, "FlappyApples")
	rl.SetTargetFPS(60)
	birdDown := rl.LoadImage("../assets/bird-down.png")
	birdUp := rl.LoadImage("../assets/bird-up.png")
	texture := rl.LoadTextureFromImage(birdUp)
	rand.Seed(time.Now().UnixNano())
	var appleLoc int = rand.Intn(660-50+1) + 50
	Apples := []Apple{}
	currentApple := Apple{screenWidth - 100, int32(appleLoc), 74, 74, rl.Red} // Adjusted apple dimensions
	Apples = append(Apples, currentApple)
	var xCoords int32 = screenWidth/2 - texture.Width/2
	var yCoords int32 = screenHeight/2 - texture.Height/2 - 40
	var score int = 0
	var lives int = 3 // Number of lives
	var name string
	saveClicked := false

	// Receive name from start page
	if len(os.Args) > 1 {
		name = os.Args[1]
		fmt.Println("Name:", name)

		// Convert lives to int
		_, err := fmt.Sscanf(os.Args[2], "%d", &lives)
		if err != nil {
			fmt.Println("Failed to parse lives:", err)
			return
		}
		fmt.Println("Lives:", lives)

		// Convert score to int
		_, err = fmt.Sscanf(os.Args[3], "%d", &score)
		if err != nil {
			fmt.Println("Failed to parse score:", err)
			return
		}
		fmt.Println("Score:", score)
	} else {
		name = "Player"
		score = 0
		lives = 3
	}

	saveButton := rl.NewRectangle(float32(screenWidth-120), 10, 110, 40)

	// Fetch high score from the database
	highScore := getHighScore()
	bgTexture := rl.LoadTexture("../assets/bgsky.png")

	// Load the apple image
	appleImage := rl.LoadImage("../apple10.png")
	appleTexture = rl.LoadTextureFromImage(appleImage)
	rl.UnloadImage(appleImage)

	for !rl.WindowShouldClose() && lives > 0 {
		rl.BeginDrawing()

		rl.ClearBackground(rl.NewColor(255, 228, 196, 255))
		rl.DrawTexture(bgTexture, 0, 0, rl.White)

		// Draw the labels for score, name, lives, and high score
		labelY := int32(0)
		rl.DrawText("Score:       Lives:       Name:       High Score:", 10, labelY, 30, rl.Black)
		labelY += 40

		// Draw the values for score, name, lives, and high score
		valueY := int32(40)
		rl.DrawText(strconv.Itoa(score), 50, valueY, 30, rl.Black)
		rl.DrawText(strconv.Itoa(lives), 220, valueY, 30, rl.Black)
		rl.DrawText(name, 360, valueY, 30, rl.Black)
		rl.DrawText(strconv.Itoa(highScore), 560, valueY, 30, rl.Black)
		// Draw the apples
		for io, currentApple := range Apples {
			rl.DrawTexture(appleTexture, currentApple.posX, currentApple.posY, rl.White)
			Apples[io].posX -= 5
			if currentApple.posX < 0 {
				Apples[io].posX = 800
				Apples[io].posY = int32(rand.Intn(int(screenHeight-40-50)) + 50)
				score--
				lives-- // Decrease lives when an apple is missed
				if score < 0 {
					score = 0
				}
			}
			birdCollisionRect := rl.NewRectangle(float32(xCoords)+15, float32(yCoords)+10, 20, 10)
			if rl.CheckCollisionRecs(birdCollisionRect, rl.NewRectangle(float32(currentApple.posX)+10, float32(currentApple.posY)+10, float32(currentApple.width)-20, float32(currentApple.height)-20)) {
				Apples[io].posX = 800
				Apples[io].posY = int32(rand.Intn(580-2+1) - 2)
				score++
				rl.PlaySound(eatNoise)
			}
		}

		// Draw the bird
		rl.DrawTexture(texture, xCoords, yCoords, rl.White)

		// Draw the save button
		rl.DrawRectangleRec(saveButton, rl.Blue)
		rl.DrawText("Save", int32(saveButton.X)+25, int32(saveButton.Y)+10, 20, rl.White)

		if rl.CheckCollisionPointRec(rl.GetMousePosition(), saveButton) && rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			saveClicked = true
		}

		if saveClicked {
			saveGameDetails(name, score, lives)
			rl.CloseWindow()
		}

		if rl.IsKeyDown(rl.KeySpace) {
			texture = rl.LoadTextureFromImage(birdUp)
			yCoords -= 5
		} else {
			texture = rl.LoadTextureFromImage(birdDown)
			yCoords += 5
		}

		// Reduce lives if bird goes below screen
		if yCoords > 700 {
			lives--
			if lives > 0 {
				yCoords = screenHeight/2 - texture.Height/2 - 40
			}
			if lives == 0 {
				rl.UnloadTexture(texture)
				Apples = nil
				rl.DrawText("Your final score is: "+strconv.Itoa(score), 30, 40, 30, rl.Red)
			}
		}

		rl.EndDrawing()
		time.Sleep(50000000)
	}

	rl.UnloadSound(eatNoise)
	rl.UnloadTexture(texture)
	rl.UnloadTexture(appleTexture)
}

func saveGameDetails(name string, score, lives int) {
	// Connect to the database
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", DBUsername, DBPassword, DBHost, DBPort, DBName))
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Insert new user record
	newUser := User{Name: name, Lives: lives, Score: score}
	err = insertUser(db, newUser)
	if err != nil {
		log.Fatal("Failed to insert user:", err)
	}
}

func insertUser(db *sql.DB, user User) error {
	query := `
        INSERT INTO User (username, lives, score)
        VALUES (?, ?, ?)
        ON DUPLICATE KEY UPDATE
        lives = VALUES(lives),
        score = VALUES(score)
    `
	_, err := db.Exec(query, user.Name, user.Lives, user.Score)
	if err != nil {
		log.Printf("Failed to execute query: %v\n", err)
	}
	return err
}

func getHighScore() int {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", DBUsername, DBPassword, DBHost, DBPort, DBName))
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	var highScore int
	err = db.QueryRow("SELECT MAX(score) FROM User").Scan(&highScore)
	if err != nil {
		log.Fatal("Failed to get high score:", err)
	}

	return highScore
}
