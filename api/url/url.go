package url

import (
	"context"
	"fmt"

	"log"
	"net/http"
	"shorturl/kitex_gen/short/url"
	"shorturl/kitex_gen/short/url/urlservice"
	"shorturl/middleware"
	"shorturl/model"
	"sort"

	"github.com/cloudwego/kitex/client"
	"github.com/gin-gonic/gin"
)


func Add(c *gin.Context)  {
	req:=url.NewUrlRequest()
	req.LongUrl=c.PostForm("longurl")
	req.ShortCode=c.PostForm("shortcode")

	//在reids查找是否存在
    val, err := model.RedisClient.Get(c, req.LongUrl).Result()
    if err == nil {
		c.JSON(http.StatusOK,gin.H{
			"message":"已存在",
			"longUrl":req.LongUrl,
			"shortCode":val,
		})
		return
	}

	//进入数据库查找
	cli,err:=urlservice.NewClient("short.url",	client.WithHostPorts("0.0.0.0:8889"))
	if err!=nil{
		log.Fatal(err)
	}
	resp,_:=cli.AddUrl(context.Background(),req)

	//添加创建人
	if resp.Massage=="生成成功"{
		u:=model.Url{}
		model.DB.Where("long_url=?",req.LongUrl).First(&u)
		model.DB.Model(&u).Select("created_user").Update(model.Url{CreatedUser: middleware.UserName})

		//加入缓存
		err = model.RedisClient.Set(c, resp.ShortCode, resp.LongUrl, 0).Err()
    if err != nil {
        log.Println("Error setting shortcode Redis cache:", err)
	}
}

	c.JSON(http.StatusOK,gin.H{
		"message":resp.Massage,
		"longUrl":resp.LongUrl,
		"shortCode":resp.ShortCode,
	})
}







func Update(c *gin.Context)  {
	req:=url.NewUrlRequest()
	req.LongUrl=c.PostForm("longurl")
	req.ShortCode=c.PostForm("shortcode")

	cli,err:=urlservice.NewClient("short.url",	client.WithHostPorts("0.0.0.0:8889"))
	if err!=nil{
		log.Fatal(err)
	}
	resp,_:=cli.Update(context.Background(),req)

	if resp.Massage=="更新成功"{
		//改变缓存中该数据
		val, err :=model.RedisClient.Exists(c,req.LongUrl).Result()          
		if err != nil {  
		log.Fatalf("Failed to check key existence: %v", err)  
		}  
		
  		if val > 0 {  
		_, err = model.RedisClient.Set(c, req.LongUrl,resp.ShortCode,0).Result()  
		if err != nil {  
			log.Fatalf("Failed to update key: %v", err)  
		}  
		}
	}
	
	c.JSON(http.StatusOK,gin.H{
		"message":resp.Massage,
		"longUrl":resp.LongUrl,
		"shortCode":resp.ShortCode,
	})
}


func Search(c *gin.Context)  {
	req:=url.NewUrlRequest()
	req.LongUrl=c.PostForm("longurl")

	cli,err:=urlservice.NewClient("short.url",	client.WithHostPorts("0.0.0.0:8889"))
	if err!=nil{
		log.Fatal(err)
	}
	resp,_:=cli.Sreach(context.Background(),req)


	//加入redis缓存
	if resp.Massage=="对应短链存在"{

		err = model.RedisClient.Set(c, resp.ShortCode, resp.LongUrl, 0).Err()
    if err != nil {
        log.Println("Error setting shortcode Redis cache:", err)
    }	
	}
	
	c.JSON(http.StatusOK,gin.H{
		"message":resp.Massage,
		"longUrl":resp.LongUrl,
		"shortCode":resp.ShortCode,
	})
}





func Delete(c *gin.Context)  {
	req:=url.NewUrlRequest()
	req.LongUrl=c.PostForm("longurl")
	req.ShortCode=c.PostForm("shortcode")
	

	cli,err:=urlservice.NewClient("short.url",	client.WithHostPorts("0.0.0.0:8889"))
	if err!=nil{
		log.Fatal(err)
	}
	resp,_:=cli.Delete(context.Background(),req)

	if resp.Massage=="删除成功"{
		//清除缓存中该数据
		val, err :=model.RedisClient.Exists(c,req.LongUrl).Result()          
		if err != nil {  
		log.Fatalf("Failed to check key existence: %v", err)  
		}  
		
  		if val > 0 {  
		_, err = model.RedisClient.Del(c, req.LongUrl).Result()  
		if err != nil {  
			log.Fatalf("Failed to delete key: %v", err)  
		}  
		}
	}

	
	c.JSON(http.StatusOK,gin.H{
		"message":resp.Massage,
		"longUrl":resp.LongUrl,
		"shortCode":resp.ShortCode,
	})
	
}




//重定向
func RedirectHandler(c *gin.Context)  {
	shortCode:=c.Param("shortCode")[1:]
	var longUrl string
	var u model.Url
	flag:=false

	//redis
	val, err:= model.RedisClient.Get(c, shortCode).Result()
	fmt.Println(shortCode,val)
    if err==nil {
		fmt.Println("redis找到")
		longUrl=val
    }else{
		result:=model.DB.Where("short_code=?",shortCode).First(&u)

		if result.RecordNotFound(){
		c.JSON(http.StatusOK,gin.H{
			"message":"未找到记录",
		})
		return 
	}
	flag=true
	longUrl=u.LongUrl
	}
	
	c.Redirect(http.StatusMovedPermanently,longUrl)

	if(!flag){
		result:=model.DB.Where("short_code=?",shortCode).First(&u)
		if result.RecordNotFound(){
		c.JSON(http.StatusOK,gin.H{
			"message":"未找到记录",
		})
		return 
	}
	}
	visited:=u.Visited+1
	model.DB.Model(&u).Select("visited").Updates(model.Url{Visited: visited})	//访问次数更新
	
}






type Result struct {
    LongURL   string
    ShortCode string
    Visited   int
}


func Rank(c *gin.Context)  {
	
	var results []Result
	res:=model.DB.Table("urls").Where("created_user=?",middleware.UserName).Select("long_url,short_code,visited").Find(&results)
	if res.RowsAffected==0{
		c.JSON(http.StatusOK,gin.H{
			"message":"还未生成短链",
		})
		return
	}

	sort.Sort(sort.Reverse(ByVisited(results)))			//对访问次数降序排序
	c.JSON(http.StatusOK,gin.H{
		"message":results,
	})

}