package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/ashikkabeer/messaging-api/config/db"
	"github.com/ashikkabeer/messaging-api/models"
	"github.com/ashikkabeer/messaging-api/queue/sender"
	"github.com/ashikkabeer/messaging-api/utils"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func SendMessage(c *gin.Context) {
    var req models.RequestBody
    
    if err := c.ShouldBindJSON(&req);

    err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
        return;
    }

    // validating the users before sending message
    senderID := req.SenderID
    receiverID := req.ReceiverID

    if senderID == "" || receiverID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Missing senderID or receiverID"})
        return
    }

    user:= isUsersExist(senderID, receiverID)
   
    if !user {
        c.JSON(http.StatusNotFound, gin.H{"error": "Users not found"})
        return
    }

    if err := sender.SendMessageToQueue(models.RequestBody{
        SenderID: req.SenderID,
        ReceiverID: req.ReceiverID,
        Content: req.Content,
    });
    err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
        return
    }
    log.Println("Message published to RabbitMQ")

    c.JSON(http.StatusOK, gin.H{"status": "Success", "message": req})
}


func RetrieveHistory(c *gin.Context) {
    firstUser := c.Query("user1")
    secondUser := c.Query("user2")
    cursor := c.Query("cursor")

    limit, errs := strconv.Atoi(c.Query("limit"))
    if errs != nil {
        return
    }

    var rows *sql.Rows
    var err error

    if firstUser == "" || secondUser == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Missing user1 or user2"})
        return
    }

    user := isUsersExist(firstUser, secondUser)

    if !user {
        c.JSON(http.StatusNotFound, gin.H{"error": "Users not found"})
        return
    }

    var timestamp time.Time

    // retrieving latest messages between the two users
    if cursor == "" {
        query := `SELECT id, senderID, receiverID, content, read, created_at 
        FROM messages 
        WHERE (senderID = $1 AND receiverID = $2) 
        OR (senderID = $2 AND receiverID = $1) 
        ORDER BY created_at DESC
        LIMIT $3;`
        rows, err = db.Query(query, firstUser, secondUser, limit)
    } else {
        // retrieving messages with cursor
        timestamp, _, err = utils.DecodeCursor(cursor)
        if err!= nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cursor"})
            return 
        }
        query := `SELECT id, senderID, receiverID, content, read, created_at 
        FROM messages 
        WHERE (senderID = $1 AND receiverID = $2) 
        OR (senderID = $2 AND receiverID = $1)
        AND created_at < $3
        ORDER BY created_at DESC
        LIMIT $4;
    `


    rows, err = db.Query(query, firstUser, secondUser, timestamp, limit)
    }

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    defer rows.Close()

    var messages []models.Message
    for rows.Next() {
        var msg models.Message
        err := rows.Scan(&msg.ID, &msg.SenderID, &msg.ReceiverID, &msg.Content, &msg.Read, &msg.CreatedAt)

        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        messages = append(messages, msg)
    }

    if len(messages) > 0 {
        lastMessage := messages[len(messages)-1]
        cursor = utils.EncodeCursor(lastMessage.CreatedAt, lastMessage.ID)
        c.JSON(http.StatusOK, gin.H{
            "messages": messages,
            "next_cursor": cursor,
            "has_more": len(messages) == int(limit),
        })
    } else {
        c.JSON(http.StatusOK, gin.H{
            "messages": messages,
            "next_cursor": nil,
            "has_more": false,
        })
    }

    if err = rows.Err(); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
}

func MarkAsRead(c *gin.Context) {
    messageId := c.Param("message_id");
    MarkAsReadAsync(messageId)
    c.JSON(http.StatusOK, gin.H{"status":"read"})
}
func MarkAsReadAsync(messageId string) {
    go func() {
        query := `UPDATE messages SET read = true WHERE id = $1`
        _,err := db.Exec(query, messageId)
        if err != nil {
            log.Printf("update failed")
        }
    }()
}

func isUsersExist(sender string, receiver string) bool {
    query := `SELECT COUNT(*)
    FROM users
    WHERE id IN ($1, $2)`
    var count int
    err := db.QueryRow(query, sender, receiver).Scan(&count)

    if err!= nil {
        return false
    }
    return count == 2
}