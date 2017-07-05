#bvpn
由于vpn的密码是由用户密码+google验证码组成的，所以每次网络不稳定的时候，都要重新输入密码，十分麻烦，所以bvpn对 openvpn 和 google 验证码封装，在vpn断线的时候自动更新google验证码，然后重新连接vpn
##安装
```
go get github.com/anoty/bvpn
```
##配置
```
ovpn = "openvpn命令行程序路径"
cfg = "openvpn配置文件路径"
pass = "openvpn密码文件路径"
username = "vpn账户"
password="vpn密码"
secretKey = "google 验证码 secret，不是手机上的那6个数字，是在第一次初始化 google 验证码的时候二维码对应的 secret"
```
##使用
bvpn需要管理员权限
```
bvpn -conf=bvpn.toml
```
