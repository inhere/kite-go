@REM
@REM use like .bashrc
@REM

@echo off

@REM add cmd alias
@REM refer https://www.cnblogs.com/mq0036/p/16255494.html
doskey ls=dir /b $*
doskey ll=dir $*

doskey pwd=cd
doskey cat=type $*
doskey rm=del $*
doskey mv=move $*
doskey cd=cd /d $*
doskey mkdir=md $*

doskey clear=cls
doskey history=doskey /history
doskey alias=doskey /macros

doskey traceroute=tracert $*
doskey tracepath=pathping $*
doskey ifconfig=ipconfig $*
doskey shell=PowerShell $*

echo The alias is loaded, input 'alias' to view