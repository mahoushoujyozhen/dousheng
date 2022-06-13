package controller

var DemoVideos = []Video{
	{
		Id:      1,
		Author:  DemoUser,
		PlayUrl: "https://www.w3schools.com/html/movie.mp4",
		//PlayUrl: "http://10.0.2.2:8080/E:/Project/douyin/demo/simple-demo/public/video/1_test.mp4",
		CoverUrl: "https://cdn.pixabay.com/photo/2016/03/27/18/10/bear-1283347_1280.jpg",
		//CoverUrl:      "http://10.0.2.2:8080/E:/Project/douyin/demo/simple-demo/public/coverImg/bbb.jpeg",
		FavoriteCount: 5,
		CommentCount:  10,
		IsFavorite:    false,
	},
}

var DemoComments = []Comment{
	{
		Id:         1,
		User:       DemoUser,
		Content:    "Test Comment",
		CreateDate: "05-01",
	},
}

var DemoUser = User{
	Id:            1,
	Name:          "TestUser",
	FollowCount:   0,
	FollowerCount: 0,
	IsFollow:      false,
}
