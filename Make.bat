@echo off

::pid.kirichok.name
::SET IP="94.158.81.107"

:: kirichok.name
SET IP="1.1.1.1"

SET PWD="LetMeIn321"

if "%1"=="clean"       Call :%1
if "%1"=="lint-front"  Call :%1
if "%1"=="lint-back"   Call :%1
if "%1"=="lint-fix"    Call :%1
if "%1"=="test-back"   Call :%1
if "%1"=="build-front" Call :%1
if "%1"=="build"       Call :%1
if "%1"=="assemble"    Call :%1
if "%1"=="run"         Call :%1
if "%1"=="run-short"   Call :%1
if "%1"=="deploy"      Call :%1
if "%1"=="deploy-config" Call :%1
if "%1"=="all"         Call :%1
exit 0

:clean
    if exist webserver.exe (
        @del   webserver.exe
    )
    if exist webserver (
            @del   webserver
    )
goto :EOF

:lint-fix
    call "C:\Program Files\nodejs\npm" run lint -- --fix

:lint-front
	call "C:\Program Files\nodejs\npm" install
	call "C:\Program Files\nodejs\npm" run lint
goto :EOF

:build-front
    call   :clean
    call   "C:\Program Files\nodejs\npm" run build
    @del   /S /Q "%CD%\internal\webserver\handlers\home\css"
    @mkdir ""%CD%\internal\webserver\handlers\home\css"
    @copy  "%CD%\frontend\styles\*" "%CD%\internal\webserver\handlers\home\css\"
    7z a -tgzip "%CD%\internal\webserver\handlers\home\index.html.gz" "%CD%\frontend\index.html"
    7z a -tgzip "%CD%\internal\webserver\handlers\home\index.js.gz" "%CD%\internal\webserver\handlers\home\index.js"
    7z a -tgzip "%CD%\internal\webserver\handlers\home\index.js.map.gz" "%CD%\internal\webserver\handlers\home\index.js.map"
goto :EOF

:lint-back
	call golangci-lint run
goto :EOF

:test-back
	call go test -v -tags=exec "%CD%\internal\..."
goto :EOF

:build-back
	SET GOOS=linux
	SET GOARCH=mipsle
	@go build -ldflags "-s -w" -o webserver "%CD%\cmd\main.go"

	SET GOOS=windows
    SET GOARCH=amd64
    @go build -ldflags "-s -w" -o webserver.exe "%CD%\cmd\main.go"
goto :EOF

:assemble
    call :build-front
    call :build-back
goto :EOF

:run
    call :lint-front
    call :build-front
    call :lint-back
    call :test-back
	@go build -o webserver.exe "%CD%\cmd\main.go"
	webserver.exe
goto :EOF

:run-short
    call :build-front
    call :lint-back
    call :test-back
	@go run "%CD%\cmd\main.go"
goto :EOF

:deploy
    call :assemble
    @echo "Deploying to remote..."
    pscp.exe -4 -P 22 -l root -pw %PWD%  -r .\\webserver %IP%:/root/webserver.new
    @echo "Restarting webserver..."
    klink -4 -P 22 -pw %PWD% -l root  %IP% "/etc/init.d/webserver stop; chmod +x webserver.new; mv webserver.new webserver; /etc/init.d/webserver start"
    @echo "Done"
goto :EOF

:deploy-config:
    pscp.exe -4 -P 22 -l root -pw %PWD%  -r .\\config-prod.json %IP%:/root/config.json
	klink -4 -P 22 -pw %PWD% -l root  %IP% "/etc/init.d/webserver restart"
goto :EOF

:all
    call :lint-front
    call :build-front
    call :lint-back
    call :test-back
    call :build-back
    call :deploy
    call :clean
goto :EOF

