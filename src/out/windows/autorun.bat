@echo off

:: 此程序会将代码放到目标机器的一定位置，然后加入到service中，开机自动运行

::强行提示用户修改conf/文件中的配制信息
echo 请先修改taskmanager/conf/app.conf文件中的配置信息 继续:Y， 退出：N
Set /p choice=请选择: 
if /i %choice%==N (
    exit 1
)

::将文件放到D盘中，如果没有D盘就放入C盘
set defaultdisc=D:
if not exist %defaultdisc% (
    %defaultdisc%=C:
)
set destPath=%defaultdisc%\taskmanager\
xcopy /y /e /i taskmanager %destPath%
echo 相关的文件已复制到了%destPath%目录下, 相关的配制信息也可以在里面修改

::加入service中，并且开机运行
sc query taskmanager > NUL
if %errorlevel% NEQ 1060 (
    sc delete "taskmanager"
) 

::ipsecpol -p myfw -r dwmrc_pass_me -f *+0:8000:tcp -n PASS -w reg -x

sc create taskmanager binPath=%destPath%\home.exe start=auto
sc start taskmanager
pause