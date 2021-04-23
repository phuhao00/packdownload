package main

import "packdownload/internal/biz"

func main()  {
	mgr:=biz.NewPackDownLoadMgr()
	mgr.Run()
}
