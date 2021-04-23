package biz

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/skip2/go-qrcode"
	"html/template"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

type PackDownLoadMgr struct {
	Mux *http.ServeMux
}

func NewPackDownLoadMgr() PackDownLoadMgr {

	return PackDownLoadMgr{
		Mux: http.DefaultServeMux,
	}

}

func (p *PackDownLoadMgr) Run()  {
	p.RegisterHandler()
	p.waitSignal()
}

func (p *PackDownLoadMgr) waitSignal() {
	var sigChan = make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGHUP, os.Interrupt, os.Kill)
	for sig := range sigChan {
		if sig == syscall.SIGHUP {
			glog.Info("Update config")
		} else {
			glog.Infof("caught signal %v, exit2\n", sig)
			break
		}
	}
}


func (p *PackDownLoadMgr) RegisterHandler()  {
	p.HandleFunc("/download", p.downloadHandler)
	p.HandleFunc("/get_qr", p.getQRCode)
	//http://192.168.20.50:8686//download
	http.Handle("/data/",http.StripPrefix("/data/",http.FileServer(http.Dir("data"))))
	fmt.Printf("http://%s:%d", gInnerIP(), 8686)
	err := http.ListenAndServe(fmt.Sprintf(":%v",8686), nil)
	if err != nil {
		glog.Errorf("ListenAndServe :", err.Error())
	}
}


func (p *PackDownLoadMgr) HandleFunc(
	pattern string, handler func(writer http.ResponseWriter, request *http.Request),
) {
	http.HandleFunc(pattern, func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
				writer.WriteHeader(http.StatusInternalServerError)
			}
		}()
		handler(writer, request)

	})
}

type User struct {
	Name string `json:"name"`
	Id   int64  `json:"id"`
}


func gInnerIP() string {
	addresses, err := net.InterfaceAddrs()
	if err != nil {

	}
	for _, address := range addresses {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				gInnerIP := ipnet.IP.String()
				return gInnerIP
			}
		}
	}
	return ""
}

func (p *PackDownLoadMgr) downloadHandler(w http.ResponseWriter, r *http.Request) {
	_,path:=GetPath()
	p.WriteFile(path+"xfgame.png",w)
}

func (p PackDownLoadMgr) WritePage(fileName string,w http.ResponseWriter)  {
	t, _ := template.ParseFiles("views/qr.html")
	err:=t.ExecuteTemplate(w, "qr.html",nil) //第二个参数表示向模版传递的数据
	if err != nil {
		glog.Error(err.Error())
	}
}

func (p *PackDownLoadMgr) WriteFile(fileName string,w http.ResponseWriter)  {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	fileNames := url.QueryEscape(fileName) // 防止中文乱码
	w.Header().Add("Content-Type", "application/octet-stream")
	w.Header().Add("Content-Disposition", "attachment; filename=\""+fileNames+"\"")

	if err != nil {
		fmt.Println("Read File Err:", err.Error())
	} else {
		_,err=w.Write(content)
		if err != nil {
			glog.Error(err.Error())
		}
	}
}

func GetPath() ( mkpath ,basePath string) {
	mkpath="./data/images/"
	basePath="/data/images/"
	if runtime.GOOS=="windows"{
		basePath="\\data\\images\\"
	}
	base,err:=os.Getwd()
	if err != nil {
		glog.Error(err.Error())
	}
	fmt.Println(base)
	basePath=base+basePath
	return mkpath,basePath

}

func (p *PackDownLoadMgr) getQRCode(res http.ResponseWriter, req *http.Request) {
	//todo
	mkpath,basePath:=GetPath()
	Mkdir(mkpath)
	qrcodeFilename:=basePath+"xfgame.png"
	p.createQRCode(qrcodeFilename)
	//p.WriteFile(qrcodeFilename,res)
	p.WritePage("",res)
}


func (p *PackDownLoadMgr) createQRCode(filename string)  {
	//filename := "example.png"
	if err :=qrcode.WriteFile("http://192.168.20.50:8686//download", qrcode.Medium, 256, filename); err != nil {
		if err = os.Remove(filename); err != nil {
			glog.Errorf(" %s",err.Error())
		}
	}
}

func Mkdir(path string)  {
	err := os.MkdirAll(path, 0666)
	if err != nil {
		fmt.Println(err)
	}
}

