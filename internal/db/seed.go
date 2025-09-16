package db

import (
	"backendwithgo/internal/store"
	"context"
	"fmt"
	"log"
	"math/rand"
	"database/sql"
)

var usernames = []string{"niga", "tom", "bob", "jerry", "alice", "eve", "mallory", "peggy", "trent", "victor", "walter"}

var titles = []string{
	"Unleashing the Future: How AI is Reshaping Our World",
	"Secrets of the Digital Underground",
	"Beyond the Horizon: Exploring the Unknown",
	"Mastering Chaos: Strategies for Thriving in Uncertain Times",
	"The Power Playbook: Unlocking Hidden Potential",
	"Zero to Hero: A Journey of Transformation",
	"Breaking the Code: Inside the Hacker’s Mind",
	"Rising Titans: The Battle for Innovation",
	"Edge of Reality: Where Science Meets Imagination",
	"The Hidden Blueprint to Success",
	"From Sparks to Flames: Igniting Your Passion",
	"Shadows of Tomorrow: Predicting the Next Big Wave",
	"The Art of Fearless Living",
	"Legends Reborn: Stories That Inspire",
	"Game Changers: Leaders Who Rewrite the Rules",
	"The Infinite Quest: Pushing Human Limits",
	"Unlocking the Vault: Secrets of High Achievers",
	"Storming the Summit: Conquering Impossible Goals",
	"Fueling the Fire: Turning Ideas into Impact",
	"Epic Journeys: Tales of Adventure and Discovery",
}

var contents = []string{
	"Content that Captivates and Converts",
	"The Secret Formula Behind Viral Content",
	"Crafting Stories that Spark Emotion",
	"Content that Builds Trust and Authority",
	"From Clicks to Conversions: Content That Sells",
	"The Hidden Psychology of Persuasive Content",
	"Unlocking Creativity: Content Ideas That Never Run Out",
	"Content that Ranks: Winning the SEO Game",
	"The Future of Interactive Content",
	"From Boring to Brilliant: Transforming Your Content",
	"Content that Inspires Action and Change",
	"How to Create Evergreen Content That Lasts Forever",
	"Content that Turns Strangers into Fans",
	"Behind the Scenes: The Making of Magnetic Content",
	"Data-Driven Content Strategies that Work",
	"Content that Sparks Conversations",
	"How to Scale Your Content Without Losing Quality",
	"Content that Dominates Social Media",
	"The Blueprint for High-Impact Content",
	"Content that Defines Brands and Moves Markets",
}

var tags = []string{"tech", "life", "music", "travel", "food", "science", "health", "fitness", "education", "finance", "history", "art", "culture", "nature", "sports", "politics", "environment", "fashion", "gaming", "movies"}

var comments = []string{
	"Great post! Really enjoyed reading it.",
	"Thanks for sharing this information.",
	"Interesting perspective, I hadn't thought of it that way.",
	"I disagree with some points, but overall a good read.",
	"Can you provide more details on this topic?",
	"Looking forward to your next post!",
	"This was very helpful, thank you!",
	"I learned something new today.",
	"Well written and easy to understand.",
	"Could you recommend further reading?",
	"Your insights are always appreciated.",
	"This topic is very relevant in today's world.",
	"Thanks for the tips, I'll definitely try them out.",
	"Fantastic article, keep up the good work!",
	"I have a question about one of your points.",
	"This is exactly what I was looking for.",
	"Your writing style is very engaging.",
}

func Seed(store store.Storage, db *sql.DB) {
	ctx := context.Background()

	users := generateUsers(100)
	tx, _ := db.BeginTx(ctx, nil)

	for _, user := range users {
		if err := store.Users.Create(ctx, tx, user); err != nil {
			_ = tx.Rollback()
			log.Println("Error creating user:", err)
			return
		}
	}

	tx.Commit()

	posts := generatePosts(200, users)
	for _, post := range posts {
		if err := store.Posts.Create(ctx, post); err != nil {
			log.Println("Error creating post:", err)
			return
		}
	}

	comments := generateComments(500, posts, users)
	for _, comment := range comments {
		if err := store.Comments.Create(ctx, comment); err != nil {
			log.Println("Error creating comment:", err)
			return
		}
	}

	log.Println("Seed Success")
}
func generateUsers(num int) []*store.User {
	users := make([]*store.User, num)

	for i := 0; i < num; i++ {
		users[i] = &store.User{
			Username: usernames[i%len(usernames)] + fmt.Sprintf("%d", i),
			Email:    usernames[i%len(usernames)] + fmt.Sprintf("%d", i) + "@example.com",
			Role: store.Role{
				Name: "user",
			},
		}
	}

	return users
}

func generatePosts(num int, users []*store.User) []*store.Post {
	posts := make([]*store.Post, num)
	for i := 0; i < num; i++ {
		user := users[rand.Intn(len(users))]
		posts[i] = &store.Post{
			UserID:  user.ID,
			Title:   titles[rand.Intn(len(titles))],
			Content: contents[rand.Intn(len(contents))],
			Tags: []string{
				tags[rand.Intn(len(tags))],
				tags[rand.Intn(len(tags))],
				tags[rand.Intn(len(tags))],
			},
		}
	}
	return posts
}

func generateComments(num int, posts []*store.Post, users []*store.User) []*store.Comment {
	cms := make([]*store.Comment, num)
	for i := 0; i < num; i++ {
		post := posts[rand.Intn(len(posts))]
		user := users[rand.Intn(len(users))]
		cms[i] = &store.Comment{
			PostID:  post.ID,
			UserID:  user.ID,
			Content: comments[rand.Intn(len(comments))],
		}
	}
	return cms
}
