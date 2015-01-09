@echo off 
fw -i=XevDemo0
IF ERRORLEVEL 1 GOTO err
fw -i=XevDemo1
IF ERRORLEVEL 1 GOTO err
fw -i=XevDemo2
IF ERRORLEVEL 1 GOTO err
fw -i=XevDemo3
IF ERRORLEVEL 1 GOTO err
fw -i=XevDemo4
IF ERRORLEVEL 1 GOTO err
fw -i=XevDemo5
IF ERRORLEVEL 1 GOTO err
fw -i=XevDemo6
IF ERRORLEVEL 1 GOTO err
fw -i=XevDemo7
IF ERRORLEVEL 1 GOTO err
fw -i=XevDemo8
IF ERRORLEVEL 1 GOTO err
fw -i=XevDemo9
IF ERRORLEVEL 1 GOTO err
fw -i=XevDemo10
IF ERRORLEVEL 1 GOTO err
fw -i=XevDemo11
IF ERRORLEVEL 1 GOTO err
fw -i=XevDemo12
IF ERRORLEVEL 1 GOTO err
fw -i=XevDemo13
IF ERRORLEVEL 1 GOTO err
fw -i=XevDemo14
IF ERRORLEVEL 1 GOTO err

GOTO ok
:err
echo FAILED
pause

:ok