## Mail-Reminder  
<a href="https://996.icu"><img src="https://img.shields.io/badge/link-996.icu-red.svg" alt="996.icu" /></a>
[![LICENSE](https://img.shields.io/badge/license-Anti%20996-blue.svg)](https://github.com/996icu/996.ICU/blob/master/LICENSE)  
A simple mail reminder.  
#### How to use?
+ **Get Mail-Reminder server**
  1. find it in ```/bin```
  2. or find it in ```release```
  3. or use ```go build``` to build ```/src/server/reminder_server.go```

+ **Make dir like this**
  + ┏ /bin
  + ┃ &nbsp;&nbsp;&nbsp;┗ reminder_server
  + ┗ /resource  
  + &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;┗ config.xml
    
+ **Start server**
  + ```./reminder_server```
  + or
  + ```nohup ./reminder_server```
  
+ **Use cli to connect server**
  + ```./reminder_cli -h``` _to get help_.
  + ``` -ip xxx.xxx.xxx.xxx``` _to set host ip_.
  + ```-p xxxx``` _to set port_.
  + ```-new``` _to add a new timed message_.
  + ```-del``` _to del a existing message_.
  + ```-ls``` _to list all messages_.
    + For example: ```./reminder_cli -ip 127.0.0.1 -p 8447 -new```
    + Or use default ip and port: ```./reminder_cli -new```
    + Default ```ip``` = ```localhost```, default ```port``` = ```8447```
    
+ **Config File**
  + ```reminder_server``` cannot get started without ```/resource/config.xml```  
  + **Template of config.xml:**
  + ```
      <config>
           <!--sender-->
           <addr>scu_sjl@outlook.com</addr>
           
           <!--auth code or password of sender-->
           <auth>password</auth>
    
           <!--smtp host of sender-->
           <host>smtp.office365.com</host>
    
           <!--port of smtp host-->
           <port>465</port>
    
           <!--max retry times when mail sending fails-->
           <retry>3</retry>
           
           <!--port for reminder_server-->
           <listen>8447</listen>
       </config>
    ```

+ **PS**
  + if you are using ```QQ-Email``` as sender, then you should enter ```auth code``` in ```<auth>``` instead of your QQ password.
  + 如果你用QQ邮箱来作为邮件的发送者，在```<auth>```标签中你应该填写QQ邮箱的授权码，而不是你的QQ密码或者QQ邮箱密码
  + I recommend you to use English for input, otherwise encoded utf-8 characters may appear.
  + 推荐你使用英文来输入信息，否则可能会出现被编码过的 utf-8 字符串.
