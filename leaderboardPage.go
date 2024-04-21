package main

import (
	"database/sql"
	"fmt"
	"log"

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
	Score int
}

func main() {
	screenWidth := int32(800)
	screenHeight := int32(600)

	rl.InitWindow(screenWidth, screenHeight, "Leaderboard")

	// Connect to the database
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", DBUsername, DBPassword, DBHost, DBPort, DBName))
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Query for leaderboard scores
	rows, err := db.Query("SELECT username, score FROM User ORDER BY score DESC ")
	if err != nil {
		log.Fatal("Failed to fetch leaderboard scores:", err)
	}
	defer rows.Close()

	// Parse leaderboard scores
	var leaderboard []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.Name, &user.Score)
		if err != nil {
			log.Fatal("Failed to scan leaderboard row:", err)
		}
		leaderboard = append(leaderboard, user)
	}
	if err := rows.Err(); err != nil {
		log.Fatal("Error while iterating over leaderboard rows:", err)
	}

	// Main loop
	for !rl.WindowShouldClose() {
		// Draw
		rl.BeginDrawing()

		rl.ClearBackground(rl.Black)

		// Draw leaderboard
		y := int32(50)
		for i, user := range leaderboard {
			rl.DrawRectangle(50, y, 700, 30, rl.Fade(rl.SkyBlue, 0.2)) // Background color with padding
			rl.DrawRectangleLines(50, y, 700, 30, rl.Black)            // Border
			rl.DrawText(fmt.Sprintf("%d. %s", i+1, user.Name), 100, y+5, 20, rl.White)
			rl.DrawText(fmt.Sprintf("%d", user.Score), 640, y+5, 20, rl.White) // Adjusted horizontal position for the score
			y += 35                                                            // Reduced vertical spacing between entries
		}

		// Draw title
		rl.DrawText("Leaderboard", screenWidth/2-rl.MeasureText("Leaderboard", 30)/2, 10, 30, rl.White)

		rl.EndDrawing()
	}

	rl.CloseWindow()
}
