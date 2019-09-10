### 百度爬虫
目录架构  
spider 爬虫代码目录  
static 爬虫任务分配和爬虫结果目录  
static-job 爬虫任务目录，文件需以job.txt结尾  
static-result 爬虫结果目录

爬取速度限制不在io,而在于网络，可根据实际情况调整  
spider.go中Schedule方法爬虫并发数量

运行`go run app.go`  
