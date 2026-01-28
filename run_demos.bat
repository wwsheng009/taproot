@echo off
echo ================================
echo Taproot TUI Framework Demo Test
echo ================================
echo.
echo This will launch three demo programs sequentially.
echo Press q or ctrl+c to exit each demo.
echo.
pause
echo.
echo [1/3] Running basic demo...
echo.
bin\demo.exe
echo.
echo [2/3] Running list demo...
echo.
bin\list.exe
echo.
echo [3/3] Running app demo (full framework)...
echo.
bin\app.exe
echo.
echo All demos completed!
pause
