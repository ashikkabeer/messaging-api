package handlers


import (
    "net/http"
    "database/sql"
    "fmt"
    "github.com/ashikkabeer/messaging-api/config/db"
    "github.com/ashikkabeer/messaging-api/queue/sender"
    "github.com/ashikkabeer/messaging-api/models"
    "github.com/gin-gonic/gin"
    _ "github.com/lib/pq"
)

func SendMessage(c *gin.Context) {
    var req models.RequestBody
    
    if err := c.ShouldBindJSON(&req);
    err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Message not available in the request body"})
        return;
    }

    // message is being sent to the Queue
    messageSender, err := sender.NewSender()
    if err!= nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error in creating a new sender"})
        return;
    }
    defer messageSender.Close()

    if err := messageSender.SendMessage(models.RequestBody{
        SenderID: req.SenderID,
        ReceiverID: req.ReceiverID,
        Content: req.Content,
    });
    err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error in sending message"})
        return
    }
    fmt.Println("pushed to queue")

    c.JSON(http.StatusOK, gin.H{"status": "Success", "message": req})
}


func RetrieveHistory(c *gin.Context) {
    firstUser := c.Query("user1")
    secondUser := c.Query("user2")
    cursor := c.Query("cursor")
    limit := 20
    var rows *sql.Rows
    var err error

    if firstUser == "" || secondUser == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "userId query parameter is required"})
        return
    }

    // check if both users exist
    query := `SELECT COUNT(*) 
    FROM users 
    WHERE id IN ($1, $2)`
    var count int
    err = db.QueryRow(query, firstUser, secondUser).Scan(&count)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    if count != 2 {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    // retrieving latest messages between the two users
    if cursor == "" {
        query = `SELECT id, senderID, receiverID, content, read, created_at 
        FROM messages 
        WHERE (senderID = $1 AND receiverID = $2) 
        OR (senderID = $2 AND receiverID = $1)
        ORDER BY created_at DESC
        LIMIT $3;`
        rows, err = db.Query(query, firstUser, secondUser, limit)
    } else {
        // retrieving messages with cursor
        query = `SELECT id, senderID, receiverID, content, read, created_at 
        FROM messages 
        WHERE (senderID = $1 AND receiverID = $2) 
        OR (senderID = $2 AND receiverID = $1)
        AND created_at < $3
        ORDER BY created_at DESC
        LIMIT $4;
    `
    rows, err = db.Query(query, firstUser, secondUser, cursor, limit)
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
        c.JSON(http.StatusOK, gin.H{
            "messages": messages,
            "next_cursor": lastMessage.CreatedAt,
            "has_more": len(messages) == limit,
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
    query := `UPDATE messages SET read = true WHERE id = $1`
    _, err := db.Exec(query, messageId)
    if err!= nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error in marking message as read"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"status": "read"})
}