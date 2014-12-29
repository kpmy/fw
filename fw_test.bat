@echo off 
fw -i=XevDemo0
IF ERRORLEVEL 1 GOTO err
fw -i=XevDemo1
fw -i=XevDemo2
fw -i=XevDemo3
fw -i=XevDemo4
fw -i=XevDemo5
fw -i=XevDemo6
fw -i=XevDemo7
fw -i=XevDemo8

GOTO ok
:err
echo FAILED
pause

:ok