package ui

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"io"
	"k8stool/util"
	"strings"
	"fmt"
	"strconv"
	"bufio"
)
var OutTE *walk.TextEdit
var Username,Password,Ip,Port *walk.LineEdit
var UserLab,PasswdLab,IpLab *walk.Label
var Login,Kubetool,Control *walk.GroupBox
var CmdReader io.Reader
var NameSpace []string
var NameSpaceSelect,PodsSelect,LogPath *walk.ComboBox
var shell_results string
var serr error
var done chan bool
var monitorbtn,cannelbtn  *walk.PushButton

func Body() {

	MainWindow{
		Icon:    Bind("'/img/logo.ico'"),
		Title:   "K8Stools 容器日志监视器",
		MinSize: Size{400, 200},
		Size: Size{1200, 800},
		Layout:  VBox{},
		Children: []Widget{
			Composite{
				Layout: VBox{},
				Children: []Widget{
					//VSpacer{Size:100},
					GroupBox{   //login ui start
						AssignTo: &Login,
						//StretchFactor: 1300,
						Title:      "第一步:建立SSH链接，刷新namespace、pods选项",
						Layout: Grid{Columns: 10},
						Children: []Widget{
							Label{
								AssignTo: &IpLab,
								Text: "SSH IP:Port",
							},
							LineEdit{
								AssignTo: &Ip,
								MaxSize:Size{100,0},
								Text: "192.168.2.137",
								//MinSize: Size{60,30},
							},
							Label{
								AssignTo: &IpLab,
								Text: ":",
							},
							LineEdit{
								AssignTo: &Port,
								MaxSize:Size{30,0},
								Text: "22",
								//MinSize: Size{60,30},
							},
							Label{
								AssignTo: &UserLab,
								Text: "用户名:",
							},
							LineEdit{
								AssignTo: &Username,
								Text: "root",
							},
							Label{
								AssignTo: &PasswdLab,
								Text: "密码:",
							},
							LineEdit{
								AssignTo: &Password,
								Text: "yunwei123",
								PasswordMode: true,
							},
							PushButton{
								Text:"确认链接",
								OnClicked:  func() {
									cannelbtn.SetEnabled(false)
									NameSpaceSelect.SetModel(nil)
									PodsSelect.SetModel(nil)
									cli := util.NewCli(Ip.Text(),Username.Text(),Password.Text(),Port.Text())
									shell_results, serr = cli.RunShell("kubectl get namespace | grep -v kube | grep -v monitoring | awk '{print $1}' | sed -n '2,$p'")
									if serr != nil{
										OutTE.SetText("远程链接失败")
									}else{
										shell_results = shell_results + "#^#"
										shell_results = strings.Replace(shell_results,"\n#^#","",1)
										NameSpace = util.StringTOSlicef(shell_results,"\n")
										NameSpaceSelect.SetModel(NameSpace)
										NameSpaceSelect.SetCurrentIndex(0)
										LogPath.SetCurrentIndex(0)
									}

									var NS string
									go func() {
										for {
											if NS != NameSpaceSelect.Text() {
												cli := util.NewCli(Ip.Text(), Username.Text(), Password.Text(), Port.Text())
												shell_results, serr = cli.RunShell("kubectl get pods -n " + NameSpaceSelect.Text() + " | awk '{print $1}' | sed -n '2,$p'")
												if serr != nil {
													OutTE.SetText("远程链接失败")
												} else {
													shell_results = shell_results + "#^#"
													shell_results = strings.Replace(shell_results, "\n#^#", "", 1)
													NameSpace = util.StringTOSlicef(shell_results, "\n")
													PodsSelect.SetModel(NameSpace)
													PodsSelect.SetCurrentIndex(0)
												}
												NS = NameSpaceSelect.Text()
											}
										}
									}()
								},
							},
							HSpacer{Size:400},
						},
					}, //login ui end

					GroupBox{ //kube ui start
						Title: "第二步:选择要监控的pods点【开始监控】",
						AssignTo: &Kubetool,
						//StretchFactor: 1300,
						Layout: Grid{Columns: 8},
						Children: []Widget{
							Label{
								Text: "NameSpace:",
							},
							ComboBox{
								AssignTo: &NameSpaceSelect,
								Value:    Bind("Namespace"),
							},
							Label{
								Text: "Pods:",
							},
							ComboBox{
								AssignTo: &PodsSelect,
								Value:    Bind("default"),
								//Model:    []string{"default"},
							},
							Label{
								Text: "LogPath:",
							},
							ComboBox{
								AssignTo: &LogPath,
								Value:    Bind("LogPath"),
								Model:    []string{"打印控制台日志","logs/catalina.out"},
							},
							PushButton{
								AssignTo: &monitorbtn,
								Text:"开始监控",
								OnClicked:  func() {
									done = make(chan bool)
									OutTE.SetText("")
									monitorbtn.SetEnabled(false)
									cannelbtn.SetEnabled(true)
									intport,_:= strconv.Atoi(Port.Text())
									session, err := util.SConnect(Username.Text(), Password.Text(),Ip.Text(), intport)
									go func() {
										if err != nil {
											fmt.Println(err)
										}
										defer session.Close()
										if LogPath.Text() == "打印控制台日志" {
											session.Run("kubectl logs -f " + PodsSelect.Text() + " -n " + NameSpaceSelect.Text() + " --tail=20")
											//session.Run("kubectl logs -f " + PodsSelect.Text() + " -n " + NameSpaceSelect.Text() + " tomcat-cbs-28082 --tail=20")
										} else {
											session.Run("kubectl exec -it " + PodsSelect.Text() + " -n " + NameSpaceSelect.Text() + " -- tail -f " + LogPath.Text())
										}
										session.Wait()
									}()
									CmdReader, err = session.StdoutPipe()
	;								if err != nil {
										OutTE.AppendText(fmt.Sprintln(err))
									}
									reader := bufio.NewReader(CmdReader)
									go func() {
										for {
											select {
											case <-done:
												OutTE.AppendText("\r\n\r\n========================== 断开监控成功，谢谢使用 =========================\r\n\r\n")
												monitorbtn.SetEnabled(true)
												return
											default:
												line, err2 := reader.ReadString('\n')
												if err2 != nil || io.EOF == err2 {
													break
												}
												OutTE.AppendText(line)
											}
										}
									}()

								},
							},
							PushButton{
								AssignTo: &cannelbtn,
								Text: "断开监控",
								OnClicked:  func() {
									OutTE.AppendText("\r\n\r\n========================== 正在断开监控，稍等片刻 =========================\r\n\r\n")
									close(done)
									cannelbtn.SetEnabled(false)
								},
							},
						},
					},//kube ui end

					GroupBox{

					},

					TextEdit{
						AssignTo: &OutTE,
						Persistent: true,
						ReadOnly: true,
						VScroll:  true,
						MaxLength: 9999999999,  //int max number 2147483647
						//DoubleBuffering: false,
					},
				},
			},
		},
	}.Run()
}
