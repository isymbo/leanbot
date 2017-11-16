package main

import (
    "net/http"
    "os"
    "os/exec"
    "time"
    "io"
    "log"
    "strings"

    "github.com/gin-gonic/gin"
    "github.com/gin-gonic/contrib/ginrus"
    "github.com/sirupsen/logrus"

    "golang.org/x/sync/errgroup"
)

var (
    g errgroup.Group
)

func accept() http.Handler {
    r := gin.New()

    // Logging to a file.
    f, _ := os.Create("leanbot.log")

    logger := logrus.New()
    logger.Level = logrus.DebugLevel
    logger.Out = io.MultiWriter(f, os.Stdout)
    r.Use(ginrus.Ginrus(logger, time.RFC3339, false))

    r.Use(gin.Recovery())

    r.POST("/accept", func(c *gin.Context) {
        user_name := c.PostForm("user_name")
        text := c.PostForm("text")

        // Remove unnesessary whitespaces
        cli := strings.TrimSpace(text)
        arg := strings.TrimLeft(cli, "!clips ")
        args := strings.Fields(arg)

        out, err := exec.Command("clips", args...).CombinedOutput()

        if err != nil {
            log.Println(err)
        }

        output := string(out)

        logger.Infoln("User: ", user_name)
        logger.Infoln("Command: ", text)
        logger.Infoln("Output: ", output)

        c.JSON(
            http.StatusOK,
            gin.H{
                "text": output,
                //"attachments": [{
                //  "title": "output",
                //  "description": "Description",
                //  "url": "http://pubu.im",
                //  "color": "info"
                //}],
                "username": "CLIPS",
                "icon_url":
                "http://open.bestv.com.cn/icons/clips.png",
            },
        )
    })

    return r
}

func main() {
    // Disable Console Color, you don't need console color when writing the
    // logs to file.
    gin.DisableConsoleColor()

    // Logging to a file.
    //f, _ := os.Create("leanbot.log")
    // gin.DefaultWriter = io.MultiWriter(f)

    // Use the following code if you need to write the logs to file and console
    // at the same time.
    //gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

    leanbot := &http.Server{
        Addr:         ":8080",
        Handler:      accept(),
        ReadTimeout:  5 * time.Second,
        WriteTimeout: 10 * time.Second,
    }

    g.Go(func() error {
        return leanbot.ListenAndServe()
    })

    if err := g.Wait(); err != nil {
        log.Fatal(err)
    }
}
