module startPage

go 1.22.2

replace gamePage => ../gamePage

require github.com/gen2brain/raylib-go/raylib v0.0.0-20240408130534-53f26d8a0802

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/ebitengine/purego v0.6.0-alpha.1.0.20231122024802-192c5e846faa // indirect
	github.com/go-sql-driver/mysql v1.8.1 // indirect
	golang.org/x/sys v0.14.0 // indirect
)

replace leaderboardPage => ../leaderboardPage.go
