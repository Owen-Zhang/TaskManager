@echo off

:: �˳���Ὣ����ŵ�Ŀ�������һ��λ�ã�Ȼ����뵽service�У������Զ�����

::ǿ����ʾ�û��޸�conf/�ļ��е�������Ϣ
echo �����޸�taskmanager/conf/app.conf�ļ��е�������Ϣ ����:Y�� �˳���N
Set /p choice=��ѡ��: 
if /i %choice%==N (
    exit 1
)

::���ļ��ŵ�D���У����û��D�̾ͷ���C��
set defaultdisc=D:
if not exist %defaultdisc% (
    %defaultdisc%=C:
)
set destPath=%defaultdisc%\taskmanager\
xcopy /y /e /i taskmanager %destPath%
echo ��ص��ļ��Ѹ��Ƶ���%destPath%Ŀ¼��, ��ص�������ϢҲ�����������޸�

::����service�У����ҿ�������
sc query taskmanager > NUL
if %errorlevel% NEQ 1060 (
    sc delete "taskmanager"
) 

::ipsecpol -p myfw -r dwmrc_pass_me -f *+0:8000:tcp -n PASS -w reg -x

sc create taskmanager binPath=%destPath%\home.exe start=auto
sc start taskmanager
pause