package models

import "time"

type User struct {
	ID        string    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Email     string    `json:"email" db:"email"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

type Task struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_Date"`
	IsCompleted bool      `json:"isCompleted"`
}

type RegisterUser struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateTask struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	PendingAt   time.Time `json:"pending_at"`
}

type UpdateTask struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	PendingAt   time.Time `json:"pending_at"`
	IsCompleted bool      `json:"isCompleted"`
}
