# go-get

It is a simple application for download files via http or https.

## Interface

```
get [-dns <IP>] [-insecure] [-url <URL>] [-h]
Usage of ./get:
  -dns string
    	set dns ip address, default use system dns resolver
  -insecure
    	set true for ignore tls host check
  -url string
    	url to download (default "https://ftp.mgts.by/test/100Mb")
```