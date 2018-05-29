@echo off

:: 编译主文件
go install home

if %errorlevel% NEQ 0 (
    echo 编译代码出错,请查看相应问题
    exit 1
)

set tempfolder= out\windows\taskmanager
if exist tempfolder (
    rd /s /Q tempfolder
)
md tempfolder

::复制相应的文件或者文件夹到temp文件夹中
xcopy /y ..\bin\home.exe %tempfolder%\
xcopy /y /e /i conf %tempfolder%\conf
xcopy /y /e /i static %tempfolder%\static
xcopy /y /e /i views %tempfolder%\views

::创建两个文件夹，一个记录本程序日记，一个用于存放用户上传的文件
md %tempfolder%\Data\Run
md %tempfolder%\Data\Temp
md %tempfolder%\logs

echo 数据生成完成, 在out/windows下面, 请将taskmanager文件夹和autorun.bat文件一起copy到要运行的机器上
echo 然后运行autorun.bat文件,请不要随便改文件夹的名字

pause