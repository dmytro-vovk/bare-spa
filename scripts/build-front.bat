@echo off

call   "C:\Program Files\nodejs\npm" install

call   "C:\Program Files\nodejs\npm" run build-guest
call   "C:\Program Files\nodejs\npm" run build-user

@del   /S /Q "%CD%\internal\webserver\handlers\home\css"

@mkdir ""%CD%\internal\webserver\handlers\home\css"
@copy  "%CD%\frontend\styles\*" "%CD%\internal\webserver\handlers\home\css\"

7z a -tgzip "%CD%\internal\webserver\handlers\home\index.html.gz" "%CD%\frontend\index.html"

powershell -Command "(Get-Content %CD%\internal\webserver\handlers\home\guest.js) -replace '=guest.js.map','=js.map' | Set-Content %CD%\internal\webserver\handlers\home\guest.js"
7z a -tgzip "%CD%\internal\webserver\handlers\home\guest.js.gz" "%CD%\internal\webserver\handlers\home\guest.js"
7z a -tgzip "%CD%\internal\webserver\handlers\home\guest.js.map.gz" "%CD%\internal\webserver\handlers\home\guest.js.map"

powershell -Command "(Get-Content %CD%\internal\webserver\handlers\home\user.js) -replace '=user.js.map','=js.map' | Set-Content %CD%\internal\webserver\handlers\home\user.js"
7z a -tgzip "%CD%\internal\webserver\handlers\home\user.js.gz" "%CD%\internal\webserver\handlers\home\user.js"
7z a -tgzip "%CD%\internal\webserver\handlers\home\user.js.map.gz" "%CD%\internal\webserver\handlers\home\user.js.map"
