@echo off

:: �������ļ�
go install home

if %errorlevel% NEQ 0 (
    echo ����������,��鿴��Ӧ����
    exit 1
)

set tempfolder= out\windows\taskmanager
if exist tempfolder (
    rd /s /Q tempfolder
)
md tempfolder

::������Ӧ���ļ������ļ��е�temp�ļ�����
xcopy /y ..\bin\home.exe %tempfolder%\
xcopy /y /e /i conf %tempfolder%\conf
xcopy /y /e /i static %tempfolder%\static
xcopy /y /e /i views %tempfolder%\views

::���������ļ��У�һ����¼�������ռǣ�һ�����ڴ���û��ϴ����ļ�
md %tempfolder%\Data\Run
md %tempfolder%\Data\Temp
md %tempfolder%\logs

echo �����������, ��out/windows����, �뽫taskmanager�ļ��к�autorun.bat�ļ�һ��copy��Ҫ���еĻ�����
echo Ȼ������autorun.bat�ļ�,�벻Ҫ�����ļ��е�����

pause