@echo off 
fw -i=XevDemo0 > .out
IF ERRORLEVEL 1 GOTO err
fw -i=XevDemo1 > .out
IF ERRORLEVEL 1 GOTO err
fw -i=XevDemo2 > .out
IF ERRORLEVEL 1 GOTO err
fw -i=XevDemo3 > .out
IF ERRORLEVEL 1 GOTO err
fw -i=XevDemo4 > .out
IF ERRORLEVEL 1 GOTO err
fw -i=XevDemo5 > .out
IF ERRORLEVEL 1 GOTO err
fw -i=XevDemo6 > .out
IF ERRORLEVEL 1 GOTO err
fw -i=XevDemo7 > .out
IF ERRORLEVEL 1 GOTO err
fw -i=XevDemo8 > .out
IF ERRORLEVEL 1 GOTO err
fw -i=XevDemo9 > .out
IF ERRORLEVEL 1 GOTO err
fw -i=XevDemo10 > .out
IF ERRORLEVEL 1 GOTO err
fw -i=XevDemo11 > .out
IF ERRORLEVEL 1 GOTO err
fw -i=XevDemo12 > .out
IF ERRORLEVEL 1 GOTO err
fw -i=XevDemo13 > .out
IF ERRORLEVEL 1 GOTO err
fw -i=XevDemo14 > .out
IF ERRORLEVEL 1 GOTO err
fw -i=XevDemo15 > .out
IF ERRORLEVEL 1 GOTO err
fw -i=XevDemo16 > .out
IF ERRORLEVEL 1 GOTO err
fw -i=XevDemo17 > .out
IF ERRORLEVEL 1 GOTO err
fw -i=XevDemo18 > .out
IF ERRORLEVEL 1 GOTO err
fw -i=XevDemo19 > .out
IF ERRORLEVEL 1 GOTO err
fw -i=XevDemo20 > .out
IF ERRORLEVEL 1 GOTO err
fw -i=XevDemo21 > .out
IF ERRORLEVEL 1 GOTO err
fw -i=XevDemo22 > .out
IF ERRORLEVEL 1 GOTO err
fw -i=XevTest0 > .out
IF ERRORLEVEL 1 GOTO err

GOTO ok
:err
echo FAILED
pause

:ok
del .out