package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"strconv"
	"time"
)

func main() {

	dsn := "go_admin:123456@tcp(127.0.0.1:13306)/go_demo?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	fmt.Println(db, err)
	sqlDB, err := db.DB()

	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(10 * time.Second)

	// 结构体
	type List struct {
		gorm.Model
		Name    string `gorm:"type:varchar(20); not null" json:"name" binding:"required"`
		State   string `gorm:"type:varchar(20); not null" json:"state" binding:"required"`
		Phone   string `gorm:"type:varchar(20); not null" json:"phone" binding:"required"`
		Email   string `gorm:"type:varchar(40); not null" json:"email" binding:"required"`
		Address string `gorm:"type:varchar(200); not null" json:"address" binding:"required"`
	}
	db.AutoMigrate(&List{})

	r := gin.Default()
	// ADD
	r.POST("/user/add", func(c *gin.Context) {
		var data List
		err := c.ShouldBindJSON(&data)
		if err != nil {
			c.JSON(200, gin.H{
				"msg":  "添加失败",
				"data": data,
				"code": 400,
			})
		} else {
			// 创建一条数据
			db.Create(&data)
			c.JSON(200, gin.H{
				"msg":  "创建成功",
				"data": data,
				"code": 200,
			})
		}
	})

	// DELETE
	r.DELETE("/user/delete/:id", func(c *gin.Context) {
		var data []List
		// 接收id
		id := c.Param("id")

		// 判断id是否存在
		db.Where("id = ?", id).Find(&data)

		// id存在就删除，不存在就报错
		if len(data) == 0 {
			c.JSON(200, gin.H{
				"msg":  "id不存在",
				"code": 400,
			})
		} else {
			db.Where("id = ?", id).Delete(&data)
			c.JSON(200, gin.H{
				"msg":  "删除成功",
				"code": 200,
			})
		}
	})

	// 改
	r.PUT("/user/update/:id", func(c *gin.Context) {
		var data List
		// 接收id
		id := c.Param("id")

		// 查找id
		db.Select("id").Where("id = ?", id).Find(&data)

		// 判断id是否存在
		if data.ID == 0 {
			c.JSON(200, gin.H{
				"msg":  "用户不存在",
				"code": 400,
			})
		} else {
			err := c.ShouldBindJSON(&data)
			if err != nil {
				c.JSON(200, gin.H{
					"mse":  "修改失败",
					"code": 400,
				})
			} else {
				// 修改数据库内容
				db.Where("id = ?", id).Updates(&data)
				c.JSON(200, gin.H{
					"msg":  "修改成功",
					"code": 200,
				})
			}
		}
	})

	// 查询
	// 条件查询
	r.GET("/user/list/:name", func(c *gin.Context) {
		// 获取路径参数
		name := c.Param("name")
		var dataList []List

		// 查询数据库
		db.Where("name = ?", name).Find(&dataList)
		// 判断是否存在数据
		if len(dataList) == 0 {
			c.JSON(200, gin.H{
				"msg":  "没有查询到数据",
				"code": 400,
				"data": gin.H{},
			})
		} else {
			c.JSON(200, gin.H{
				"msg":  "查询成功",
				"code": 200,
				"data": dataList,
			})
		}

	})
	// 全部查询
	r.GET("user/list", func(c *gin.Context) {
		var dataList []List

		pageNum, _ := strconv.Atoi(c.Query("pageNum")) // 转化成int
		pageSize, _ := strconv.Atoi(c.Query("pageSize"))

		// 判断是否需要分页
		if pageSize == 0 {
			pageSize = -1
		}

		if pageNum == 0 {
			pageNum = -1
		}

		offsetVal := (pageNum - 1) * pageSize // 固定写法
		if pageNum == -1 && pageSize == -1 {
			offsetVal = -1
		}

		// 查询全部数据，查询分页数据
		var total int64
		db.Model(dataList).Count(&total).Limit(pageSize).Offset(offsetVal).Find(&dataList)
		if len(dataList) == 0 {
			c.JSON(200, gin.H{
				"msg":  "没有查询到数据",
				"code": 400,
				"data": gin.H{},
			})
		} else {
			c.JSON(200, gin.H{
				"msg":  "查询成功",
				"code": 200,
				"data": gin.H{
					"list":     dataList,
					"total":    total,
					"pageNum":  pageNum,
					"pageSize": pageSize,
				},
			})
		}
	})
	r.Run(": 8082") // listen and serve on 0.0.0.0:8080
}
