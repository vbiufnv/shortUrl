package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"shorturl/kitex_gen/short/url"
	"shorturl/model"
	"strconv"
	"time"
)

// 哈希算法
func ShortUrl(url string) []string {
	key := "miku"

	chars := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z", "0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

	//MD5加密
	md5Hash := md5.New()
	md5Hash.Write([]byte(key + url))
	hashResult := hex.EncodeToString(md5Hash.Sum(nil))

	
	resUrl := make([]string, 4)


	for i := 0; i < 4; i++ {
		sTempSubString := hashResult[i*8 : i*8+8]
		//转换为整数
		lHexLong, _ := strconv.ParseInt(sTempSubString, 16, 64)
		lHexLong = 0x3FFFFFFF & lHexLong	//30位
		outChars := ""


		//每个短码6位
		for j := 0; j < 6; j++ {
			index := 0x0000003D & lHexLong	//0~61	
			outChars += chars[index]
			lHexLong = lHexLong >> 5
		}
		resUrl[i] = outChars	
	}
	return resUrl
}



// UrlServiceImpl implements the last service interface defined in the IDL.
type UrlServiceImpl struct{}


// AddUrl implements the UrlServiceImpl interface.
func (s *UrlServiceImpl) AddUrl(ctx context.Context, req *url.UrlRequest) (resp *url.UrlResponse, err error) {
	// TODO: Your code here...
	longUrl := req.LongUrl
	resp = url.NewUrlResponse()
	resp.LongUrl = longUrl

	//检查mysql
	url := model.Url{}
	result := model.DB.Where("long_url=?", longUrl).First(&url)
	if result.RowsAffected == 1 {
		resp.Massage = "该URL已存在"
		resp.ShortCode = url.ShortCode
		return resp, nil
	}

	shortUrls := ShortUrl(longUrl)

	//随机
	rand.Seed(time.Now().UnixNano())
	randomNumber := rand.Intn(4)
	shortUrl := shortUrls[randomNumber]

	resp.ShortCode = shortUrl


	//存入数据库
	u:= model.Url{LongUrl:longUrl, ShortCode: shortUrl,Visited: 0}
	if err:=model.DB.Create(&u).Error;err!=nil{
		resp.Massage = "生成失败"
		return resp, err
	}else{
		resp.Massage = "生成成功"
	return resp, nil
	}

}



// Update implements the UrlServiceImpl interface.
func (s *UrlServiceImpl) Update(ctx context.Context, req *url.UrlRequest) (resp *url.UrlResponse, err error) {
	// TODO: Your code here...

	longUrl := req.LongUrl
	resp = url.NewUrlResponse()
	resp.LongUrl = longUrl


	//检查mysql
	url := model.Url{}
	result := model.DB.Where("long_url=?", longUrl).First(&url)
	if result.RecordNotFound() {
		resp.Massage = "该长链不存在"
		resp.ShortCode = ""
		return resp, nil
	}

	shortUrls := ShortUrl(longUrl)

	//随机取
	rand.Seed(time.Now().UnixNano())
	randomNumber := rand.Intn(4)
	shortUrl := shortUrls[randomNumber]
	resp.ShortCode = shortUrl

	//取不同的短链
	if shortUrl != req.ShortCode {
		result := model.DB.Model(&url).Select("short_code").Update(model.Url{LongUrl: req.LongUrl, ShortCode: shortUrl})
		if result.Error != nil {
			resp.Massage = "更新失败"
			return resp,err
		}
	}

	resp.Massage = "更新成功"
	return resp, nil

}







// Delete implements the UrlServiceImpl interface.
func (s *UrlServiceImpl) Delete(ctx context.Context, req *url.UrlRequest) (resp *url.UrlResponse, err error) {
	// TODO: Your code here...
	longUrl := req.LongUrl
	resp = url.NewUrlResponse()
	resp.LongUrl = longUrl


	//检查存在
	url := model.Url{}
	result := model.DB.Where("long_url=?", longUrl).First(&url)
	if result.RecordNotFound() {
		resp.Massage = "该长链不存在"
		resp.ShortCode = ""
		return resp, nil
	}

	//删除
	if err := model.DB.Where("long_url=?", longUrl).Unscoped().Delete(&url).Error; err != nil {
		resp.Massage = "删除失败"
		resp.ShortCode = ""
		return resp, nil
	}

	resp.Massage = "删除成功"
	return resp, nil
}



// Sreach implements the UrlServiceImpl interface.
func (s *UrlServiceImpl) Sreach(ctx context.Context, req *url.UrlRequest) (resp *url.UrlResponse, err error) {
	// TODO: Your code here...
	longUrl := req.LongUrl
	resp = url.NewUrlResponse()
	resp.LongUrl = longUrl

	//检查存在
	u := model.Url{}
	result := model.DB.Where("long_url=?", longUrl).First(&u)
	
	//不存在
	if result.RecordNotFound() {
		resp.Massage = "该长链还未生成"
		resp.ShortCode = ""
		return resp, nil
	}

	resp.Massage = "对应短链存在"
	resp.ShortCode = u.ShortCode
	return resp, nil

}




func (s *UrlServiceImpl) RedictUrl(ctx context.Context, req *url.UrlRequest) (resp *url.UrlResponse, err error) {
	shortCode := req.ShortCode
	resp = url.NewUrlResponse()
	resp.ShortCode = shortCode

	u := model.Url{}
	result := model.DB.Where("short_code=?", shortCode).First(&u)
	//不存在
	if result.RecordNotFound() {
		resp.Massage = "该短链不存在"
		resp.ShortCode = ""
		return resp, nil
	}

	return resp, nil
}


