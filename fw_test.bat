@echo off 
fw -i=%1
IF ERRORLEVEL 1 GOTO err

GOTO ok
:err
echo FAILED %1
pause

:ok
exit
