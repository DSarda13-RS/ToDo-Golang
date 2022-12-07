package handler

import (
	"ToDo/database"
	"ToDo/database/dbHelper"
	"ToDo/middlewares"
	"ToDo/models"
	"ToDo/utils"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	body := models.RegisterUser{}
	if parseErr := utils.ParseBody(r.Body, &body); parseErr != nil {
		utils.RespondError(w, http.StatusBadRequest, parseErr, "Failed to parse request body.")
		return
	}
	if len(body.Password) < 6 {
		utils.RespondError(w, http.StatusBadRequest, nil, "Password must contain at-least 6 Characters.")
		return
	}
	exists, existsErr := dbHelper.IsUserExist(body.Email)
	if existsErr != nil {
		utils.RespondError(w, http.StatusBadRequest, existsErr, "Failed to check users' existence.")
		return
	}
	if exists {
		utils.RespondError(w, http.StatusBadRequest, nil, "User already exists.")
		return
	}
	hashedPassword, hashErr := utils.HashPassword(body.Password)
	if hashErr != nil {
		utils.RespondError(w, http.StatusBadRequest, hashErr, "Failed to secure password.")
		return
	}
	sessionToken := utils.HashString(body.Email + time.Now().String())
	txErr := database.Tx(func(tx *sqlx.Tx) error {
		userID, saveErr := dbHelper.CreateUser(tx, body.Name, body.Email, hashedPassword)
		if saveErr != nil {
			return saveErr
		}
		sessionErr := dbHelper.CreateUserSession(tx, userID, sessionToken)
		if sessionErr != nil {
			return sessionErr
		}
		return nil
	})
	if txErr != nil {
		utils.RespondError(w, http.StatusBadRequest, txErr, "Failed to create User.")
		return
	}
	utils.RespondJSON(w, http.StatusOK, struct {
		Token string `json:"token"`
	}{
		Token: sessionToken,
	})
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	body := models.LoginUser{}
	if parseErr := utils.ParseBody(r.Body, &body); parseErr != nil {
		utils.RespondError(w, http.StatusBadRequest, parseErr, "Failed to parse request body.")
		return
	}
	userID, userErr := dbHelper.GetUserIDByPassword(body.Email, body.Password)
	if userErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, userErr, "Failed to find user.")
		return
	}
	sessionToken := utils.HashString(body.Email + time.Now().String())
	sessionErr := dbHelper.CreateUserSession(database.Todo, userID, sessionToken)
	if sessionErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, sessionErr, "Failed to create user session.")
		return
	}
	utils.RespondJSON(w, http.StatusOK, struct {
		Token string `json:"token"`
	}{
		Token: sessionToken,
	})
}

func GetUserInfo(w http.ResponseWriter, r *http.Request) {
	userCtx := middlewares.UserContext(r)
	utils.RespondJSON(w, http.StatusOK, userCtx)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("x-api-key")
	userCtx := middlewares.UserContext(r)
	err := dbHelper.DeleteSessionToken(userCtx.ID, token)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to logout user")
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	body := models.CreateTask{}
	userCtx := middlewares.UserContext(r)
	if parseErr := utils.ParseBody(r.Body, &body); parseErr != nil {
		utils.RespondError(w, http.StatusBadRequest, parseErr, "Failed to parse request body.")
		return
	}
	taskErr := dbHelper.CreateTask(database.Todo, userCtx.ID, body.Name, body.Description, body.PendingAt)
	if taskErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, taskErr, "Failed to create Task.")
		return
	}
	utils.RespondJSON(w, http.StatusOK, struct {
		Msg string `json:"msg"`
	}{
		Msg: "Task created.",
	})
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	body := models.UpdateTask{}
	taskID := r.URL.Query().Get("id")
	userCtx := middlewares.UserContext(r)
	if parseErr := utils.ParseBody(r.Body, &body); parseErr != nil {
		utils.RespondError(w, http.StatusBadRequest, parseErr, "Failed to parse request body.")
		return
	}
	task, err := dbHelper.GetTask(userCtx.ID, taskID)
	if err != nil || task == nil {
		logrus.WithError(err).Errorf("Failed to get task with ID: %s", taskID)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var updateName string
	var updateDescription string
	var updatePendingAt time.Time
	var updateIsCompleted bool
	if body.Name != "" {
		updateName = body.Name
	} else {
		updateName = task.Name
	}
	if body.Description != "" {
		updateDescription = body.Description
	} else {
		updateDescription = task.Description
	}
	if time.Time.IsZero(body.PendingAt) {
		updatePendingAt = task.DueDate
	} else {
		updatePendingAt = body.PendingAt
	}
	if body.IsCompleted != false {
		updateIsCompleted = body.IsCompleted
	} else {
		updateIsCompleted = task.IsCompleted
	}
	updateErr := dbHelper.UpdateTask(database.Todo, updateName, updateDescription, userCtx.ID, taskID, updatePendingAt, updateIsCompleted)
	if updateErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, updateErr, "Failed to update Task.")
		return
	}
	utils.RespondJSON(w, http.StatusOK, struct {
		Msg string `json:"msg"`
	}{
		Msg: "Task updated.",
	})
}

func GetTask(w http.ResponseWriter, r *http.Request) {
	taskID := r.URL.Query().Get("id")
	userCtx := middlewares.UserContext(r)
	task, err := dbHelper.GetTask(userCtx.ID, taskID)
	if err != nil || task == nil {
		logrus.WithError(err).Errorf("Failed to get task with ID: %s", taskID)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	utils.RespondJSON(w, http.StatusOK, task)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	taskID := r.URL.Query().Get("id")
	userCtx := middlewares.UserContext(r)
	deleteErr := dbHelper.DeleteTask(database.Todo, userCtx.ID, taskID)
	if deleteErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, deleteErr, "Failed to delete Task.")
		return
	}
	utils.RespondJSON(w, http.StatusOK, struct {
		Msg string `json:"msg"`
	}{
		Msg: "Task deleted.",
	})
}
