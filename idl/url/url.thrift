namespace go short.url

struct UrlRequest{
    1:string LongUrl
    2:string ShortCode
}


struct UrlResponse{
    1:string LongUrl 
    2:string ShortCode
    3:string massage
}


service UrlService{
    UrlResponse AddUrl(1:UrlRequest req)   
    UrlResponse Update(1:UrlRequest req)
    UrlResponse Delete(1:UrlRequest req)
    UrlResponse Sreach(1:UrlRequest req)
}