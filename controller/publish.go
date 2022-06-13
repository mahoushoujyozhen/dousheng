package controller

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
)

//	VideoListResponse  a response struct for publishList api
type VideoListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}

var dbConn *sql.DB

//func ReadFrameAsJpeg(inFileName string, frameNum int) []byte {
//	// Returns specified frame as []byte
//	buf := bytes.NewBuffer(nil)
//	err := ffmpeg.Input(inFileName).
//		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", frameNum)}).
//		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
//		WithOutput(buf, os.Stdout).Run()
//	if err != nil {
//		fmt.Println(err)
//		panic(err)
//	}

func initDB() (err error) {

	//xxxxxx is password
	dsn := "root:xxxxxx@tcp(175.178.106.176:3306)/Publish?charset=utf8mb4"
	// open函数只是验证格式是否正确，并不是创建数据库连接
	dbConn, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	//dbConn.SetMaxOpenConns(10)
	//dbConn.SetMaxIdleConns(5)
	//dbConn.SetConnMaxLifetime(time.Minute * 60)

	// 与数据库建立连接
	err = dbConn.Ping()
	if err != nil {
		return err
	}
	return nil
}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {
	//token := c.PostForm("token")
	//
	//if _, exist := usersLoginInfo[token]; !exist {
	//	//自己测试，没有tocken，自己写一个*************
	//	c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	//	return
	//}

	//使用c获取表单文件
	data, err := c.FormFile("data")

	//获取失败！
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	title := c.PostForm("title")
	//获取失败！
	if err != nil {
		fmt.Println("111")
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	fmt.Printf("title:%v\n", title)

	filename := filepath.Base(data.Filename)
	//user := usersLoginInfo[token]

	//userId := user.Id
	userId := 1

	//finalName := fmt.Sprintf("%d_%s", user.Id, filename)
	finalName := fmt.Sprintf("%d_%s", 1, filename)

	//这里将文件名and路径写入对应用户的数据库，这里设置不能够重名
	//这里需要获取上传用户的id

	err = initDB()
	if err != nil {
		fmt.Printf("err:%v\n", err)
		fmt.Println("mysql DB connect failed!!")
		return
	}
	//================================
	//后面模块化设计时i，要把defer改了，不能这样关闭
	defer dbConn.Close()

	coverUrl := "E:\\Project\\douyin\\demo\\simple-demo\\public\\coverImg\\bbb.jpeg"
	favoriteCount := 10
	isFavorite := true
	fmt.Println(filename)
	fmt.Println(userId)

	//查看该视频是否发布过，如果发布过，则不上传
	s := `select file_id from video where author_id=? and fileName=?`

	buf, err := dbConn.Exec(s, userId, finalName)

	if err != nil {

		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	//如果文件已经存在，提示不要重复发布
	if buf != nil {
		fmt.Printf("buf:%v\n", buf)
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "请勿重复发布视频",
		})
		return
	}

	//保存文件放在检测文件是否存在之后
	saveFile := filepath.Join("E:\\Project\\douyin\\demo\\simple-demo\\public\\video", finalName)
	fmt.Println(saveFile)
	err = c.SaveUploadedFile(data, saveFile)
	if err != nil {
		fmt.Println("23232323")
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	//var t assert.TestingT
	//获取视频第1帧作为图片
	//err = ffmpeg.Input("..\\public\\video\\1_test.mp4", ffmpeg.KwArgs{"ss": 1}).
	//	Output("..\\public\\coverImg\\1_test.mp4", ffmpeg.KwArgs{"t": 1}).OverWriteOutput().Run()
	//assert.Nil(t, err)
	//err = ffmpeg.Input("./sample_data/in1.mp4", ffmpeg.KwArgs{"ss": 1}).
	//	Output("./sample_data/out1.mp4", ffmpeg.KwArgs{"t": 1}).OverWriteOutput().Run()
	//assert.Nil(t, err)

	s = `insert into video (author_id,fileName,play_url,cover_url,favorite_count,is_favorite,title) values(?,?,?,?,?,?,?)`

	_, err = dbConn.Exec(s, userId, finalName, saveFile, coverUrl, favoriteCount, isFavorite, title)

	if err != nil {
		fmt.Printf("loaction:publish.go->85 Err:%v\n", err)
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error() + "\n请勿重复提交",
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  finalName + " uploaded successfully",
	})
}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {
	var err error
	err = initDB()
	if err != nil {
		fmt.Printf("err:%v\n", err)
		fmt.Println("mysql DB connect failed!!")
		return
	}
	//================================
	//后面模块化设计时i，要把defer改了，不能这样关闭
	defer dbConn.Close()

	userId := c.PostForm("user_id")

	fmt.Printf("111%v", userId)

	if userId == "" {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "请登录",
		})
	}

	//0长度的切片
	var videoList = make([]Video, 10)
	userTemp := User{
		Id:            0,
		Name:          "",
		FollowCount:   0,
		FollowerCount: 0,
		IsFollow:      false,
	}
	videoTemp := Video{
		Id:            0,
		Author:        userTemp,
		PlayUrl:       "",
		CoverUrl:      "",
		FavoriteCount: 0,
		CommentCount:  0,
		IsFavorite:    false,
		Title:         "",
	}
	fmt.Println(videoTemp)

	//读author_id=userId的author的所有信息，赋值给userTemp
	s := `select author_id,name,follow_count,follower_count,is_follow from author where author_id=?`
	row := dbConn.QueryRow(s, userId)
	err = row.Scan(&userTemp.Id, &userTemp.Name, &userTemp.FollowerCount, &userTemp.FollowerCount, &userTemp.IsFollow)
	if err != nil {
		fmt.Println("Scan author msg err!")
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
	}

	//查找author_id=userId的用户发布的所有视频信息
	s = `select file_id,play_url,cover_url,favorite_count,comment_count,is_favorite,title 
			from video where author_id=?`
	rc, err := dbConn.Query(s, userId)
	if err != nil {
		fmt.Println("Query video msg err!")
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
	}
	for rc.Next() {
		err = rc.Scan(&videoTemp.Id, &videoTemp.PlayUrl, &videoTemp.CoverUrl, &videoTemp.FavoriteCount, &videoTemp.CommentCount, &videoTemp.IsFavorite, &videoTemp.Title)
		if err != nil {
			fmt.Println("line 226 rc.Scan err! ")
			c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			})
		}
		videoList = append(videoList, videoTemp)
		fmt.Println(videoTemp)
	}

	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: DemoVideos,
	})
}
