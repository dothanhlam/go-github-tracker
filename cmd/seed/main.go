package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/dothanhlam/go-github-tracker/internal/database"
	"github.com/dothanhlam/go-github-tracker/internal/store"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	dbFile := "dummy.db"
	_ = os.Remove(dbFile)

	db, err := database.Connect("sqlite3", dbFile)
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}
	defer db.Close()

	if err := db.RunMigrations(); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	// Insert Dummy Data using Raw SQL for Teams
	_, err = db.Exec(`INSERT INTO teams (id, name, description, created_at, updated_at) VALUES 
		(1, 'Backend Protectors', 'Core backend systems team', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
		(2, 'Frontend Wizards', 'User interface squad', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
		(3, 'Data Architects', 'Data engineering team', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`)
	if err != nil {
		log.Fatalf("failed to insert teams: %v", err)
	}

	// Insert Members
	_, err = db.Exec(`INSERT INTO team_memberships (team_id, github_username, allocation_weight, joined_at, created_at) VALUES 
		(1, 'alice', 1.0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
		(1, 'bob', 1.0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
		(2, 'charlie', 1.0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
		(2, 'diana', 1.0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
		(3, 'eve', 1.0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
		(3, 'frank', 1.0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`)
	if err != nil {
		log.Fatalf("failed to insert members: %v", err)
	}

	st := store.New(db)
	now := time.Now()
	rand.Seed(time.Now().UnixNano())

	// Insert 100 random PR metrics across the 3 teams over the past 90 days
	for i := 1; i <= 100; i++ {
		teamID := rand.Intn(3) + 1
		daysAgo := rand.Intn(90)
		createdAt := now.AddDate(0, 0, -daysAgo)
		
		cycleHours := rand.Intn(168) + 1 // 1 hour to 1 week
		mergedAt := createdAt.Add(time.Duration(cycleHours) * time.Hour)
		
		firstReviewAfter := rand.Intn(48) + 1
		firstReviewAt := createdAt.Add(time.Duration(firstReviewAfter) * time.Hour)
		turnaroundHours := firstReviewAfter
		
		comments := rand.Intn(20)
		conversations := rand.Intn(5)
		reviewers := rand.Intn(3) + 1
		
		pr := &database.PRMetric{
			TeamID:                 teamID,
			PRNumber:               1000 + i,
			Repository:             fmt.Sprintf("repo-%d", teamID),
			Author:                 "user",
			Title:                  fmt.Sprintf("Dummy PR %d", i),
			CreatedAt:              createdAt,
			MergedAt:               &mergedAt,
			ClosedAt:               &mergedAt,
			CycleTimeHours:         &cycleHours,
			State:                  "merged",
			CreatedDate:            &createdAt,
			FirstReviewAt:          &firstReviewAt,
			ReviewTurnaroundHours:  &turnaroundHours,
			ReviewCommentsCount:    comments,
			ConversationCount:      conversations,
			ChangesRequestedCount:  rand.Intn(2),
			ApprovedCount:          1,
			ReviewersCount:         reviewers,
			ExternalReviewersCount: rand.Intn(2),
			ReviewersList:          `["reviewer1"]`,
		}
		
		if err := st.UpsertPRMetric(pr); err != nil {
			log.Fatalf("Failed to insert PR %d: %v", i, err)
		}
	}

	fmt.Printf("Successfully created %s with dummy data!\n", dbFile)
	fmt.Printf("\nTo test the TUI with this DB, add this to your .env.local:\n")
	fmt.Printf("DB_DRIVER=sqlite3\n")
	fmt.Printf("DB_URL=dummy.db\n")
}
